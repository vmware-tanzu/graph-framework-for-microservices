package main

import (
	dircompare "graph-framework-for-microservices/install-validator/internal/dir-compare"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		panic("no directory provided")
	}
	b, t, e := dircompare.CheckDir(os.Args[1])
	if e != nil {
		panic(e)
	}
	if b {
		panic(t)
	}
}
