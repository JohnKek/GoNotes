В программировании контекст (context) — это информация о среде, в которой существует объект или выполняется функция. В Go под контекстом обычно имеют в виду интерфейс Context из пакета context. Исходно его придумали, чтобы облегчить работу с HTTP-запросами. Но контексты можно использовать и в обычном многозадачном коде. Давайте посмотрим, как именно.

## Отмена операции через канал
Рассмотрим функцию execute(), которая умеет запустить переданную функцию и поддерживает отмену:
```go
func execute(cancel <-chan struct{}, fn func() int) (int, error) {
    ch := make(chan int, 1)

    go func() {
        ch <- fn()
    }()

    select {
    case res := <-ch:
        return res, nil
    case <-cancel:
        return 0, errors.New("canceled")
    }
}
```
Здесь все знакомо:

- функция принимает канал, через который может получить сигнал отмены;
- запускает fn() в отдельной горутине;
- через select дожидается выполнения fn() либо прерывается по отмене, смотря что наступит раньше.

Напишем клиента, который отменяет операцию с 50% вероятностью:
```go
func main() {
    rand.Seed(time.Now().Unix())

    // работает в течение 100 мс
    work := func() int {
        time.Sleep(100 * time.Millisecond)
        fmt.Println("work done")
        return 42
    }

    // ждет 50 мс, после этого
    // с вероятностью 50% отменяет работу
    maybeCancel := func(cancel chan struct{}) {
        time.Sleep(50 * time.Millisecond)
        if rand.Float32() < 0.5 {
            close(cancel)
        }
    }

    cancel := make(chan struct{})

    go maybeCancel(cancel)

    res, err := execute(cancel, work)
    fmt.Println(res, err)
}
```
А теперь сделаем то же самое через контекст.

# Отмена операции через контекст
Основное назначение контекста в Go — отмена операций.

Повторим через контекст то, что мы только что сделали через канал отмены. Функция execute() будет принимать контекст ctx вместо канала cancel:

```go
// выполняет функцию fn с учетом контекста ctx
func execute(ctx context.Context, fn func() int) (int, error) {
    ch := make(chan int, 1)

    go func() {
        ch <- fn()
    }()

    select {
    case res := <-ch:
        return res, nil
    case <-ctx.Done():       // (1)
        return 0, ctx.Err()  // (2)
    }
}
```
Код почти не изменился:
- вместо канала cancel сигнал об отмене может прийти из канала ctx.Done() ➊
- вместо ручного создания ошибки при отмене возвращаем ctx.Err() ➋

Клиент тоже несколько меняется:
```go
func main() {
    // ...
    
    // работает в течение 100 мс
    work := func() int {
        // ...
    }

    // ждет 50 мс, после этого
    // с вероятностью 50% отменяет работу
    maybeCancel := func(cancel func()) {
        time.Sleep(50 * time.Millisecond)
        if rand.Float32() < 0.5 {
            cancel()
        }
    }

    ctx := context.Background()              // (1)
    ctx, cancel := context.WithCancel(ctx)   // (2)
    defer cancel()                           // (3)

    go maybeCancel(cancel)                   // (4)

    res, err := execute(ctx, work)           // (5)
    fmt.Println(res, err)
}
```
Вот что здесь происходит:

➊ через context.Background() создаем пустой контекст;

➋ через context.WithCancel() на базе пустого контекста создаем новый, с возможностью ручной отмены;

➌ планируем отложенную отмену контекста при выходе из main();

➍ отменяем контекст с 50% вероятностью;

➎ передаем контекст в функцию execute().

context.WithCancel() возвращает сам контекст и функцию cancel для его отмены. Вызов cancel() освобождает занятые контекстом ресурсы и закрывает канал ctx.Done() — этот эффект мы и используем для прерывания execute(). Если контекст отменен, то ctx.Err() возвращает ошибку с причиной отмены (context.Canceled в нашем случае).

