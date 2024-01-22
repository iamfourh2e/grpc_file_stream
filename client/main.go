package main

import (
	"context"
	"grpc_client/pb"
	"io"
	"log"
	"os"

	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc"
)

const chunkSize = 4096 // Adjust chunk size as needed

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewFileServiceClient(conn)
	ctx := context.Background()
	md := metadata.Pairs("filename", "gg.dmg")

	ctx = metadata.NewOutgoingContext(ctx, md)
	stream, err := client.Upload(ctx)
	if err != nil {
		log.Fatalf("Failed to create stream: %v", err)
	}

	file, err := os.Open("./gg.dmg")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	chunkNumber := 0
	for {
		chunk := make([]byte, chunkSize)
		n, err := file.Read(chunk)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Failed to read chunk: %v", err)
		}

		chunkNumber++
		err = stream.Send(&pb.FileUploadRequest{
			Filename:    "file.dmg",
			Chunk:       chunk[:n],
			ChunkNumber: int32(chunkNumber),
		})
		if err != nil {
			log.Fatalf("Failed to send chunk: %v", err)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Failed to receive response: %v", err)
	}

	log.Printf("File uploaded successfully: %s (%d bytes)", res.Filename, res.Size)
}
