package worker

var counter int

func Increment() {
	go func() {
		counter++
	}()
}
