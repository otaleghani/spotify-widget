package server

import (
	"log"
	"net/http"
	"time"

	"github.com/otaleghani/spotify-widget/internal/database"
)

func Serve() {
	http.HandleFunc("/image", serveImage)

	srv := &http.Server{
		Addr:         ":8081",
		Handler:      nil,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	go func() {
		log.Println("Server is starting at :8081")
		log.Fatal(srv.ListenAndServe())
	}()
}

func serveImage(w http.ResponseWriter, r *http.Request) {
	filepath, err := database.WidgetImageFilepath()
	if err != nil {
		return
	}
  w.Header().Set("Cache-Control", "public, max-age=0")
  w.Header().Set("Expires", time.Now().Add(time.Hour).Format(http.TimeFormat))
  w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))

	http.ServeFile(w, r, filepath)
}
