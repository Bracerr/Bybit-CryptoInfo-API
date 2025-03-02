package main

import (
	"API-CRYPT/src/handlers"
	"API-CRYPT/src/services"
	"fmt"
	"net/http"
)

func main() {
	klinesService := services.NewKlinesService()
	klinesHandler := handler.NewKlinesHandler(klinesService)

	http.HandleFunc("/klines", klinesHandler.GetKlines)
	http.HandleFunc("/klines/csv", klinesHandler.GetKlinesCSV)
	http.HandleFunc("/symbols", klinesHandler.GetAvailableSymbols)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../static"))))

	fmt.Println("Сервер запущен на http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
	}
}
