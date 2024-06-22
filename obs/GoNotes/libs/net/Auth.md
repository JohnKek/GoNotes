Если обычные веб-страницы открыты всем желающим, то API (особенно коммерческие) требуют от клиента «представиться» — аутентифицироваться. Для этого клиент передает в запросе свой идентификатор. Рассмотрим пару популярных способов это сделать.

### Basic-аутентификация

Basic-аутентификация — это когда в заголовке `Authorization` на сервер уходит пара «логин-пароль». Предварительно она кодируется по стандарту Base64. От вас в данном случае требуется только задать логин и пароль через метод запроса `SetBasicAuth()`:

```go
const uri = "https://httpbingo.org/get"
req, err := http.NewRequest(http.MethodGet, uri, nil)
if err != nil {
panic(err)
}

req.SetBasicAuth("anton", "secret")
fmt.Println(req.Header["Authorization"])
// [Basic YW50b246c2VjcmV0]

resp, err := client.Do(req)
// ..
```

### Аутентификация по токену

При аутентификации по токену придется явно указать этот самый токен в заголовке `Authorization`:

```go
const uri = "https://httpbingo.org/get"
req, err := http.NewRequest(http.MethodGet, uri, nil)
if err != nil {
panic(err)
}

token := "1234567890"
req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
fmt.Println(req.Header["Authorization"])
// [Bearer 1234567890]

resp, err := client.Do(req)
// ...
```

Обычно токен, логин и пароль не «зашивают» в коде, а берут из настроек или какого-то доверенного источника.

[песочница](https://go.dev/play/p/ksGhfBNx-rs)

> В песочнице запрещены внешние вызовы, потому код примеров выдает там ошибки. Запускайте локально, если хотите попробовать.