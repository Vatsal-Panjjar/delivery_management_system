package workers

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Vatsal-Panjjar/delivery_management_system/internal/cache"
	"github.com/Vatsal-Panjjar/delivery_management_system/internal/db"
)

type OrderTracker struct {
	store *db.Store
	cache *cache.Cache

	queue chan int
	wg    sync.WaitGroup
	quit  chan struct{}
}

func NewOrderTracker(store *db.Store, cacheClient *cache.Cache) *OrderTracker {
	return &OrderTracker{
		store: store,
		cache: cacheClient,
		queue: make(chan int, 200),
		quit:  make(chan struct{}),
	}
}

func (t *OrderTracker) Start() {
	go func() {
		for {
			select {
			case id := <-t.queue:
				t.wg.Add(1)
				go t.processOrder(id)
			case <-t.quit:
				return
			}
		}
	}()
}

func (t *OrderTracker) Stop() {
	close(t.quit)
	t.wg.Wait()
}

func (t *OrderTracker) Enqueue(orderID int) {
	select {
	case t.queue <- orderID:
	default:
		fmt.Printf("order tracker queue full; dropping %d\n", orderID)
	}
}

func (t *OrderTracker) processOrder(orderID int) {
	defer t.wg.Done()

	steps := []struct {
		status string
		delay  time.Duration
	}{
		{"accepted", 3 * time.Second},
		{"in_transit", 5 * time.Second},
		{"delivered", 4 * time.Second},
	}

	for _, s := range steps {
		time.Sleep(s.delay)

		tx := t.store.DB.MustBegin()
		var cur string
		err := tx.Get(&cur, "SELECT status FROM orders WHERE id=$1 FOR UPDATE", orderID)
		if err != nil {
			tx.Rollback()
			return
		}
		if cur == "cancelled" || cur == "delivered" {
			tx.Rollback()
			return
		}
		_, err = tx.Exec("UPDATE orders SET status=$1, updated_at=now() WHERE id=$2", s.status, orderID)
		if err != nil {
			tx.Rollback()
			return
		}
		if err := tx.Commit(); err != nil {
			return
		}
		_ = t.cache.Set(context.Background(), fmt.Sprintf("order:%d:status", orderID), s.status, time.Hour)
	}
}
