### Запросы и статусы

За HTTP-запросы отвечает структура `http.Client`:

```go
client := http.Client{Timeout: 3 * time.Second}
```

Таймаут — это максимальное время, которое клиент готов ждать ответа от сервера, прежде чем вернет ошибку. По умолчанию таймаут не задан, то есть ждать клиент будет до бесконечности. Поэтому всегда явно указывайте его.

Клиента достаточно создать один раз и дальше использовать для всех запросов.

Выполним GET-запрос и посмотрим на результат:

```go
const uri = "https://httpbingo.org/status/200"
resp, err := client.Get(uri)  // (1)
if err != nil {
    panic(err)
}

fmt.Printf("GET %v\n", uri)
fmt.Println(resp.Status)      // (2)
fmt.Println()
/*
    GET https://httpbingo.org/status/200
    200 OK
*/
```

Метод клиента `Get()` выполнил запрос на указанный URL ➊. Запрос успешно прошел и вернул ответ со статусом `200 OK` ➋

Вообще статусов существует [великое множество](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status), но чаще всего встречаются такие:

-   `200 OK` — запрос успешно отработал.
-   `301 Moved Permanently` и `302 Found` — ресурс, к которому вы обращаетесь, переехал. Придется сделать еще один запрос по новому адресу.
-   `4xx Client Error` (включает любые статусы, которые начинаются на 4) — ошибка клиента (проблема на вашей стороне).
-   `5xx Server Error` (включает любые статусы, которые начинаются на 5) — ошибка сервера (проблема на стороне сервера).

### Параметры запроса

В GET-запросах часто передают параметры, например:

```http
GET https://httpbingo.org/get?id=42
```

Вот как это делается в Go:

```go
const uri = "https://httpbingo.org/get"
req, err := http.NewRequest(http.MethodGet, uri, nil)  // (1)
if err != nil {
panic(err)
}

params := url.Values{}              // (2)
params.Add("id", "42")
req.URL.RawQuery = params.Encode()  // (3)

resp, err := client.Do(req)         // (4)
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
```

Поскольку мы хотим модифицировать запрос перед отправкой, метод клиента `Get()` больше не подходит. Вместо этого сначала явно создадим запрос через `http.NewRequest()` ➊, а затем выполним его через метод клиента `Do()` ➍.

HTTP-метод указывается первым параметром `http.NewRequest()`. В нашем случае это константа `http.MethodGet` (`GET`), но аналогично можно указать любой другой метод (`http.MethodPost`, `http.MethodDelete`, и так далее).

`url.Values` ➋ — это карта, в которую можно добавлять, изменять и удалять параметры:

```go
Has(key string) bool 
Get(key string) string 
Add(key, value string) 
Set(key, value string) 
Del(key string)
```

Чтобы добавить параметры в запрос, придется превратить карту в строку через метод `Encode()` и присвоить ее свойству `URL.RawQuery` ➌

### Заголовки запроса

Чтобы добавить в запрос заголовки, используют свойство `Header`. Это карта с примерно таким же набором методов, как было у параметров:

```go
const uri = "https://httpbingo.org/headers"
req, err := http.NewRequest(http.MethodGet, uri, nil)
if err != nil {
panic(err)
}

req.Header.Add("Accept", "application/json")  // (1)
req.Header.Add("X-Request-Id", "42")          // (2)

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
```

Здесь, к счастью, не требуется возиться с кодированием — достаточно явно добавить заголовки в свойство `Header` ➊ ➋

### «Простые» запросы

В обучающих статьях запросы часто выполняют через функции `http.Get()` и `http.Post()`:

```go
// "простой" GET-запрос
const uri = "https://httpbingo.org/status/200"
resp, err := http.Get(uri)
fmt.Println(resp.Status, err)
// 200 OK <nil>

// "простой" POST-запрос
const uri = "https://httpbingo.org/status/200"
body := []byte("hello")
resp, err := http.Post(uri, "text/plain", bytes.NewBuffer(body))
fmt.Println(resp.Status, err)
// 200 OK <nil>
```

Кажется удобным, что не приходится отдельно создавать `http.Client`. Но рекомендую вам не применять эти функции в реальной жизни. Они используют умолчательный HTTP-клиент (`http.DefaultClient`), у которого не задан таймаут — то есть клиент ждет ответа до скончания времен.

[песочница](https://go.dev/play/p/ZlLGPZHDRd-)

> В песочнице запрещены внешние вызовы, потому код примеров выдает там ошибки. Запускайте локально, если хотите попробовать.