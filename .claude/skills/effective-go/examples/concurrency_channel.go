package worker

func Work(ch chan int) {
	go func() {
		ch <- 42
	}()
}
