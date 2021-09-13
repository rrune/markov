package main

import (
	"fmt"
	"rune/markov/markov"
)

type Chain map[string]string

func main() {
	c := markov.New()
	c.Read("./db.txt")
	fmt.Println(c.GetChain(20))

	//c.Train("./text/ted.txt")
	//c.Write("./db.txt")
}
