package main

// Sieve of Atkin

import (
	"log"
	"math"
	"os"
	"strconv"

	"github.com/willf/bitset"
	"github.com/wtnb75/intseq"
)

func dumpfile(fn string) {
	var primes = intseq.NewIntSeq()
	fp, err := os.Open(fn)
	if err != nil {
		log.Panic("open failed", err)
	}
	primes.Read(fp)
	log.Println(primes)
}

func main() {
	var maxv = uint(65536)
	args := os.Args[1:len(os.Args)]
	log.Println("args", args)
	if len(args) != 0 {
		if _, err := os.Stat(args[0]); err == nil {
			dumpfile(args[0])
			return
		}
		maxv64, _ := strconv.ParseUint(args[0], 0, 64)
		maxv = uint(maxv64)
	}
	var maxvv = uint(math.Sqrt(float64(maxv)))
	log.Println("calculate primes", maxv, maxvv)
	var bs = bitset.New(maxv)
	for x := uint(1); x < maxvv; x++ {
		for y := uint(1); y < maxvv; y++ {
			var k = uint(4*x*x + y*y)
			if k < maxv && (k%12 == 1 || k%12 == 5) {
				bs.Flip(k)
			}
			k = uint(3*x*x + y*y)
			if k < maxv && (k%12 == 7) {
				bs.Flip(k)
			}
			if x > y {
				k = 3*x*x - y*y
				if k < maxv && k%12 == 11 {
					bs.Flip(k)
				}
			}
		}
	}
	cnt := bs.Count()
	log.Println("bs.Count=", cnt, "/", bs.Len(), bs.Len()/(bs.Len()-cnt))
	bs.Set(2)
	bs.Set(3)
	for n := uint(5); n < maxvv; n++ {
		if bs.Test(n) {
			n2 := n * n
			for k := n2; k < maxv; k += n2 {
				bs.Clear(k)
			}
		}
	}
	fp, err := os.Create("prime.dat")
	if err != nil {
		log.Panic("file open error")
	}
	defer fp.Close()
	var primes = intseq.NewIntSeq()
	for n := uint(2); n < maxv; n++ {
		if bs.Test(n) {
			primes.Add(uint64(n))
		}
	}
	primes.Write(fp)
}