### Контекст — это матрешка. 
Объект контекста неизменяемый. Чтобы добавить контексту новые свойства, создают новый контекст («дочерний») на основе старого («родительского»). Поэтому мы сначала создали пустой контекст, а затем новый (с возможностью отмены) на его основе:
```go
// родительский контекст
ctx := context.Background()

// дочерний
ctx, cancel := context.WithCancel(ctx)
Если отменить родительский контекст — отменятся и все дочерние (но не наоборот).

// родительский контекст
parentCtx, parentCancel := context.WithCancel(context.Background())

// дочерний контекст
childCtx, childCancel := context.WithCancel(parentCtx)

// parentCancel() отменит parentCtx и childCtx
// childCancel() отменит только childCtx
```

Многократная отмена безопасна. Если два раза вызвать close() на канале, получим панику. А вот вызывать cancel() контекста можно сколько угодно. Первая отмена сработает, а остальные будут проигнорированы. Это удобно, потому что можно сразу после создания контекста запланировать отложенный cancel(), плюс явно отменить контекст при необходимости (как мы сделали в функции maybeCancel). С каналом бы так не получилось.

# Таймаут
Настоящая сила контекста в том, что его можно использовать как для ручной отмены, так и для отмены по таймауту. Следите за руками:

```go
func execute(ctx context.Context, fn func() int) (int, error) {
    // код не меняется
}

func main() {
    // ...

    // работает в течение 100 мс
    work := func() int {
        // ...
    }

    // возвращает случайный аргумент из переданных
    randomChoice := func(arg ...int) int {
        i := rand.Intn(len(arg))
        return arg[i]
    }

    // случайный таймаут - 50 мс либо 150 мс
    timeout := time.Duration(randomChoice(50, 150)) * time.Millisecond
    ctx, cancel := context.WithTimeout(context.Background(), timeout)    // (1)
    defer cancel()

    res, err := execute(ctx, work)
    fmt.Println(res, err)
}
```

Функция execute() вообще не изменилась, а в main() вместо context.WithCancel() теперь context.WithTimeout() ➊. Этого достаточно, чтобы execute() теперь отваливалась по таймауту в половине случаев (ошибка context.DeadlineExceeded):

Классно, вот твой текст в Markdown:

---

## Родительский и дочерний таймауты

Допустим, у нас есть все та же функция `execute()` и две функции, которые она может выполнить — быстрая `work()` и медленная `slow()`:

```go
// выполняет функцию fn с учетом контекста ctx
func execute(ctx context.Context, fn func() int) (int, error) {
    ch := make(chan int, 1)

    go func() {
        ch <- fn()
    }()

    select {
    case res := <-ch:
        return res, nil
    case <-ctx.Done():
        return 0, ctx.Err()
    }
}

// работает в течение 100 мс
work := func() int {
    time.Sleep(100 * time.Millisecond)
    return 42
}

// работает в течение 300 мс
slow := func() int {
    time.Sleep(300 * time.Millisecond)
    return 13
}
```

Пусть таймаут по умолчанию составляет 200 мс:

```go
// возвращает контекст с умолчательным таймаутом 200 мс
getDefaultCtx := func() (context.Context, context.CancelFunc) {
    const timeout = 200 * time.Millisecond
    return context.WithTimeout(context.Background(), timeout)
}
```

Тогда `work()` с умолчательным контекстом успеет выполниться:

```go
// таймаут 200 мс
ctx, cancel := getDefaultCtx()
defer cancel()

// успеет выполниться
res, err := execute(ctx, work)
fmt.Println(res, err) // 42 <nil>
```

А `slow()` — не успеет:

```go
// таймаут 200 мс
ctx, cancel := getDefaultCtx()
defer cancel()

// НЕ успеет выполниться
res, err := execute(ctx, slow)
fmt.Println(res, err) // 0 context deadline exceeded
```

Мы можем создать дочерний контекст, чтобы задать более жесткий таймаут. Тогда применится именно он, а не родительский:

```go
// родительский контекст с таймаутом 200 мс
parentCtx, cancel := getDefaultCtx()
defer cancel()

// дочерний контекст с таймаутом 50 мс
childCtx, cancel := context.WithTimeout(parentCtx, 50*time.Millisecond)
defer cancel()

// теперь work НЕ успеет выполниться
res, err := execute(childCtx, work)
fmt.Println(res, err) // 0 context deadline exceeded
```

