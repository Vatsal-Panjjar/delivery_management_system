package handlers

import (
	"fmt"
	"log"
	"time"

	"github.com/Vatsal-Panjjar/delivery_management_system/internal/db"
)

type Tracker struct {
	store *db.Store
}

func NewTracker(store *db.Store) *Tracker {
	return &Tracker{store: store}
}

func (t *Tracker) StartOrderTracking(orderID int) {
	go func() {
		statuses := []string{"pending", "in_transit", "delivered"}
		for _, s := range statuses {
			time.Sleep(5 * time.Second)
			_, err := t.store.DB.Exec(`UPDATE orders SET status=$1, updated_at=NOW() WHERE id=$2`, s, orderID)
			if err != nil {
				log.Println("Error updating order:", err)
				return
			}
			fmt.Printf("📦 Order %d updated to %s\n", orderID, s)
		}
	}()
}
