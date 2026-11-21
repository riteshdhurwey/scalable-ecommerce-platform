package order

import (
	"context"

	"time"

	orderpb "github.com/riteshdhurwey/scalable-ecommerce-platform/proto/order"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Order struct {
	conn *grpc.ClientConn
	client orderpb.OrderServiceClient
}

func InitOrderClient(addr string)(*Order,error){
	connection,err := grpc.NewClient(addr,grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil,err
	}

	return &Order{conn: connection,client: orderpb.NewOrderServiceClient(connection)},nil
}

func (o *Order) GetOrderStatus(ctx context.Context,orderIDs []int64)(*orderpb.GetOrderStatusResponse,error){
	ctx,cancel := context.WithTimeout(ctx, 5 * time.Second)

	defer cancel()

	request := make([]*orderpb.OrderID,0,len(orderIDs))

	for _,order := range orderIDs {

		request = append(request, &orderpb.OrderID{OrderID: order})
	}

	return o.client.GetOrderStatus(ctx,&orderpb.GetOrderStatusRequest{OrderIDs: request})
}

func (o *Order) Close()(error){
	return o.conn.Close()
}