## Что такое каналы?
Канал — это объект связи, с помощью которого горутины обмениваются данными. Технически это конвейер (или труба), откуда можно считывать или помещать данные. То есть одна горутина может отправить данные в канал, а другая — считать помещенные в этот канал данные.

## Создание канала
Go для создания канала предоставляет ключевое слово `chan`. Канал может передавать данные только одного типа, данные других типов через это канал передавать невозможно.
```go
package main
import "fmt"
func main() {    
    var c chan int    
    fmt.Println(c)
} 
```

## Запись и чтение данных
Go предоставляет простой синтаксис для чтения `<-` и записи в канал
```go
c <- data
```
В этом примере мы передаем данные в канал `c`. Направление стрелки указывает на то, что мы извлекаем данные из `data` и помещаем в канал `c`.
```go
<- c
```
А здесь мы считываем данные с канала `c`. Эта операция не сохраняет данные в переменную и она является корректной. Если вам необходимо сохранить данные с канала в переменную, вы можете использовать следующий синтаксис:
```go
var data int
data = <- c
```
Теперь данные из канала `c`, который имеет тип `int`, могут быть записаны в переменную `data`. Так же можно упростить запись, используя короткий синтаксис:
```go
data := <- c
```
Go определит тип данных, передаваемый каналу `c`, и предоставит `data` корректный тип данных.
## *Важно!*
Все вышеобозначенные операции с каналом являются блокируемыми. Когда вы помещаете данные в канал, горутина блокируется до тех пор, пока данные не будут считаны другой горутиной из этого канала. В то же время операции канала говорят планировщику о планировании другой горутины, поэтому программа не будет заблокирована полностью. Эти функции весьма полезны, так как отпадает необходимость писать блокировки для взаимодействия горутин.

## Как закрыть канал?

В Go есть механизм, который решает ровно эту задачу:

- писатель может _закрыть_ (close) канал;
- читатель может понять, что канал закрыт.

Писатель закрывает канал функцией `close()`:

```go
in := make(chan string)
go func() {
    words := strings.Split(str, ",")
    for _, word := range words {
        in <- word
    }
    close(in)
}()
```

Читатель проверяет статус канала вторым значением при считывании:

```go
for {
    word, ok := <-in
    if !ok {
        break
    }
    if word != "" {
        fmt.Printf("%s ", word)
    }
}
```

Допустим, в канал передают строки `one`, `two`, после чего закрывают его. Вот что получит при этом читатель:

```go
// in <- "one"
word, ok := <-in
// word = "one", ok = true

// in <- "two"
word, ok := <-in
// word = "two", ok = true

// close(in)
word, ok := <-in
// word = "", ok = false

word, ok := <-in
// word = "", ok = false

word, ok := <-in
// word = "", ok = false
```

Пока канал открыт, читатель получает очередное значение и статус `true`. Если канал закрыт, читатель получает нулевое значение (для строки это `""`) и статус `false`.

Как видно из примера, читать из закрытого канала можно сколько угодно — каждый раз вернется нулевое значение и статус `false`. Это неспроста — через несколько шагов разберемся, зачем так сделано.

Закрыть канал можно только один раз. Повторное закрытие приведет к панике:

```go
in := make(chan string)
close(in)
close(in)
// panic: close of closed channel
```

Записать в закрытый канал тоже не получится:

```go
in := make(chan string)
go func() {
    in <- "hi"
    close(in)
}()
fmt.Println(<-in)
// hi

in <- "bye"
// panic: send on closed channel
```

Отсюда два важных правила:

1. _Закрыть канал имеет право только писатель, но не читатель_. Если читатель закроет канал, то писатель словит панику при следующей записи.
2. _Писатель имеет право закрыть канал, только если владеет им единолично_. Если писателей несколько, и один из них закроет канал, то остальные словят панику при следующей записи или попытке закрыть канал со своей стороны.

|                                                                                                                                                                                                |
| ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Всегда ли закрывать канал**<br><br>Закрывать канал стоит с единственной целью — сообщить его читателям, что все данные отправлены. Если читателей у канала нет, то и закрывать его не нужно. |

