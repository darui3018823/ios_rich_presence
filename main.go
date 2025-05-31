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
		log.Println("Unauthorized (/clear-rpc)")
		http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
		return
	}

	app := payload.Data.App
	if app == "" {
		http.Error(w, "App name is required", http.StatusBadRequest)
		return
	}

	// 明示的にRPC削除
	err := exec.Command("./python/del_rpc.exe", app).Run()
	if err != nil {
		log.Println("clear_rpc 実行エラー:", err)
		http.Error(w, "Failed to clear RPC", http.StatusInternalServerError)
	}

	if pid, ok := appPIDs[app]; ok {
		err := exec.Command("taskkill", "/PID", fmt.Sprint(pid), "/F").Run()
		if err != nil {
			log.Println("taskkillエラー:", err)
		} else {
			log.Printf("プロセス %s (PID %d) を終了しました\n", app, pid)
			delete(appPIDs, app)
		}
	} else {
		log.Println("記録されたPIDがありません")
	}

	w.Write([]byte("RPC cleared"))
}

func main() {
	fmt.Println("iOS ShortCut DiscordRP Server v2.0.1")
	http.HandleFunc("/set-rpc", handleSetRPC)
	http.HandleFunc("/clear-rpc", handleClearRPC)
	log.Println("サーバー起動中 (http://localhost:8080)")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
