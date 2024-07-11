package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/time", handler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Ошибка запуска сервера", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now().Format(time.RFC1123)
	fmt.Fprint(w, "Текущее время: ", currentTime)
}
