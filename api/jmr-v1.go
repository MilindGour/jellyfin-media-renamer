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
	"github.com/gorilla/mux"
)

type JmrAPI struct {
	// any future api keys or such will
	// appear here
	*mux.Router
	configProvider     config.ConfigProvider
	fileSystemProvider filesystem.FileSystemProvider

	sourcesWithID []DirConfigWithID
}

func NewJmrApi(isDevEnv bool) *JmrAPI {
	var targetCfg config.ConfigProvider
	var targetFSProvider filesystem.FileSystemProvider

	if isDevEnv {
		// DEV Environment ONLY
		log.Println("JMR API starting in DEVELOPER environment")
		targetCfg = config.NewDevJmrConfig()
		targetFSProvider = filesystem.NewMockJmrFS()

	} else {
		// PROD Environment ONLY
		log.Println("JMR API starting in PRODUCTION environment")
		targetCfg = config.NewJmrConfig()
		targetFSProvider = filesystem.NewJmrFS()
	}

	jmrApi := JmrAPI{
		Router:             mux.NewRouter(),
		configProvider:     targetCfg,
		fileSystemProvider: targetFSProvider,
	}

	// register all the api routes
	jmrApi.RegisterAPIRoutes()

	if jmrApi.configProvider == nil {
		log.Fatal("Please place config before running the server")
	}

	jmrApi.populateSourcesWithID()

	log.Printf("Starting server on port %s", jmrApi.GetPort())
	return &jmrApi
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
func (j *JmrAPI) GetPort() string {
	return j.configProvider.GetPort()
}

func (j *JmrAPI) RegisterAPIRoutes() {
	j.Use(middlewares.CorsMW)
	j.Use(middlewares.LogMW)
	j.Use(middlewares.Json)

	j.HandleFunc("/api/sources", j.Get_Sources()).Methods("GET")
	j.HandleFunc("/api/sources/{id}", j.Get_SourceByID()).Methods("GET")
}

func (j *JmrAPI) Get_Sources() APIHandlerFunction {
	return func(w http.ResponseWriter, r *http.Request) {
		raw := NewSourcesResponse(j.sourcesWithID)
		w.Write(ToJSON(raw))
	}
}
func (j *JmrAPI) Get_SourceByID() APIHandlerFunction {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, hasId := vars["id"]
		if !hasId {
			j.HandleAPIError(w, r, http.StatusBadRequest, errors.New("ID is required"))
			return
		}
		idInt, err := strconv.Atoi(id)
		if err != nil {
			j.HandleAPIError(w, r, http.StatusBadRequest, err)
			return
		}
		src := util.Find[DirConfigWithID](j.sourcesWithID, func(dcwi DirConfigWithID) bool {
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
		msg += fmt.Sprint(" Unknown error occured.")
	}
	log.Println(msg)
	w.WriteHeader(errorCode)
	w.Write([]byte(msg))
}
