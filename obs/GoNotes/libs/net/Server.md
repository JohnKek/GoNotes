### httptest.Server

В начале урока я сказал, что разбирать HTTP-серверы мы не будем. Но сделаем небольшое исключение: посмотрим, как написать заглушку сервера для тестов. Иначе пришлось бы использовать реальное внешнее API — это не подходит для тестирования.

Допустим, мы хотим получать запросы такого вида:

```http
GET /hello?name=Alice
```

И отвечать на них статусом `200 OK` с телом ответа вида:

```go
Hello, Alice!
```

Напишем обработчик, который это сделает:

```go
func helloHandler(w http.ResponseWriter, r *http.Request) { params := r.URL.Query() // (1) name := params.Get("name") // (2) body := fmt.Sprintf("Hello, %s", name) // (3) w.Header().Set("Content-Type", "text/plain") // (4) w.WriteHeader(http.StatusOK) // (5) w.Write([]byte(body)) // (6) }
```

Обработчик — это функция, которая принимает запрос и писатель ответа. Вот что в ней происходит:

1.  Получаем карту с параметрами запроса.
2.  Извлекаем значение параметра `name`.
3.  Формируем строку с будущим телом ответа.
4.  Устанавливаем заголовок ответа `Content-Type`.
5.  Устанавливаем статус ответа `200 OK`.
6.  Записываем тело ответа.

Теперь создадим фейковый HTTP-сервер, который принимает обращения к `/hello` и отправляет их в обработчик:

```go
func startServer() *httptest.Server { mux := http.NewServeMux() mux.HandleFunc("/hello", helloHandler) server := httptest.NewServer(mux) return server }
```

Наконец, запустим сервер:

```go
server := startServer() defer server.Close()
```

Теперь можно отправлять запросы обычным способом. Единственный нюанс — клиента получаем через `server.Client()`:

```go
client := server.Client() uri := server.URL + "/hello?name=Alice" resp, err := client.Get(uri) if err != nil { panic(err) } defer resp.Body.Close() body, err := io.ReadAll(resp.Body) if err != nil { panic(err) } fmt.Printf("GET %v\n", uri) fmt.Println(resp.Status) fmt.Println(string(body)) fmt.Println() /* GET http://127.0.0.1:56131/hello?name=Alice 200 OK Hello, Alice */
```

Можно создать разные обработчики для разных урлов сервера:

```go
func handlerOne(w http.ResponseWriter, r *http.Request) { // ... } func handlerTwo(w http.ResponseWriter, r *http.Request) { // ... } func handlerThree(w http.ResponseWriter, r *http.Request) { // ... } func startServer() *httptest.Server { mux := http.NewServeMux() mux.HandleFunc("/one", handlerOne) mux.HandleFunc("/two", handlerTwo) mux.HandleFunc("/three", handlerThree) server := httptest.NewServer(mux) return server }
```

[песочница](https://go.dev/play/p/daX7Slrn0FU)

> В песочнице запрещены внешние вызовы, потому код примеров выдает там ошибки. Запускайте локально, если хотите попробовать.