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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r, err := client.SayHello(ctx, &greetv1.HelloRequest{Name: "Gofer"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Unary Response: %s", r.GetMessage())

	// server streaming call
	stream, err := client.ManyHellos(ctx, &greetv1.ManyHellosRequest{
		Name:       "StreamFan",
		NumReplies: 5,
	})
	if err != nil {
		log.Printf("error starting steam: %v", err)
	}
	log.Println("--- Starting Stream Reseption ---")
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			break //поток закончился успешно!
		}
		if err != nil {
			log.Fatalf("error receiving stream: %v", err)
		}
		log.Printf("Stream Message: %s", in.GetMessage())
	}
	log.Println("--- Stream Finish ---")
}
