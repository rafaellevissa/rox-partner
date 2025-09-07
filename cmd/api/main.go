package main

import (
	"database/sql"
	"log"

	"github.com/rafaellevissa/rox-partner/internal/api"
	"github.com/spf13/viper"

	_ "github.com/lib/pq"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./configs")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error reading config file: %v", err)
	}

	dsn := viper.GetString("database.dsn")
	if dsn == "" {
		log.Fatal("database DSN not configured")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	r := api.NewRouter(db)
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
