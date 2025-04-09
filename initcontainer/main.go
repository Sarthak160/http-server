// app1.go
package main

import (
	"fmt"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from app1! PID: %d\n", os.Getpid())
}

func main() {
	http.HandleFunc("/", handler)
	go func() {
		fmt.Printf("App1 PID: %d\n", os.Getpid())
	}()

	fmt.Println("App1 is running on port 8080...")
	http.ListenAndServe(":8080", nil)
}
