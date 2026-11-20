package catalogclient

import (
	"context"
	"time"

	catalogpb "github.com/riteshdhurwey/scalable-ecommerce-platform/proto/catalog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn *grpc.ClientConn
	client catalogpb.CatalogServiceClient
}

func InitCatalogClient(addr string)(*Client,error){
	connection,err := grpc.NewClient(addr,grpc.WithTransportCredentials(insecure.NewCredentials()),)

	if err != nil {
		return nil,err
	}

	return &Client{conn: connection,client: catalogpb.NewCatalogServiceClient(connection)},nil

}

func (c *Client) GetProduct(ctx context.Context,productID int64)(*catalogpb.Product,error){

	ctx,cancel := context.WithTimeout(ctx,10 * time.Second)

	defer cancel()

	return c.client.GetProduct(ctx,&catalogpb.GetProductRequest{Id: productID})
}

func (c *Client) GetProductByIDs(ctx context.Context,productIDs[]int64)(*catalogpb.GetProductsResponse,error){
	ctx,cancel := context.WithTimeout(ctx,5 * time.Second)

	defer cancel()

	return c.client.GetProducts(ctx,&catalogpb.GetProductsRequest{Id: productIDs})
}

func (c *Client)Close()error{
	return c.conn.Close()
}