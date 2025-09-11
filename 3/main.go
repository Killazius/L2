package main

import (
	"fmt"
	"os"
)

// Foo возращает интерфейс ошибки
//
//	type iface struct {
//		tab  *itab // метеданные, там как раз будет наш тип *os.PathError
//		data unsafe.Pointer // nil
//	}
//
// поэтому err == nil вернет false, так как tab != nil
//
// fmt.Println(err) выведет nil, потому что внутри интерфейса лежит nil-указатель (data == nil), а тип (tab) задан.

// пустой интерфейс interface{} имеет такую структуру:
// у него нет метаданных *itab с методами, но есть тип и данные
// поэтому пустой интерфейс с nil-указателем внутри будет равен nil
//
//	type eface struct {
//		_type *_type
//		data  unsafe.Pointer
//	}
func Foo() error {
	var err *os.PathError = nil
	return err
}

func main() {
	err := Foo()
	fmt.Println(err)        // nil
	fmt.Println(err == nil) // false
}
