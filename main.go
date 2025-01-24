package main

import (
	"flag"
	"fmt"
	"hot-coffee/add"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	port := flag.Int("port", 8080, "Port number to listen on")
	dir := flag.String("dir", "data", "Base directory for storage")
	help := flag.Bool("help", false, "Show help")
	flag.Parse()

	if *help {
		add.PrintUsage()
		return
	}

	http.HandleFunc("/", handleRequests) // Настройка маршрутов

	// Создаем директорию для данных, если она не существует
	if err := os.MkdirAll(*dir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	addr := fmt.Sprintf(":%d", *port)

	// Запуск браузера
	go add.OpenBrowser(addr)

	// Запуск HTTP сервера
	log.Printf("The server is running on the port %s...\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
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
