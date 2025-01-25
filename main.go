package main

import (
	"hot-coffee/internal/handler"
	"log"
	"net/http"
)

func main() {
	// Регистрация маршрутов для создания заказов
	http.HandleFunc("/orders", handler.CreateOrder)

	// Запуск HTTP-сервера
	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
