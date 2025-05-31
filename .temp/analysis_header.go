package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("=== 新しいリクエスト ===")
	fmt.Println("メソッド:", r.Method)
	fmt.Println("URL:", r.URL.Path)

	fmt.Println("--- ヘッダー一覧 ---")
	for name, values := range r.Header {
		for _, value := range values {
			fmt.Printf("%s: %s\n", name, value)
		}
	}

	fmt.Println("--- ボディ ---")
	defer r.Body.Close()
	buf := make([]byte, r.ContentLength)
	r.Body.Read(buf)
	fmt.Println(string(buf))

	w.Write([]byte("OK"))
}

func main() {
	http.HandleFunc("/", handler)
	log.Println("テストサーバー起動中 (http://localhost:8080)")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
