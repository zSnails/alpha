package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/zSnails/alpha/parser"
	"github.com/zSnails/alpha/tokenizer"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "error: missing filename")
		return
	}

	tok, err := tokenizer.FromFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}

	parser, err := parser.NewParser(tok)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}

	node, err := parser.Program()
	if err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}
	// fmt.Printf("%s\n", node)
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "    ")
	_ = encoder.Encode(node)
}
