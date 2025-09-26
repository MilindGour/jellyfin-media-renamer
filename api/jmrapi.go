package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/MilindGour/jellyfin-media-renamer/config"
	"github.com/MilindGour/jellyfin-media-renamer/filesystem"
	mediainfoprovider "github.com/MilindGour/jellyfin-media-renamer/mediaInfoProvider"
	"github.com/MilindGour/jellyfin-media-renamer/middlewares"
	"github.com/MilindGour/jellyfin-media-renamer/renamer"
	"github.com/MilindGour/jellyfin-media-renamer/util"
)

type JmrAPI struct {
	// any future api keys or such will
	// appear here
	serveMux           *http.ServeMux
	configProvider     config.ConfigProvider
	fileSystemProvider filesystem.FileSystemProvider
	ren                renamer.Renamer
	mip                mediainfoprovider.MediaInfoProvider

	configResponse *ConfigResponse
	allowedExts    []string
}

func NewJmrApi(
	configProvider config.ConfigProvider,
	filesystemProvider filesystem.FileSystemProvider,
	ren renamer.Renamer,
	mip mediainfoprovider.MediaInfoProvider,
) *JmrAPI {
	jmrApi := JmrAPI{
		configProvider:     configProvider,
		fileSystemProvider: filesystemProvider,
		ren:                ren,
		mip:                mip,
	}

	return &jmrApi
}

func (j *JmrAPI) Initialize(enableCors bool) {
	j.RegisterAPIRoutes()

	if j.configProvider == nil {
		log.Fatal("Please place config before running the server")
	}

	j.populateConfig()
	j.startListenAndServe(enableCors)
}

func (j *JmrAPI) startListenAndServe(enableCors bool) {
	addr := ":" + j.configProvider.GetPort()

	var mwStack middlewares.Middleware
	if enableCors {
		mwStack = middlewares.Pipe(
			middlewares.Cors,
			middlewares.LogMW,
			middlewares.Json,
		)
	} else {
		mwStack = middlewares.Pipe(
			middlewares.LogMW,
			middlewares.Json,
		)
	}

	log.Printf("Starting JMR on port %s\n", addr)
	http.ListenAndServe(addr, mwStack(j.serveMux))
}
func (j *JmrAPI) populateConfig() {
	// get the sources by ID and cache it in memory
	config := j.configProvider.GetConfig()
	j.configResponse = NewConfigResponse(config)
}

func (j *JmrAPI) RegisterAPIRoutes() {
	j.serveMux = http.NewServeMux()

	// ping APIs
	j.serveMux.HandleFunc("GET /api/ping", j.Get_Ping())

	// select source page APIs
	j.serveMux.HandleFunc("GET /api/config", j.Get_Config())
	j.serveMux.HandleFunc("GET /api/sources", j.Get_Sources())
	j.serveMux.HandleFunc("GET /api/sources/{id}", j.Get_SourceByID())

	// identify page APIs
	j.serveMux.HandleFunc("POST /api/media/identify-names", j.Post_IdentifyNames())
	j.serveMux.HandleFunc("POST /api/media/identify-info", j.Post_IdentifyMediaInfo())

	// rename page APIs
	j.serveMux.HandleFunc("POST /api/media/rename", j.Post_Rename())
}

func (j *JmrAPI) Get_Config() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg := NewConfigResponse(j.configProvider.GetConfig())
		w.Write(ToJSON(cfg))
	}
}

func (j *JmrAPI) Get_Ping() APIHandlerFn {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("online"))
	}
}
func (j *JmrAPI) Get_Sources() APIHandlerFn {
	return func(w http.ResponseWriter, r *http.Request) {
		raw := NewSourcesResponse(j.configResponse.Source)
		w.Write(ToJSON(raw))
	}
}
func (j *JmrAPI) Get_SourceByID() APIHandlerFn {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		hasId := len(id) > 0

		if !hasId {
			j.HandleAPIError(w, r, http.StatusBadRequest, errors.New("ID is required"))
			return
		}
		idInt, err := strconv.Atoi(id)
		if err != nil {
			j.HandleAPIError(w, r, http.StatusBadRequest, err)
			return
		}
		src := util.Find(j.configResponse.Source, func(dcwi DirConfigWithID) bool {
			return dcwi.ID == idInt
		})
		if src == nil {
			j.HandleAPIError(w, r, http.StatusNotFound, errors.New("Cannot find id "+id))
			return
		}

		w.Write(ToJSON(NewSourceByIDResponse(*src, j.fileSystemProvider.ScanDirectory(src.Path, j.configProvider.GetAllowedExtensions()))))
	}
}

