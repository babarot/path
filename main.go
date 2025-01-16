package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/hashicorp/logutils"
)

type CLI struct {
	Stdin          io.Reader
	Stdout, Stderr io.Writer

	// options
	isRange bool
	// flags
	countDirHierarchy bool
}

func main() {
	log.SetOutput(logOutput())
	cli := &CLI{
		Stdin:             os.Stdin,
		Stdout:            os.Stdout,
		Stderr:            os.Stderr,
		isRange:           false,
		countDirHierarchy: false,
	}
	flag.BoolVar(&cli.countDirHierarchy, "c", false, "show a count of directory hierarchy")
	flag.Parse()
	if err := cli.main(flag.Args()); err != nil {
		fmt.Fprintln(cli.Stderr, err)
		os.Exit(1)
	}
}

func (c *CLI) main(args []string) error {
	var nums []int
	for _, arg := range args {
		switch {
		case regexp.MustCompile(`^\d+\.\.(-?\d+)?$`).MatchString(arg): // range number
			c.isRange = true
			var start, end int
			ranges := strings.Split(arg, "..")
			start, err := strconv.Atoi(ranges[0])
			if err != nil {
				return err
			}
			if ranges[1] != "" {
				end, err = strconv.Atoi(ranges[1])
				if err != nil {
					return err
				}
			}
			nums = []int{start, end}
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
			return fmt.Errorf("%s: invalid arguments", arg)
		}
	}
	if err := c.build(nums...); err != nil {
		return err
	}
	return nil
}

func (c *CLI) build(nums ...int) error {
	log.Printf("[DEBUG] build: args: %#v\n", nums)
	if len(nums) == 0 && !c.countDirHierarchy {
		return fmt.Errorf("too few arguments")
	}
	var r io.Reader = c.Stdin
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		input := scanner.Text()
		if c.countDirHierarchy {
			dirs := slices.DeleteFunc(strings.Split(input, "/"), func(s string) bool {
				return s == "" || s == "."
			})
			fmt.Println(len(dirs))
			continue
		}
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
		if c.isRange {
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
					return fmt.Errorf("index out of range (range should be -%d..%d)",
						len(dirs1), len(dirs1))
				}
			case end == 0:
				end = len(dirs1)
			case end > len(dirs1):
				log.Printf("[ERROR] LEAKED! %d is bigger than length of given path (%d)", end, len(dirs1))
				continue
			}
			for i := start; i <= end; i++ {
				indexes = append(indexes, i)
			}
		} else {
			for _, i := range nums {
				switch {
				case i < 0:
					j := i + len(dirs1) + 1
					log.Printf("[DEBUG] arg %d is negative number, so calculated from backward: %d", i, j)
					if j < 0 {
						log.Printf("[ERROR] %d: LEAKED! calculated result from backward is still negative, index out of range", j)
						continue
					}
					indexes = append(indexes, j)
				case i > len(dirs1):
					log.Printf("[ERROR] LEAKED! %d is bigger than length of given path (%d)", i, len(dirs1))
					continue
				default:
					indexes = append(indexes, i)
				}
			}
		}
		log.Printf("[DEBUG] indexes: %#v", indexes)
		for _, idx := range indexes {
			idx -= 1 // handle 1-origin
			if dirs1[idx] == "" {
				log.Printf("[DEBUG] removing a space from paths: %#v", dirs1)
				continue // remove a space
			}
			dirs2 = append(dirs2, dirs1[idx])
		}
		fmt.Fprintf(c.Stdout, "%s\n", strings.Join(dirs2, "/"))
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}
	return nil
}

func logOutput() io.Writer {
	levels := []logutils.LogLevel{"TRACE", "DEBUG", "INFO", "WARN", "ERROR"}
	minLevel := os.Getenv("LOG_LEVEL")
	if len(minLevel) == 0 {
		minLevel = "INFO" // default log level
	}

	// default log writer is null
	writer := ioutil.Discard
	if minLevel != "" {
		writer = os.Stderr
	}

	filter := &logutils.LevelFilter{
		Levels:   levels,
		MinLevel: logutils.LogLevel(strings.ToUpper(minLevel)),
		Writer:   writer,
	}

	return filter
}
