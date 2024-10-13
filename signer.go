package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

var md5Mutex sync.Mutex

// SingleHash вычисляет crc32(data) и crc32(md5(data))
func SingleHash(in, out chan interface{}) {
	var wg sync.WaitGroup

	for data := range in {
		wg.Add(1)
		go func(data interface{}) {
			defer wg.Done()

			value := fmt.Sprintf("%v", data)
			fmt.Println(data, "SingleHash", "data", value)

			// Параллельно считаем crc32(data) и crc32(md5(data))
			crc32Chan := make(chan string)
			md5Crc32Chan := make(chan string)

			// crc32(data)
			go func() {
				crc32Hash := DataSignerCrc32(value)
				fmt.Println(data, "SingleHash", "crc32(data)", crc32Hash)
				crc32Chan <- crc32Hash
			}()

			// crc32(md5(data))
			go func() {
				md5Mutex.Lock()
				md5Hash := DataSignerMd5(value)
				md5Mutex.Unlock()
				fmt.Println(data, "SingleHash", "md5(data)", md5Hash)

				crc32Md5Hash := DataSignerCrc32(md5Hash)
				fmt.Println(data, "SingleHash", "crc32(md5(data))", crc32Md5Hash)
				md5Crc32Chan <- crc32Md5Hash
			}()

			crc32Result := <-crc32Chan
			md5Crc32Result := <-md5Crc32Chan

			// Объединяем результаты
			result := crc32Result + "~" + md5Crc32Result
			fmt.Println(data, "SingleHash", "result", result)
			out <- result
		}(data)
	}

	wg.Wait()
	//close(out)
}

// MultiHash вычисляет crc32(th+data) для th от 0 до 5
func MultiHash(in, out chan interface{}) {
	var wg sync.WaitGroup

	for data := range in {
		wg.Add(1)
		go func(data interface{}) {
			defer wg.Done()

			input := fmt.Sprintf("%v", data)
			var result [6]string
			var mhWg sync.WaitGroup

			// Параллельно вычисляем crc32(th+data)
			for th := 0; th < 6; th++ {
				mhWg.Add(1)
				go func(th int) {
					defer mhWg.Done()
					thStr := fmt.Sprintf("%d%s", th, input)
					result[th] = DataSignerCrc32(thStr)
					fmt.Println(data, "MultiHash: crc32(th+step1))", th, result[th])
				}(th)
			}

			mhWg.Wait()
			concatenatedResult := strings.Join(result[:], "")
			fmt.Println(data, "MultiHash result:", concatenatedResult)
			out <- concatenatedResult
		}(data)
	}

	wg.Wait()
	//close(out)
}

// CombineResults получает все результаты, сортирует их и объединяет через _
func CombineResults(in, out chan interface{}) {
	var results []string

	for data := range in {
		result := fmt.Sprintf("%v", data)
		results = append(results, result)
	}

	// Сортируем результаты
	sort.Strings(results)
	finalResult := strings.Join(results, "_")
	fmt.Println("CombineResults:", finalResult)
	out <- finalResult

	//close(out)
}

// ExecutePipeline организует последовательное выполнение job'ов через каналы
func ExecutePipeline(jobs ...job) {
	var wg sync.WaitGroup
	in := make(chan interface{})

	for _, j := range jobs {
		out := make(chan interface{})
		wg.Add(1)

		go func(in, out chan interface{}, j job) {
			defer wg.Done()
			j(in, out)
			close(out) // Закрываем канал после завершения работы
		}(in, out, j)

		in = out // Переносим out на следующий шаг
	}

	wg.Wait()
}

func main() {
	fmt.Println("Pipeline example")

	// Пример работы ExecutePipeline с SingleHash, MultiHash и CombineResults
	ExecutePipeline(
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
	)
}
