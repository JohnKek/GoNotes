Принципы SOLID - это пять основных принципов объектно-ориентированного программирования, которые были представлены
Робертом Мартином в начале 2000-х годов. Каждая буква в слове "SOLID" представляет собой один из этих принципов:

S - Принцип единственной ответственности (Single Responsibility Principle): Каждый класс должен быть ответственен только
за одну важную часть функциональности программы.

O - Принцип открытости/закрытости (Open/Closed Principle): Классы должны быть открыты для расширения, но закрыты для
модификации.

L - Принцип подстановки Барбары Лисков (Liskov Substitution Principle): Объекты в программе должны быть заменяемыми на
экземпляры их подтипов без изменения правильности выполнения программы.

I - Принцип разделения интерфейса (Interface Segregation Principle): Клиенты не должны зависеть от интерфейсов, которые
они не используют. Более специфические интерфейсы предпочтительнее более общих.

D - Принцип инверсии зависимостей (Dependency Inversion Principle): Модули верхнего уровня не должны зависеть от модулей
нижнего уровня. Оба типа модулей должны зависеть от абстракций.

SOLID
Single Responsibility Principle гласит, что класс или модуль должен иметь только одну причину для изменения. Корочег
говоря - каждый класс или функция должны решать лишь одну задачу, не более. Если у вас есть функция или класс, который
меняется по нескольким причинам, это первый звоночек, что вы нарушаете SRP.

Когда есть компонент, который в ответе только за одну задачу, его намного проще изменять, не затрагивая остальные части
системы.

Неправильное применение SRP:

```go
package main

import (
	"net/http"
)

type User struct {
	ID        int
	FirstName string
	LastName  string
}

// функция обрабатывает HTTP-запросы И управляет пользователями
func (u *User) HandleRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// получение данных пользователя
	case "POST":
		// создание нового пользователя
	}
}
```

HandleRequest класса User выполняет две задачи: обрабатывает HTTP-запросы и управляет пользователями, это большая ошибка

Правильное применение SRP:

```go
package main

import (
	"net/http"
)

type User struct {
	ID        int
	FirstName string
	LastName  string
}

type UserHandler struct {
	// ...
}

// UserHandler отвечает только за обработку HTTP-запросов
func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// получение данных пользователя
	case "POST":
		//создание нового пользователя
	}
}

```

User хранит данные о пользователе, а UserHandler управляет HTTP-запросами. Каждый класс фокусируется на своей уникальной
задаче. Если потребуется изменить логику обработки HTTP-запросов, можно это сделать в UserHandler, не затрагивая класс
User.

----
O
pen/closed Principle - программные сущности (классы, модули, функции и т.д.) должны быть открыты для расширения, но
закрыты для изменения. Нужно свой код таким образом, чтобы для добавления новой функциональности не требовалось менять
существующий код. Соблюдение этого уменьшает вероятность возникновения багов, т.к вам не нужно трогать уже работающий
код

Пример без использования OCP:

```go
package main

import "fmt"

type Printer struct{}

func (p *Printer) Print(data string) {
	fmt.Println("Data: ", data)
}

// Допустим, нам нужно добавить функционал печати в HTML
// придется измнить класс Printer, что нарушает OCP
func (p *Printer) PrintHTML(data string) {
	fmt.Println("<habr>" + data + "</habr>")
}

func main() {
	printer := Printer{}
	printer.Print("Hello, World!")
	printer.PrintHTML("Hello, HABR World!")
}
```

Для добавления новой функциональности (печати в HTML), мы изменили класс Printer. Это нарушает OCP.

Пример с использованием OCP:
```go
package main

import "fmt"

type Printer interface {
	Print(data string)
}

type TextPrinter struct{}

func (p *TextPrinter) Print(data string) {
	fmt.Println("Data: ", data)
}

type HTMLPrinter struct{}

func (h *HTMLPrinter) Print(data string) {
	fmt.Println("<html>" + data + "</html>")
}

func main() {
	var printer Printer

	printer = &TextPrinter{}
	printer.Print("Hello, World!")

	printer = &HTMLPrinter{}
	printer.Print("Hello, HTML World!")

}
```

Вместо изменения существующего кода, мы расширили функциональность системы, добавив новую реализацию интерфейса Printer.
Соблюдаем OCP: существующий код не изменяется, а новый функционал добавляется через новые реализации.

LSP гласит, что объекты в программе должны быть заменяемыми на экземпляры их подтипов без изменения правильности работы
программы. Это звучит как-то непонятно, но на самом деле всё просто: если у вас есть класс-родитель и класс-потомок, то
любой код, который использует родительский класс, должен работать так же хорошо и с объектами дочернего класса.

