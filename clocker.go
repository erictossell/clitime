package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func parse_args(args []string) string {
	return args[0]
}

func main() {
	args := os.Args[1:]
	start := time.Now()
	fmt.Println(start)
	fmt.Println(parse_args(args))
	fmt.Println(len(args))
	fmt.Println(strings.Join(args, " - "))
}
