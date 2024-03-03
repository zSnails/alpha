package main

import (
	"encoding/json"
	"lang_test/parser"
	"lang_test/tokenizer"
	"log"
	"os"
)

const CONTENT = `variable = 44 + ("hola" == 44*(22-variableDeNumero+11))`

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
	_ = enc.Encode(node)
}
