package main
import (
	"github.com/ericlagergren/decimal"
	"fmt"
	"flag"
	"time"
)
var (
	precision int
)
func main() {
	start := time.Now()
	precPtr := flag.Int("p", 10001, "Precision for calculations")
	iterPtr := flag.Uint64("i", 1625, "Value of infinity")
	flag.Parse()
	precision = *precPtr
	fmt.Println(series(0,*iterPtr))
	fmt.Print("Time: ")
	fmt.Println(time.Now().Sub(start))
}
func series(lower, upper uint64) (res *decimal.Big) {
	res = decimal.WithPrecision(precision).SetUint64(0)
	for n := lower; n < upper; n++ {
		add := decimal.WithPrecision(precision).SetUint64(((2*n)+2))
		add.Quo(add, factorial((2*n)+1))
		res.Add(res, add)
	}
	return
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