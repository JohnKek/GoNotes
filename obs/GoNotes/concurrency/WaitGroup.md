```go
import "sync"

// ...

func main() {
    var wg sync.WaitGroup

    wg.Add(1)
    go func() {
        time.Sleep(50 * time.Millisecond)
        fmt.Print(".")
        wg.Done()
    }()

    wg.Wait()
    fmt.Println("all done")
}
```
wg.Done() кладем в deffer в горутине

# Значение vs указатель
Еще важный нюанс реализации. Передавать группу ожидания в функции следует как указатель *WaitGroup, а не как значение WaitGroup. Иначе каждая функция получит свою копию, а не оригинальную группу. В результате у них будут разные счетчики, и ничего работать не будет.
```go
Вот передача по значению:

func runWork(wg sync.WaitGroup) {
    wg.Add(1)
    go func() {
        time.Sleep(50 * time.Millisecond)
        fmt.Println("work done")
        wg.Done()
    }()
}

func main() {
    var wg sync.WaitGroup
    runWork(wg)
    wg.Wait()
    fmt.Println("all done")
}
```
runWork получила копию группы, и увеличила ее счетчик через Add. В main же своя копия с нулевым счетчиком, поэтому Wait не заблокировала выполнение. В результате main завершилась, не ожидая завершения горутины из runWork:

А вот передача по указателю:
```go
func runWork(wg *sync.WaitGroup) {
    wg.Add(1)
    go func() {
        time.Sleep(50 * time.Millisecond)
        fmt.Println("work done")
        wg.Done()
    }()
}

func main() {
    var wg sync.WaitGroup
    runWork(&wg)
    wg.Wait()
    fmt.Println("all done")
}
```

# Инкапсуляция
В Go считается хорошим тоном прятать потроха синхронизации от клиентов (не тех клиентов, что деньги платят, а программного кода, который вызывает ваш код). Коллеги не скажут спасибо, если вы будете заставлять их жонглировать группами ожидания. Лучше скрыть логику синхронизации в отдельной функции или типе, а вовне предоставить удобный интерфейс.

## Функции-обертки
Допустим, я написал функцию RunConc, которая выполняет переданные функции одновременно:
```go
func RunConc(wg *sync.WaitGroup, funcs ...func()) {
    wg.Add(len(funcs))
    for _, fn := range funcs {
        fn := fn
        go func() {
            defer wg.Done()
            fn()
        }()
    }
}
```
И предлагаю вызывать ее так:
```go
work := func() {
    time.Sleep(50 * time.Millisecond)
    fmt.Print(".")
}

var wg sync.WaitGroup
RunConc(&wg, work, work, work)
wg.Wait()
```
Удобно ли это, если клиент просто хочет выполнить функции параллельно и дождаться их завершения? Не слишком.

Лучше спрятать wait-группу внутрь функции:
```go
func RunConc(funcs ...func()) {
    var wg sync.WaitGroup
    wg.Add(len(funcs))
    for _, fn := range funcs {
        fn := fn
        go func() {
            defer wg.Done()
            fn()
        }()
    }
    wg.Wait()
}
```
# Типы-обертки
Коллеги попробовали RunConc, и им не понравилось. Говорят, не хотят передавать все функции одним куском и сразу запускать, а вместо этого хотят добавлять их постепенно по одной, и потом когда-нибудь запустить все разом. И еще хотелось бы возможность запустить набор функций несколько раз.
```go
// ConcRunner выполняет заданные функции одновременно.
type ConcRunner struct {
    wg    sync.WaitGroup
    funcs []func()
}

// NewConcRunner создает новый экземпляр ConcRunner.
func NewConcRunner() *ConcRunner {
    return &ConcRunner{wg: sync.WaitGroup{}}
}

// Add добавляет функцию, не выполняя ее.
func (cg *ConcRunner) Add(fn func()) {
    cg.funcs = append(cg.funcs, fn)
}

// Run выполняет функции одновременно и дожидается их окончания.
func (cg *ConcRunner) Run() {
    cg.wg.Add(len(cg.funcs))
    for _, fn := range cg.funcs {
        fn := fn
        go func() {
            defer cg.wg.Done()
            fn()
        }()
    }
    cg.wg.Wait()
}
```