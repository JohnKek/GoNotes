Сделать запрос — в лучшем случае половина дела. Давайте посмотрим, как обработать ответ от сервера.

### Код ответа

Помните статус ответа (`200 OK`)? Он состоит из кода ответа (`200`) и описания (`OK`). Часто работать с кодом удобнее, потому что это число. Поэтому код можно получить отдельно от статуса через поле ответа `StatusCode`:

```go
const uri = "https://httpbingo.org/status/400"
resp, err := client.Get(uri)
if err != nil {
panic(err)
}

fmt.Printf("GET %v\n", uri)
fmt.Printf("%s, code = %d\n", resp.Status, resp.StatusCode)
fmt.Println()
/*
   GET https://httpbingo.org/status/400
   400 Bad Request, code = 400
*/
```

### Тело ответа

Тело ответа доступно через поле `Body`. Но там лежат не сами данные, а структура, через которую можно их прочитать (как мы это делали на уроке «Чтение и запись»). В большинстве случаев удобно считать весь ответ в один прием через `io.ReadAll()`:

```go
const uri = "https://httpbingo.org/get"
resp, err := client.Get(uri)
if err != nil {
panic(err)
}

defer resp.Body.Close()             // (1)
body, err := io.ReadAll(resp.Body)  // (2)

fmt.Printf("GET %v\n", uri)
fmt.Println(resp.Status)
fmt.Printf("%T, %d bytes, err = %v\n", body, len(body), err)
fmt.Println(string(body))
fmt.Println()
/*
   GET https://httpbingo.org/get
   200 OK
   []uint8, 649 bytes, err = <nil>
   {"args":{},...}
*/
```

Тело ответа после начитывания из `Body` — это срез байт ➋. Само `Body` обязательно надо закрыть, чтобы освободить занятые ресурсы ➊

### Заголовки ответа

У HTTP-ответа, как и у запроса, есть заголовки. Они доступны через карто-подобное свойство `Header`:

```go
const uri = "https://httpbingo.org/get"
resp, err := http.Get(uri)
if err != nil {
panic(err)
}

fmt.Printf("GET %v\n", uri)
fmt.Println(resp.Status)
for key, val := range resp.Header {
fmt.Printf("%v = %v\n", key, val)
}
fmt.Println()
/*
   GET https://httpbingo.org/get
   200 OK
   Access-Control-Allow-Credentials = [true]
   Access-Control-Allow-Origin = [*]
   Content-Type = [application/json; encoding=utf-8]
   Date = [Sun, 11 Sep 2022 16:59:16 GMT]
   ...
*/
```

### Ответ в JSON

Если вы вызываете чье-то API, ответ в большинстве случаев будет в JSON. Чтобы преобразовать тело ответа из набора байт в объект, используют знакомую нам по уроку о JSON функцию `json.Unmarshal()`:

```go
const uri = "https://httpbingo.org/json"
resp, err := http.Get(uri)
if err != nil {
    panic(err)
}

defer resp.Body.Close()
body, err := io.ReadAll(resp.Body)
if err != nil {
    panic(err)
}

var data map[string]any
err = json.Unmarshal(body, &data)
if err != nil {
    panic(err)
}

fmt.Printf("GET %v\n", uri)
fmt.Println(resp.Status)
fmt.Println(data)
fmt.Println()
/*
    GET https://httpbingo.org/json
    200 OK
    map[slideshow:map[author:Yours Truly date:date of publication ...]]
*/
```

Или `json.Decoder`:

```go
const uri = "https://httpbingo.org/json"
resp, err := http.Get(uri)
if err != nil {
panic(err)
}

defer resp.Body.Close()

var data map[string]any
err = json.NewDecoder(resp.Body).Decode(&data)
// ...
```

### Игнорирование ответа

Бывает, что тело ответа вам вообще не интересно — достаточно кода. Но даже в этом случае всегда начитывайте и закрывайте поле `Body`, чтобы освободить занятые ресурсы. Начитывать можно в «никуда» (`io.Discard`):

```go
const uri = "https://httpbingo.org/get"
resp, err := client.Get(uri)
if err != nil {
panic(err)
}

defer resp.Body.Close()
_, err = io.Copy(io.Discard, resp.Body)
if err != nil {
panic(err)
}
```

[песочница](https://go.dev/play/p/Ip3xU4qOgs0)

> В песочнице запрещены внешние вызовы, потому код примеров выдает там ошибки. Запускайте локально, если хотите попробовать.