## Итерирование по каналу

На предыдущем шаге мы заставили читателя постоянно проверять, открыт ли канал:

```go
for {
    word, ok := <-in
    if !ok {
        break
    }
    if word != "" {
        fmt.Printf("%s ", word)
    }
}
```

Довольно громоздко. Чтобы не делать этого вручную, Go поддерживает конструкцию `range` для чтения из канала:

```go
for word := range in {
    if word != "" {
        fmt.Printf("%s ", word)
    }
}
```

`range` автоматически считывает очередное значение из канала и проверяет, не закрыт ли тот. Если закрыт — выходит из цикла. Удобно, не правда ли?

Обратите внимание, что `range` по каналу возвращает одно значение, а не пару, в отличие от `range` по срезу. Сравните:

```go
// срез
words := []string{"1", "2", "3"}
for idx, val := range words {
    fmt.Println(idx, val)
}

// канал
in := make(chan string)
go func() {
    in <- "1"
    in <- "2"
    in <- "3"
    close(in)
}()
for val := range in {
    fmt.Println(val)
}
```

## Направление каналов
Уменьшить путаницу с каналами можно, если задать им _направление_ (direction). Каналы бывают:

- `chan`: для чтения и записи (по умолчанию);
- `chan<-` : только для записи (send-only);
- `<-chan`: только для чтения (receive-only).

Функции `submit()` подойдет канал «только для записи»:

```go
func submit(str string, stream chan<- string) {  // (1)
    words := strings.Split(str, ",")
    for _, word := range words {
        stream <- word
    }
    // <-stream                                  // (2)
    close(stream)
}
```

В сигнатуре функции ➊ мы указали, что канал только для записи, так что теперь прочитать из него не получится. Если раскомментировать строчку ➋, получим ошибку при компиляции:

```bash
invalid operation: cannot receive from send-only channel stream
```

Функции `print()` подойдет канал «только для чтения»:

```
func print(stream <-chan string) {  // (1)
    for word := range stream {
        if word != "" {
            fmt.Printf("%s ", word)
        }
    }
    // stream <- "oops"             // (2)
    // close(stream)                // (3)
    fmt.Println()
}
```

В сигнатуре функции ➊ мы указали, что канал только для чтения. Записать в него не получится. Если раскомментировать строчку ➋, получим ошибку при компиляции:

```bash
invalid operation: cannot send to receive-only channel stream
```

Закрыть receive-only канал тоже нельзя. Если раскомментировать строчку ➌, получим ошибку при компиляции:

```bash
invalid operation: cannot close receive-only channel stream
```

Задать направление канала можно и при инициализации. Но толку от этого мало:

```go
func main() {
    str := "one,two,,four"
    stream := make(chan<- string)  // (!)
    go submit(str, stream)
    print(stream)
}
```

Здесь `stream` объявлен только для записи, так что для функции `print()` он больше не подходит. А если объявить только для чтения — не подойдет для `submit()`. Поэтому обычно каналы инициализируют для чтения и записи, а в параметрах конкретных функций заявляют как однонаправленные. Go конвертирует обычный канал в направленный автоматически:

```go
stream := make(chan int)

go func(in chan<- int) {
    in <- 42
}(stream)

func(out <-chan int) {
    fmt.Println(<-out)
}(stream)
// 42
```

Старайтесь всегда указывать направление канала в параметрах функции, чтобы застраховаться от ошибок во время выполнения программы.

### Канал завершения

Есть функция, которая произносит фразу по словам (с некоторыми задержками):

```
func say(id int, phrase string) {
    for _, word := range strings.Fields(phrase) {
        fmt.Printf("Worker #%d says: %s...\n", id, word)
        dur := time.Duration(rand.Intn(100)) * time.Millisecond
        time.Sleep(dur)
    }
}
```

Запускаем несколько одновременных болтушек, по одной на каждую фразу:

```go
func main() {
    phrases := []string{
        "go is awesome",
        "cats are cute",
        "rain is wet",
        "channels are hard",
        "floor is lava",
    }
    for idx, phrase := range phrases {
        go say(idx+1, phrase)
    }
}
```

