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

type CLI struct {
	Stdin          io.Reader
	Stdout, Stderr io.Writer
}

func (c *CLI) construct(sequential bool, nums ...int) error {
	if len(nums) == 0 {
		return fmt.Errorf("invalid usage")
	}
	var r io.Reader = c.Stdin
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		input := scanner.Text()
		var indexes []int
		var dirs1 []string = strings.Split(input, "/")
		var dirs2 []string
		switch dirs1[0] {
		case "": // a case of `/root/1/2/3`
			// remove from dirs1 to treat it without blank in operating
			// but join to output dirs
			dirs1 = dirs1[1:]
			if nums[0] == 1 {
				dirs2 = []string{""}
			}
		case ".": // a case of `./local/1/2/3`
			// remove from dirs1 to treat it without blank in operating
			// but join to output dirs
			dirs1 = dirs1[1:]
			if nums[0] == 1 {
				dirs2 = []string{"."}
			}
		}
		if sequential {
			start := nums[0]
			end := nums[1]
			switch {
			case 0 < end && end < start:
				return fmt.Errorf(
					"On positive number the right side (%d) should be smaller than the left side (%d)",
					start, end)
			case end < 0:
				end = len(dirs1) + 1 + end
				if end < start {
					return fmt.Errorf(
						"index out of range (range should be -%d..%d)\n%v",
						len(dirs1), len(dirs1), dirs1)
				}
			case end == 0:
				end = len(dirs1)
			}
			for i := start; i <= end; i++ {
				indexes = append(indexes, i)
			}
		} else {
			for _, i := range nums {
				if i < 0 {
					indexes = append(indexes, len(dirs1)+1+i)
				} else {
					indexes = append(indexes, i)
				}
			}
		}
		var leakage bool
		for _, idx := range indexes {
			idx -= 1 // handle 1-origin
			if len(dirs1)-1 < idx {
				leakage = true
				continue // avoid runtime error index out of range
			}
			if dirs1[idx] == "" {
				continue // remove a space
			}
			dirs2 = append(dirs2, dirs1[idx])
		}
		if leakage {
			// do not print if leakage
			continue
		}
		fmt.Fprintf(c.Stdout, "%s\n", strings.Join(dirs2, "/"))
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}
	return nil
}

func (c *CLI) run(args []string) error {
	var sequential bool
	var nums []int
	for _, arg := range args {
		switch {
		case regexp.MustCompile(`^\d+\.\.(-?\d+)?$`).MatchString(arg): // range number
			ranges := strings.Split(arg, "..")
			start, err := strconv.Atoi(ranges[0])
			if err != nil {
				return err
			}
			end := 0
			if ranges[1] != "" {
				end, err = strconv.Atoi(ranges[1])
				if err != nil {
					return err
				}
			}
			nums = []int{start, end}
			sequential = true
		case regexp.MustCompile(`^-?\d+$`).MatchString(arg): // single number
			num, err := strconv.Atoi(arg)
			if err != nil {
				continue
			}
			if num == 0 {
				return errors.New("cannot use 0 as a index because of 1-origin")
			}
			nums = append(nums, num)
		default:
			return fmt.Errorf("%s: invalid argument type", arg)
		}
	}
	if err := c.construct(sequential, nums...); err != nil {
		return err
	}
	return nil
}

func main() {
	cli := &CLI{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	flag.Parse()
	if err := cli.run(flag.Args()); err != nil {
		fmt.Fprintln(cli.Stderr, err)
		os.Exit(1)
	}
}
