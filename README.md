# easy-outbox

`easy-outbox` is a Go library that provides a simple and reliable implementation of the transactional outbox pattern. It's designed to help you send messages idempotently without worrying about the underlying implementation details.

The outbox pattern is crucial for services that need to publish messages or events reliably. It ensures that messages are sent "at least once" by persisting them in your application's database before notifying a message broker or other services. `easy-outbox` automates this process for you.

## Installation

To install `easy-outbox`, use `go get`:

```bash
go get github.com/hertzcodes/easy-outbox/go
```

## Usage

Here's a quick example of how to use `easy-outbox`:

### Basic Usage: Manual Fetching

You can manually fetch and delete messages from the outbox.

```go
package main

import (
	"fmt"
	"log"
	"github.com/hertzcodes/easy-outbox/go/outbox"
)

func main() {
	// Initialize the outbox with PebbleDB
	box, err := outbox.New(outbox.DBTypePebble, "./outbox_db")
	if err != nil {
		log.Fatalf("failed to create outbox: %v", err)
	}

	// Add messages to the outbox
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("message_%d", i)
		// The value can be any data you want to store.
		// For this example, we'll use a simple string.
		value := fmt.Sprintf("This is message %d", i)
		if err := box.SetMessage(key, []byte(value)); err != nil {
			log.Printf("failed to set message: %v", err)
		}
	}

	// Get messages from the outbox
	messages, err := box.GetMessages(10)
	if err != nil {
		log.Fatalf("failed to get messages: %v", err)
	}

	fmt.Println("Processing messages...")
	for _, msgKey := range messages {
		fmt.Printf("  - Processing message: %s\n", msgKey)
		// Here you would send the message to your message broker (e.g., Kafka, RabbitMQ)
		// ...

		// After successful processing, delete the message from the outbox
		if err := box.Delete(msgKey); err != nil {
			log.Printf("failed to delete message: %v", err)
		}
	}
	fmt.Println("Done.")
}
```

### Automated Usage: Interval-Based Fetching

For a more hands-off approach, you can configure `easy-outbox` to fetch messages at a regular interval and send them to a channel for processing.

```go
package main

import (
	"fmt"
	"log"
	"time"
	"github.com/hertzcodes/easy-outbox/go/outbox"
)

func main() {
	// Initialize the outbox with PebbleDB
	box, err := outbox.New(outbox.DBTypePebble, "./outbox_db_interval")
	if err != nil {
		log.Fatalf("failed to create outbox: %v", err)
	}

	// Add some messages
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("message_%d", i)
		value := fmt.Sprintf("This is message %d", i)
		if err := box.SetMessage(key, []byte(value)); err != nil {
			log.Printf("failed to set message: %v", err)
		}
	}

	// Create a channel to receive messages from the outbox
	messageChannel := make(chan []string)

	// Set up the interval-based fetching
	// This will fetch 10 messages every 2 seconds
	box.SetInterval(messageChannel, 2*time.Second, 10)

	fmt.Println("Listening for messages...")
	// Process messages as they arrive on the channel
	for messages := range messageChannel {
		fmt.Printf("Received a batch of %d messages.\n", len(messages))
		for _, msgKey := range messages {
			fmt.Printf("  - Processing message: %s\n", msgKey)
			// ... send to message broker ...
			box.Delete(msgKey)
		}
	}
}
```

## Supported Databases

`easy-outbox` is designed to be extensible with different database backends.

### Currently Supported
- **Pebble**: A lightweight, embedded key-value store.

### Adding a New Database
To add support for a new database, you need to implement the `bindings.DB` interface:

```go
// located at /internal/bindings/contract.go
package bindings

type DB interface {
	Write(key string, value any) error
	Read(key string) (interface{}, error)
	Delete(key string) error
	ReadBulkKeys(amount int) []string
	PrintMetrics()
}
```
Then, you can add it to the `New` function in `outbox.go`.

## Running Tests

To run the test suite:

```bash
go test ./...
```

## Contributing

Contributions are welcome! Please feel free to submit a pull request.

## License

This project is licensed under the GNU GENERAL PUBLIC LICENSE.