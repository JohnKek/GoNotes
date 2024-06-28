Go обладает превосходной поддержкой кодирования и декодирования этих форматов, предоставленной в пакетах стандартной библиотеки encoding/ j son, encoding/ xml, encoding/asnl и других, и все эти пакеты имеют схожие API.

Рассмотрим приложение, которое собирает кинообзоры и предлагает рекомендации. Его тип данных Movie и типичный список значений объявляются ниже. (Строковые литералы после объявлений полей Year и Color являются дескрипторами полей (field tags); мы объясним их чуть позже.)
aooi.io/ch4/movie 
```go
type Movie struct {
Title string
Year int 'json:"released"'
Color bool 'j son:"color,omitempty"'
Actors []string
}
```
```go
var movies = []Movie{
{Title: "Casablanca", Year: 1942, Color: false,
Actors: []string{"Humphrey Bogart", "Ingrid Bergman"}},
{Title: "Cool Hand Luke", Year: 1967, Color: true,
Actors: []string{"Paul Newman"}},
{Title: "Bullitt", Year: 1968, Color: true,
Actors: []string{"Steve McQueen", "Jacqueline Bisset"}},
// ...
}
```

Такие структуры данных отлично подходят для JSON и легко конвертируются в обоих направлениях. Преобразование структуры данных Go наподобие movies в JSON называется маршачингом (marshaling). Маршалинг осуществляется с помощью json.Marshal:

```go
data, err := json.Marshal(movies) 
if err != nil {
    log.Fatalf("Сбой маршалинга JSON: %s", err)
}
fmt.Printf("%s\n", data)
```


Для чтения человеком можно использовать функцию json. Marshallndent, которая дает аккуратно отформатированное представление с отсту
пами. Два дополнительных аргумента определяют префикс для каждой строки вывода и строку для каждого уровня отступа:
```go
data, err := json.MarshalIndent(movies, "", " ") 
if err != nil {
    log.Fatalf("Сбой маршалинга DSON: %s", err)
}
fmt.Printf("%s\n", data)
```

Обратная к маршалингу операция, декодирования JSON и заполнения структуры данных Go, называется восстановлением объекта или демаршалингом (unmarshaling) и выполняется с помощью json.Unmarshal
```go
var titles []struct{ Title string } 
if err := json.Unmarshal(data, Stitles); err != nil 
{ 
	log.Fatalf("Сбой демаршалинга JSON: %s", err)
}
fmt.Println(titles) // "[{Casablanca} {Cool Hand Luke} {Bullitt}]"
```