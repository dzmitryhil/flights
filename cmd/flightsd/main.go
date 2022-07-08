package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/dzmitryhil/flights/docs"
	"github.com/dzmitryhil/flights/finder"
	"github.com/dzmitryhil/flights/handler"
)

func main() {
	const addr = "0.0.0.0:8080"

	router := mux.NewRouter()

	router.Handle("/openapi/openapi.yml", http.FileServer(http.FS(docs.Docs)))
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/openapi/openapi.yml"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	pathFinder := finder.NewFlightPath()
	appHandlers := handler.NewFlightPathHandler(pathFinder)
	router.HandleFunc("/flights/path", appHandlers.PostFlightPath())

	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  time.Second * 15, //nolint:gomnd // constant for the main
		WriteTimeout: time.Second * 15, //nolint:gomnd // constant for the main
		IdleTimeout:  time.Second * 60, //nolint:gomnd // constant for the main
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	log.Printf("The service is listening...")
	log.Printf("Open http://%s/swagger/ to access the Swagger UI", addr)

	<-c
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15) //nolint:gomnd // constant for the main
	defer cancel()
	defer os.Exit(0)
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("shutdown with error, %s", err)
		return
	}

	log.Println("shutting down")
}
