У wire есть 2 основные концепции
- провайдер
- инжектор
Основным механизмом wire является провайдер. Функция которая возвращает значение
Пример
```go
package foobarbaz

type Foo struct {
    X int
}

// ProvideFoo returns a Foo.
func ProvideFoo() Foo {
    return Foo{X: 42}
}
```

Провайдер должен быть экспортированным
Также провайдер может иметь параметры 
```go
package foobarbaz

// ...

type Bar struct {
    X int
}

// ProvideBar returns a Bar: a negative Foo.
func ProvideBar(foo Foo) Bar {
    return Bar{X: -foo.X}
}
```

Также допускается вовзращение **error** из провайдера

```go
package foobarbaz

import (
    "context"
    "errors"
)

// ...

type Baz struct {
    X int
}

// ProvideBaz returns a value if Bar is not zero.
func ProvideBaz(ctx context.Context, bar Bar) (Baz, error) {
    if bar.X == 0 {
        return Baz{}, errors.New("cannot provide baz when bar is zero")
    }
    return Baz{X: bar.X}, nil
}
```

Есть возможность группировки провайдеров в группы. Полезно когда несколько провайдеров используются совместно.
Чтобы объединить провайдеры используйте метод wire.NewSet
```go
package foobarbaz

import (
    // ...
    "github.com/google/wire"
)

// ...

var SuperSet = wire.NewSet(ProvideFoo, ProvideBar, ProvideBaz)
```

После объявления можно дополнять сет создавая на основе него новый
```go
package foobarbaz

import (
    // ...
    "example.com/some/other/pkg"
)

// ...

var MegaSet = wire.NewSet(SuperSet, pkg.OtherSet)
```

# Внедрение

Приложение внедреняет провайдеров с помощью инжектора.
Инжектор - функция которая вызывает провайдеров в необходимом порядке
Инжектор объявляется путем записи объявления функции, тело которой является вызовом `wire.Build`
```go
// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package main

import (
    "context"

    "github.com/google/wire"
    "example.com/foobarbaz"
)

func initializeBaz(ctx context.Context) (foobarbaz.Baz, error) {
    wire.Build(foobarbaz.MegaSet)
    return foobarbaz.Baz{}, nil
}
```
Как и поставщики, инжекторы могут быть параметризованы на входах (которые затем отправляются поставщикам) и могут возвращать ошибки. Аргументы для `wire.Build`те же, что и `wire.NewSet`: они формируют набор поставщиков. Это набор поставщиков, который используется во время генерации кода для этого инжектора.
## Расширенные функции

### Интерфейсы привязки

