package main

import ( 
        "fmt"
        "math/big"
			)
/*			
func sum(a [] int, c chan int) {
	sum := 0
	for _, v := range a {
		sum += v
		time.Sleep(100*time.Millisecond)
		c <- sum
	}
}

func main() {
	a := [] int { 7,10,2,234,2,54,103,-123,234,-1234,1231}
	c := make(chan int)
	go sum(a[:len(a)/2], c)
	go sum(a[len(a)/2:], c)
	x, y := <-c, <-c
	
	fmt.Println(x, y, x+y)
}
*/
			
			
/* prime numbers */
			
func isPrime(arg2 *big.Int) (ret bool) {
	arg1 := arg2
	fmt.Println("Start ", arg1)
	ret = true
	four := new(big.Int)
	four.SetString("4",10)
	two := new(big.Int)
	two.SetString("2",10)
	zero := new(big.Int)
	zero.SetString("0",10)
	if arg1.Cmp(four) < 0 {
		return
	}
	modres := arg1.Mod(arg1, two)
	
	if modres.Cmp(zero) == 0 {
		ret = false
		return
	}
	var itr big.Int
	arg1 = arg2
	fmt.Println("Before ", arg1)
	for itr.SetString("3",10); itr.Cmp(arg1.Div(arg1, two)) < 0; itr.Add(&itr,two) {
		arg1=arg2
		modres = arg1.Mod(arg1, &itr)
		fmt.Println("Inside ",itr, arg1, modres)
		if modres.Cmp(zero) == 0 {
			ret = false
			return
		}
	}
	return
}	

func main() {
	// body
	i := new(big.Int)
	i.SetString("15",10)
	fmt.Println(i, isPrime(i))
}		
			
