Конечно, быстрее выполнять работу параллельно, в N «обработчиков». Тогда логика будет такой:

если есть свободный обработчик — отдать запрос ему;
иначе ждать, пока кто-нибудь освободится.
Сделаем обертку withWorkers(n, fn), которая обеспечивает одновременное выполнение. Для этого заведем канал free и будем следить, чтобы запускались не более n функций fn() одновременно:
```go
func withWorkers(n int, fn func()) (handle func(), wait func()) {
    // канал с токенами
    free := make(chan struct{}, n)
    for i := 0; i < n; i++ {
        free <- struct{}{}
    }

    // выполняет fn, но не более n одновременно
    handle = func() {
        <-free
        go func() {
            fn()
            free <- struct{}{}
        }()
    }

    // ожидает, пока все запущенные fn отработают
    wait = func() {
        for i := 0; i < n; i++ {
            <-free
        }
    }

    return handle, wait
}
```
Теперь клиент вызывает функцию work() не напрямую, а через обертку:
```go
func main() {
    work := func() {
        time.Sleep(100 * time.Millisecond)
    }

    handle, wait := withWorkers(2, work)

    start := time.Now()

    handle()
    handle()
    handle()
    handle()
    wait()

    fmt.Println("4 calls took", time.Since(start))
}
```

Получилось так:

первый и второй вызовы сразу пошли в обработку;
третий и четвертый ждали, пока предыдущие два закончат обрабатываться.
В результате при двух обработчиках 4 вызова выполняются за 200 мс.

Схема обработки с ожиданием отлично подходит, когда степень параллелизма n и индивидуальное время работы work() примерно соответствуют интенсивности, с которой вызывают handle(). Тогда у каждого вызова есть хороший шанс сразу пойти в обработку.

Если же вызовов сильно больше, чем способны «прожевать» обработчики — система начнет «тормозить». Каждый отдельный work() по-прежнему будет отрабатывать за 100 мс. Но вызовы handle() будут подвисать, ведь каждому придется ждать свободного токена. Для обработки данных в конвейере это не страшно, а вот для онлайн-запросов может быть не слишком хорошо.

Бывает, что клиент предпочел бы сразу получить ошибку в ответ на handle(), если все обработчики заняты. Тут схема с ожиданием уже не подойдет.

---

# Обработка без ожидания

Изменим логику withWorkers():

если есть свободный токен — выполнить функцию;
иначе сразу вернуть ошибку.
Так клиенту не придется ждать «подвисшего» вызова.

Нам в очередной раз поможет инструкция select:
```go
handle = func() error {
    select {
    case <-free:
        go func() {
            fn()
            free <- struct{}{}
        }()
        return nil
    default:
        return errors.New("busy")
    }
}
```

Вспомним алгоритм работы селекта:

1. Проверяет, какие ветки не заблокированы.
2. Если таких веток несколько, выбирает одну из них случайным образом и выполняет ее.
3. Если все ветки заблокированы, блокирует выполнение, пока хотя бы одна ветка не разблокируется.

На самом деле, третий пункт разбивается на два:

- Если все ветки заблокированы и нет блока default — блокирует выполнение, пока хотя бы одна ветка не разблокируется.
- Если все ветки заблокированы и есть блок default — выполняет его.

default идеально подходит для нашей ситуации:

- если есть свободный токен в канале free — запускаем fn;
- иначе ничего не ждем и возвращаем ошибку busy.

Посмотрим на клиента:
```go
func main() {
    work := func() {
        time.Sleep(100 * time.Millisecond)
    }

    handle, wait := withWorkers(2, work)

    start := time.Now()

    err := handle()
    fmt.Println("1st call, error:", err)

    err = handle()
    fmt.Println("2nd call, error:", err)

    err = handle()
    fmt.Println("3rd call, error:", err)

    err = handle()
    fmt.Println("4th call, error:", err)

    wait()

    fmt.Println("4 calls took", time.Since(start))
}
```

```
1st call, error: <nil>
2nd call, error: <nil>
3rd call, error: busy
4th call, error: busy
4 calls took 100ms
```

Первые два вызова выполнились одновременно (по 100 мс каждый), а третий и четвертый моментально получили ошибку. В итоге все вызовы обработались за 100 мс.

Конечно, при таком подходе требуется некоторая «сознательность» клиента. Клиент должен понять, что ошибка «busy» означает перегруз, и отложить дальнейшие вызовы handle() на некоторое время, или уменьшить их частоту.

---

# Таймаут операции

Есть функция, которая обычно отрабатывает за 10 мс, но в 20% случаев занимает 200 мс:

