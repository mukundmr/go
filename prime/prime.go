package main

import (
	"fmt"
	"math/big"
	//	"sync"
	//	"time"
)

type splitMap map[string]*big.Int

/* some handy constants. big.Int isn't friendly */
var zero = big.NewInt(int64(0))
var one = big.NewInt(int64(1))
var two = big.NewInt(int64(2))
var three = big.NewInt(int64(3))
var four = big.NewInt(int64(4))
var eight = big.NewInt(int64(8))

/* create splits to use */
func get_splits(num *big.Int) (splits splitMap) {
	var one_eight, two_eight, three_eight, four_eight big.Int
	make_odd(one_eight.Div(num, eight))
	make_odd(two_eight.Mul(&one_eight, two))
	make_odd(three_eight.Mul(&one_eight, three))
	make_odd(four_eight.Mul(&one_eight, four))
	splits = make(splitMap)
	splits[three.String()] = &one_eight // 3..x/8
	splits[one_eight.String()] = &two_eight
	splits[two_eight.String()] = &three_eight
	splits[three_eight.String()] = &four_eight
	return
}

/* helpers */
func isDivisible(num_to_check *big.Int, divisor *big.Int) (isDivisible bool) {
	var tmp big.Int
	isDivisible = (tmp.Rem(num_to_check, divisor).Cmp(zero) == 0)
	return
}

/* make odd */
func make_odd(num_to_check *big.Int) {
	if isDivisible(num_to_check, two) {
		num_to_check.Sub(num_to_check, one)
	}
}

/* evaluation loop */
func looper(num_to_check *big.Int, start, end *big.Int, quit chan bool) (ret chan *big.Int) {
	ret = make(chan *big.Int)
	fmt.Println(start, end)
	go func() {
		var interrupt = false
		var itr big.Int
		for itr.Set(start); itr.Cmp(end) < 0; itr.Add(&itr, two) {
			select {
			default:
				if isDivisible(num_to_check, &itr) {
					interrupt = true
					quit <- true
					var val big.Int
					val.Set(&itr)
					ret <- &val
					fmt.Println("returning ", itr.String(), " from ", start, " to ", end)
					break
				}
			case <-quit:
				interrupt = true
				ret <- zero
				quit <- true
				fmt.Println("interupt due to quit from ", start, " to ", end)
				break
			}
		}
		if !interrupt {
			ret <- zero
			fmt.Println("returning ", zero.String())
		}
	}()
	return
}

/* mux */
func mux4to1(a, b, c, d <-chan *big.Int) (ret chan *big.Int) {
	ret = make(chan *big.Int)
	go func() {
		for {
			select {
			case z := <-a:
				ret <- z
			case z := <-b:
				ret <- z
			case z := <-c:
				ret <- z
			case z := <-d:
				ret <- z
			}
		}
	}()
	return
}

/* prime numbers */

func isPrime(num_to_check *big.Int) (isPrime bool) {
	isPrime = true
	if num_to_check.Cmp(four) < 0 {
		return
	} else if isDivisible(num_to_check, two) {
		isPrime = false
		return
	}
	splits := get_splits(num_to_check)
	var stream = make(chan (chan *big.Int), len(splits))
	var quit = make(chan bool)
	var start,stop big.Int
	for key, value := range splits {
		start.SetString(key, 10)
		stop.Set(value)
		stream <- looper(num_to_check, &start, &stop, quit)
	}
	resultChan := mux4to1(<-stream, <-stream, <-stream, <-stream)
	for waitloop := 0; waitloop < len(splits); waitloop++ {
		select {
		case val := <-resultChan:
			if val.Cmp(zero) != 0 {
				fmt.Println(num_to_check, " divisible by ", val)
				isPrime = false
			}
		}
	}
	return
}

func main() {
	var i big.Int
	i.SetString("9513", 10)
	if isPrime(&i) {
		fmt.Println(i.String(), "is a prime number")
	} else {
		fmt.Println(i.String(), "is not a prime number")
	}
}
