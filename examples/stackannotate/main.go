package main

func main() {
	var n int
	for i := 0; i < 10; i++ {
		n += foo(1)
	}
	println(n)
	<-chan struct{}(nil)
}

func foo(a int) int {
	return bar(a, 2)
}

func bar(a int, b int) int {
	s := 3
	for i := 0; i < 100; i++ {
		s += a * b
	}
	return s
}
