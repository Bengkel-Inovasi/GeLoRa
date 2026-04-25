package main

import (
	"os"

	domainappseeder "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/app/seeder"
)

func main() {
	ec := domainappseeder.Run()
	os.Exit(ec)
}
