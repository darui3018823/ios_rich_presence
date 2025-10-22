// 2025 iOS ShortCut DiscordRPC Server: darui3018823 All rights reserved.
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
	"sync"
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

var (
	appMu   sync.Mutex
	appPIDs = make(map[string][]int)
)

func addAppPID(app string, pid int) {
	appMu.Lock()
	appPIDs[app] = append(appPIDs[app], pid)
	appMu.Unlock()
}

func popAppPIDs(app string) []int {
	appMu.Lock()
	pids := appPIDs[app]
	delete(appPIDs, app)
	appMu.Unlock()
	return pids
}

func takeAllAppPIDs() map[string][]int {
	appMu.Lock()
	defer appMu.Unlock()
	all := appPIDs
	appPIDs = make(map[string][]int)
	return all
}

func killProcesses(app string, pids []int) {
	for _, pid := range pids {
		if pid == 0 {
			continue
		}
		if err := exec.Command("taskkill", "/PID", fmt.Sprint(pid), "/f", "/t").Run(); err != nil {
			log.Printf("taskkillエラー: app=%s pid=%d err=%v\n", app, pid, err)
		} else {
			log.Printf("プロセス %s (PID %d) を終了しました\n", app, pid)
		}
	}
}

func handleSetRPC(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)

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

	app := payload.Data.App

	// 同じアプリ名で既に起動しているRPCを終了
	if existing := popAppPIDs(app); len(existing) > 0 {
		log.Printf("既存のRPCを終了します: app=%s pids=%v\n", app, existing)
		killProcesses(app, existing)
	}

	cmd := exec.Command("./python/set_rpc.exe", app, payload.Data.Device, payload.Data.User)
	if err := cmd.Start(); err != nil {
		log.Println("set_rpc起動失敗:", err)
		http.Error(w, "Failed to start RPC", http.StatusInternalServerError)
		return
	}

	// 起動成功 → PID記録
	addAppPID(app, cmd.Process.Pid)
	log.Printf("RPC started for %s with PID %d\n", app, cmd.Process.Pid)
	w.Write([]byte("OK"))
}

func handleClearRPC(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)

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

	if pids := popAppPIDs(app); len(pids) > 0 {
		killProcesses(app, pids)
	} else {
		log.Println("記録されたPIDがありません")
	}

	w.Write([]byte("RPC cleared"))
}

func handleAlldelRPC(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)

	var payload Payload
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		fmt.Println("JSON解析エラー:", err)
		return
	}

	if payload.Token != expectedToken {
		log.Println("Unauthorized (/all_del)")
		http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
		return
	}

	all := takeAllAppPIDs()
	for app, pids := range all {
		if err := exec.Command("./python/del_rpc.exe", app).Run(); err != nil {
			log.Printf("clear_rpc 実行エラー(app=%s): %v\n", app, err)
		}
		killProcesses(app, pids)
	}

	w.Write([]byte("All RPCs cleared"))
}

func main() {
	fmt.Println("iOS ShortCut DiscordRP Server v2.1.0")
	http.HandleFunc("/set-rpc", handleSetRPC)
	http.HandleFunc("/clear-rpc", handleClearRPC)
	http.HandleFunc("/all_del", handleAlldelRPC)
	log.Println("サーバー起動中 (http://localhost:8080)")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
