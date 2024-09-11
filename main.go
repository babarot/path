package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
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
			var dirs1 []string = strings.Split(input, "/")
			var dirs2 []string
			switch dirs1[0] {
			case "":
				dirs1 = dirs1[1:]
				dirs2 = []string{""}
			case ".":
				dirs1 = dirs1[1:]
				dirs2 = []string{"."}
			}
			for _, order := range orders {
				order -= 1 // handle 1 origin
				if len(dirs1)-1 < order {
					// avoid runtime error index out of range
					continue
				}
				dirs2 = append(dirs2, dirs1[order])
			}
			fmt.Fprintf(os.Stdout, "%s\n", strings.Join(dirs2, "/"))
		}
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}
	return nil
}

func run() error {
	var orders []int
	flag.Parse()
	for _, arg := range flag.Args() {
		order, err := strconv.Atoi(arg)
		if err != nil {
			continue
		}
		if order == 0 {
			return errors.New("Cannot use 0 as index because of 1 origin")
		}
		orders = append(orders, order)
	}
	if err := construct(orders...); err != nil {
		return err
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
