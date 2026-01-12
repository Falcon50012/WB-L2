// Что выведет программа?

package main

import (
	"fmt"
	"os"
)

func Foo() error {
	var err *os.PathError = nil
	return err
}

func main() {
	err := Foo()
	fmt.Println(err)        // <nil>
	fmt.Println(err == nil) // false
}

//Объяснение внутреннего устройства интерфейсов и их отличия от пустых интерфейсов:
//
//Непустой интерфейс представлен структурой iface:
//type iface struct {
//	tab  *itab           -- указатель на itab структуру описывающую реализацию конкретного интерфейса конкретным типом
//	data unsafe.Pointer  -- указатель на данные
//}
// Использует методы
// Интерфейс считается nil только если tab == nil && data == nil
//
//Пустой интерфейс представлен структурой eface:
//type eface struct {
//	_type *_type          -- указатель на описание типа
//	data  unsafe.Pointer  -- указатель на данные
//}
// Не использует методы
// Интерфейс считается nil только если _type == nil && data == nil
