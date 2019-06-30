package main

import (
	"fmt"

	"github.com/regex-compiler/compiler"
)

func main() {
	regex := "abc*(ab)"
	str := "abccc"

	postfix := compiler.Re2post(regex)
	nfa := compiler.Post2nfa([]rune(postfix))
	match := compiler.MatchRe(nfa, str)

	fmt.Println("match => ", match)
}
