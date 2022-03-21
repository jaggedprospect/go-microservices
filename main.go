package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/jaggedprospect/go-microservices/product-api/handlers"
)

func main() {
	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	// create product handlers
	ph := handlers.NewProducts(l)

	// create new serve mux and register handlers
	sm := mux.NewRouter()
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", ph.GetProducts)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProducts)
	putRouter.Use(ph.MiddlewareValidateProduct)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", ph.AddProduct)
	postRouter.Use(ph.MiddlewareValidateProduct)

	// create new server
	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	}

	// start the server
	go func() {
		l.Println("Starting server on port :9090")

		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	sig := <-c
	l.Println("Received terminate, graceful shutdown", sig)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
