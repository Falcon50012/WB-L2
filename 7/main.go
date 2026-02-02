// Что выведет программа?

package main

import (
	"fmt"
	"math/rand"
	"time"
)

func asChan(vs ...int) <-chan int {
	c := make(chan int)
	go func() {
		for _, v := range vs {
			c <- v
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
		close(c)
	}()
	return c
}

func merge(a, b <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		for {
			select {
			case v, ok := <-a:
				if ok {
					c <- v
				} else {
					a = nil
				}
			case v, ok := <-b:
				if ok {
					c <- v
				} else {
					b = nil
				}
			}
			if a == nil && b == nil {
				close(c)
				return
			}
		}
	}()
	return c
}

func main() {
	rand.Seed(time.Now().Unix())
	a := asChan(1, 3, 5, 7)
	b := asChan(2, 4, 6, 8)
	c := merge(a, b)
	for v := range c {
		fmt.Print(v)
	}
}

// Объяснение работы конвейера с использованием select:
//
// В данном коде реализован конвейер типа fan-in.
// Два отправителя асинхронно отправляют значения в свои каналы.
// Функция merge с помощью select конкурентно читает из обоих каналов
// и пересылает значения в один выходной канал.
// Закрытые каналы исключаются из select через присваивание nil,
// а выходной канал закрывается после завершения всех входных потоков.
// Порядок вывода значений не детерминирован.
