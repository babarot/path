package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func construct(nums ...int) error {
	var r io.Reader = os.Stdin
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		input := scanner.Text()
		switch len(nums) {
		case 0:
			fmt.Fprintf(os.Stdout, "%s\n", input)
			continue
		}
		var indexes []int
		var dirs1 []string = strings.Split(input, "/")
		var dirs2 []string
		switch dirs1[0] {
		case "":
			dirs1 = dirs1[1:]
			if nums[0] == 1 {
				dirs2 = []string{""}
			}
		case ".":
			dirs1 = dirs1[1:]
			if nums[0] == 1 {
				dirs2 = []string{"."}
			}
		}
		switch last := nums[len(nums)-1]; {
		case last < 0:
			start := nums[0]
			end := len(dirs1) + nums[len(nums)-1]
			if end < start {
				return fmt.Errorf("Index last %d (%d) is smaller than index start %d", nums[len(nums)-1], end, start)
			}
			for i := start; i <= end; i++ {
				indexes = append(indexes, i)
			}
		case last == 0:
			start := nums[0]
			end := len(dirs1)
			for i := start; i <= end; i++ {
				indexes = append(indexes, i)
			}
		default:
			indexes = nums
		}
		var leakage bool
		for _, idx := range indexes {
			idx -= 1 // handle 1 origin
			if len(dirs1)-1 < idx {
				leakage = true
				continue // avoid runtime error index out of range
			}
			if dirs1[idx] == "" {
				continue // remove a space
			}
			dirs2 = append(dirs2, dirs1[idx])
		}
		// do not print if leakage
		if leakage {
			continue
		}
		fmt.Fprintf(os.Stdout, "%s\n", strings.Join(dirs2, "/"))
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}
	return nil
}

func run() error {
	var nums []int
	flag.Parse()
	for _, arg := range flag.Args() {
		switch {
		case regexp.MustCompile(`^\d+\.\.$`).MatchString(arg):
			ranges := strings.Split(arg, "..")
			start, err := strconv.Atoi(ranges[0])
			if err != nil {
				return err
			}
			nums = append(nums, start)
			nums = append(nums, 0)
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
				return errors.New("left..right: left should be positive number")
			}
			if end < 0 {
				nums = append(nums, start)
				nums = append(nums, end)
			} else {
				for num := start; num <= end; num++ {
					nums = append(nums, num)
				}
			}
		case regexp.MustCompile(`^\d+$`).MatchString(arg): // positive number
			num, err := strconv.Atoi(arg)
			if err != nil {
				continue
			}
			if num == 0 {
				return errors.New("Cannot use 0 as index because of 1 origin")
			}
			nums = append(nums, num)
		default:
			return fmt.Errorf("%s: invalid arguments", arg)
		}
	}
	if err := construct(nums...); err != nil {
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
