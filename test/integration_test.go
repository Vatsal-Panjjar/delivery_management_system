package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestAuthCreateAndProgress(t *testing.T) {
	client := &http.Client{}
	// adjust host if needed
	base := "http://localhost:8080"

	// register
	reg := map[string]string{"username":"testuser","email":"test@example.com","password":"pass123","role":"customer"}
	b,_:=json.Marshal(reg)
	resp,err:=client.Post(base+"/api/auth/register","application/json",bytes.NewReader(b))
	if err!=nil { t.Fatal(err) }
	if resp.StatusCode!=201 && resp.StatusCode!=200 { t.Fatalf("register status %d", resp.StatusCode) }

	// login
	login := map[string]string{"email":"test@example.com","password":"pass123"}
	b,_=json.Marshal(login)
	resp,err=client.Post(base+"/api/auth/login","application/json",bytes.NewReader(b))
	if err!=nil { t.Fatal(err) }
	defer resp.Body.Close()
	var data map[string]string
	json.NewDecoder(resp.Body).Decode(&data)
	token := data["token"]
	if token=="" { t.Fatal("no token") }

	// create order
	order := map[string]string{"description":"integration order"}
	b,_=json.Marshal(order)
	req, _ := http.NewRequest("POST", base+"/api/orders", bytes.NewReader(b))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp,err=client.Do(req)
	if err!=nil { t.Fatal(err) }
	if resp.StatusCode!=200 && resp.StatusCode!=201 { t.Fatalf("create order status %d", resp.StatusCode) }

	// wait for progression (tracker uses seconds)
	time.Sleep(13 * time.Second)

	// fetch orders
	req, _ = http.NewRequest("GET", base+"/api/orders", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp,err = client.Do(req)
	if err!=nil { t.Fatal(err) }
	var orders []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&orders)
	if len(orders)==0 { t.Fatal("no orders") }
	// ensure status progressed to at least "in_transit" or "delivered"
	status := orders[0]["status"].(string)
	if status == "pending" { t.Fatalf("status did not progress, still %s", status) }
}
