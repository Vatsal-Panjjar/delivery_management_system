package main
"github.com/jmoiron/sqlx"
"github.com/go-redis/redis/v8"


"github.com/Vatsal-Panjjar/delivery_management_system/internal/cache"
"github.com/Vatsal-Panjjar/delivery_management_system/internal/db"
"github.com/Vatsal-Panjjar/delivery_management_system/internal/handlers"
"github.com/Vatsal-Panjjar/delivery_management_system/internal/tracker"
)


func main() {
dbURL := os.Getenv("DATABASE_URL")
redisAddr := os.Getenv("REDIS_ADDR")
port := os.Getenv("PORT")
if port == "" {
port = "8080"
}


if dbURL == "" || redisAddr == "" {
log.Fatal("DATABASE_URL and REDIS_ADDR must be set")
}


sqlxDB, err := db.NewPostgres(dbURL)
if err != nil {
log.Fatal(err)
}


// Load schema
schema, err := os.ReadFile("./migrations/init.sql")
if err == nil {
db.ExecSchema(sqlxDB, string(schema))
}


rclient := redis.NewClient(&redis.Options{Addr: redisAddr})
if _, err := rclient.Ping(rclient.Context()).Result(); err != nil {
log.Fatalf("redis ping failed: %v", err)
}


trackerSvc := tracker.NewTracker(sqlxDB, rclient, 4)
cacheClient := cache.NewRedis(redisAddr)


srv := &handlers.Server{DB: sqlxDB, Tracker: trackerSvc, Cache: cacheClient}


r := chi.NewRouter()
fs := http.FileServer(http.Dir("./web/static"))
r.Handle("/static/*", http.StripPrefix("/static/", fs))


r.Get("/", srv.IndexPage)
r.Get("/login", srv.LoginPage)
r.Get("/register", srv.RegisterPage)


r.Post("/api/register", srv.Register)
r.Post("/api/login", srv.Login)


r.Group(func(r chi.Router) {
r.Use(srv.AuthMiddleware)
r.Post("/api/orders", srv.CreateOrder)
r.Post("/api/orders/{id}/cancel", srv.CancelOrder)
r.Get("/api/orders/{id}", srv.GetOrder)
r.Get("/api/orders", srv.ListOrders)
r.Get("/dashboard", srv.DashboardPage)
r.Get("/admin", srv.AdminPage)
})


srvAddr := fmt.Sprintf(":%s", port)
s := &http.Server{
Addr: srvAddr,
Handler: r,
ReadTimeout: 15 * time.Second,
WriteTimeout: 15 * time.Second,
}


fmt.Printf("listening on %s\n", srvAddr)
log.Fatal(s.ListenAndServe())
}
