package async

import (
	"log"
	"testing"
	"time"
)

func TestAsync(t *testing.T) {

	for i := 0; i < 30; i++ {

		num := i
		fn := func() {
			log.Println("put ", num)
			time.Sleep(time.Microsecond * 200)
		}

		Async(fn)
	}

}

func TestAsyncWithPanic(t *testing.T) {

	for i := 0; i < 30; i++ {

		num := i
		fn := func() {
			log.Panicf("panic %d !\n", num)
			time.Sleep(time.Microsecond * 200)
		}

		Async(fn)
	}

	pool.Stop()
}
