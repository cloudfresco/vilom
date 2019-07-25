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

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"

	"github.com/cloudfresco/vilom/routes"
	"github.com/cloudfresco/vilom/search/searchservices"
)

/* error message range: 100-249 */

// SetUpLogging - start the logging sub-system
func setUpLogging(logFile string, logLevel log.Level) error {
	var err error
	var f *os.File

	// open the log file
	f, err = os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 100,
		}).Error(err)
		return err
	}

	log.SetOutput(f)
	log.SetFormatter(&log.JSONFormatter{})

	log.SetLevel(logLevel)
	return nil
}

func getKeys(caCertPath string, certPath string, keyPath string) *tls.Config {

	caCert, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 101,
		}).Error(err)
	}

	caCertpool := x509.NewCertPool()
	caCertpool.AppendCertsFromPEM(caCert)

	// LoadX509KeyPair reads files, so we give it the paths
	serverCert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 102,
		}).Error(err)
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

	err := setUpLogging(logFile, logLevel)
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 103,
		}).Error(err)
		os.Exit(1)
	}

	appState = &routes.AppState{}
	devFlag = false
	err = appState.Init(devFlag)
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 103,
		}).Error(err)
		os.Exit(1)
	}

	appState.SearchIndex = searchservices.InitSearch("", appState.DBService.DB)
	mux := appState.CreateRoutes()

	if appState.ServerOptions.ServerTLS == "true" {
		var caCertPath, certPath, keyPath string
		var tlsConfig *tls.Config
		caCertPath = pwd + filepath.FromSlash(appState.ServerOptions.CaCertPath)
		certPath = pwd + filepath.FromSlash(appState.ServerOptions.CertPath)
		keyPath = pwd + filepath.FromSlash(appState.ServerOptions.KeyPath)

		tlsConfig = getKeys(caCertPath, certPath, keyPath)

		srv := &http.Server{
			Addr:      appState.ServerOptions.ServerAddr,
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
				log.WithFields(log.Fields{
					"msgnum": 104,
				}).Warn("HTTP server Shutdown")
			}
			close(idleConnsClosed)
		}()

		if err := srv.ListenAndServeTLS(certPath, keyPath); err != http.ErrServerClosed {
			// Error starting or closing listener:
			log.WithFields(log.Fields{
				"msgnum": 105,
			}).Warn("HTTP server ListenAndServeTLS")
		}
		log.WithFields(log.Fields{
			"msgnum": 106,
		}).Warn("Server shutting down")

		<-idleConnsClosed
	} else {

		srv := &http.Server{
			Addr:    appState.ServerOptions.ServerAddr,
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
				log.WithFields(log.Fields{
					"msgnum": 107,
				}).Warn("HTTP server Shutdown")
			}
			close(idleConnsClosed)
		}()

		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// Error starting or closing listener:
			log.WithFields(log.Fields{
				"msgnum": 108,
			}).Warn("HTTP server ListenAndServe")
		}

		log.Error("")
		log.WithFields(log.Fields{
			"msgnum": 109,
		}).Warn("Server shutting down")

		<-idleConnsClosed
	}
}
