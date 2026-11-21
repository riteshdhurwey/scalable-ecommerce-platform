package server

import (
	"log"
	"net"

	pb "github.com/riteshdhurwey/scalable-ecommerce-platform/proto/order"
	"google.golang.org/grpc"
)

func StartGrpcServer (addr string,orderHandler pb.OrderServiceServer) error {

	listen,err := net.Listen("tcp",addr)

	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()

	pb.RegisterOrderServiceServer(grpcServer,orderHandler)

	log.Printf("order gRPC server listening on %s",addr)

	return grpcServer.Serve(listen)
}
