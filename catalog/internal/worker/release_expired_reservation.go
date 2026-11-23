package worker

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	orderClient "github.com/riteshdhurwey/scalable-ecommerce-platform/services/catalog/internal/clients/order"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/catalog/internal/models"
)

// Returns true if expired reservations were found and processed, false if none exist
func ReleaseExpiredOrders(ctx context.Context, db *pgxpool.Pool, client *orderClient.Order) bool {
	log.Print("worker started")
	var expiredReservations []models.ExpiredReservation

	query := `SELECT id, order_id, product_id, quantity FROM inventory_reservations WHERE status = 'RESERVED' AND expires_at < NOW() ORDER BY expires_at LIMIT 50 FOR UPDATE SKIP LOCKED;`

	tx, err := db.Begin(ctx)
	if err != nil {
		log.Printf("error starting transaction: %v", err)
		return false
	}
	defer tx.Rollback(ctx)

	results, err := tx.Query(ctx, query)
	if err != nil {
		log.Printf("error collecting expired reservations: %v", err)
		return false
	}

	for results.Next() {
		var r models.ExpiredReservation
		if err := results.Scan(&r.ID, &r.OrderID, &r.ProductID, &r.Quantity); err != nil {
			results.Close()
			log.Printf("error scanning expired reservations: %v", err)
			return false
		}
		expiredReservations = append(expiredReservations, r)
	}

	if err = results.Err(); err != nil {
		log.Printf("error iterating expired reservations: %v", err)
		return false
	}
	results.Close()

	// Nothing to process -> tell caller to freeze for 5 minutes
	if len(expiredReservations) == 0 {
		return false
	}

	uniqueOrderIDs := make([]int64, 0, len(expiredReservations))
	seen := make(map[int64]struct{}, len(expiredReservations))
	for _, reservation := range expiredReservations {
		if _, exists := seen[reservation.OrderID]; !exists {
			seen[reservation.OrderID] = struct{}{}
			uniqueOrderIDs = append(uniqueOrderIDs, reservation.OrderID)
		}
	}

	response, err := client.UpdateOrderStatus(ctx, uniqueOrderIDs)
	
	if err != nil {
		log.Printf("error verifying expired reservations: %v", err)
		return false
	}

	if !response.StatusUpdated {
		log.Println("order status not updated")
		return false
	}

	if response.OrderIDs == nil {
		return false
	}

	orderExpired := make(map[int64]bool)

	for _,id := range response.OrderIDs {

		orderExpired[id.OrderID] = true
	}
	

	updateQuery := `UPDATE inventory_reservations SET status = 'RELEASED' WHERE id = $1 AND status = 'RESERVED'`
	releaseStockQuery := `UPDATE products SET stock_quantity = stock_quantity + $1 WHERE id = $2;`

	for _, item := range expiredReservations {
		if !orderExpired[item.OrderID] {
			continue
		}

		res, err := tx.Exec(ctx, updateQuery, item.ID)
		if err != nil {
			log.Printf("error updating reservation status: %v", err)
			return false
		}

		if res.RowsAffected() == 1 {
			stockRes, err := tx.Exec(ctx, releaseStockQuery, item.Quantity, item.ProductID)
			if err != nil {
				log.Printf("error releasing stock: %v", err)
				return false
			}
			if stockRes.RowsAffected() != 1 {
				log.Print("product not found while releasing stock")
				return false
			}
		}
	}

	if err = tx.Commit(ctx); err != nil {
		log.Printf("error committing release of expired reservations: %v", err)
		return false
	}

	log.Print("worker done")
	return true
}

func StartReservationWorker(ctx context.Context, db *pgxpool.Pool, client *orderClient.Order) {
	for {
		if ctx.Err() != nil {
			return
		}

		hasMoreWork := ReleaseExpiredOrders(ctx, db, client)

		if !hasMoreWork {
			log.Println("No expired reservations found. Freezing worker for 5 minutes...")
			timer := time.NewTimer(5 * time.Minute)
			select {
			case <-timer.C:
				continue
			case <-ctx.Done():
				timer.Stop()
				return
			}
		}
	}
}
