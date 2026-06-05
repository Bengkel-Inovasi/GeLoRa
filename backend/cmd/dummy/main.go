package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Bengkel-Inovasi/GeLoRa/backend/internal/config"
	domainappcore "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/app/core"
)

func main() {
	ctx := context.Background()
	_ = config.LoadEnv()

	// Using existing core logic to set up infra
	core := &domainappcore.Core{}
	if err := core.NewInfrastructure(ctx); err != nil {
		fmt.Printf("Failed infra: %v\n", err)
		os.Exit(1)
	}
	if err := core.NewWiring(ctx); err != nil {
		fmt.Printf("Failed wiring: %v\n", err)
		os.Exit(1)
	}

	// 1. Create Dummy Node
	mid := "DUMMY-NODE-001"
	name := "John Doe (Dummy)"
	nodeId, err := core.Wiring().RepoNode.CreateNode(ctx, mid, name)
	if err != nil {
		fmt.Printf("Node might already exist: %v\n", err)
		node, _ := core.Wiring().RepoNode.ReadNodeByMid(ctx, mid)
		if node != nil {
			nodeId = node.Id
		}
	}

	// 2. Validate Node (by Super User ID 1)
	_ = core.Wiring().RepoNode.UpdateNodeValidation(ctx, nodeId, true, 1)

	// 3. Create Active Session
	sessionId, err := core.Wiring().RepoSession.CreateSession(ctx, 1, nodeId)
	if err != nil {
		fmt.Printf("Session might already exist: %v\n", err)
	} else {
		// 4. Add Dummy Record
		hr := 75.5
		temp := 36.6
		lat := -6.2088
		lon := 106.8456
		_, _ = core.Wiring().RepoRecord.CreateRecord(ctx, sessionId, time.Now(), &hr, &temp, &lat, &lon)
	}

	fmt.Println("Successfully created dummy node, session, and record.")
}
