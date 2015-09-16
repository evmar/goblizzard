// mpqtool is a tool for dumping mpq files.
package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"blizzard/mpq"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("must specify input path")
	}
	if len(os.Args) < 3 {
		log.Fatalf("must specify command")
	}
	path, command, args := os.Args[1], os.Args[2], os.Args[3:]

	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	r := mpq.NewReader(f)

	switch command {
	case "ls":
		for _, f := range r.GetFileList() {
			fmt.Printf("%s\n", f)
		}
	case "cat":
		if len(args) < 1 {
			log.Fatalf("must specify file")
		}
		path = args[0]

		f := r.OpenFile(path)
		_, err := io.Copy(os.Stdout, f)
		if err != nil {
			panic(err)
		}
	}
}
