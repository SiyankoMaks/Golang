package atomic

import (
	"fmt"
	"sync/atomic"
	"time"
)

func Atomic() {
	var done int32 = 0
	go func() {
		fmt.Println("hello world")
		atomic.StoreInt32(&done, 1)
	}()
	for atomic.LoadInt32(&done) == 0 {
		time.Sleep(3 * time.Millisecond)
	}
}