package server

import (
	"log"
	"net"

	catalogpb "github.com/riteshdhurwey/scalable-ecommerce-platform/proto/catalog"
	"google.golang.org/grpc"
)


func StartGRPCServer(
	addr string,
	catalogHandler catalogpb.CatalogServiceServer,
) error {

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()

	catalogpb.RegisterCatalogServiceServer(
		grpcServer,
		catalogHandler,

	)

	log.Printf("Catalog gRPC server listening on %s", addr)

	return grpcServer.Serve(listener)
}
