package atomic

import (
	"sync/atomic"
	"testing"
	"time"
)

// Тест на корректность работы функции Atomic
func TestAtomic(t *testing.T) {
	var done int32 = 0

	go func() {
		time.Sleep(10 * time.Millisecond) // Симуляция работы горутины
		atomic.StoreInt32(&done, 1)
	}()

	start := time.Now()
	for atomic.LoadInt32(&done) == 0 {
		time.Sleep(3 * time.Millisecond)
	}
	elapsed := time.Since(start)

	if atomic.LoadInt32(&done) != 1 {
		t.Errorf("Ожидалось done=1, получено done=%d", done)
	}
	if elapsed > 100*time.Millisecond {
		t.Errorf("Выполнение заняло слишком много времени: %v", elapsed)
	}
}

// Бенчмарк для измерения производительности функции Atomic
func BenchmarkAtomic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var done int32 = 0

		go func() {
			atomic.StoreInt32(&done, 1)
		}()

		for atomic.LoadInt32(&done) == 0 {
			time.Sleep(3 * time.Millisecond)
		}
	}
}