```go
func work() int {
    if rand.Intn(10) < 8 {
        time.Sleep(10 * time.Millisecond)
    } else {
        time.Sleep(200 * time.Millisecond)
    }
    return 42
}
```

Мы в принципе не хотим ждать дольше, скажем, 50 мс. Поэтому установим таймаут (timeout) — максимальное время, в течение которого готовы ждать ответ. Если операция не уложилась в таймаут, будем считать это ошибкой.

делаем обертку, которая выполняет переданную функцию с указанным таймаутом, и будем вызывать ее вот так:

```go
func withTimeout(fn func() int, timeout time.Duration) (int, error) {
    // ...
}

func main() {
    for i := 0; i < 10; i++ {
        start := time.Now()
        timeout := 50 * time.Millisecond
        if answer, err := withTimeout(work, timeout); err != nil {
            fmt.Printf("Took longer than %v. Error: %v\n", time.Since(start), err)
        } else {
            fmt.Printf("Took %v. Result: %v\n", time.Since(start), answer)
        }
    }
}
```

Идея работы withTimeout() следующая:

- Запускаем переданную fn() в отдельной горутине.
- Ждем в течение timeout времени.
- Если fn() вернула ответ — возвращаем его.
- Если не успела — возвращаем ошибку.

Вот как можно это реализовать:
```go
func withTimeout(fn func() int, timeout time.Duration) (int, error) {
    var result int

    done := make(chan struct{})
    go func() {
        result = fn()
        close(done)
    }()

    select {
    case <-done:
        return result, nil
    case <-time.After(timeout):
        return 0, errors.New("timeout")
    }
}
```

Здесь все знакомо, кроме time.After(). Эта библиотечная функция возвращает канал, который изначально пуст, а через timeout времени отправляет в него значение. Благодаря этому select выберет нужный вариант:
- ветку <-done, если fn() успела до таймаута (вернет ответ);
- ветку <-time.After(), если не успела (вернет ошибку).

---

# Таймер
Бывает, хочется выполнить действие не прямо сейчас, а через какое-то время. В Go для этого предусмотрен инструмент таймер (timer):
```go
func main() {
    work := func() {
        fmt.Println("work done")
    }

    var eventTime time.Time

    start := time.Now()
    timer := time.NewTimer(100 * time.Millisecond)    // (1)
    go func() {
        eventTime = <-timer.C                         // (2)
        work()
    }()

    // достаточно времени, чтобы сработал таймер
    time.Sleep(150 * time.Millisecond)
    fmt.Printf("delayed function started after %v\n", eventTime.Sub(start))
}
```

time.NewTimer() создает новый таймер ➊, который сработает через указанный промежуток времени. Таймер — это структура с каналом C, в который он запишет текущее время, когда сработает ➋. Благодаря этому, функция work() выполнится только после того, как таймер сработает.

Таймер можно остановить, тогда в канал С значение не придет, и work() не запустится:
```go
func main() {
    // ...

    start := time.Now()
    timer := time.NewTimer(100 * time.Millisecond)
    go func() {
        <-timer.C
        work()
    }()

    time.Sleep(10 * time.Millisecond)
    fmt.Println("10ms has passed...")
    // таймер еще не успел сработать
    if timer.Stop() {
        fmt.Printf("delayed function canceled after %v\n", time.Since(start))
    }
}
```
Если остановить таймер слишком поздно, Stop() вернет false:

Для отложенного запуска функции не обязательно вручную возиться с созданием таймера и считыванием из канала. Есть удобная обертка time.AfterFunc():
```go
func main() {
    work := func() {
        fmt.Println("work done")
    }

    timer := time.AfterFunc(100*time.Millisecond, work)

    time.Sleep(10 * time.Millisecond)
    fmt.Println("10ms has passed...")
    // таймер еще не успел сработать
    if timer.Stop() {
        fmt.Println("execution canceled")
    }
}
```
---

# Тикер
Бывает, хочется выполнять какое-то действие с определенной периодичностью. В Go и для этого есть инструмент — тикер (ticker). Тикер похож на таймер, только срабатывает не один раз, а регулярно, пока не остановите:
```go
func main() {
    work := func(at time.Time) {
        fmt.Printf("%s: work done\n", at.Format("15:04:05.000"))
    }

    ticker := time.NewTicker(50 * time.Millisecond)
    defer ticker.Stop()

    go func() {
        for {
            at := <-ticker.C
            work(at)
        }
    }()

    // хватит на 5 тиков
    time.Sleep(260 * time.Millisecond)
}
```

NewTicker(d) создает тикер, который с периодичностью d отправляет текущее время в канал С. Тикер обязательно надо рано или поздно остановить через Stop(), чтобы он освободил занятые ресурсы.