Мы начали урок с того, что за HTTP-запросы отвечает _клиент_ — структура `http.Client`. Давайте рассмотрим его более пристально.

### Таймауты

Если сервер не успел вернуть ответ за отведенный таймаут, то метод клиента — `Get()`, `Post()`, `PostForm()` или `Do()` — вернет ошибку:

```go
// таймаут
client := http.Client{Timeout: 1 * time.Second}

const uri = "https://httpbingo.org/delay/5"
_, err := client.Get(uri)
fmt.Println(err)
/*
   Get "https://httpbingo.org/delay/5":
   context deadline exceeded (Client.Timeout exceeded while awaiting headers)
*/
```

### Отмена запроса

Чтобы отменить запрос вручную, используют знакомый нам механизм контекстов (`context.Context`). Запрос в этом случае создают через `http.NewRequestWithContext()`:

```go
const uri = "https://httpbingo.org/delay/5"
ctx, cancel := context.WithCancel(context.Background())
req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
if err != nil {
panic(err)
}

time.AfterFunc(200*time.Millisecond, func() {
cancel()
})
_, err = http.DefaultClient.Do(req)
fmt.Println(err)
/*
   Get "https://httpbingo.org/delay/5": context canceled
*/
```

Вызов `cancel()` контекста отменяет запрос.

### Редиректы

Редирект — это ответ от сервера с кодом 301 или 302 (и еще некоторые 30x коды). Вместе с таким кодом сервер присылает новый URL, на который переехал запрошенный ресурс. Дальше клиент обращается по новому урлу, если снова получил редирект — обращается по следующему урлу, и так далее — пока не получит ответ с кодом 2xx, 4xx или 5xx.

Обычно задумываться о редиректах не приходится, потому что стандартный `http.Client` выполняет их автоматически и прозрачно для программиста. По умолчанию поддерживается до 10 редиректов — хватит с запасом. Если будет больше — клиент вернет ошибку.

Если вы зачем-то хотите отключить автоматические редиректы — перекройте свойство клиента `CheckRedirect`:

```go
// не следовать редиректам
client := http.Client{Timeout: 3 * time.Second}
client.CheckRedirect = func(req *http.Request, via []*http.Request) error {  // (1)
return http.ErrUseLastResponse
}

const uri = "https://httpbingo.org/redirect/5"
resp, err := client.Get(uri)
if err != nil {    // (2)
panic(err)
}

fmt.Printf("GET %v\n", uri)
fmt.Println(resp.Status)
fmt.Println()
/*
   GET https://httpbingo.org/redirect/5
   302 Found, code = 302
*/
```

`CheckRedirect()` ➊ — это функция, которую клиент вызывает при каждом редиректе. Специальная ошибка `ErrUseLastResponse`, которую она здесь возвращает, заставляет клиента не обрабатывать редирект, но и не возвращать ошибку (паника в ➋ не срабатывает). В результате возвращается ответ с кодом 302 и пустой ошибкой.

[песочница](https://go.dev/play/p/CoA-w32auta)

> В песочнице запрещены внешние вызовы, потому код примеров выдает там ошибки. Запускайте локально, если хотите попробовать.