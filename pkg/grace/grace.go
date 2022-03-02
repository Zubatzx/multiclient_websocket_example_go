package grace

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"
)

// Serve will run HTTP server with graceful shutdown capability
func Serve(port string, h http.Handler) error {

	// create new http server object
	server := &http.Server{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      h,
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	idleConnsClosed := make(chan struct{})
	go func() {

		signals := make(chan os.Signal, 1)

		signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
		<-signals

		// We received an os signal, shut down.
		if err := server.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server shutdown error: %v", err)
		}

		close(idleConnsClosed)
	}()

	log.Println("HTTP server running on port", port)
	if err := server.Serve(lis); err != http.ErrServerClosed {
		// Error starting or closing listener:
		return err
	}

	<-idleConnsClosed
	log.Println("HTTP server shutdown gracefully")
	return nil
}

// ServeTLS will run HTTPS server with graceful shutdown capability
func ServeTLS(h http.Handler) error {
	port := ":443"

	// create new http server object
	server := &http.Server{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      h,
	}

	idleConnsClosed := make(chan struct{})
	go func() {

		signals := make(chan os.Signal, 1)

		signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
		<-signals

		// We received an os signal, shut down.
		if err := server.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTPS server shutdown error: %v", err)
		}

		close(idleConnsClosed)
	}()

	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	cert := filepath.Join(basepath, "../..", "files", "ssl", "certificate.crt")
	key := filepath.Join(basepath, "../..", "files", "ssl", "private.key")

	if _, err := os.Stat(cert); err != nil {
		cert = "/certificate.crt"
	}

	if _, err := os.Stat(key); err != nil {
		key = "/private.key"
	}

	log.Println("HTTPS server running on port", port)
	if err := server.ListenAndServeTLS(cert, key); err != http.ErrServerClosed {
		// Error starting or closing listener:
		return err
	}

	<-idleConnsClosed
	log.Println("HTTPS server shutdown gracefully")
	return nil
}
