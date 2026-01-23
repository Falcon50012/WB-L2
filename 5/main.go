// Что выведет программа?

package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	// ... do something
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

// Объяснение вывода программы:

// Программа выведет error,
// так как функция test возвращает nil указатель типа *customError
// в переменную интерфейсного типа, поле data интерфейса указывает
// на nil значение customError, однако при этом поле itab данного интерфейса
// указывает на таблицу реализации интерфейса error типом customError и не является nil,
// поэтому интерфейсное значение является непустым,
// и при сравнении err с nil сравнивается сам интерфейс,
// вследствие чего и выполняется ветка вывода error.
