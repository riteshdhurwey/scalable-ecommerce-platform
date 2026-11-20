package server

import (
	"log"
	"net"

	cartpb "github.com/riteshdhurwey/scalable-ecommerce-platform/proto/cart"
	"google.golang.org/grpc"
)

func StartGrpcServer(addr string,cartHandler cartpb.CartServiceServer) error {
	listener,err := net.Listen("tcp",addr)

	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()

	cartpb.RegisterCartServiceServer(grpcServer,cartHandler)

	log.Printf("cart grpc server Listening on %s",addr)

	return grpcServer.Serve(listener)
}