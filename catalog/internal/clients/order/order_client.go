package order

import (
	"context"
	"log"

	"time"

	orderpb "github.com/riteshdhurwey/scalable-ecommerce-platform/proto/order"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Order struct {
	conn *grpc.ClientConn
	client orderpb.OrderServiceClient
}

func InitOrderClient(addr string) (*Order, error) {
    target := "dns:///" + addr

    log.Printf("Initializing Order gRPC client: %q", target)

    connection, err := grpc.NewClient(
        target,
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    )

    if err != nil {
        return nil, err
    }

    return &Order{
        conn:   connection,
        client: orderpb.NewOrderServiceClient(connection),
    }, nil
}

func (o *Order) UpdateOrderStatus(ctx context.Context,orderIDs []int64)(*orderpb.GetOrderStatusResponse,error){
	ctx,cancel := context.WithTimeout(ctx, 5 * time.Second)

	defer cancel()

	request := make([]*orderpb.OrderID,0,len(orderIDs))

	for _,order := range orderIDs {

		request = append(request, &orderpb.OrderID{OrderID: order})
	}

	return o.client.UpdateExpiredOrderStatus(ctx,&orderpb.GetOrderStatusRequest{OrderIDs: request})
}

func (o *Order) Close()(error){
	return o.conn.Close()
}