package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	calculator "explore-grpc/proto"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type server struct {
	calculator.UnimplementedCalculatorServer
}

func (s *server) SquareRoot(ctx context.Context, req *calculator.SquareRootRequest) (*calculator.SquareRootResponse, error) {
	num := req.GetNumber()
	return &calculator.SquareRootResponse{Result: int32(num * num)}, nil
}

func (s *server) Sum(ctx context.Context, req *calculator.SumRequest) (*calculator.SumResponse, error) {
	num1, num2 := req.GetNum1(), req.GetNum2()
	return &calculator.SumResponse{Sum: num1 + num2}, nil
}

func main() {
	// command-line options:
	// gRPC server endpoint
	grpcServerEndpoint := flag.String("grpc-server-endpoint", "localhost:50051", "gRPC server endpoint")
	flag.Parse()
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	listner, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("unable to listen to the port")
	}

	srv := grpc.NewServer()
	calculator.RegisterCalculatorServer(srv, &server{})
	fmt.Println("Server is listening on port 50051")

	go func() {
		if err := srv.Serve(listner); err != nil {
			log.Fatalf("Failed to server: %v", err)
		}
	}()

	fmt.Println("are you still running bruh!!!!")

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := calculator.RegisterCalculatorHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts); err != nil {
		fmt.Println("error", err)
	}

	if err := http.ListenAndServe(":8085", mux); err != nil {
		fmt.Println("unable to get listen and serve http", err)
	}
}
