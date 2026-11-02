package main

import (
	"context"
	"io"
	"log"
	"time"

	greetv1 "github.com/andrey-korolev-ms/data-provider-service/gen/greet/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := greetv1.NewGreeterServiceClient(conn)

	// ==========================================
	// ТЕСТ 1: Unary Call (SayHello)
	// ==========================================
	log.Println("--- Test 1: Unary SayHello ---")
	ctxUnary, cancelUnary := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelUnary() // Обязательно освобождаем ресурсы

	r, err := client.SayHello(ctxUnary, &greetv1.HelloRequest{Name: "Gofer"})
	if err != nil {
		log.Printf("Error calling SayHello: %v", err)
	} else {
		log.Printf("Success! Response: %s", r.GetMessage())
	}

	// ==========================================
	// ТЕСТ 2: Server Streaming (ManyHellos)
	// ==========================================
	log.Println("\n--- Test 2: Server Streaming ManyHellos ---")

	// Создаем НОВЫЙ контекст с большим таймаутом для стрима
	// Если сервер шлет 5 сообщений по 1 сек каждое, нам нужно минимум 6-7 сек
	ctxStream, cancelStream := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelStream()

	stream, err := client.ManyHellos(ctxStream, &greetv1.ManyHellosRequest{
		Name:       "StreamFan",
		NumReplies: 5,
	})
	if err != nil {
		log.Fatalf("Failed to start stream: %v", err)
	}

	log.Println("Starting reception...")
	count := 0
	for {
		in, err := stream.Recv()

		// Конец потока — штатная ситуация
		if err == io.EOF {
			break
		}

		// Любая другая ошибка — проблема сети или сервера
		if err != nil {
			log.Fatalf("Error receiving from stream: %v", err)
		}

		count++
		log.Printf("[%d] Received: %s", count, in.GetMessage())
	}

	log.Printf("--- Stream Finished successfully after %d messages ---", count)
}
