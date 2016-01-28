# Integer Sequence Array

[API reference](http://godoc.org/github.com/wtnb75/intseq)

# Usage

- `go get github.com/wtnb75/intseq`

```golang
import "github.com/wtnb75/intseq"
  :
var intarray = intseq.NewIntSeq()
  :
intarray.Add(1)
intarray.Add(2)  // put data
  :
if intarray.Contains(1) {
  // true
}
fmt.Println(intarray)  // convert to string  [1,2]
intarray.Write(os.Stdout)  // dump data to io.Writer
intarray.Read(os.Stdin)  // read data from io.Reader
```

# Examples

## primes

- source -> [primes.go](examples/primes/primes.go)
    - go run primes.go 65536
        - save prime numbers to file "prime.dat"
    - go run primes.go prime.dat
        - show primes
