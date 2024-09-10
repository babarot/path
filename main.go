package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func construct(orders ...int) error {
	var r io.Reader = os.Stdin
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		input := scanner.Text()
		switch len(orders) {
		case 0:
			fmt.Fprintf(os.Stdout, "%s\n", input)
		default:
			dirs1 := strings.Split(input, "/")
			dirs2 := []string{}
			for _, order := range orders {
				if len(dirs1)-1 < order {
					// avoid runtime error index out of range
					continue
				}
				dirs2 = append(dirs2, dirs1[order])
			}
			log.Println(dirs1)
			log.Println(dirs2)
			fmt.Fprintf(os.Stdout, "%s\n", strings.Join(dirs2, "/"))
		}
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}
	return nil
}

func main() {
	var orders []int
	flag.Parse()
	for _, arg := range flag.Args() {
		order, err := strconv.Atoi(arg)
		if err != nil {
			continue
		}
		orders = append(orders, order)
	}
	if err := construct(orders...); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
