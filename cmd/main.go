package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"os"
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

		s := &http.Server{
			Addr:      appState.ServerAddr,
			Handler:   mux,
			TLSConfig: tlsConfig,
		}
		log.Fatal(s.ListenAndServeTLS(certPath, keyPath))
	} else {
		err = http.ListenAndServe(appState.ServerAddr, mux)
		if err != nil {
			log.Error(stacktrace.Propagate(err, ""))
		}
	}
}
