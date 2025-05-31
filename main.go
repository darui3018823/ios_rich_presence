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
var appPIDs = make(map[string]int)

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

	cmd := exec.Command("./python/set_rpc.exe", payload.Data.App, payload.Data.Device, payload.Data.User)
	if err := cmd.Start(); err != nil {
		log.Println("set_rpc起動失敗:", err)
		http.Error(w, "Failed to start RPC", http.StatusInternalServerError)
		return
	}

	// 起動成功 → PID記録
	appPIDs[payload.Data.App] = cmd.Process.Pid
	log.Printf("RPC started for %s with PID %d\n", payload.Data.App, cmd.Process.Pid)
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

	pid, ok := appPIDs[payload.Data.App]
	if !ok {
		http.Error(w, "No RPC running for this app", http.StatusNotFound)
		return
	}

	// taskkill実行（Windows専用）
	err := exec.Command("taskkill", "/PID", fmt.Sprint(pid), "/F").Run()
	if err != nil {
		log.Println("taskkillエラー:", err)
		http.Error(w, "Failed to kill RPC process", http.StatusInternalServerError)
		return
	}

	delete(appPIDs, payload.Data.App)
	log.Printf("RPC process for %s (PID %d) terminated.\n", payload.Data.App, pid)
	w.Write([]byte("Cleared"))
}

func main() {
	fmt.Println("iOS ShortCut DiscordRP Server v1.7.0")
	http.HandleFunc("/set-rpc", handleSetRPC)
	http.HandleFunc("/clear-rpc", handleClearRPC)
	log.Println("サーバー起動中 (http://localhost:8080)")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
