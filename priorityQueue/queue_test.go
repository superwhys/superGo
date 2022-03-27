package priorityQueue

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

func main() {

	var wg sync.WaitGroup

	pq := Init(5000)
	wg.Add(4)

	go func() {
		fmt.Println("start push")
		for i := 1; i <= 50000; i++ {
			key := fmt.Sprintf("why_%s", strconv.Itoa(i))
			e := &Entry{
				Key:      key,
				Priority: rand.Intn(5),
			}
			pq.PushQueue(e)
			//fmt.Printf("push_data: %+v\n", e)
		}
		wg.Done()
	}()

	time.Sleep(time.Second * 5)

	for i := 0; i < 3; i++ {
		go func() {
			fmt.Println("start pop")
			for {
				data := pq.PopQueue()
				fmt.Printf("---------pop_data: %+v\n", data)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
