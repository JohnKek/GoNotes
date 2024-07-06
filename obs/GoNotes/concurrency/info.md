горутина легковесный поток 2Кб Треб 8мб 
![img_7.png](img_7.png)
Горутина управляет планировщик ГО

А потоками ОС
 
```go
fmt.Println(runtime.NumCPU())
```

```go
runtime.GOMAXPROCS(1)
```

```go
runtime.Gosched() 
```

defer выполняется в обратном порядке 

![img_8.png](img_8.png)

![img_9.png](img_9.png)
