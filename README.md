# Regex compiler

Thompson's algorithm for parsing regular expressions

Steps:
* parse infix to postfix
* compile postfix to NFA
* walk nfa with input

Supported grammars:
  * `|`
  * `(` `)`
  * `*`

Inspired by <a href="https://swtch.com/~rsc/regexp/regexp1.html">Regular Expression Matching Can Be Simple And Fast </a> by Russ Cox