Часто внедрение зависимости используется для привязки конкретной реализации для интерфейса. Wire сопоставляет входы с выходами через [type identity](https://golang.org/ref/spec#Type_identity) , поэтому может возникнуть желание создать поставщика, который возвращает тип интерфейса. Однако это не будет идиоматичным, поскольку наилучшей практикой Go является возврат [конкретных типов](https://github.com/golang/go/wiki/CodeReviewComments#interfaces) . Вместо этого вы можете объявить привязку интерфейса в наборе поставщиков:
```go
type Fooer interface {
    Foo() string
}

type MyFooer string

func (b *MyFooer) Foo() string {
    return string(*b)
}

func provideMyFooer() *MyFooer {
    b := new(MyFooer)
    *b = "Hello, World!"
    return b
}

type Bar string

func provideBar(f Fooer) string {
    // f will be a *MyFooer.
    return f.Foo()
}

var Set = wire.NewSet(
    provideMyFooer,
    wire.Bind(new(Fooer), new(*MyFooer)),
    provideBar)
```
Первый аргумент — `wire.Bind`это указатель на значение желаемого типа интерфейса, а второй аргумент — это указатель на значение типа, реализующего интерфейс. Любой набор, включающий привязку интерфейса, должен также иметь поставщика в том же наборе, который предоставляет конкретный тип.

**Когда wire будет создавать зависимости, он будет знать, что любой раз, когда требуется Fooer, он должен предоставить указатель на MyFooer типо того получается?**

### Поставщики структур
	Структуры могут быть созданы с использованием предоставленных типов. Используйте `wire.Struct`функцию для создания типа структуры и сообщите инжектору, какие поля должны быть внедрены. Инжектор заполнит каждое поле, используя поставщика для типа поля. Для полученного типа структуры `S`, `wire.Struct`предоставляет `S`и `*S`. Например, если даны следующие поставщики:

```go
package main  
  
import (  
    "fmt"  
    "github.com/google/wire")  
  
type Foo int  
type Bar int  
  
func ProvideFoo() Foo { return 10 }  
  
func ProvideBar() Bar { return 20 }  
  
type FooBar struct {  
    MyFoo Foo  
    MyBar Bar  
}  
  
var Set = wire.NewSet(  
    ProvideBar,  
    ProvideFoo,  
    wire.Struct(new(FooBar), "MyFoo", "MyBar"))  
  
func main() {  
    result := InitializeFooBar()  
    fmt.Println(result)  
  
}
```

Иногда полезно предотвратить заполнение определенных полей инжектором, особенно при передаче `*`в `wire.Struct`. Вы можете пометить поле тегом , `` `wire:"-"` ``чтобы Wire игнорировал такие поля. Например:
```go
type Foo struct {
    mu sync.Mutex `wire:"-"`
    Bar Bar
}
```
Специальная строка `"*"`может использоваться как сокращение, чтобы сообщить инжектору о необходимости внедрить все поля. Так что `wire.Struct(new(FooBar), "*")`получается тот же результат, что и выше.
### Обязательные значения
Иногда бывает полезно привязать базовое значение (обычно `nil`) к типу. Вместо того, чтобы инжекторы зависели от функции-поставщика одноразового использования, вы можете добавить выражение значения к набору поставщиков.
```go
type Foo struct {
    X int
}

func injectFoo() Foo {
    wire.Build(wire.Value(Foo{X: 42}))
    return Foo{}
}
```

### Использование полей структуры в качестве поставщиков
Иногда поставщики, которые нужны пользователю, являются некоторыми полями структуры. Если вы обнаружите, что пишете поставщика, как `getS`в примере ниже, чтобы продвигать поля структуры в предоставленные типы:
```go
type Foo struct {
    S string
    N int
    F float64
}

func getS(foo Foo) string {
    // Bad! Use wire.FieldsOf instead.
    return foo.S
}

func provideFoo() Foo {
    return Foo{ S: "Hello, World!", N: 1, F: 3.14 }
}

func injectedMessage() string {
    wire.Build(
        provideFoo,
        getS)
    return ""
}
```

Вместо этого вы можете `wire.FieldsOf`использовать эти поля напрямую, без написания `getS`:

```go
func injectedMessage() string {
    wire.Build(
        provideFoo,
        wire.FieldsOf(new(Foo), "S"))
    return ""
}
```

### Функции очистки

Если поставщик создает значение, которое необходимо очистить (например, закрытие файла), то он может вернуть замыкание для очистки ресурса. Инжектор будет использовать это либо для возврата агрегированной функции очистки вызывающему объекту, либо для очистки ресурса, если поставщик, вызванный позже в реализации инжектора, вернет ошибку.

```go
func provideFile(log Logger, path Path) (*os.File, func(), error) {
    f, err := os.Open(string(path))
    if err != nil {
        return nil, nil, err
    }
    cleanup := func() {
        if err := f.Close(); err != nil {
            log.Log(err)
        }
    }
    return f, cleanup, nil
}
```

Функция очистки гарантированно вызывается перед функцией очистки любых входных данных поставщика и должна иметь сигнатуру `func()`.

### Синтаксис альтернативного инжектора

[](https://github.com/google/wire/blob/main/docs/guide.md#alternate-injector-syntax)

Если вам надоело писать `return foobarbaz.Foo{}, nil`в конце объявления функции инжектора, вы можете написать его более кратко с помощью `panic`:

```go
func injectFoo() Foo {
    panic(wire.Build(/* ... */))
}
```