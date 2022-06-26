package main

import (
	"SYBD/internal/api"
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	defaultAddress = "0.0.0.0"
	defaultPort    = "8080"
)

func main() {
	// -------------------- Set up viper (config) -------------------- //

	viper.AutomaticEnv()

	viper.SetConfigFile("/cmd/configs/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config file: %s\n", err)
	}

	viper.SetDefault("service.bind.address", defaultAddress)
	viper.SetDefault("service.bind.port", defaultPort)

	// -------------------- Set up logging -------------------- //

	log := logrus.New()

	//formatter := logrus.JSONFormatter{
	//	TimestampFormat: time.RFC3339,
	//}

	switch viper.GetString("logging.level") {
	case "warning":
		log.SetLevel(logrus.WarnLevel)
	case "notice":
		log.SetLevel(logrus.InfoLevel)
	//case "debug":
	//	log.SetLevel(logrus.DebugLevel)
	//	formatter.PrettyPrint = true
	default:
		log.SetLevel(logrus.InfoLevel)
	}

	//	log.SetFormatter(&formatter)

	//log.Infof("log level: %s", log.Level.String())

	// -------------------- Set up database -------------------- //

	db, err := pgxpool.Connect(context.Background(), viper.GetString("db.connection_string"))
	if err != nil {
		log.Fatalf("unable to connect to database: %s", err)
	}
	defer db.Close()

	// -------------------- Set up service -------------------- //
	svc, err := api.NewAPIService(logrus.NewEntry(log), db)
	if err != nil {
		log.Fatalf("error creating service instance: %s", err)
	}

	go svc.Serve()

	// -------------------- Listen for Interruption signal and shutdown -------------------- //

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second*time.Duration(viper.GetInt("service.shutdown_timeout")),
	)
	defer cancel()

	if err := svc.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
