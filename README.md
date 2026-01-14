# This fork

Upgrades the original to Mem0 V2 API.

**NOTE**: this shall be regarded as an MVP and not as a fully-fledged Mem0 Go client.

## Installation
**THIS WON'T WORK**
```
go get github.com/intx4/mem0-go
```
Instead, fetch the latest version from the output of the command above (e.g. `v0.0.0-20260114100254-688eb5c13010`)
and use it to insert a `replace` directive in your `go.mod` as such:
```
replace github.com/bytectlgo/mem0-go => github.com/intx4/mem0-go v0.0.0-20260113080456-7f54441e2c2b
```
then simply:
```
go get github.com/bytectlgo/mem0-go
```


# Mem0 Go Client

[中文文档 (Chinese Documentation)](README_ZH.md)

Mem0 Go Client is a Go language client library for interacting with the Mem0 API.

## Installation

```bash
go get github.com/bytectlgo/mem0-go
```

## Requirements

- Go 1.18 or higher

## Quick Start

```go
package main

import (
	"fmt"
	"log"

	"github.com/bytectlgo/mem0-go/client"
	"github.com/bytectlgo/mem0-go/types"
)

func main() {
	// Create client
	mem0, err := client.NewMemoryClient(client.ClientOptions{
		APIKey: "your-api-key",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Add memory
	memories, err := mem0.Add("Hello, World!", types.MemoryOptions{
		UserID: "user-123",
		Metadata: map[string]interface{}{
			"source": "example",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Added memory: %+v\n", memories[0])

	// Search memory
	results, err := mem0.Search("Hello", &types.SearchOptions{
		Limit: 10,
		Threshold: 0.8,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Search results: %+v\n", results)

	// Get project info
	project, err := mem0.GetProject(types.ProjectOptions{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Project info: %+v\n", project)
}
```

## Features

- Memory Management
  - Add memory
  - Update memory
  - Get memory
  - Search memory
  - Delete memory
  - Batch operations
- User Management
  - Get user list
  - Delete user
- Project Management
  - Get project info
  - Update project settings
- Webhook Management
  - Create webhook
  - Update webhook
  - Delete webhook
  - Get webhook list
- Feedback
  - Submit feedback

## API Documentation

### Client Initialization

```go
client, err := client.NewMemoryClient(client.ClientOptions{
	APIKey:          "your-api-key",
	Host:            "https://api.mem0.ai", // Optional, defaults to https://api.mem0.ai
	OrganizationName: "org-name",           // Optional
	ProjectName:     "project-name",        // Optional
	OrganizationID:  "org-id",             // Optional
	ProjectID:       "project-id",         // Optional
})
```

### Memory Operations

#### Add Memory

```go
memories, err := client.Add("Hello, World!", types.MemoryOptions{
	UserID: "user-123",
	Metadata: map[string]interface{}{
		"source": "example",
	},
})
```

#### Update Memory

```go
memories, err := client.Update("memory-id", "Updated content")
```

#### Get Memory

```go
memory, err := client.Get("memory-id")
```

#### Search Memory

```go
results, err := client.Search("query", &types.SearchOptions{
	Limit: 10,
	Threshold: 0.8,
})
```

#### Delete Memory

```go
err := client.Delete("memory-id")
```

### User Management

#### Get User List

```go
users, err := client.Users()
```

#### Delete User

```go
err := client.DeleteUser("user-id")
```

### Project Management

#### Get Project Info

```go
project, err := client.GetProject(types.ProjectOptions{})
```

#### Update Project Settings

```go
err := client.UpdateProject(types.PromptUpdatePayload{
	CustomInstructions: "New instructions",
})
```

### Webhook Management

#### Create Webhook

```go
webhook, err := client.CreateWebhook(types.WebhookPayload{
	EventTypes: []types.WebhookEvent{types.MemoryAdded},
	ProjectID:  "project-id",
	Name:       "My Webhook",
	URL:        "https://example.com/webhook",
})
```

#### Update Webhook

```go
err := client.UpdateWebhook(types.WebhookPayload{
	WebhookID:  "webhook-id",
	EventTypes: []types.WebhookEvent{types.MemoryAdded, types.MemoryUpdated},
	Name:       "Updated Webhook",
	URL:        "https://example.com/webhook",
})
```

#### Delete Webhook

```go
err := client.DeleteWebhook("webhook-id")
```

#### Get Webhook List

```go
webhooks, err := client.GetWebhooks("project-id")
```

### Feedback

```go
err := client.Feedback(types.FeedbackPayload{
	MemoryID:      "memory-id",
	Feedback:      types.Positive,
	FeedbackReason: "Helpful response",
})
```

## Error Handling

All API methods may return errors. Error types include:

- `APIError`: Returned when API request fails
- Other standard Go errors

```go
memory, err := client.Get("memory-id")
if err != nil {
	if apiErr, ok := err.(*client.APIError); ok {
		fmt.Printf("API Error: %s\n", apiErr.Message)
	} else {
		fmt.Printf("Error: %v\n", err)
	}
	return
}
```

## License

MIT

## FAQ

### How to handle API rate limiting?

When encountering API rate limiting, it's recommended to implement an exponential backoff retry mechanism. Example code:

```go
func retryWithBackoff(fn func() error) error {
	var err error
	for i := 0; i < 3; i++ {
		err = fn()
		if err == nil {
			return nil
		}
		if apiErr, ok := err.(*client.APIError); ok && apiErr.StatusCode == 429 {
			time.Sleep(time.Duration(math.Pow(2, float64(i))) * time.Second)
			continue
		}
		return err
	}
	return err
}
```

### How to batch process memories?

Use the `AddBatch` method to add multiple memories at once:

```go
memories := []string{"memory1", "memory2", "memory3"}
results, err := client.AddBatch(memories, types.MemoryOptions{})
```

## Contributing

We welcome contributions of any kind! Before submitting a Pull Request, please ensure:

1. Code follows Go standard formatting
2. Necessary tests are added
3. Relevant documentation is updated
4. Commit messages are clear and descriptive

## Changelog

### v0.1.0 (2025-04-18)
- Initial version release
- Support for basic memory management features
- Support for user management features
- Support for project management features
- Support for webhook management features
- Support for feedback features