Пример:
```go
package main

import "fmt"

// Bird базовый тип
type Bird struct{}

func (b *Bird) Fly() {
	fmt.Println("Птица летит")
}

// Penguin - подтип Bird, но не может летать
type Penguin struct {
	Bird
}

func main() {
	var bird = &Bird{}
	bird.Fly()

	var penguin = &Penguin{}
	penguin.Fly() // Нарушение LSP, т.к. пингвины не летают

}

```

Penguin наследуется от Bird, но не соответствует поведению, ожидаемому от Bird, что нарушает LSP.

В данном случае, так как пингвины не умеют летать (или все же умеют?), нам следует отделить способность летать от
базового класса Bird:
```go
package main

import "fmt"

// Bird базовый тип
type Bird struct{}

func (b *Bird) MakeSound() {
	fmt.Println("Птица издает звук")
}

// FlyingBird интерфейс для летающих птиц
type FlyingBird interface {
	Fly()
}

// Sparrow подтип Bird, который умеет летать
type Sparrow struct {
	Bird
}

func (s *Sparrow) Fly() {
	fmt.Println("Воробей летит")
}

// Penguin подтип Bird, но не реализует интерфейс FlyingBird
type Penguin struct {
	Bird
}

func main() {
	var sparrow FlyingBird = &Sparrow{}
	sparrow.Fly()

	var penguin = &Penguin{}
	penguin.MakeSound() // Penguin может издавать звук, но не летать

}

```

Bird остается базовым классом для всех птиц, обеспечивая общее поведение (например, издавать звук). Создается интерфейс
FlyingBird для птиц, которые могут летать. Sparrow реализует интерфейс FlyingBird, так как воробьи умеют летать. Penguin
является подтипом Bird, но не реализует интерфейс FlyingBird, поскольку пингвины не летают.

----

ISP утверждает, что юзеры не должны быть вынуждены зависеть от интерфейсов, которые они не используют. Это означает, что
вместо одного наполненного интерфейса лучше иметь несколько тонких и специализированных

Пример:
```go
package main

type Printer interface {
	Print(document string)
}

type Scanner interface {
	Scan(document string)
}

// MultiFunctionDevice наследует от обоих интерфейсов
type MultiFunctionDevice interface {
	Printer
	Scanner
}

// класс, реализующий только функцию печати
type SimplePrinter struct{}

func (p *SimplePrinter) Print(document string) {
	// реализация печати
}

// класс, реализующий обе функции
type AdvancedPrinter struct{}

func (p *AdvancedPrinter) Print(document string) {
}

func (p *AdvancedPrinter) Scan(document string) {
}

```

Не заставляем SimplePrinter реализовывать функции сканирования, которые он не использует, соблюдая ISP.

DIP гласит, что высокоуровневые модули не должны зависеть от низкоуровневых модулей. Оба типа модулей должны зависеть от
абстракций.

Приме:

package main

import "fmt"

// Интерфейс для абстракции хранения данных
type DataStorage interface {
Save(data string)
}

// Низкоуровневый модуль для хранения данных в файле
type FileStorage struct {}

func (fs *FileStorage) Save(data string) {
fmt.Println("Сохранение данных в файл:", data)
}

// Высокоуровневый модуль, не зависит напрямую от FileStorage
type DataManager struct {
storage DataStorage // зависит от абстракции
}

func (dm *DataManager) SaveData(data string) {
dm.storage.Save(data) // делегирование сохранения
}

func main() {
fs := &FileStorage{}
dm := DataManager{storage: fs}
dm.SaveData("Тестовые данные")
}
DataManager не зависит напрямую от FileStorage. Вместо этого он использует интерфейс DataStorage, что позволяет легко
заменить способ хранения данных без изменения DataManager.

DRY
Каждый кусочек знаний в системе должен иметь единственное, недвусмысленное, авторитетное представление в рамках системы.
Проще говоря, надо избегать повторения одного и того же кода в разных частях вашей программы. Когда логика дублируется,
любое изменение в ней требует обновления во всех местах, где она встречается.

Нарушение DRY:

package main

import "fmt"

type User struct {
Name string
Age int
}

func (u User) PrintName() {
fmt.Println(u.Name)
}

func (u User) PrintAge() {
fmt.Println(u.Age)
}

func main() {
user := User{Name: "Alex", Age: 30}
user.PrintName()
user.PrintAge()
}
Соблюдение DRY:

package main

import "fmt"

type User struct {
Name string
Age int
}

func (u User) PrintInfo() {
fmt.Printf("Name: %s, Age: %d\n", u.Name, u.Age)
}

func main() {
user := User{Name: "Alex", Age: 30}
user.PrintInfo()
}
В первом примере мы видм, что методы PrintName и PrintAge дублируют логику вывода информации о пользователе. Во втором
примере мы исправляем это, объединяя логику в одном методе PrintInfo.

