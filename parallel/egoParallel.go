package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ericlagergren/decimal"
)

var (
	precision  int
	iterations uint64
	hard       bool
	channel    chan *decimal.Big
)

func main() {
	// Options
	precPtr := flag.Int("p", 10001, "Precision for calculations")
	iterPtr := flag.Uint64("i", 1625, "Value of infinity")
	hard := flag.Bool("hard", false, "Stress your hardware more, more iterations! Forces set iterations and precison, overiding any set.")
	debug := flag.Bool("debug", false, "Used for debugging. This will write to log.txt")
	flag.Parse()

	// Iterations
	precision = *precPtr
	iterations = *iterPtr
	if *hard {
		iterations = 4288
		precision = 30001
	}
	start := time.Now()
	channel = make(chan *decimal.Big, iterations)
	//go series(0, *iterPtr)
	var answer = decimal.WithPrecision(precision).SetUint64(0)
	for i := uint64(1); i < iterations; i++ {
		go series(i-1, i)
	}
	for counter := uint64(0); counter < iterations-1; counter++ {
		answer = answer.Add(<-channel, answer)
		//fmt.Print(".")
		//time.Sleep(time.Millisecond*5)
	}

	// Logging. Only creates log.txt with -debug option
	if *debug {
		f, err := os.OpenFile("log.txt",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		defer f.Close()
		logger := log.New(f, "eGoDecimal ", log.LstdFlags)
		// Add things to log for debug here
		logger.Println(answer)
		logger.Printf("\nRun Time: %v\n", time.Now().Sub(start))
	}
	// Print result to console
	fmt.Printf("Run Time: %v\n", time.Now().Sub(start))

}
func series(lower, upper uint64) {
	var res = decimal.WithPrecision(precision).SetUint64(0)
	for n := lower; n < upper; n++ {
		add := decimal.WithPrecision(precision).SetUint64(((2 * n) + 2))
		add.Quo(add, factorial((2*n)+1))
		res.Add(res, add)
	}
	channel <- res
}

func factorial(x uint64) (fact *decimal.Big) {
	fact = decimal.WithPrecision(precision).SetUint64(1)
	//fmt.Println("Prec",fact.Precision())
	for i := x; i > 0; i-- {
		fact.Mul(fact, decimal.New((int64(i)), 0))
	}
	//fmt.Println("ActualPrec:",fact.Precision())
	return
}
