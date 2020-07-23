package main

import (
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

var s512 = flag.Bool("512", false, "print sha512")
var s384 = flag.Bool("384", false, "print sha384")

func main() {
	flag.Parse()
	input := flag.Args()
	var tokens = make(chan struct{}, 200)
	var ch = make(chan struct{})

	for _, file := range input {
		go func(f string) {
			tokens <- struct{}{}
			b, err := ioutil.ReadFile(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "shasum: %v\n", err)
				<-tokens
				ch <- struct{}{}
				return
			}
			switch {
			case *s512:
				digest := sha512.Sum512(b)
				fmt.Printf("%x\t%s\n", digest, f)
			case *s384:
				digest := sha512.Sum384(b)
				fmt.Printf("%x\t%s\n", digest, f)
			default:
				digest := sha256.Sum256(b)
				fmt.Printf("%x\t%s\n", digest, f)
			}
			<-tokens
			ch <- struct{}{}
		}(file)
	}
	for range input {
		<-ch
	}
}
