Что выведет программа? Объяснить вывод программы.

```go
package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	{
		// do something
	}
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}
```

**Ответ:**

Вывод:
```
error
```

Для объяснения вывода этой задачи, воспользуемся [решением третьей задачи этого же уровня](listing03.md). В функции `main` мы 
создаем переменную `err` интерфейсного типа `error`. Функция `test` возвращает указателя на `customError` - структуру, 
реализующую интерфейс `error`. То есть значение конкретного типа переменной err - `nil`, конкретный 
тип - `*customError`. Далее мы сравниваем `err` с `nil`, в случае неравности мы возвращаем `"error"` и завершаем 
программу, иначе возвращаем `"ok"`. Так как сравнение идет с переменной интерфейсного типа, то равенство с `nil` будет 
только тогда, когда значение и тип - `nil`. Так как это не так, то вернется `"error"`.