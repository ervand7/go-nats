package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	role := flag.String("role", "pub", "role: pub | sub | req | rep")
	natsURL := flag.String("nats", "nats://nats1:4222", "NATS server URL")
	flag.Parse()

	nts, err := nats.Connect(*natsURL, nats.Name("NATS Demo"))
	if err != nil {
		log.Fatalf("failed to connect to NATS: %v", err)
	}
	defer func() {
		if err := nts.Drain(); err != nil {
			log.Printf("error during drain: %v", err)
		}
	}()

	// Setup context with cancellation on SIGINT/SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	switch *role {
	case "sub":
		_, err := nts.Subscribe("updates", func(msg *nats.Msg) {
			log.Printf("received: %s", msg.Data)
		})
		if err != nil {
			log.Fatalf("failed to subscribe: %v", err)
		}
		// Ensures the subscription is registered with the NATS server before proceeding
		if err := nts.Flush(); err != nil {
			log.Fatalf("flush error: %v", err)
		}
		log.Println("subscribed to 'updates' – waiting …")
		<-ctx.Done() // Wait until Ctrl+C or SIGTERM

	case "pub":
		const total = 100_000
		start := time.Now()
		for i := 0; i < total; i++ {
			if err := nts.Publish("updates", []byte(fmt.Sprintf("msg %d", i))); err != nil {
				log.Fatalf("publish failed at msg %d: %v", i, err)
			}
		}
		// Ensures all published messages are sent to the NATS server
		if err := nts.Flush(); err != nil {
			log.Fatalf("flush error: %v", err)
		}
		dur := time.Since(start)
		log.Printf("sent %d msgs in %v (%.0f msg/s)", total, dur,
			float64(total)/dur.Seconds())

	case "rep":
		_, err := nts.Subscribe("ping", func(msg *nats.Msg) {
			if err := nts.Publish(msg.Reply, []byte("pong")); err != nil {
				log.Printf("failed to reply: %v", err)
			}
		})
		if err != nil {
			log.Fatalf("failed to subscribe: %v", err)
		}
		// Ensures the reply subscription is active before handling incoming requests
		if err := nts.Flush(); err != nil {
			log.Fatalf("flush error: %v", err)
		}
		log.Println("replying on subject 'ping'")
		<-ctx.Done() // Block until shutdown signal

	case "req":
		resp, err := nts.Request("ping", []byte("ping"), time.Second)
		if err != nil {
			log.Fatalf("request failed: %v", err)
		}
		log.Printf("request/reply → %q", resp.Data)

	default:
		log.Fatalf("unknown role %q", *role)
	}
}
