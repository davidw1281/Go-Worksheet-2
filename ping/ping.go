package main

import (
	"fmt"
	"log"
	"os"
	"runtime/trace"
	"time"
)

func foo(channel chan string) {
	for {
		fmt.Println("Foo is sending: ping")
		channel <- "ping"

		<-channel
		fmt.Println("Foo has received: Pong")
	}
}

func bar(channel chan string) {
	for {
		<-channel
		fmt.Println("Bar has received: ping")

		fmt.Println("Bar is sending: pong")
		channel <- "pong"
	}
}

func pingPong() {
	channel := make(chan string)
	go foo(channel) // Nil is similar to null. Sending or receiving from a nil chan blocks forever.
	go bar(channel)
	time.Sleep(500 * time.Millisecond)
}

func main() {
	f, err := os.Create("trace.out")
	if err != nil {
		log.Fatalf("failed to create trace output file: %v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatalf("failed to close trace file: %v", err)
		}
	}()

	if err := trace.Start(f); err != nil {
		log.Fatalf("failed to start trace: %v", err)
	}
	defer trace.Stop()

	pingPong()
}
