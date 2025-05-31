// 2025 iOS ShortCut DiscordRP: darui3018823 All rights reserved.
// All works created by darui3018823 associated with this repository are the intellectual property of darui3018823.
// Packages and other third-party materials used in this repository are subject to their respective licenses and copyrights.

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

type InnerData struct {
	App    string `json:"app"`
	Device string `json:"device"`
	User   string `json:"user"`
}

type Payload struct {
	Token string    `json:"token"`
	Data  InnerData `json:"data"`
}

var expectedToken = os.Getenv("RPC_AUTH_TOKEN")

func handleSetRPC(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	fmt.Println("受信Rawボディ (/set-rpc):")
	fmt.Println(string(body))

	var payload Payload
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		fmt.Println("JSON解析エラー:", err)
		return
	}

	if payload.Token != expectedToken {
		payloadJson, _ := json.MarshalIndent(payload, "", "  ")
		log.Println("Unauthorized (body token mismatch)")
		log.Println("Request payload:\n" + string(payloadJson))
		http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
		return
	}

	fmt.Printf("RPC update requested: app=%s, device=%s, user=%s\n", payload.Data.App, payload.Data.Device, payload.Data.User)
	cmd := exec.Command("./python/set_rpc.exe", payload.Data.App, payload.Data.Device, payload.Data.User)
	err := cmd.Run()
	if err != nil {
		log.Println("Pythonバイナリ実行エラー:", err)
		http.Error(w, "RPC update failed", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("OK"))
}

func handleClearRPC(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	fmt.Println("受信Rawボディ (/clear-rpc):")
	fmt.Println(string(body))

	var payload Payload
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		fmt.Println("JSON解析エラー:", err)
		return
	}

	if payload.Token != expectedToken {
		payloadJson, _ := json.MarshalIndent(payload, "", "  ")
		log.Println("Unauthorized (clear token mismatch)")
		log.Println("Request payload:\n" + string(payloadJson))
		http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
		return
	}

	fmt.Printf("RPC clear requested: app=%s\n", payload.Data.App)
	cmd := exec.Command("./python/clear_rpc.exe", payload.Data.App)
	err := cmd.Run()
	if err != nil {
		log.Println("Pythonクリアバイナリ実行エラー:", err)
		http.Error(w, "RPC clear failed", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("RPC Cleared"))
}

func main() {
	fmt.Println("iOS ShortCut DiscordRP Server v1.6.0")
	http.HandleFunc("/set-rpc", handleSetRPC)
	http.HandleFunc("/clear-rpc", handleClearRPC)
	log.Println("サーバー起動中 (http://localhost:8080)")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
