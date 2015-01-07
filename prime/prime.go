package main

import (
	"fmt"
	"math/big"
	"sync"
	"time"
)

/* prime numbers */

func isPrime(num_to_check *big.Int) (ret bool) {
	// pathetic algorithm.  not meant for anything other than messing with language constructs
	var tmp1, zero, two, three, four, eight big.Int
	var modres *big.Int = big.NewInt(int64(0))

	tmp1.Set(num_to_check)
	ret = true
	four.SetInt64(int64(4))
	three.SetInt64(int64(3))
	two.SetInt64(int64(2))
	zero.SetInt64(int64(0))
	eight.SetInt64(int64(8))

	if tmp1.Cmp(&four) < 0 && tmp1.Cmp(&zero) > 0 {
		return
	}

	modres = modres.Mod(&tmp1, &two)

	if modres.Cmp(&zero) == 0 {
		fmt.Println("it is a even number ")
		ret = false
		return
	}
	// split the range of 3..num/2 into 4 (1/8, 2/8, 3/8, 4/8)
	var splits map[string]*big.Int = make(map[string]*big.Int)
	var one_eight, two_eight, three_eight, four_eight big.Int
	one_eight.Div(&tmp1, &eight)
	two_eight.Div(&tmp1, &eight)
	two_eight.Mul(&two_eight, &two)
	three_eight.Div(&tmp1, &eight)
	three_eight.Mul(&three_eight, &three)
	four_eight.Div(&tmp1, &eight)
	four_eight.Mul(&four_eight, &four)
	splits[three.String()] = &one_eight // 3..x/8
	splits[one_eight.String()] = &two_eight
	splits[two_eight.String()] = &three_eight
	splits[three_eight.String()] = &four_eight

	resChan := make(chan bool, 5)
	resChan <- true
	divByChan := make(chan string, 1)
	var wg sync.WaitGroup
	wg.Add(len(splits))
	fmt.Println(len(splits), " tracks.")
	for key, val := range splits {
		go func(start string, end *big.Int) {
			var local_itr big.Int
			local_itr.SetString(start, 10)
			var local_tmp big.Int
			defer wg.Done()
			fmt.Println("Running from ", local_itr.String(), " to ", (*end).String())
			for ; local_itr.Cmp(end) < 0; local_itr.Add(&local_itr, &two) {
				if local_tmp.Mod(&tmp1, &local_itr).Cmp(&zero) == 0 {
					divByChan <- local_itr.String()
					resChan <- false
					wg.Done()
					wg.Done()
					wg.Done()
					break
				}
			}
			fmt.Println("Running from ", start, " to ", (*end).String(), " done.")
		}(key, val)
	}
	wg.Wait() // fake wait. if one routine finishes, the wg is freed. we don't care about other routines
	for _ = range resChan {
		select {
		case ret = <-resChan:
			close(resChan) // assigned true first time. rest depends on Go routine
		case <-time.After(time.Millisecond): // one millisecond timeout.  wonder if one microsecond makes sense
			continue
		}
	}
	if ret == false {
		fmt.Println("Divisible by ", <-divByChan)
	}
	return
}

func main() {
	var i big.Int
	i.SetString("1347852092143232009187340981273401274309812734901234079109283881", 10)
	if isPrime(&i) {
		fmt.Println(i.String(), "is a prime number")
	} else {
		fmt.Println(i.String(), "is not a prime number")
	}
}
