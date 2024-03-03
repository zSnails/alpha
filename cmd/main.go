package main

import (
	"encoding/json"
	"fmt"
	"lang_test/parser"
	"lang_test/tokenizer"
	"log"
	"os"
)

const CONTENT = `if condition then print() else print(22) end`

func main() {
	lex := tokenizer.NewTokenizer(CONTENT)
	parser, err := parser.NewParser(lex)
	if err != nil {
		panic(err)
	}

	node, err := parser.SingleCommand()
	if err != nil {
		log.Fatal(err)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	fmt.Print("Ast ")
	_ = enc.Encode(node)
}
