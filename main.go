package main

import (
	"fmt"

	"github.com/regex-compiler/compiler"
)

func main() {
	str := "ab"

	postfix := compiler.Re2post(str)
	nfa := compiler.Post2nfa([]rune(postfix))

	match := compiler.MatchRe(nfa, str)

	fmt.Println("match => ", match)
}
