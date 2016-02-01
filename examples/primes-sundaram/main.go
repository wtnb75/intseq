package main

// Sieve of Sundaram

import (
	"log"
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
	var maxvv = uint((maxv-2-1)/2 + 1)
	log.Println("calculate primes", maxv, maxvv)
	// var maxv uint = 65536
	var bs = bitset.New(maxvv)
	for j := uint(1); j < maxvv; j++ {
		for i := uint(1); i <= j; i++ {
			if i+j+2*i*j > maxvv {
				break
			}
			bs.Set(i + j + 2*i*j)
		}
	}
	cnt := bs.Count()
	log.Println("bs.Count=", cnt, "/", bs.Len(), bs.Len()/(bs.Len()-cnt))
	fp, err := os.Create("prime.dat")
	if err != nil {
		log.Panic("file open error")
	}
	defer fp.Close()
	var primes = intseq.NewIntSeq()
	primes.Add(2)
	for i := uint(1); i < maxvv; i++ {
		if !bs.Test(i) {
			primes.Add(uint64(i*2 + 1))
		}
	}
	primes.Write(fp)
}