Программа, конечно, ничего не печатает — функция `main()` завершается до того, как отработает хотя бы одна болтушка:

```bash
$ go run say.go
<пусто>
```

Раньше мы использовали `[sync.WaitGroup](https://pkg.go.dev/sync#WaitGroup)`, чтобы дождаться завершения горутин. А можно использовать прием «канал завершения» (done channel):

```go
func say(done chan<- struct{}, id int, phrase string) {
    for _, word := range strings.Fields(phrase) {
        fmt.Printf("Worker #%d says: %s...\n", id, word)
        dur := time.Duration(rand.Intn(100)) * time.Millisecond
        time.Sleep(dur)
    }
    done <- struct{}{}                     // (1)
}

func main() {
    phrases := []string{
        "go is awesome",
        "cats are cute",
        "rain is wet",
        "channels are hard",
        "floor is lava",
    }

    done := make(chan struct{})            // (2)

    for idx, phrase := range phrases {
        go say(done, idx+1, phrase)        // (3)
    }

    // wait for goroutines to finish
    for i := 0; i < len(phrases); i++ {    // (4)
        <-done
    }
}
```

Вот что здесь происходит:

- создаем отдельный канал ➋ и передаем его в каждую горутину ➌;
- в горутине записываем значение в канал по окончании работы ➊;
- в основной функции ждем, пока каждая горутина запишет в канал ➍.

Чтобы это работало, основная функция должна точно знать, сколько горутин запущено (в нашем случае — по одной на каждую исходную строку). Иначе непонятно, сколько значений читать из `done`.

Теперь все в порядке:

```no-highlight
$ go run say.go
Worker #5 says: floor...
Worker #1 says: go...
Worker #4 says: channels...
Worker #3 says: rain...
Worker #2 says: cats...
Worker #4 says: are...
Worker #3 says: is...
Worker #4 says: hard...
Worker #2 says: are...
Worker #5 says: is...
Worker #5 says: lava...
Worker #3 says: wet...
Worker #1 says: is...
Worker #2 says: cute...
Worker #1 says: awesome...
```

Если прием с каналом завершения вам не по душе, вместо него всегда можно использовать `sync.WaitGroup`.

### Буферизованные каналы

Есть горутина `send()`, которая передает значение горутине `receive()` через канал `stream`:

```go
var wg sync.WaitGroup
wg.Add(2)

stream := make(chan bool)

send := func() {
    defer wg.Done()
    fmt.Println("Sender ready to send...")
    stream <- true                                // (1)
    fmt.Println("Sent!")
}

receive := func() {
    defer wg.Done()
    fmt.Println("Receiver not ready yet...")
    time.Sleep(500 * time.Millisecond)
    fmt.Println("Receiver ready to receive...")
    <-stream                                      // (2)
    fmt.Println("Received!")
}

go send()
go receive()
wg.Wait()
```

`send()` сразу после запуска хочет передать значение в канал, но `receive()` пока не готова. Поэтому `send()` вынуждена заблокироваться в точке ➊ и ждать 500 миллисекунд, пока `receive()` не придет в точку ➋ и не согласится принять значение из канала. Получается, что горутины _синхронизируются_ в точке приема/передачи:

```bash
Receiver not ready yet...
Sender ready to send...
Receiver ready to receive...
Received!
Sent!
```

Чаще всего такое поведение нас устраивает. Но что делать, если мы хотим, чтобы отправитель не ждал получателя? Хотим, чтобы он отправил значение в канал и занимался своими делами. А получатель пусть заберет, когда будет готов. Ах, если бы только в канал можно было сложить значение, как в очередь! И как хорошо, что Go предоставляет ровно такую возможность:

```go
// второй аргумент - размер буфера канала
// то есть количество значений, которые он может хранить
stream := make(chan int, 3)
// ⬜ ⬜ ⬜

stream <- 1
// 1️⃣ ⬜ ⬜

stream <- 2
// 1️⃣ 2️⃣ ⬜

stream <- 3
// 1️⃣ 2️⃣ 3️⃣

fmt.Println(<-stream)
// 1
// 2️⃣ 3️⃣ ⬜

fmt.Println(<-stream)
// 2
// 3️⃣ ⬜ ⬜

stream <- 4
stream <- 5
// 3️⃣ 4️⃣ 5️⃣

stream <- 6
// в канале больше нет места,
// горутина блокируется
```

