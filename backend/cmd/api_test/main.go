package main

import (
	"os"

	domainappapitest "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/app/api_test"
)

func main() {
	ec := domainappapitest.Run()
	os.Exit(ec)
}
