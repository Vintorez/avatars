package main

import (
	"net/http"

	"github.com/drone/config"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
)

const (
	BaseUrl    = `/`
	ApiUrl     = `api/v1/`
	BaseApiUrl = BaseUrl + ApiUrl
)

var (
	Listen = config.String("http", "0.0.0.0:4567")
)

func main() {
	config.SetPrefix("AV_")
	config.Parse("")

	mux := web.New()
	mux.Use(SetHeaders)
	mux.Use(middleware.Logger)
	mux.Use(Options)

	// NOTE: "RouterWithId" router needs for using URL parameters in CheckId middleware.
	// Goji can't bind URL parameters until after the middleware stack runs.
	// https://github.com/zenazn/goji/issues/32#issuecomment-46124240
	RouterWithId := web.New()
	RouterWithId.Use(CheckId)
	RouterWithId.Post(BaseApiUrl+"file/:id", UploadFile)
	RouterWithId.Put(BaseApiUrl+"file/:id", UpdateFile)
	RouterWithId.Patch(BaseApiUrl+"file/:id", ChangeMask)
	RouterWithId.Delete(BaseApiUrl+"file/:id", DeleteFile)
	RouterWithId.Get(BaseApiUrl+"file/:id", GetResizedFile)
	RouterWithId.Get(BaseApiUrl+"file/:id/raw", GetOriginalFile)

	mux.Handle(BaseApiUrl+"file/:id", RouterWithId)
	mux.Handle(BaseApiUrl+"file/:id/*", RouterWithId)

	http.Handle(BaseApiUrl, mux)

	http.Handle(BaseUrl, http.FileServer(http.Dir("app")))

	panic(http.ListenAndServe(*Listen, nil))
}
