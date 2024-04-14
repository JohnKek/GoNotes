## Goroutine
Функции, запущенные через `go`, называются _горутинами_ (goroutine). Планировщик Go распределяет их по _потокам_ операционной системы (threads), которые, в свою очередь, выполняются на разных _ядрах CPU_ (CPU cores). По сравнению с потоками ОС горутины очень легкие, так что их можно создавать сотнями и тысячами.
```go
func main() {
    go say(1, "go is awesome") 
    go say(2, "cats are cute") 
    time.Sleep(500 * time.Millisecond) 
}
```
Если убрать строку `time.Sleep(500 * time.Millisecond)`, то в поток вывода ничего не напечатается, так как горутины работают независимо и не ожидают друг друга, плюс сама точка входа программы `main` является горутиной.
Правильный способ ожидания горутин:
```go
func main() {
    var wg sync.WaitGroup

    wg.Add(1)
    go say(&wg, 1, "go is awesome")

    wg.Add(1)
    go say(&wg, 2, "cats are cute")

    wg.Wait()
}

func say(wg *sync.WaitGroup, id int, phrase string) {
    for _, word := range strings.Fields(phrase) {
        fmt.Printf("Worker #%d says: %s...\n", id, word)
        dur := time.Duration(rand.Intn(100)) * time.Millisecond
        time.Sleep(dur)
    }
    wg.Done()
}
```
WaitGroup под капотом реализована как атомарный счётчик `Add` увеличивает на 1, `Done` уменьшает на 1, `Wait` ждёт пока значение счётчика не станет равным 0