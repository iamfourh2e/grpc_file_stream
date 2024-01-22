package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"stream_server/pb"

	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc"
)

type fileServiceServer struct {
	pb.UnimplementedFileServiceServer
}

func (s *fileServiceServer) Upload(stream pb.FileService_UploadServer) error {
	md, _ := metadata.FromIncomingContext(stream.Context())
	filename := md["filename"]
	file, err := os.Create(md["filename"][0])
	fmt.Printf("%v", md)
	if err != nil {
		return err
	}
	defer file.Close()

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		fmt.Printf("%v", req.Chunk)
		_, err = file.Write(req.Chunk)

		if err != nil {
			return err
		}
	}

	return stream.SendAndClose(&pb.UploadResponse{
		Filename: filename[0],
		Size:     2,
	})
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterFileServiceServer(s, &fileServiceServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatal("unable to start server", err)
	}
}
