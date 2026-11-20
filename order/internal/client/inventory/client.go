package inventory

import (
	"context"
	"time"

	"github.com/riteshdhurwey/scalable-ecommerce-platform/order/internal/models"
	catalogpb "github.com/riteshdhurwey/scalable-ecommerce-platform/proto/catalog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)
type Inventory struct {
	conn *grpc.ClientConn
	client catalogpb.CatalogServiceClient
}

func InitInventoryClient(addr string) (*Inventory,error){
	connection,err := grpc.NewClient(addr,grpc.WithTransportCredentials(insecure.NewCredentials()),)

	if err != nil {
		return nil,err
	}

	return &Inventory{conn: connection,client: catalogpb.NewCatalogServiceClient(connection)},nil
}

func (i *Inventory) ReverseInventory(ctx context.Context,req models.ReservationRequest)(*catalogpb.ReserveInventoryResponse,error){

	ctx,cancel := context.WithTimeout(ctx,5 * time.Second)
	defer cancel()

	items := make([]*catalogpb.ReservationItem, 0, len(req.Items))

	for _, item := range req.Items {
		items = append(items, &catalogpb.ReservationItem{
			ProductId: item.ProductID,
			Quantity:  int32(item.Quantity),
		})
	}

	return i.client.ReserveInventory(ctx,&catalogpb.ReserveInventoryRequest{OrderId: req.OrderID,ReservationItems: items})
}

func (i *Inventory) Close() {
	i.conn.Close()
}