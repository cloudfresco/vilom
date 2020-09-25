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

	"github.com/throttled/throttled/v2/store/goredisstore"

	"github.com/cloudfresco/vilom/common"
	"github.com/cloudfresco/vilom/msg/msgservices"
	"github.com/cloudfresco/vilom/search/searchservices"
	"github.com/cloudfresco/vilom/user/userservices"

	"github.com/cloudfresco/vilom/msg/msgcontrollers"
	"github.com/cloudfresco/vilom/search/searchcontrollers"
	"github.com/cloudfresco/vilom/user/usercontrollers"
)

/* error message range: 100-249 */

func getConfigOpt() (*common.DBOptions, *common.RedisOptions, *common.MailerOptions, *common.ServerOptions, *common.RateOptions, *common.JWTOptions, *common.OauthOptions, *common.UserOptions, *common.RoleOptions, *common.LogOptions) {

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

	roleOpt, err := common.GetRoleConfig(v)
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

	return dbOpt, redisOpt, mailerOpt, serverOpt, rateOpt, jwtOpt, oauthOpt, userOpt, roleOpt, logOpt
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
	var err error

	dbOpt, redisOpt, mailerOpt, serverOpt, rateOpt, jwtOpt, _, userOpt, roleOpt, logOpt := getConfigOpt()

	common.SetUpLogging(logOpt)
	common.SetJWTOpt(jwtOpt)

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

	store, err := goredisstore.New(redisService.RedisClient, "throttled:")
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 754,
		}).Error(err)
		os.Exit(1)
	}

	searchIndex := searchservices.InitSearch("", dbService.DB)

	/*authEnforcer, err := casbin.NewEnforcer("./auth_model.conf", "./policy.csv")
		if err != nil {
				log.WithFields(log.Fields{
				"msgnum": 755,
			}).Error(err)
	    os.Exit(1)
		}*/
	authEnforcer, err := common.LoadEnforcer(dbService, roleOpt)

	userService := userservices.NewUserService(dbService, redisService, mailerService, jwtOpt, userOpt, authEnforcer)
	ugroupService := userservices.NewUgroupService(dbService, redisService)
	ubadgeService := userservices.NewUbadgeService(dbService, redisService)

	workspaceService := msgservices.NewWorkspaceService(dbService, redisService)
	channelService := msgservices.NewChannelService(dbService, redisService)
	msgService := msgservices.NewMessageService(dbService, redisService)

	searchService := searchservices.NewSearchService(dbService, redisService, searchIndex)

	mux := http.NewServeMux()

	usercontrollers.Init(userService, ugroupService, ubadgeService, rateOpt, jwtOpt, mux, store)
	msgcontrollers.Init(workspaceService, channelService, msgService, userService, rateOpt, jwtOpt, mux, store)
	searchcontrollers.Init(searchService, userService, rateOpt, jwtOpt, mux, store)

	if serverOpt.ServerTLS == "true" {
		var caCertPath, certPath, keyPath string
		var tlsConfig *tls.Config
		pwd, _ := os.Getwd()
		caCertPath = pwd + filepath.FromSlash(serverOpt.CaCertPath)
		certPath = pwd + filepath.FromSlash(serverOpt.CertPath)
		keyPath = pwd + filepath.FromSlash(serverOpt.KeyPath)

		tlsConfig = getKeys(caCertPath, certPath, keyPath)

		srv := &http.Server{
			Addr:      serverOpt.ServerAddr,
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
			Addr:    serverOpt.ServerAddr,
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
