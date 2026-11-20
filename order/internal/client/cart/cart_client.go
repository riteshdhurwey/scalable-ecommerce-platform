package cart

import (
	"context"
	"time"

	cartpb "github.com/riteshdhurwey/scalable-ecommerce-platform/proto/cart"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Cart struct {
	conn   *grpc.ClientConn
	client cartpb.CartServiceClient
}

func InitCartClient(addr string) (*Cart, error) {

	connection, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	return &Cart{conn: connection, client: cartpb.NewCartServiceClient(connection)}, nil
}

func (c *Cart) GetCartItems(ctx context.Context, userID int64) (*cartpb.GetCartResponse, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	

	defer cancel()

	return c.client.GetCart(ctx, &cartpb.GetCartRequest{UserID: userID})
}

func (c *Cart) ClearCart(ctx context.Context,userID int64)(*cartpb.ClearCartResponse,error){

	ctx , cancel := context.WithTimeout(ctx, 5 * time.Second)

	defer cancel()

	return c.client.ClearCart(ctx,&cartpb.ClearCartRequest{UserID:userID})
}

func (c *Cart) Close() error {
	return c.conn.Close()
}
