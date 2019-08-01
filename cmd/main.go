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

	"github.com/cloudfresco/vilom/common"
	"github.com/cloudfresco/vilom/routes"
	"github.com/cloudfresco/vilom/search/searchservices"
)

/* error message range: 100-249 */

func getConfigOpt() (*common.DBOptions, *common.RedisOptions, *common.MailerOptions, *common.ServerOptions, *common.RateOptions, *common.JWTOptions, *common.OauthOptions, *common.UserOptions, *common.LogOptions) {

	v, err := common.GetViper()
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 103,
		}).Error(err)
		os.Exit(1)
	}

	dbOpt, err := common.GetDbConfig(v)
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 103,
		}).Error(err)
		os.Exit(1)
	}

	redisOpt, err := common.GetRedisConfig(v)
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 103,
		}).Error(err)
		os.Exit(1)
	}

	mailerOpt, err := common.GetMailerConfig(v)
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 103,
		}).Error(err)
		os.Exit(1)
	}

	serverOpt, err := common.GetServerConfig(v)
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 103,
		}).Error(err)
		os.Exit(1)
	}

	rateOpt, err := common.GetRateConfig(v)
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 103,
		}).Error(err)
		os.Exit(1)
	}

	jwtOpt, err := common.GetJWTConfig(v)
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 103,
		}).Error(err)
		os.Exit(1)
	}

	oauthOpt, err := common.GetOauthConfig(v)
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 103,
		}).Error(err)
		os.Exit(1)
	}

	userOpt, err := common.GetUserConfig(v)
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 103,
		}).Error(err)
		os.Exit(1)
	}

	logOpt, err := common.GetLogConfig(v)
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 103,
		}).Error(err)
		os.Exit(1)
	}

	return dbOpt, redisOpt, mailerOpt, serverOpt, rateOpt, jwtOpt, oauthOpt, userOpt, logOpt
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
	var appState *routes.AppState
	var err error

	dbOpt, redisOpt, mailerOpt, serverOpt, rateOpt, jwtOpt, oauthOpt, userOpt, logOpt := getConfigOpt()

	common.SetUpLogging(logOpt)

	dbService, err := common.CreateDBService(dbOpt)
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 750,
		}).Error(err)
		os.Exit(1)
	}

	redisService, err := common.CreateRedisService(redisOpt)
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 750,
		}).Error(err)
		os.Exit(1)
	}

	mailerService, err := common.CreateMailerService(mailerOpt)
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 750,
		}).Error(err)
		os.Exit(1)
	}

	appState = &routes.AppState{}
	err = appState.Init(dbService, redisService, mailerService, serverOpt, rateOpt, jwtOpt, oauthOpt, userOpt)
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
		pwd, _ := os.Getwd()
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
