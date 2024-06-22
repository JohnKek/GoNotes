Делать GET-запросы — это, конечно, здорово. Но часто клиенты передают на сервер данные другими методами — через POST,
PUT или PATCH. Тогда данные идут не в параметрах URL, а в теле запроса.

Самые распространенные форматы тела запроса — это HTML-форма и JSON. Давайте посмотрим, как с ними работать.

### Отправка формы

Форма отправляется на сервер в виде строки, в которой пары «ключ-значение» объединены через амперсанд `&`. Например:

```http
brand=lg&category=tv
```

Одному ключу можно сопоставить несколько значений:

```http
brand=lg&category=tv&category=notebook
```

Самостоятельно собирать такую строку не требуется — для этого подойдет уже знакомая нам структура`url.Values`:

```go
const uri = "https://httpbingo.org/post"

data := url.Values{}
data.Add("brand", "lg")
data.Add("category", "tv")
data.Add("category", "notebook")

resp, err := client.PostForm(uri, data)  // (1)
if err != nil {
panic(err)
}

fmt.Printf("POST %v\n", uri)
fmt.Println(resp.Status)
fmt.Println()
/*
   POST https://httpbingo.org/post
   200 OK
*/
```

Метод клиента `PostForm()` ➊ принимает структуру с формой вторым параметром, кодирует ее в строку и отправляет на сервер
через метод `POST`. Заодно он устанавливает подходящий заголовок `Content-Type` (`application/x-www-form-urlencoded`),
так что нам делать этого не придется.

### Отправка JSON

Чтобы отправить JSON на сервер, его придется предварительно кодировать в набор байт. Для этого используют знакомую нам
по уроку о JSON функцию `json.Marshal()`:

```go
const uri = "https://httpbingo.org/post"

data := map[string]any{
"brand":    "lg",
"category": []string{"tv", "notebook"},
}

b, err := json.Marshal(data)
if err != nil {
panic(err)
}

req, err := http.NewRequest(http.MethodPost, uri, bytes.NewReader(b))  // (1)
if err != nil {
panic(err)
}

req.Header.Add("Content-Type", "application/json")  // (2)
req.Header.Add("Accept", "application/json")        // (3)

resp, err := client.Do(req)
if err != nil {
panic(err)
}

fmt.Printf("POST %v\n", uri)
fmt.Println(resp.Status)
fmt.Println()
/*
   POST https://httpbingo.org/post
   200 OK
*/
```

Обратите внимание:

- ➊ Конструктор запроса`http.NewRequest()` принимает третьим параметром тело запроса — но не срез байт, а читателя.
- ➋ Раз тело запроса в JSON — следует указать соответствующий заголовок `Content-Type`.
- ➌ Если тело ответа тоже ожидается в JSON — стоит указать и заголовок `Accept`.

Если в запросе не нужны другие заголовки помимо `Content-Type` — вместо`http.NewRequest` +`client.Do` можно
использовать `client.Post`:

```go
const uri = "https://httpbingo.org/post"

data := map[string]any{
"brand":    "lg",
"category": []string{"tv", "notebook"},
}

b, err := json.Marshal(data)
if err != nil {
panic(err)
}

resp, err := client.Post(uri, "application/json", bytes.NewReader(b))
// ...
```

[песочница](https://go.dev/play/p/SsfDZlwCWaQ)

> В песочнице запрещены внешние вызовы, потому код примеров выдает там ошибки. Запускайте локально, если хотите
> попробовать.