**Одиночка** — это порождающий паттерн проектирования,
который гарантирует, что у класса есть только один
экземпляр, и предоставляет к нему глобальную
точку доступа.

![[Pasted image 20240820161742.png]]

**Одиночка** определяет статический метод getInstance ,
который возвращает единственный экземпляр своего
класса.
Конструктор одиночки должен быть скрыт от клиентов.
Вызов метода getInstance должен стать единственным
способом получить объект этого класса.

**Реализация через sync.Once**

```go
var once sync.Once  
  
type single struct {  
}  
  
var singleInstance *single  
  
func getInstance() *single {  
    if singleInstance == nil {  
       once.Do(  
          func() {  
             fmt.Println("Creating single instance now.")  
             singleInstance = &single{}  
          })  
    } else {  
       fmt.Println("Single instance already created.")  
    }  
  
    return singleInstance  
}
```

**Через поле** 

```go
var lock = &sync.Mutex{}  
  
type single struct {  
}  
  
var singleInstance *single  
  
func getInstance() *single {  
    if singleInstance == nil {  
       lock.Lock()  
       defer lock.Unlock()  
       if singleInstance == nil {  
          fmt.Println("Creating single instance now.")  
          singleInstance = &single{}  
       } else {  
          fmt.Println("Single instance already created.")  
       }  
    } else {  
       fmt.Println("Single instance already created.")  
    }  
  
    return singleInstance  
}
```


```go
package main

import (
	"fmt"
	"sync"
)

type singleton struct {
	value int
	mu    sync.Mutex // Мьютекс для обеспечения потокобезопасности
}

var instance *singleton
var once sync.Once

// GetInstance возвращает единственный экземпляр синглтона
func GetInstance() *singleton {
	once.Do(func() {
		fmt.Println("Creating single instance now.")
		instance = &singleton{}
	})
	return instance
}

// SetValue устанавливает значение в синглтоне
func (s *singleton) SetValue(val int) {
	s.mu.Lock()         // Блокируем мьютекс перед изменением
	defer s.mu.Unlock() // Освобождаем мьютекс после изменения
	s.value = val
}

// GetValue возвращает текущее значение синглтона
func (s *singleton) GetValue() int {
	s.mu.Lock()         // Блокируем мьютекс перед чтением
	defer s.mu.Unlock() // Освобождаем мьютекс после чтения
	return s.value
}

func main() {
	s1 := GetInstance()
	s1.SetValue(42)

	s2 := GetInstance()
	fmt.Println("Value from s2:", s2.GetValue()) // Вывод: Value from s2: 42

	// Проверяем, что s1 и s2 указывают на один и тот же экземпляр
	fmt.Println("s1 and s2 are the same instance:", s1 == s2) // true
}

```