package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Payload struct {
	Token  string `json:"token"`
	App    string `json:"app"`
	Device string `json:"device"`
	User   string `json:"user"`
}

var expectedToken = os.Getenv("RPC_AUTH_TOKEN")

func handler(w http.ResponseWriter, r *http.Request) {
	var payload Payload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if payload.Token != expectedToken {
		// payloadを見やすくJSON整形して出力
		payloadJson, _ := json.MarshalIndent(payload, "", "  ")
		log.Println("Unauthorized (body token mismatch)")
		log.Println("Request payload:\n" + string(payloadJson))
		http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
		return
	}

	fmt.Printf("RPC update requested: app=%s, device=%s, user=%s\n", payload.App, payload.Device, payload.User)
	w.Write([]byte("OK"))
}

func main() {
	http.HandleFunc("/", handler)
	log.Println("サーバー起動中 (http://localhost:8080)")
	http.ListenAndServe(":8080", nil)
}
