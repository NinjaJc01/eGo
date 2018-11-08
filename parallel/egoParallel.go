package main
import (
	"github.com/ericlagergren/decimal"
	"fmt"
	"flag"
	"time"
)
var (
	precision int
	iterations uint64
	channel chan *decimal.Big
)

func main() {
	start := time.Now()
	precPtr := flag.Int("p", 10001, "Precision for calculations")
	iterPtr := flag.Uint64("i", 1625, "Value of infinity")
	flag.Parse()
	precision = *precPtr
	iterations = *iterPtr
	channel = make(chan *decimal.Big, iterations)
	//go series(0,*iterPtr)
	var answer = decimal.WithPrecision(precision).SetUint64(0)
	for i := uint64(1); i < iterations; i++ {
		go series(i-1,i)
	}
	for counter := uint64(0); counter < iterations-1; counter++ {
		answer = answer.Add(<- channel, answer)
		//fmt.Print(".")
		//time.Sleep(time.Millisecond*5)
	}
	fmt.Println()
	fmt.Println(answer)
	fmt.Print("Time: ")
	fmt.Println(time.Now().Sub(start))
}
func series(lower, upper uint64) {
	var res  = decimal.WithPrecision(precision).SetUint64(0)
	for n := lower; n < upper; n++ {
		add := decimal.WithPrecision(precision).SetUint64(((2*n)+2))
		add.Quo(add, factorial((2*n)+1))
		res.Add(res, add)
	}
	channel <- res
}

func factorial(x uint64) (fact *decimal.Big) {
	fact = decimal.WithPrecision(precision).SetUint64(1)
	//fmt.Println("Prec",fact.Precision())
	for i := x; i > 0; i-- {
		fact.Mul(fact, decimal.New((int64(i)),0))
	}
	//fmt.Println("ActualPrec:",fact.Precision())
	return
}