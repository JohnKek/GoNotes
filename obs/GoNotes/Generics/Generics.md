Дженерики позволяют единообразно работать с данными разных типов. На самом деле, вы уже знакомы с дженериками и не раз их использовали, потому что встроенные функции append и len — generic-функции. Судите сами:
```go
// срез целых чисел
intSlice := []int{1, 2, 3}

// срез строк
strSlice := []string{"a", "b", "c"}
```
Generic-функция

Создать собственную generic-функцию очень просто:
```go
// Reverse инвертирует порядок элементов в срезе.
func Reverse[T any](s []T) {
    for i := 0; i < len(s)/2; i++ {
        s[i], s[len(s)-1-i] = s[len(s)-1-i], s[i]
    }
}
```
T — это параметр типа (type parameter), а []T указывает, что параметр функции s — это срез объектов типа T.

any говорит, что T может быть абсолютно любым типом (пока не обращайте внимания, разберемся с этим на следующих шагах).

Вместо T можно указать любую другую букву или слово, это ни на что не влияет:

1 задание 
```go
func Produce[T any](val T, n int) []T {
    vals := make([]T, n)
    for i := range n {
        vals[i] = val
    }
    return vals
}
```

Ограничение типа
Например, функции или срезы сравнивать через == нельзя. А поскольку у нас T может быть абсолютно любым (any), то в общем случае == работать не будет — вот компилятор и ругается.

Решение — ограничить множество типов T, разрешив только те типы, с которыми работают операторы == и != (так называемые comparable-типы):
```go
// Find находит элемент elem в срезе slice и возвращает его индекс.
func Find[T comparable](slice []T, elem T) int {
    for i, v := range slice {
        if v == elem {
            return i
        }
    }
    return -1
}
```
comparable — это специальное ограничение типа (type constraint), которое включает все сравнимые типы:

- числа, строки, логические значения,
- массивы,
- указатели,
- структуры (если все поля структуры тоже сравнимые),
- интерфейсы,
- и некоторые другие. https://go.dev/ref/spec#Comparison_operators


Упорядоченные типы
Есть множество типов, значения которых можно сравнивать через операторы < <= >= >. Такие типы называются упорядоченными (ordered). Они включают целые и вещественные числа, а также строки. Чтобы разрешить только упорядоченные типы, укажите cmp.Ordered в качестве ограничения:
```go
// Max возвращает максимальное значение из a и b.
func Max[T cmp.Ordered](a, b T) T {
    if a > b {
        return a
    }
    return b
}
```

Список разрешенных типов

Можно явно перечислить разрешенные типы через |. Например, разрешить искать только по целым и вещественным числам:

```go
func Find[T int | float64](slice []T, elem T) int {
    // реализация не меняется
}
Или по строкам и логическим значениям (сомнительное желание, но допустим):

func Find[T string | bool](slice []T, elem T) int {
    // реализация не меняется
}
```

Такой способ работает только для встроенных типов. Если, к примеру, мы захотим искать людей (тип Person) и домашних питомцев (тип Pet) по полю Name, которое есть у обоих типов — ничего не получится:

Здесь придется использовать следующий способ ограничить типы.

Интерфейс

В ограничении типа можно указать интерфейс. Воспользуемся этим, чтобы реализовать поиск по имени среди людей и животных.

Создадим интерфейс Named с методом Name:
```go
// Named - кто-то, у кого есть имя.
type Named interface {
    Name() string
}
```
```go
// Find находит элемент с именем name в срезе slice и возвращает его индекс.
func Find[T Named](slice []T, name string) int {
    for i, v := range slice {
        if v.Name() == name {
            return i
        }
    }
    return -1
}
```

Этих возможностей вам хватит почти всегда:

- any для любых типов;
- comparable для типов, которые поддерживают == !=;
- cmp.Ordered для типов, которые поддерживают < <= >= >;
- T1 | T2 | ... чтобы явно перечислить разрешенные типы (только для встроенных);
- I чтобы разрешить только типы, которые реализуют интерфейс I.


Как вы знаете, у функции может быть несколько обычных параметров. Точно так же и параметров типа может быть несколько. Например, мы можем написать собственный вариант встроенной функции make для карт:

```go
// MakeMap создает карту указанной емкости.
func MakeMap[K comparable, V any](size int) map[K]V {
    m := make(map[K]V, size)
    return m
}
```
Когда параметров типа несколько, они перечисляются через запятую (аналогично обычным параметрам). После названия параметра указывается ограничение типа (аналогично тому, как указываются типы у обычных параметров):

Пример использования:
```go
m := MakeMap[string, int](5)
m["one"] = 1
m["two"] = 2
m["three"] = 3
fmt.Println(m)
// map[one:1 two:2 three:3]
```

Здесь компилятор не может вывести конкретные K и V автоматически, поэтому мы указываем их явно:
```go
m := MakeMap[string, int](5)
```

Чтобы работал автоматический вывод типов, K и V должны использоваться среди обычных параметров:

```go
// SingleMap создает карту с одной записью.
func SingleMap[K comparable, V any](key K, val V) map[K]V {
    return map[K]V{key: val}
}
m := SingleMap("one", 1)
fmt.Println(m)
// map[one:1]
```

2 task
```go
func ZipMap[K comparable, V any](keys []K, vals []V) map[K]V {
    size := min(len(keys), len(vals))
    m := make(map[K]V, size)
    for i := range size {
        m[keys[i]] = vals[i]
    }
    return m
}
```

Generic-срез
```go
// Slice - обобщенный срез, который
// работает со значениями типа T.
type Slice[T any] []T
```

Как и в случае с обобщенной функцией, T — это параметр типа, а []T указывает, что наш Slice создан на основе среза значений типа T.

```go
// Reverse инвертирует порядок элементов в срезе.
func (s Slice[T]) Reverse() {
    for i := 0; i < len(s)/2; i++ {
        s[i], s[len(s)-1-i] = s[len(s)-1-i], s[i]
    }
}
```

Параметр T — часть типа Slice, поэтому получатель метода (s Slice[T]), а не просто (s Slice).

Наш срез можно использовать для чисел и строк:
```go
// срез целых чисел
intSlice := Slice[int]{1, 2, 3}           // (1)
intSlice.Reverse()
fmt.Println(intSlice)
// [3 2 1]

// срез строк
strSlice := Slice[string]{"a", "b", "c"}  // (2)
strSlice.Reverse()
fmt.Println(strSlice)
// [c b a]
```

При создании конкретного значения типа Slice в ➊ и ➋ мы указываем конкретный тип (int или string) вместо T. Кажется, что компилятор мог бы вывести тип самостоятельно (как было с generic-функцией), но вот не умеет. Так что указывайте явно.

Наш срез подойдет и для любых других значений:

```go
type Person struct {
    Name string
}

// срез значений типа Person
personSlice := Slice[Person]{
    {Name: "Alice"}, {Name: "Bob"}, {Name: "Cindy"},
}
personSlice.Reverse()
fmt.Println(personSlice)
// [{Cindy} {Bob} {Alice}]
```

Generic-тип

Параметризовать типом можно не только срез. Вот, например, структура с двумя полями типа T:
```go
// Pair - пара значений.
type Pair[T any] struct {
    first  T
    second T
}
```
Метод Swap меняет местами значения в паре:
```go
// Swap меняет местами значения в паре.
func (p *Pair[T]) Swap() {
    p.first, p.second = p.second, p.first
}
```