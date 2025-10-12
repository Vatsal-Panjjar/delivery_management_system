module delivery_management_system

go 1.19

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-redis/redis/v8 v8.11.5
	github.com/gorilla/mux v1.8.0
	github.com/lib/pq v1.10.4
	github.com/joho/godotenv v1.4.0
)

replace github.com/dgrijalva/jwt-go => github.com/golang-jwt/jwt v4.0.0
