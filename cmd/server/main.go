package main

import (
	"image-service/internal/pkg/service"
	"log"
	"net/http"
)

func registerWatermarkController(mux *http.ServeMux) {
	controller := service.NewWatermarkController(service.ImageTransformer{}, 1024, 768, "/tmp")
	mux.Handle("/watermark", controller)
}

func main() {
	muxer := http.NewServeMux()
	registerWatermarkController(muxer)
	server := http.Server{
		Addr:    ":3210",
		Handler: muxer,
	}
	log.Fatal(server.ListenAndServe())

}
