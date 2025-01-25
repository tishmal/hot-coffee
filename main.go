package main

import (
	"flag"
	"fmt"
	"hot-coffee/add"
	"hot-coffee/internal/handler"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {
	port := flag.Int("port", 8080, "Port number to listen on")
	help := flag.Bool("help", false, "Show help")
	flag.Parse()

	if *help {
		add.PrintUsage()
		return
	}

	http.HandleFunc("/orders", handler.GetAllOrders)

	addr := fmt.Sprintf(":%d", *port)

	// Запуск браузера
	go add.OpenBrowser(addr)

	// Запуск HTTP сервера
	log.Printf("The server is running on the port %s...\n", addr)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// Обработчик запросов
func handleRequests(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем заголовки для правильного отображения HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Чтение HTML файла
	html, err := ioutil.ReadFile("index.html")
	if err != nil {
		http.Error(w, "Failed to load HTML file", http.StatusInternalServerError)
		return
	}

	// Вставляем путь из URL в HTML-страницу
	username := r.URL.Path[1:]
	if username == "" {
		username = "Guest"
	}

	// Заменяем метку <username> на значение из URL
	htmlStr := strings.Replace(string(html), "<username>", username, 1)

	// Отправляем страницу клиенту
	fmt.Fprint(w, htmlStr)
}
