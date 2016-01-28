// calculate primes and write them to file
package main

import (
	"log"
	"math"
	"os"
	"strconv"

	"github.com/wtnb75/intseq"
)

var primes = intseq.NewIntSeq()

func isPrime(n uint64) bool {
	var m = uint64(math.Sqrt(float64(n)))
	var ret = true
	primes.Each(func(v uint64) bool {
		if n%v == 0 {
			ret = false
			return true
		}
		if v > m {
			return true
		}
		return false
	})
	return ret
}

func dumpfile(fn string) {
	fp, err := os.Open(fn)
	if err != nil {
		log.Panic("open failed", err)
	}
	primes.Read(fp)
	log.Println(primes)
}

// Usage1: $0 [maxvalue]
// Usage2: $0 filename
func main() {
	var maxv = uint64(65536)
	args := os.Args[1:len(os.Args)]
	log.Println("args", args)
	if len(args) != 0 {
		if _, err := os.Stat(args[0]); err == nil {
			dumpfile(args[0])
			return
		}
		maxv, _ = strconv.ParseUint(args[0], 0, 64)
	}
	log.Println("calculate primes", maxv)
	fp, err := os.Create("prime.dat")
	if err != nil {
		log.Panic("open failed", err)
	}
	defer fp.Close()
	primes.Add(2)
	for i := uint64(3); i < maxv; i++ {
		if isPrime(i) {
			primes.Add(uint64(i))
			if primes.Count()%(1024*128) == 0 {
				log.Println(i, primes.Count(), uint64(i)/primes.Count())
			}
		}
	}
	primes.Write(fp)
	log.Println("wrote", primes.Count(), "primes")
}
