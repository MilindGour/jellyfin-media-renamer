package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/MilindGour/jellyfin-media-renamer/config"
	"github.com/MilindGour/jellyfin-media-renamer/filesystem"
	"github.com/MilindGour/jellyfin-media-renamer/middlewares"
	"github.com/MilindGour/jellyfin-media-renamer/util"
)

type JmrAPI struct {
	// any future api keys or such will
	// appear here
	serveMux           *http.ServeMux
	configProvider     config.ConfigProvider
	fileSystemProvider filesystem.FileSystemProvider

	sourcesWithID []DirConfigWithID
}

func NewJmrApi(
	configProvider config.ConfigProvider,
	filesystemProvider filesystem.FileSystemProvider,
) *JmrAPI {
	jmrApi := JmrAPI{
		configProvider:     configProvider,
		fileSystemProvider: filesystemProvider,
	}

	return &jmrApi
}

func (j *JmrAPI) Initialize(enableCors bool) {
	j.RegisterAPIRoutes()

	if j.configProvider == nil {
		log.Fatal("Please place config before running the server")
	}

	j.populateSourcesWithID()
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
func (j *JmrAPI) populateSourcesWithID() {
	// get the sources by ID and cache it in memory
	sources := j.configProvider.GetSourceList()
	j.sourcesWithID = []DirConfigWithID{}
	id := 0
	for _, src := range sources {
		id = id + 1
		j.sourcesWithID = append(j.sourcesWithID, DirConfigWithID{
			DirConfig: src,
			ID:        id,
		})
	}

}

func (j *JmrAPI) RegisterAPIRoutes() {
	j.serveMux = http.NewServeMux()

	j.serveMux.HandleFunc("GET /api/sources", j.Get_Sources())
	j.serveMux.HandleFunc("GET /api/sources/{id}", j.Get_SourceByID())
}

func (j *JmrAPI) Get_Sources() APIHandlerFn {
	return func(w http.ResponseWriter, r *http.Request) {
		raw := NewSourcesResponse(j.sourcesWithID)
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
		src := util.Find(j.sourcesWithID, func(dcwi DirConfigWithID) bool {
			return dcwi.ID == idInt
		})
		if src == nil {
			j.HandleAPIError(w, r, http.StatusNotFound, errors.New("Cannot find id "+id))
			return
		}

		w.Write(ToJSON(NewSourceByIDResponse(*src, j.fileSystemProvider.ScanDirectory(src.Path))))
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
