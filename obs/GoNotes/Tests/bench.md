### Бенчмарки

Сравнить производительность в Go помогут _бенчмарки_ (benchmark).

Бенчмарки похожи на обычные тесты, только начинаются со слова `Benchmark` вместо `Test`, и принимают параметр `testing.B` вместо `testing.T`.

Поскольку сравниваем два варианта — `MatchContains()` и `MatchRegexp()` — подготовим по бенчмарку на каждого:

```go
// длинная строка
const src = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."

// слова из середины строки
const pattern = "commodo"

func BenchmarkMatchContains(b *testing.B) {
for n := 0; n < b.N; n++ {
MatchContains(pattern, src)
}
}

func BenchmarkMatchRegexp(b *testing.B) {
for n := 0; n < b.N; n++ {
MatchRegexp(pattern, src)
}
}

```

Бенчмарк всегда устроен одинаково — внутри у него цикл от `0` до `b.N`, а в теле цикла вызывается целевая функция. Количество итераций `b.N` Go определяет самостоятельно — так, чтобы получились статистически значимые результаты.

Запускаем и смотрим:


Интуиция не подвела. Проверка по регулярке работает аж в 25 раз медленнее 😱
 песочнице бенчмарки не работают. Поэтому, если захотите их запустить — придется делать это локально.

### Варьируем размер

Метод `IntSet.Contains()` отлично работает на множестве из 3 или 10 элементов. Но как он поведет себя на 1000 элементов? 10000? 100000? Ответить помогут бенчмарки.

Возможно, вы помните, как на прошлом уроке мы выполняли табличные тесты — один большой тест, внутри которого запускается много маленьких. Когда тестируют производительность одной и той же функции на разных входных данных — поступают аналогично:

```go
func BenchmarkIntSet(b *testing.B) {
for _, size := range []int{1, 10, 100, 1000, 10000, 100000} {
set := randomSet(size)
name := fmt.Sprintf("Contains-%d", size)
b.Run(name, func (b *testing.B) {
for n := 0; n < b.N; n++ {
elem := rand.Intn(100000)
set.Contains(elem)
}
})

}
}

var rnd = rand.New(rand.NewSource(42))

func randomSet(size int) IntSet {
set := MakeIntSet()
for i := 0; i < size; i++ {
n := rnd.Intn(100000)
set.Add(n)
}
return set
}

```

`BenchmarkIntSet()` работает так:

1.  Берет конкретный размер множества (1, 10, 100...).
2.  Создает множество указанного размера и заполняет его случайными числами.
3.  Запускает бенчмарк, который меряет производительность `IntSet.Contains()` на этом множестве.
4.  Переходит к следующему размеру и повторяет шаги 2–3.

Таким образом, выполняется по одному независимому бенчмарку для каждого размера множества.

Посмотрим на результаты:

```http
$ go test -bench="." BenchmarkIntSet/Contains-1-8 56203182 19.70 ns/op BenchmarkIntSet/Contains-10-8 47717664 22.26 ns/op BenchmarkIntSet/Contains-100-8 17775404 62.53 ns/op BenchmarkIntSet/Contains-1000-8 3852343 309.0 ns/op BenchmarkIntSet/Contains-10000-8 418706 2645 ns/op BenchmarkIntSet/Contains-100000-8 85210 13574 ns/op
```

Не слишком радужно. С увеличением размера множества скорость `Contains()` падает в сотни раз. Хорошие алгоритмы так себя не ведут.

[исходники](https://gist.github.com/nalgeon/6507509659efc0c1d1c46baf8def2058)