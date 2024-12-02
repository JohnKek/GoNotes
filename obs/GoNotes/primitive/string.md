Для преобразования целого числа в строку можно воспользоваться функцией fmt. Sprintf; другой вариант — функция strconv. Itoa (“integer to ASCII”):
x := 123
у := fmt.Sprintf("%d", x)
fmt.Println(y, strconv.Itoa(x)) // "123 123"

Для форматирования чисел в другой системе счисления можно использовать функции Formatlnt и FormatUint:
fmt.Println(strconv.FormatInt(int64(x), 2)) // "1111011"

Для анализа строкового представления целого числа используйте такие функции пакета strconv, как Atoi или Parselnt, а для беззнаковых чисел — ParseUint:
х, err := strconv.Atoi("123") // x имеет тип int
// В десятичной системе счисления, до 64 битов: у, err := strconv.Parselnt("123", 10, 64)

```
//  Преобразование числа в строку
	num := 123
	numStr := strconv.Itoa(num)
	fmt.Println(numStr) // Выводит: "123"

	// Преобразование строки в число
	numStr := "456"
	num, err := strconv.Atoi(numStr)
	if err == nil {
		fmt.Println(num) // Выводит: 456
	// В десятичной системе счисления, до 64 битов: у, err := strconv.Parselnt("123", 10, 64)
```

