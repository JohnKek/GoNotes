package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	name := params.Get("name")
	body := fmt.Sprintf("Hello, %s", name)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(body))
}

func handlerOne(w http.ResponseWriter, r *http.Request) {
	// ...
}

func handlerTwo(w http.ResponseWriter, r *http.Request) {
	// ...
}

func handlerThree(w http.ResponseWriter, r *http.Request) {
	// ...
}

func startServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/one", handlerOne)
	mux.HandleFunc("/two", handlerTwo)
	mux.HandleFunc("/three", handlerThree)
	server := httptest.NewServer(mux)
	return server
}

func main() {
	server := startServer()
	defer server.Close()

	client := server.Client()
	uri := server.URL + "/hello?name=Alice"
	resp, err := client.Get(uri)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("GET %v\n", uri)
	fmt.Println(resp.Status)
	fmt.Println(string(body))
	fmt.Println()
	/*
		GET http://127.0.0.1:56131/hello?name=Alice
		200 OK
		Hello, Alice
	*/
}
