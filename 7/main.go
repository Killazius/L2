package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Функция asChan преобразует слайс целых чисел в канал, отправляя каждый элемент с случайной задержкой.
func asChan(vs ...int) <-chan int {
	c := make(chan int)
	// при создании канала, запускаем горутину, которая будет писать в канал
	// и закрывать его по завершении
	// еще это важно, чтобы сразу из функции отдать канал, а не ждать пока все запишется
	go func() {
		for _, v := range vs {
			c <- v
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
		close(c)
	}()
	return c
}

// Функция merge объединяет два канала в один, читая из обоих и отправляя значения в результирующий канал.
// Когда оба входных канала закрываются, результирующий канал также закрывается.
func merge(a, b <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		// В бесконечном цикле читаем из обоих каналов
		for {
			// рандомно выбирает из какого канала читать
			// если канал закрыт, то присваивает nil, чтобы не читать из него больше, это мы понимаем при чтении
			// когда канал == nil, то в ok будет false
			// type scase struct {
			//	c    *hchan         // chan
			//	elem unsafe.Pointer // data element
			// }
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
			// когда оба канала будут nil, то завершаем горутину и закрываем результирующий канал
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
	// создаем два канала с разными наборами чисел
	a := asChan(1, 3, 5, 7)
	b := asChan(2, 4, 6, 8)
	// объединяем два канала в один
	c := merge(a, b)
	// читаем из канала пока он не закроется
	for v := range c {
		fmt.Print(v)
	}

}
