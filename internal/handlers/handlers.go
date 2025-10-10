package handlers

import (
    "encoding/json"
    "net/http"
    "sync"
    "time"

    "github.com/Vatsal-Panjiar/delivery_management_system/internal/repo"
    "github.com/Vatsal-Panjiar/delivery_management_system/internal/cache"
)





// DeliveryHandler implements endpoints and uses a sync.Mutex map to handle some concurrency safely
type DeliveryHandler struct{
Repo *repo.DeliveryRepo
Cache *cache.RedisCache
mu sync.Mutex // coarse-grained for demo; can be per-delivery map for finer control
}


func (h *DeliveryHandler) ListByStatus(w http.ResponseWriter, r *http.Request) {
    status := r.URL.Query().Get("status")
    if status == "" {
        status = "pending"
    }

    cacheKey := "deliveries:" + status
    if cached, err := h.Cache.Get(cacheKey); err == nil {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(cached))
        return
    }

    ds, err := h.Repo.ListByStatus(status, 100, 0) // example limit/offset
    if err != nil {
        http.Error(w, "failed to fetch deliveries", http.StatusInternalServerError)
        return
    }

    b, _ := json.Marshal(ds)
    _ = h.Cache.Set(cacheKey, b, 5*time.Minute)

    w.Header().Set("Content-Type", "application/json")
    w.Write(b)
}
