package main

import (
	"context"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
	"homework9/internal/adapters/adrepo"
	"homework9/internal/adapters/userrepo"
	"homework9/internal/app"
	grpcPort "homework9/internal/ports/grpc"
	"homework9/internal/ports/httpgin"
	"log"
	"net"
)

func main() {
	server := httpgin.NewHTTPServer(":18080", app.NewApp(adrepo.New(), userrepo.New()))
	err := server.Listen()
	if err != nil {
		panic(err)
	}

	port := ":50054"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("fail")
	}
	service := grpcPort.NewService()
	serverGrpc := grpc.NewServer(grpc.ChainUnaryInterceptor(LoggerInterceptor, grpc_recovery.UnaryServerInterceptor()))
	defer serverGrpc.GracefulStop()

	grpcPort.RegisterAdServiceServer(serverGrpc, service)

	if err = serverGrpc.Serve(lis); err != nil {
		log.Fatal("failed to serve")
	}
}

func LoggerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println(info.FullMethod)
	return handler(ctx, req)
}
