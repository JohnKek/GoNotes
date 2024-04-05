Чтобы создать собственный тип ошибки, достаточно реализовать метод Error().
Он будет вызываться по дефолту если мы попытаемся напечатать полученную ошибку

```go
// фрагмент кода стандартной библиотеки
type error interface {
Error() string
}
```
Создадим ошибку, которая описывает проблему поиска substr в строке src:
```go
type lookupError struct {
src    string
substr string
}

func (e lookupError) Error() string {
return fmt.Sprintf("'%s' not found in '%s'", e.substr, e.src)
}
```
Напишем функцию indexOf(), которая возвращает индекс вхождения подстроки substr в строку src. Если вхождения нет, возвращает ошибку типа lookupError:
````go

func indexOf(src string, substr string) (int, error) {
idx := strings.Index(src, substr)
if idx == -1 {
// Создаем и возвращаем ошибку типа `lookupError`.
return -1, lookupError{src, substr}
}
return idx, nil
}
````
Проверим работу indexOf() для разных подстрок.
```go
src := "go is awesome"
for _, substr := range []string{"go", "js"} {
if res, err := indexOf(src, substr); err != nil {
fmt.Printf("indexOf(%#v, %#v) failed: %v\n", src, substr, err)
} else {
fmt.Printf("indexOf(%#v, %#v) = %v\n", src, substr, res)
}
}
// indexOf("go is awesome", "go") = 0
// indexOf("go is awesome", "js") failed: 'js' not found in 'go is awesome'

```
Поскольку indexOf() возвращает общий тип error, чтобы получить доступ к конкретному объекту ошибки, придется использовать приведение типа:
```go
_, err := indexOf(src, "js")
if err, ok := err.(lookupError); ok {
fmt.Println("err.src:", err.src)
fmt.Println("err.substr:", err.substr)
}
// err.src: go is awesome
// err.substr: js
```