func (j *JmrAPI) Post_IdentifyNames() APIHandlerFn {
	return func(w http.ResponseWriter, r *http.Request) {
		var request IdentifyNameRequest
		err := json.NewDecoder(r.Body).Decode(&request)

		if err != nil {
			j.HandleAPIError(w, r, http.StatusBadRequest, err)
			return
		}
		if len(request) == 0 {
			j.HandleAPIError(w, r, http.StatusBadRequest, errors.New("Atleast 1 request object is required"))
			return
		}

		out := IdentifyMediaResponse{}
		for _, requestItem := range request {
			rawFilename := requestItem.Entry.Name
			nameAndYear := j.ren.GetMediaNameAndYear(rawFilename)
			out = append(out, IdentifyMediaResponseItem{
				SourceDirectory:      requestItem,
				IdentifiedMediaName:  nameAndYear.Name,
				IdentifiedMediaYear:  nameAndYear.Year,
				IdentifiedMediaInfos: []mediainfoprovider.MediaInfo{},
			})
		}
		w.Write(ToJSON(out))
	}
}

func (j *JmrAPI) Post_IdentifyMediaInfo() APIHandlerFn {
	return func(w http.ResponseWriter, r *http.Request) {
		var request IdentifyMediaRequest
		err := json.NewDecoder(r.Body).Decode(&request)

		if err != nil {
			j.HandleAPIError(w, r, http.StatusBadRequest, err)
			return
		}
		if len(request) == 0 {
			j.HandleAPIError(w, r, http.StatusBadRequest, errors.New("Atleast 1 request object is required"))
			return
		}

		out := IdentifyMediaResponse{}
		var wg sync.WaitGroup
		startTime := time.Now()
		for _, requestItem := range request {
			wg.Add(1)
			go func(requestItem IdentifyMediaResponseItem) {
				term := requestItem.IdentifiedMediaName
				year := requestItem.IdentifiedMediaYear
				requestItem.IdentifiedMediaInfos = j.mip.SearchMediaInfo(term, year, requestItem.SourceDirectory.Type)
				out = append(out, requestItem)
				wg.Done()
			}(requestItem)

		}
		wg.Wait()
		duration := time.Since(startTime)
		log.Printf("Fetched %d mediaInfo in %s.\n", len(request), duration.String())

		w.Write(ToJSON(out))
	}
}

func (j *JmrAPI) Post_Rename() APIHandlerFn {
	return func(w http.ResponseWriter, r *http.Request) {
		var request RenameMediaRequest
		err := json.NewDecoder(r.Body).Decode(&request)

		if err != nil {
			j.HandleAPIError(w, r, http.StatusBadRequest, err)
			return
		}
		if len(request) == 0 {
			j.HandleAPIError(w, r, http.StatusBadRequest, errors.New("Atleast 1 request object is required"))
			return
		}

		out := RenameMediaResponse{}
		for _, reqItem := range request {
			targetInfo := util.Filter(reqItem.IdentifiedMediaInfos, func(x mediainfoprovider.MediaInfo) bool {
				return x.MediaID == reqItem.IdentifiedMediaId
			})
			var children []filesystem.DirEntry
			if reqItem.SourceDirectory.Entry.IsDirectory {
				children = j.fileSystemProvider.ScanDirectory(reqItem.SourceDirectory.Entry.Path, j.configProvider.GetAllowedExtensions())
			}
			entry := filesystem.DirEntry{
				Name:        reqItem.SourceDirectory.Entry.Name,
				Path:        reqItem.SourceDirectory.Entry.Path,
				Size:        reqItem.SourceDirectory.Entry.Size,
				IsDirectory: reqItem.SourceDirectory.Entry.IsDirectory,
				Children:    children,
			}
			entriesAndIgnores := j.ren.SelectEntriesForRename(entry, reqItem.SourceDirectory.Type)
			resItem := RenameMediaResponseItem{
				Info:              targetInfo[0],
				Type:              reqItem.SourceDirectory.Type,
				Entry:             entry,
				EntriesAndIgnores: entriesAndIgnores,
			}
			out = append(out, resItem)
		}

		w.Write(ToJSON(out))
	}
}

func (j *JmrAPI) HandleAPIError(w http.ResponseWriter, r *http.Request, errorCode int, err error) {
	msg := fmt.Sprintf("[HTTP %d]", errorCode)
	if err != nil {
		msg += fmt.Sprintf(" %s", err.Error())
	} else {
		msg += " Unknown error occured."
	}
	log.Println(msg)
	w.WriteHeader(errorCode)
	w.Write([]byte(msg))
}
