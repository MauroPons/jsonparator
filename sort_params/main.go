package main

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"time"
)

var (
	numberLonger = int64(534234234234232654)
	numberSmaller = int64(324)
)

func main() {
	start1 := time.Now()
	fmt.Println(GetFirstDigitsOrLessFromInt64_v1(numberLonger,6))
	fmt.Println(GetFirstDigitsOrLessFromInt64_v1(numberSmaller,6))
	elapsed1 := time.Since(start1)
	log.Printf("v1 took %dms", elapsed1.Nanoseconds())


	start2 := time.Now()
	fmt.Println(GetFirstDigitsOrLessFromInt64_v2(numberLonger,6))
	fmt.Println(GetFirstDigitsOrLessFromInt64_v2(numberSmaller,6))
	elapsed2 := time.Since(start2)
	log.Printf("v2 took %dms", elapsed2.Nanoseconds())

	start3 := time.Now()
	fmt.Println(GetFirstDigitsOrLessFromInt64_v3(numberLonger,6))
	fmt.Println(GetFirstDigitsOrLessFromInt64_v3(numberSmaller,6))
	elapsed3 := time.Since(start3)
	log.Printf("v2 took %dms", elapsed3.Nanoseconds())
}

func GetFirstDigitsOrLessFromInt64_v3(n int64, qty int) int64 {
	s := fmt.Sprintf("%v", n)
	if len(s) > qty {
		fd := s[:qty]
		result, _ := strconv.ParseInt(fd, 10, 64)
		return result
	}
	return n
}


func GetFirstDigitsOrLessFromInt64_v2(n int64, qty int) int64 {
	nf := float64(n)
	pq := math.Pow10(qty)
	for nf > pq {
		nf = nf / 10
	}
	return int64(nf)
}

func GetFirstDigitsOrLessFromInt64_v1(n int64, qty int) int64 {
	s := strconv.FormatInt(n, 10)
	if len(s) > qty {
		fd := s[:qty]
		truncated, _ := strconv.ParseInt(fd, 10, 64)
		return truncated
	}
	return n
}
