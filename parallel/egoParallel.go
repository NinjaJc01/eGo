package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
	"sync"	//WaitGroup for thread synchronization

	"github.com/ericlagergren/decimal"
)

var (
	precision  int
	iterations uint64
	increment uint64
	hard       bool
	channel    chan *decimal.Big
	factorial_buf  []*decimal.Big	//store all needed factorials
	wg sync.WaitGroup	//synchronize factorial calculation threads
)

func main() {
	// Options
	precPtr := flag.Int("p", 10001, "Precision for calculations")
	iterPtr := flag.Uint64("i", 1625, "Value of infinity")
	incrementPtr := flag.Uint64("increment", 64, "increment size")
	hard := flag.Bool("hard", false, "Stress your hardware more, more iterations! Forces set iterations and precison, overiding any set.")
	debug := flag.Bool("debug", false, "Used for debugging. This will write to log.txt")
	flag.Parse()

	// Iterations
	precision = *precPtr
	iterations = *iterPtr
	increment = *incrementPtr
	if *hard {
		iterations = 4288
		precision = 30001
	}
	start := time.Now()
	channel = make(chan *decimal.Big, iterations)
	var answer = decimal.WithPrecision(precision).SetUint64(0)
	factorial_buf = make([]*decimal.Big, 2*(iterations+1),2*(iterations+1))
	factorial_buf[0] = decimal.WithPrecision(precision).SetUint64(1)

	var iteration_overflow uint64 = iterations%increment
	iterations -= iteration_overflow

	//factorial_initialize sets values of factorial_buf[i] = i!/(lower-1)! by calculating partial product
	for i := uint64(1); i < 2*(iterations); i += increment {
		wg.Add(1)
		go factorial_initialize(i,i + increment)
	}
	wg.Add(1)
	go factorial_initialize(2*iterations+1,2*(iterations+iteration_overflow+1))
	wg.Wait()

	//factorial_calc_final sets values of factorial_buf[i] = i! by multiplying by the precalculated value
	//the last value in each series must be pre-propagated to ensure the next series is correct
	//increment 1 is a special case as no parallel work can be completed
	if increment != 1 {
		for i := uint64(1); i < 2*(iterations); i+= increment {
			factorial_buf[i+increment-1].Mul(factorial_buf[i-1],factorial_buf[i+increment-1])
			wg.Add(1)
			go factorial_calc_final(i,i+ increment-1)
		}
		wg.Add(1)
		go factorial_calc_final(2*iterations+1,2*(iterations+iteration_overflow+1))
		wg.Wait()
	} else {
		for i := uint64(1); i < 2*(iterations); i++ {
			factorial_buf[i].Mul(factorial_buf[i-1],factorial_buf[i])
		}
	}

	//calculate partial sums of e
	for i := uint64(0); i < iterations; i+= increment {
		go series(i, i+ increment)
	}
	go series(iterations, iterations+iteration_overflow)

	//finish adding the sum
	for counter := uint64(0); counter < iterations; counter+= increment {
		answer = answer.Add(<-channel, answer)
	}
	answer = answer.Add(<-channel, answer)

	end := time.Now()

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
		logger.Printf("\nRun Time: %vs\n", end.Sub(start).Seconds())
	}
	// Print running time to console
	fmt.Printf("Run Time: %vs\n", end.Sub(start).Seconds())

}

//set values of factorial_buf[i] = i!/(lower-1)! by calculating partial product (lower*(lower+1)...*upper)
func factorial_initialize(lower, upper uint64) {
	defer wg.Done()
	factorial_buf[lower] = decimal.WithPrecision(precision).SetUint64(lower)
	for i := lower+1; i < upper; i++ {
		factorial_buf[i] = decimal.WithPrecision(precision).SetUint64(i)
		factorial_buf[i].Mul(factorial_buf[i-1],factorial_buf[i])
	}
}

//finishes calculating factorial_buf[i] = i! by multiplying by the precalculated (lower-1)!
func factorial_calc_final(lower, upper uint64) {
	defer wg.Done()
	for i := lower; i < upper; i++ {
		factorial_buf[i].Mul(factorial_buf[lower-1],factorial_buf[i])
	}
}

//partial sum of e
func series(lower, upper uint64) {
	var res = decimal.WithPrecision(precision).SetUint64(0)
	for n := lower; n < upper; n++ {
		add := decimal.WithPrecision(precision).SetUint64(((2 * n) + 2))
		add.Quo(add, factorial_buf[(2*n)+1])
		res.Add(res, add)
	}
	channel <- res
}
