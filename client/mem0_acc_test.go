//go:build testacc

////////////////////////////////////////////////////////////
// This file is used to test the mem0 API in an acceptance test environment.
// It is not meant to be run as a unit test.
////////////////////////////////////////////////////////////

package client_test

import (
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/bytectlgo/mem0-go/client"
	"github.com/bytectlgo/mem0-go/types"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

var (
	mem0ApiKey string
	memID      string
)

const (
	userID        = "test-gosdk-user"
	agentID       = "test-gosdk-agent"
	appID         = "test-gosdk-app"
	runID         = "test-gosdk-run"
	metadataKeyID = "test-gosdk-metadata-key-id"
)

func TestMain(m *testing.M) {
	if err := godotenv.Load(); err != nil {
		envPath := filepath.Join("..", ".env")
		if err := godotenv.Load(envPath); err != nil {
			log.Printf("Warning: could not load .env file: %v", err)
		}
	}

	mem0ApiKey = os.Getenv("TEST_MEM0_API_KEY")
	if mem0ApiKey == "" {
		log.Fatal("TEST_MEM0_API_KEY environment variable is required")
	}
	memID = os.Getenv("TEST_MEM0_MEMORY_ID")
	if memID == "" {
		log.Fatal("TEST_MEM0_MEMORY_ID environment variable is required")
	}
	os.Exit(m.Run())
}

func TestAccAddMemoryAndRetrieveEventAsync(t *testing.T) {
	client, err := client.NewMemoryClient(client.ClientOptions{
		APIKey: mem0ApiKey,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	message := types.Message{
		Role:    "user",
		Content: "Client was onboarded on 2025-02-03 and is a new customer. He likes churros and ice cream, but not burgers. He LOVES pizza tho",
	}

	events, err := client.AddAsync(message, types.MemoryOptions{
		UserID:  userID,
		AgentID: agentID,
		AppID:   appID,
		RunID:   runID,
		Metadata: map[string]any{
			"timestamp":       time.Now().Unix(),
			"metadata_key_id": metadataKeyID,
		},
		CustomCategories: types.CustomCategories{
			{
				CategoryName:        "payment_terms",
				CategoryDescription: "Customer preferences on payment terms",
			},
			{
				CategoryName:        "delivery_preferences",
				CategoryDescription: "Customer preferences on delivery time",
			},
			{
				CategoryName:        "product_preferences",
				CategoryDescription: "Customer preferences on product",
			},
		},
	})
	if err != nil {
		t.Fatalf("Failed to add memory: %v", err)
	}
	if len(events) == 0 {
		t.Fatalf("No add responses")
	}
	eventID := events[0].EventID

	event, err := client.GetEvent(eventID)
	if err != nil {
		t.Fatalf("Failed to get event: %v", err)
	}
	t.Logf("Event: %v", event)
}

func TestAccAddMemoryAndRetrieveEventSync(t *testing.T) {
	client, err := client.NewMemoryClient(client.ClientOptions{
		APIKey: mem0ApiKey,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	message := types.Message{
		Role:    "user",
		Content: "Client enjoys the product and is satisfied with the purchase. Client has 5 delivery orders in the last 12 months. Client has 2 delivery addresses in the last 12 months.",
	}

	memories, err := client.Add(message, types.MemoryOptions{
		UserID:  userID,
		AgentID: agentID,
		AppID:   appID,
		RunID:   runID,
		Metadata: map[string]any{
			"timestamp":       time.Now().Unix(),
			"metadata_key_id": metadataKeyID,
		},
		Version: types.V2,
	})
	if err != nil {
		t.Fatalf("Failed to add memory: %v", err)
	}
	if len(memories) == 0 {
		t.Fatalf("No add responses")
	}
	t.Logf("Memory: %v", memories[0])
}

func TestAccGetEvents(t *testing.T) {
	client, err := client.NewMemoryClient(client.ClientOptions{
		APIKey: mem0ApiKey,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	events, err := client.GetEvents("")
	if err != nil {
		t.Fatalf("Failed to get events: %v", err)
	}
	t.Logf("Events: %v", events)
}

func TestAccGetMemory(t *testing.T) {
	memoryID := memID
	client, err := client.NewMemoryClient(client.ClientOptions{
		APIKey: mem0ApiKey,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	memory, err := client.Get(memoryID)
	if err != nil {
		t.Fatalf("Failed to get memory: %v", err)
	}
	t.Logf("Memory: %v", memory)
}

func TestAccGetAllMemories(t *testing.T) {
	client, err := client.NewMemoryClient(client.ClientOptions{
		APIKey: mem0ApiKey,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	memories, err := client.GetAll(&types.SearchOptions{
		MemoryOptions: types.MemoryOptions{
			PageSize: 1,
			Filters: map[string]any{
				"user_id": userID,
				"metadata": map[string]any{
					"metadata_key_id": metadataKeyID,
				},
				"created_at": map[string]any{
					"gte": time.Now().Add(-1 * time.Hour * 24 * 30 * 12).Format(time.RFC3339),
					"lte": time.Now().Format(time.RFC3339),
				},
			},
		},
		Categories: []string{
			"product_preferences",
		},
	})
	if err != nil {
		t.Fatalf("Failed to get all memories: %v", err)
	}
	require.Equal(t, 1, len(memories))
	t.Logf("Memories: %v", memories)
}

func TestAccSearchMemories(t *testing.T) {
	client, err := client.NewMemoryClient(client.ClientOptions{
		APIKey: mem0ApiKey,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	memories, err := client.Search("Does the client enjoy the product?", &types.SearchOptions{
		TopK: 1,
		MemoryOptions: types.MemoryOptions{
			PageSize: 2,
			Filters: map[string]any{
				"user_id": userID,
				"metadata": map[string]any{
					"metadata_key_id": metadataKeyID,
				},
				"created_at": map[string]any{
					"gte": time.Now().Add(-1 * time.Hour * 24 * 30 * 12).Format(time.RFC3339),
					"lte": time.Now().Format(time.RFC3339),
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("Failed to get all memories: %v", err)
	}
	require.Equal(t, 1, len(memories))
	t.Logf("Memories: %v", memories)
}
