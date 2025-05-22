package main

import (
	// "context"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /product/{id}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("request received",r.PathValue("id"))
		time.Sleep(10*time.Second)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Inbound Mediation successful"}`))
		fmt.Println("response sent", r.PathValue("id"))
	})

	mux.HandleFunc("POST /product/{id}", func(w http.ResponseWriter, r *http.Request) {
    	productId := r.PathValue("id")
    	fmt.Fprintf(w, "creating product with id %s", productId)
	})
	
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	waitGrpoup := &sync.WaitGroup{}


	
	// Start the server in a goroutine
	waitGrpoup.Add(1)
	go func() {
		defer waitGrpoup.Done()
		slog.Info("Starting HTTP server on :8080")
		if err := server.ListenAndServe(); err != nil  {
			slog.Error("HTTP Inbound listener error", "error", err)
		}
	}()
	// Start a goroutine to monitor context cancellation and shut down server
	waitGrpoup.Add(1)
	go func() {
		defer waitGrpoup.Done()
		<-ctx.Done()
		slog.Info("Shutting down HTTP server...")
		// Shutdown the server gracefully
        // Calculate the shutdown time
		shutdownStart := time.Now()
		err := server.Shutdown(ctx); 
		if err != nil {
			slog.Error("Error shutting down HTTP server", "error", err.Error())
		}
		shutdownDuration := time.Since(shutdownStart)
		slog.Info("Server shutdown completed", "duration", shutdownDuration)
		slog.Info("waiting after shutting down HTTP server")
		time.Sleep(20 * time.Second)
	}()
	waitGrpoup.Wait()
	slog.Info("HTTP server shut down gracefully")
}




