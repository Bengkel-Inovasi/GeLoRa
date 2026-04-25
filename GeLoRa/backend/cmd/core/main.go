package main

import (
	"os"

	_ "github.com/Bengkel-Inovasi/GeLoRa/backend/docs/swagger"
	domainappcore "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/app/core"
)

// @title GeLoRa API
// @version 1.0
// @description The core API documentation for GeLoRa backend endpoints.
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer " followed by a space and JWT token.
func main() {
	ec := domainappcore.Run()
	os.Exit(ec)
}
