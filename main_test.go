package main

import "testing"


func TestFibonacci(t *testing.T) {
	ans := fibonacci(5)
	fib := [6]int{0, 1, 1, 2, 3, 5}
	if ans[4] != fib[4] {
		t.Errorf("fibonacci(5) = %d; want {0, 1, 1, 2, 3, 5}", ans)
	}
}

//func TestFibonacciFail(t *testing.T) {
//	ans := fibonacci(5)
//	fib := [6]int{0, 1, 1, 2, 3, 5}
//	if ans[4] != fib[5] {
//		t.Errorf("fibonacci(5) = %d; want {0, 1, 1, 2, 3, 5}", ans)
//	}
//}

func TestGetFibonacciNumberForTesting(t *testing.T) {
	ans := getFibonacciNumberForTesting(11)
	if ans != 89 {
		t.Errorf("getFibonacciNumberForTesting(11) = %d; want 89", ans)
	}
}

//func TestGetFibonacciNumberForTestingFail(t *testing.T) {
//	ans := getFibonacciNumberForTesting(11)
//	if ans != 144 {
//		t.Errorf("getFibonacciNumberForTesting(11) = %d; want 89", ans)
//	}
//}

func TestGetNumbersLessThanForTesting(t *testing.T) {
	ans := getNumbersLessThanForTesting(120)
	if ans != 12 {
		t.Errorf("getNumbersLessThanForTesting(120) = %d; want 12", ans)
	}
}

//func TestGetNumbersLessThanForTestingFail(t *testing.T) {
//	ans := getNumbersLessThanForTesting(120)
//	if ans != 10 {
//		t.Errorf("getNumbersLessThanForTesting(120) = %d; want 12", ans)
//	}
//}

