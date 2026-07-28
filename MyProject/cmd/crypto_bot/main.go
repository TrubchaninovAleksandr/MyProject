package main

import (
	"MyProject/deploy/config"
	"MyProject/internal/application"
	"log"
)

// @title           MyProject Crypto API
// @version         1.0
// @description     API for retrieving cryptocurrency exchange rates and statistics on average, minimum, and maximum values.

// @host      localhost:8080
// @BasePath  /v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

func main() {
	cfg, err := config.LoadConfig("./deploy/config/config.yaml")
	log.Print(cfg)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	application := application.NewApp(cfg)
	application.Run()

}