А если создать дочерний контекст с более мягким ограничением — он окажется бесполезен. Таймаут родительского контекста сработает раньше:

```go
// родительский контекст с таймаутом 200 мс
parentCtx, cancel := getDefaultCtx()
defer cancel()

// дочерний контекст с таймаутом 500 мс
childCtx, cancel := context.WithTimeout(parentCtx, 500*time.Millisecond)
defer cancel()

// slow все равно НЕ успеет выполниться
res, err := execute(childCtx, slow)
fmt.Println(res, err) // 0 context deadline exceeded
```

Получается вот что:

- Из таймаутов, наложенных родительским и дочерним контекстами, всегда срабатывает более жесткий.
- Дочерние контексты могут только ужесточить таймаут родительского, но не ослабить его.

--- 

Надеюсь, это поможет в понимании работы с контекстами в Go! 👩‍💻🚀

Спасибо за информацию! Вот текст о дедлайнах и использовании контекстов в Go:

---

## Дедлайн

Помимо таймаута, контекст поддерживает дедлайн (deadline) — это когда операция отменяется не через N секунд, а в конкретный момент времени.

```go
// выполняет функцию fn с учетом контекста ctx
func execute(ctx context.Context, fn func() int) (int, error) {
    // без изменений
}

func main() {
    // ...

    // работает в течение 100 мс
    work := func() int {
        // ...
    }

    // возвращает случайный аргумент из переданных
    randomChoice := func(arg ...int) int {
        // ...
    }

    // случайный дедлайн - +50 мс либо +150 мс
    // от текущего времени
    timeout := time.Duration(randomChoice(50, 150)) * time.Millisecond
    deadline := time.Now().Add(timeout)
    ctx, cancel := context.WithDeadline(context.Background(), deadline)  // (1)
    defer cancel()

    res, err := execute(ctx, work)
    fmt.Println(res, err)
}

// примеры вывода
// work выполнен
// 42 <nil>

// дедлайн превышен
// 0 context deadline exceeded

// дедлайн превышен
// 0 context deadline exceeded

// work выполнен
// 42 <nil>

// песочница
```

Из примера видно, что `context.WithDeadline()` работает так же, как и `context.WithTimeout()`, но принимает значение `time.Time` вместо `time.Duration`.

Более того, `context.WithTimeout()` является оберткой над `context.WithDeadline()`:

```go
// фрагмент кода стандартной библиотеки
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) {
    return WithDeadline(parent, time.Now().Add(timeout))
}
```

Каждый контекст всегда оперирует конкретным дедлайном, который можно получить через метод `Deadline()`:

```go
now := time.Now()
fmt.Println(now)
// 2009-11-10 23:00:00

ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
deadline, ok := ctx.Deadline()
fmt.Println(deadline, ok)
// 2009-11-10 23:00:05 true

plus5s := now.Add(5 * time.Second)
ctx, _ = context.WithDeadline(context.Background(), plus5s)
deadline, ok = ctx.Deadline()
fmt.Println(deadline, ok)
// 2009-11-10 23:00:05 true
```

Второе значение при вызове `Deadline()` показывает, установлен ли дедлайн:

- Для контекстов, созданных через `WithTimeout()` и `WithDeadline()`, это значение равно `true`.
- Для `WithCancel()` и `Background()` — `false`.

```go
ctx, _ := context.WithCancel(context.Background())
deadline, ok := ctx.Deadline()
fmt.Println(deadline, ok)
// 0001-01-01 00:00:00 false

ctx = context.Background()
deadline, ok = ctx.Deadline()
fmt.Println(deadline, ok)
// 0001-01-01 00:00:00 false
```

Надеюсь, эта информация о дедлайнах и контекстах в Go будет полезной для вас! 🕰️👩‍💻🚀

--- 

Если у вас есть дополнительные вопросы или нужна дальнейшая помощь, не стесняйтесь спрашивать!
