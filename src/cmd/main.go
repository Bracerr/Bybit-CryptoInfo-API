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
	http.HandleFunc("/symbols", klinesHandler.GetAvailableSymbols)

	fmt.Println("Сервер запущен на http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
	}
}
