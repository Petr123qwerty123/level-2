Что выведет программа? Объяснить вывод программы.

```go
package main

import (
    "fmt"
)

func main() {
    a := [5]int{76, 77, 78, 79, 80}
    var b []int = a[1:4]
    fmt.Println(b)
}
```

Ответ:
```
Вывод: [77 78 79]

Тут мы создали массив a, передали в b слайс, ссылающийся на a. Числа внутри квадратных скобок - 
Simple slice expression (т.к. 2 аргумента - low = 1, high = 4).

low - с какого элемента обрезать массив, high - low - сколько элементов будет в слайсе, причем high <= len(a),
таким образом 4 - 1 = 3 элемента в слайсе, начиная со второго (индексация с 0, значит 1 - второй), то есть 77, 78, 79.
```