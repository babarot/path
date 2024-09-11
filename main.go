package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
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
				if orders[0] == 1 {
					dirs2 = []string{""}
				}
			case ".":
				dirs1 = dirs1[1:]
				if orders[0] == 1 {
					dirs2 = []string{"."}
				}
			}
			if orders[len(orders)-1] < 0 {
				// if last order is negative number
				var newOrders []int
				start := orders[0]
				end := len(dirs1) + orders[len(orders)-1]
				if end < start {
					return fmt.Errorf("Index last %d (%d) is smaller than index start %d", orders[len(orders)-1], end, start)
				}
				for i := start; i <= end; i++ {
					newOrders = append(newOrders, i)
				}
				orders = newOrders
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
		log.Println(arg)
		switch {
		case regexp.MustCompile(`^\d+\.\.-?\d+$`).MatchString(arg):
			ranges := strings.Split(arg, "..")
			start, err := strconv.Atoi(ranges[0])
			if err != nil {
				return err
			}
			end, err := strconv.Atoi(ranges[1])
			if err != nil {
				return err
			}
			if start < 0 {
				return errors.New("left..right: left should be possitive number")
			}
			if end < 0 {
				orders = append(orders, start)
				orders = append(orders, end)
			} else {
				for order := start; order <= end; order++ {
					orders = append(orders, order)
				}
			}
		case regexp.MustCompile(`^\d+$`).MatchString(arg):
			order, err := strconv.Atoi(arg)
			if err != nil {
				continue
			}
			if order == 0 {
				return errors.New("Cannot use 0 as index because of 1 origin")
			}
			orders = append(orders, order)
		default:
			return errors.New("error")
		}
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
