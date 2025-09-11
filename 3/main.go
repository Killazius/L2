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
func Foo() error {
	var err *os.PathError = nil
	return err
}

func main() {
	err := Foo()
	fmt.Println(err)        // nil
	fmt.Println(err == nil) // false
}
