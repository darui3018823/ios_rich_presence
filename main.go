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

type OuterPayload struct {
	Data string `json:"data"`
}

var expectedToken = "Bearer " + os.Getenv("RPC_AUTH_TOKEN")

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != expectedToken {
		http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
		return
	}

	var outer OuterPayload
	bodyBytes, _ := io.ReadAll(r.Body)

	// とりあえず保存（後で検証用に使える）
	fmt.Println("--- Raw JSON ---")
	fmt.Println(string(bodyBytes))

	err := json.Unmarshal(bodyBytes, &outer)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// そのままPythonに渡す
	cmd := exec.Command("python3", "dispatch.py", outer.Data)
	err = cmd.Run()
	if err != nil {
		log.Println("Pythonエラー:", err)
		http.Error(w, "RPC更新失敗", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("OK"))
}

func main() {
	http.HandleFunc("/", handler)
	log.Println("サーバー起動中 (http://localhost:8080)")
	http.ListenAndServe(":8080", nil)
}
