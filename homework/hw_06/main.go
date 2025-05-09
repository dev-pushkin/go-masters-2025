package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/fibonacci", handler)
	fmt.Println("Server is running on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Query().Get("N")

	n, err := strconv.Atoi(param)
	if err != nil {
		http.Error(w, "invalid parameter", http.StatusBadRequest)
		return
	}

	if n > 40 {
		http.Error(w, "parameter too large", http.StatusBadRequest)
		return
	}

	result := fibonacci(n)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"result": %d}`, result)))
}

func fibonacci(n int) int {
	if n <= 0 {
		return 0
	} else if n == 1 {
		return 1
	}
	return fibonacci(n-1) + fibonacci(n-2)
}
