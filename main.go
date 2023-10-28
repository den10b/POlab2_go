package main

import (
	"crypto/sha256"
	"fmt"
	"sync"
	"time"
)

func main() {

	var nRoutines int
	fmt.Println("Введите количество горутин: ")
	_, err := fmt.Scanf("%d\n", &nRoutines)
	if err != nil {
		fmt.Printf("Ошибка ввода %v\n", err)
	}

	resArr := new([]string)
	enumerateStrings("", 5, resArr) // Создаём список всех возможных комбинаций
	toFind := []string{
		"1115dd800feaacefdf481f1f9070374a2a81e27880f187396db67958b207cbad",
		"3a7bd3e2360a3d29eea436fcfb7e44c735d117c42d1c1835420b6b9942dd4f1b",
		"74e1bb62f8dabb8125a58852b63bdf6eaef667cb56ac7f7cdba6d7305c50a22f",
	}

	//nRoutines := 4
	useAsync := true
	resChan := make(chan string, len(toFind))
	funcWG := new(sync.WaitGroup)
	beginTime := time.Now()
	for i := 0; i < nRoutines; i++ {
		funcWG.Add(1)
		if useAsync {
			// асинхронный запуск горутины
			go brute(resArr,
				i*len(*resArr)/nRoutines,
				(i+1)*len(*resArr)/nRoutines,
				&toFind, &resChan, funcWG)
		} else {
			// вызов функции, с ожиданием её выполнения
			brute(resArr,
				i*len(*resArr)/nRoutines,
				(i+1)*len(*resArr)/nRoutines,
				&toFind, &resChan, funcWG)
		}
	}
	funcWG.Wait()
	fmt.Println("Полученные перебором пароли: ")
	for i := 0; i < len(toFind); i++ {
		fmt.Println(<-resChan)
	}
	fmt.Println("Время выполнения: ", time.Since(beginTime))
}

// brute
func brute(passws *[]string, bI, eI int, toFind *[]string, resChan *chan string, wg *sync.WaitGroup) {
	for i := bI; i < eI; i++ {
		shaCode := sha256.Sum256([]byte((*passws)[i]))
		for _, rree := range *toFind {
			if fmt.Sprintf("%x", shaCode) == rree {
				*resChan <- (*passws)[i]
			}
		}

	}
	wg.Done()
}

// enumerateStrings перебираент строки заданной длинны
func enumerateStrings(prefix string, length int, resArr *[]string) {
	if length == 0 {
		*resArr = append(*resArr, prefix)
		return
	}

	for c := 'a'; c <= 'z'; c++ {
		enumerateStrings(prefix+string(c), length-1, resArr)
	}
}
