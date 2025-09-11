package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

// возращает конкретный тип кастомной ошибки
func test() *customError {
	// ... do something
	return nil
}

func main() {
	var err error   // data = nil, itab = nil
	err = test()    // data = nil, itab = *customError
	if err != nil { // сравнение с nil вернет true, так как err != nil (в itab лежит *customError)
		println("error")
		return
	}
	println("ok") // не будет вызвано, так как err != nil
}
