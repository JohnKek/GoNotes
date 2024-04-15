Приведенная ниже программа осуществляет обратный отсчет для запуска ракеты. Функция time. Tic к возвращает канал, по которому она периодически отправляет события, действуя как метроном. Каждое событие представляет собой значение момента времени, но оно не так интересно, как сам факт его доставки.
```go
func main() {
	fmt.Println("Commencing countdown.")
	tick := time.Tick(1 * time.Second)
	for countdown := 10; countdown > 0; countdown-- {
		fmt.Println(countdown)
		<-tick
	}
	launch()
}

//!-

func launch() {
	fmt.Println("Lift off!")
}
```
Давайте теперь добавим возможность прервать последовательность запуска, нажав клавишу <Enter> во время обратного отсчета. Сначала мы запустим go-подпрограмму, которая пытается прочитать один байт из стандартного ввода и, если это удастся, отправляет значение в канал, который называется abort.
```go
func main() {
	abort := make(chan struct{})
	go func() {
		os.Stdin.Read(make([]byte, 1)) // read a single byte
		abort <- struct{}{}
	}()
	fmt.Println("Commencing countdown.  Press return to abort.")
	select {
	case <-time.After(10 * time.Second):
		// Do nothing.
	case <-abort:
		fmt.Println("Launch aborted!")
		return
	}
	launch()
}

func launch() {
	fmt.Println("Lift off!")
}
```
Давайте сделаем нашу программу запуска выводящей обратный отсчет. Инструкция select ниже приводит к тому, что на каждой итерации цикла выполняется ожидание сигнала прерывания в течение секунды, но не дольше.
```go
func main() {
	abort := make(chan struct{})
	go func() {
		os.Stdin.Read(make([]byte, 1)) // read a single byte
		abort <- struct{}{}
	}()

	fmt.Println("Commencing countdown.  Press return to abort.")
	tick := time.Tick(1 * time.Second)
	for countdown := 10; countdown > 0; countdown-- {
		fmt.Println(countdown)
		select {
		case <-tick:
		case <-abort:
			fmt.Println("Launch aborted!")
			return
		}
	}
	launch()
}
func launch() {
	fmt.Println("Lift off!")
}
```
Канал отмены
Для начала создадим отдельный канал отмены (cancel channel), через который main() будет сигнализировать rangeGen(), что пора завершаться:
```go
func main() {
    cancel := make(chan struct{})    // (1)
    defer close(cancel)              // (2)

    generated := rangeGen(cancel, 41, 46)    // (3)
    for val := range generated {
        fmt.Println(val)
        if val == 42 {
            break
        }
    }
}
```
Мы создаем канал cancel ➊ и сразу настраиваем отложенный вызов close(cancel) ➋. Так часто делают, чтобы не отслеживать по коду все места, в которых нужно закрыть канал. defer гарантирует, что канал в любом случае закроется при выходе из функции, так что о нем можно не беспокоиться.

Затем передаем канал cancel в горутину ➌. Теперь, когда канал закроется, горутина должна как-то понять это и завершить работу. Хотелось бы добавить примерно такую проверку:

```go
func rangeGen(cancel <-chan struct{}, start, stop int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for i := start; i < stop; i++ {
            out <- i
            if <-cancel == struct{}{} {    // (1)
                return
            }
        }
    }()
    return out
}
```
Если cancel закрыт, то проверка ➊ пройдет (закрытый канал всегда возвращает нулевое значение, помните?), и горутина завершит работу. Но вот беда: если cancel не закрыт, то горутина заблокируется и на следующую итерацию цикла не пойдет.

Нам нужна другая, неблокирующая логика:

если cancel закрыт, выйти из горутины;
иначе отправить очередное значение в out.
В Go для этого существует инструкция select:
```go
func rangeGen(cancel <-chan struct{}, start, stop int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for i := start; i < stop; i++ {
            select {
            case out <- i:    // (1)
            case <-cancel:    // (2)
                return
            }
        }
    }()
    return out
}
```