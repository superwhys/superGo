package goroutinePool

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	routinePool := NewPool(2)
	for i := 0; i < 10; i++ {
		routinePool.Add(1)
		go func(i int) {
			time.Sleep(time.Second)
			fmt.Println("the NumGoroutine continue is:",runtime.NumGoroutine())
			routinePool.Done()
		}(i)
	}
	routinePool.Wait()
}
