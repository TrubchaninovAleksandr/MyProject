package main

import (
	"MyProject/internal/application"
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
	application := application.NewApp()
	application.Run()
}
