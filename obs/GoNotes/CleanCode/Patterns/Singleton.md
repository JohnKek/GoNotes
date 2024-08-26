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