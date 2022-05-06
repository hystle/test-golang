package testGoroutinePkg

import (
	"fmt"
)


func FibWithChan(c, q chan int){
	x, y := 0, 1
	for {
		select {
		case c <- x:
			fmt.Printf("x: %d, y: %d\n", x, y)
			x, y = y, x+y
		case rc := <-q:
			fmt.Printf("Quit: %d\n", rc)
			return
		}
	}
}