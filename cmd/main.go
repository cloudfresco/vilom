package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/palantir/stacktrace"

	"github.com/cloudfresco/vilom/routes"
	"github.com/cloudfresco/vilom/search/searchservices"
)

// SetUpLogging - start the logging sub-system
func setUpLogging(logFile string, logLevel log.Level) error {
	var err error
	var f *os.File

	// open the log file
	f, err = os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		return err
	}

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	log.SetOutput(f)

	log.SetLevel(logLevel)
	return nil
}

func getKeys(caCertPath string, certPath string, keyPath string) *tls.Config {

	caCert, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		log.Fatal(err)
	}

	caCertpool := x509.NewCertPool()
	caCertpool.AppendCertsFromPEM(caCert)

	// LoadX509KeyPair reads files, so we give it the paths
	serverCert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		log.Fatal(err)
	}

	tlsConfig := tls.Config{
		ClientCAs:    caCertpool,
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}

	return &tlsConfig

}

func main() {
	var devFlag bool
	var logFile string
	var logLevel log.Level
	var appState *routes.AppState

	pwd, _ := os.Getwd()
	logFile = pwd + filepath.FromSlash("/files/log/app.log")
	logLevel = log.InfoLevel

	stacktrace.DefaultFormat = stacktrace.FormatFull

	err := setUpLogging(logFile, logLevel)
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
	}

	appState = &routes.AppState{}
	devFlag = false
	appState.Init(devFlag)
	appState.SearchIndex = searchservices.InitSearch("", appState.Db)
	mux := appState.RoutesInit()
	log.Info("Server Started")

	if appState.ServerTLS == "true" {
		var caCertPath, certPath, keyPath string
		var tlsConfig *tls.Config
		caCertPath = pwd + filepath.FromSlash(appState.KeyOptions.CaCertPath)
		certPath = pwd + filepath.FromSlash(appState.KeyOptions.CertPath)
		keyPath = pwd + filepath.FromSlash(appState.KeyOptions.KeyPath)

		tlsConfig = getKeys(caCertPath, certPath, keyPath)

		srv := &http.Server{
			Addr:      appState.ServerAddr,
			Handler:   mux,
			TLSConfig: tlsConfig,
		}

		idleConnsClosed := make(chan struct{})
		go func() {
			sigint := make(chan os.Signal, 1)
			signal.Notify(sigint, os.Interrupt)
			<-sigint

			// We received an interrupt signal, shut down.
			if err := srv.Shutdown(context.Background()); err != nil {
				// Error from closing listeners, or context timeout:
				log.Error("HTTP server Shutdown:", err)
			}
			close(idleConnsClosed)
		}()

		if err := srv.ListenAndServeTLS(certPath, keyPath); err != http.ErrServerClosed {
			// Error starting or closing listener:
			log.Error("HTTP server ListenAndServeTLS:", err)
		}
		log.Error("Server shutting down")

		<-idleConnsClosed
	} else {

		srv := &http.Server{
			Addr:    appState.ServerAddr,
			Handler: mux,
		}

		idleConnsClosed := make(chan struct{})
		go func() {
			sigint := make(chan os.Signal, 1)
			signal.Notify(sigint, os.Interrupt)
			<-sigint

			// We received an interrupt signal, shut down.
			if err := srv.Shutdown(context.Background()); err != nil {
				// Error from closing listeners, or context timeout:
				log.Error("HTTP server Shutdown:", err)
			}
			close(idleConnsClosed)
		}()

		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// Error starting or closing listener:
			log.Error("HTTP server ListenAndServe:", err)
		}

		log.Error("Server shutting down")

		<-idleConnsClosed

	}
}
