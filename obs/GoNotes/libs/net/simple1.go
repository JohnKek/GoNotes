package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func main() {
	client := http.Client{Timeout: 3 * time.Second}

	{
		// GET-запрос
		const uri = "https://httpbingo.org/status/200"
		resp, err := client.Get(uri)
		if err != nil {
			panic(err)
		}

		fmt.Printf("GET %v\n", uri)
		fmt.Println(resp.Status)
		fmt.Println()
		/*
			GET https://httpbingo.org/status/200
			200 OK
		*/
	}

	{
		// GET-запрос с параметром
		const uri = "https://httpbingo.org/get"
		req, err := http.NewRequest(http.MethodGet, uri, nil)
		if err != nil {
			panic(err)
		}

		params := url.Values{}
		params.Add("id", "42")
		req.URL.RawQuery = params.Encode()

		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%v %v\n", req.Method, req.URL)
		fmt.Println(resp.Status)
		fmt.Println()

		/*
			GET https://httpbingo.org/get?id=42
			200 OK
		*/
	}

	{
		// GET-запрос с несколькими параметрами
		const uri = "https://httpbingo.org/get"
		req, err := http.NewRequest(http.MethodGet, uri, nil)
		if err != nil {
			panic(err)
		}

		params := url.Values{}
		params.Add("brand", "lg")
		params.Add("category", "tv")
		params.Add("category", "notebook")
		req.URL.RawQuery = params.Encode()

		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%v %v\n", req.Method, req.URL)
		fmt.Println(resp.Status)
		fmt.Println()
		/*
			GET https://httpbingo.org/get?brand=lg&category=tv&category=notebook
			200 OK
		*/
	}

	{
		// Заголовки запроса
		const uri = "https://httpbingo.org/headers"
		req, err := http.NewRequest(http.MethodGet, uri, nil)
		if err != nil {
			panic(err)
		}

		req.Header.Add("Accept", "application/json")
		req.Header.Add("X-Request-Id", "42")

		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%v %v\n", req.Method, req.URL)
		fmt.Println(resp.Status)
		fmt.Println()
		/*
			GET https://httpbingo.org/headers
			200 OK
		*/
	}

	{
		// "простой" GET-запрос
		const uri = "https://httpbingo.org/status/200"
		resp, err := http.Get(uri)
		fmt.Println(resp.Status, err)
		// 200 OK <nil>
	}

	{
		// "простой" POST-запрос
		const uri = "https://httpbingo.org/status/200"
		body := []byte("hello")
		resp, err := http.Post(uri, "text/plain", bytes.NewBuffer(body))
		fmt.Println(resp.Status, err)
		// 200 OK <nil>
	}

}