Такие каналы называются _буферизованными_ (buffered), потому что у них есть собственный буфер фиксированного размера, в котором можно хранить значения. По умолчанию, если не указать размер буфера, будет создан канал с буфером размера 0 — именно с такими каналами мы работали до сих пор:

```go
// канал без буфера
unbuffered := make(chan int)

// канал с буфером
buffered := make(chan int, 3)
```

На буферизованных каналах работают встроенные функции `len()` и `cap()`:

- `cap()` возвращает общую емкость канала;
- `len()` – количество значений в канале.

```go
stream := make(chan int, 2)
fmt.Println(cap(stream), len(stream))
// 2 0

stream <- 7
fmt.Println(cap(stream), len(stream))
// 2 1

stream <- 7
fmt.Println(cap(stream), len(stream))
// 2 2

<-stream
fmt.Println(cap(stream), len(stream))
// 2 1
```

Чтобы отвязать `send()` от `receive()` с помощью буферизованного канала, достаточно изменить единственную строчку, оставив остальной код без изменений:

```go
// создаем канал с буфером 1
// вместо обычного
stream := make(chan bool, 1)

send := func() {
    // ...
}

receive := func() {
    // ...
}

go send()
go receive()
```

Теперь отправитель не ждет получателя:

```bash
Receiver not ready yet...
Sender ready to send...
Sent!
Receiver ready to receive...
Received!
```

## Закрыть буферизованный канал

Как мы знаем, обычный канал после закрытия отдает нулевое значение и признак `false`:

```go
stream := make(chan int)
close(stream)

val, ok := <-stream
fmt.Println(val, ok)
// 0 false

val, ok = <-stream
fmt.Println(val, ok)
// 0 false

val, ok = <-stream
fmt.Println(val, ok)
// 0 false
```

Канал с буфером ведет себя точно так же, если буфер пуст. А вот если в буфере есть значения — иначе:

```go
stream := make(chan int, 2)
stream <- 1
stream <- 2
close(stream)

val, ok := <-stream
fmt.Println(val, ok)
// 1 true

val, ok = <-stream
fmt.Println(val, ok)
// 2 true

val, ok = <-stream
fmt.Println(val, ok)
// 0 false
```

Пока в буфере есть значения, канал отдает их и признак `true`. Когда все значения выбраны — отдает нулевое значение и признак `false`, как обычный канал.

Благодаря этому отправитель может в любой момент закрыть канал, не задумываясь о том, остались в нем значения или нет. Получатель в любом случае их считает:

```go
stream := make(chan int, 3)

go func() {
    fmt.Println("Sending...")
    stream <- 1
    stream <- 2
    stream <- 3
    close(stream)
    fmt.Println("Sent and closed!")
}()

time.Sleep(500 * time.Millisecond)
fmt.Println("Receiving...")
for val := range stream {
    fmt.Printf("%v ", val)
}
fmt.Println()
fmt.Println("Received!")
```

```bash
Sending...
Sent and closed!
Receiving...
1 2 3 
Received!
```

## nil-канал

Как у любого типа в Go, у каналов тоже есть нулевое значение. Это `nil`:

```go
var stream chan int
fmt.Println(stream)
// <nil>
```

nil-канал — малоприятная штука:

- Запись в nil-канал навсегда блокирует горутину.
- Чтение из nil-канала навсегда блокирует горутину.
- Закрытие nil-канала приводит к панике.

```go
var stream chan int

go func() {
    stream <- 1
}()

<-stream

// fatal error: all goroutines are asleep - deadlock!
```

```go
var stream chan int
close(stream)

// panic: close of nil channel
```

У nil-каналов есть некоторые очень специфические сценарии использования. Один из них мы рассмотрим на следующем уроке. В целом — старайтесь избегать nil-каналов до тех пор, пока не почувствуете, что никак не можете без них обойтись.