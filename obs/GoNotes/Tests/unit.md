Тесты пишут в отдельном файле с суффиксом _test. Например, если IntMin() определена в файле ints.go, то тесты будут в
том же пакете, но в отдельном файле ints_test.go.

```go
func IntMin(a, b int) int {
if a < b {
return a
}
return b
}
```
```go
func TestIntMin(t *testing.T) {
    got := IntMin(2, -2)
    want := -2
    if got != want {
        t.Errorf("got %d; want %d", got, want)
    }
}
```
Тест всегда принимает указатель на testing.T — это такой сборник полезных функций. Например, t.Errorf() выводит сообщение об ошибке, не прекращая выполнение теста.

### Табличные тесты

Чтобы проверить граничные условия и различные сочетания параметров, нам пришлось бы написать целую кучу тестовых функций. Обычно вместо этого используют _табличные_ (table-driven) тесты:

-   отдельно описывают входные данные и ожидаемые значения;
-   последовательно выполняют каждую проверку с помощью `t.Run()`.

```go
func TestIntMin(t *testing.T) {
var tests = []struct {
a, b int
want int
}{
{0, 1, 0},
{1, 0, 0},
{1, 1, 1},
}

for _, test := range tests {
name := fmt.Sprintf("case(%d,%d)", test.a, test.b)
t.Run(name, func (t *testing.T) {
got := IntMin(test.a, test.b)
if got != test.want {
t.Errorf("got %d, want %d", got, test.want)
}
})
}
}

```

Срез `tests` — наша таблица с отдельными проверками. У каждой проверки есть входные параметры (`a`, `b`) и ожидаемый результат (`want`). `t.Run()` запускает конкретную проверку и выводит ошибку, если она не прошла.

### Подготовка и завершение

Бывает, хочется выполнить какой-то код до старта тестов (setup), а какой-то — после их заверешения (teardown). Например, перед тестами подготовить данные, а после — почистить их.

В других языках для этого часто используют отдельные функции. В Go же обходятся одной — `TestMain()`:

```go
func TestMain(m *testing.M) {
fmt.Println("Setup tests...")
start := time.Now()

m.Run()

fmt.Println("Teardown tests...")
end := time.Now()
fmt.Println("Tests took", end.Sub(start))
}
```

`m.Run()` запускает тесты (функции вида `Test*()`). Получается, что любой код до него — это setup, а после — teardown.



### Избирательный запуск

Иногда хочется запустить не все тесты, а только часть. Например, если некоторые тесты медленные, и каждый раз гонять их не хочется.

Наша старая знакомая — функция, которая суммирует целые числа:

```go
func Sum(nums ...int) int { total := 0 for _, num := range nums { total += num } return total }
```

И пара тестов:

```go
func TestSumFew(t *testing.T) { if Sum(1, 2, 3, 4, 5) != 15 { t.Errorf("Expected Sum(1, 2, 3, 4, 5) == 15") } } func TestSumN(t *testing.T) { n := 1_000_000_000 nums := make([]int, n) for i := 0; i < n; i++ { nums[i] = i + 1 } got := Sum(nums...) want := n * (n + 1) / 2 if got != want { t.Errorf("Expected sum[i=1..n](i) == n*(n+1)/2") } }
```

Прогоним тесты:

```http
$ go test -v === RUN TestSum --- PASS: TestSum (0.00s) === RUN TestSumN --- PASS: TestSumN (3.78s)
```

Ни малейшего желания каждый раз ждать 4 секунды, пока выполняется `TestSumN` на миллиарде слагаемых. Воспользуемся так называемым short-режимом, который делит все тесты на «короткие» и «длинные». `TestSumFew` будет коротким тестом, а `TestSumN` — длинным:

```go
func TestSumN(t *testing.T) { if testing.Short() { t.Skip("skipping test in short mode.") } // сам тест }
```

Теперь тесты в коротком режиме будут игнорировать `TestSumN` и отрабатывать моментально:

```http
$ go test -v -short === RUN TestSumFew --- PASS: TestSumFew (0.00s) === RUN TestSumN selective_test.go:21: skipping test in short mode. --- SKIP: TestSumN (0.00s)
```

[песочница](https://go.dev/play/p/FzWJvGubMve)

Альтернативный вариант — указать маску названия теста. Так выполнится только `TestSumFew`:

```http
$ go test -v -run Few
```

А так — только `TestSumN`:

```http
$ go test -v -run N
```
