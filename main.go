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
	version           bool
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
	flag.BoolVar(&cli.version, "v", false, "show a version")
	flag.Parse()
	if err := cli.main(flag.Args()); err != nil {
		fmt.Fprintln(cli.Stderr, err)
		os.Exit(1)
	}
}

func (c *CLI) main(args []string) error {
	if c.version {
		fmt.Fprintf(os.Stdout, "%s %s (%s)\n", appName, version, revision)
		return nil
	}
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
	var r io.Reader = c.Stdin
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		input := scanner.Text()
		if err := c.build(input, nums...); err != nil {
			return err
		}
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}
	return nil
}

func (c *CLI) build(path string, nums ...int) error {
	log.Printf("[DEBUG] build: args: %#v\n", nums)
	if len(nums) == 0 && !c.countDirHierarchy {
		return fmt.Errorf("too few arguments")
	}
	if c.countDirHierarchy {
		dirs := slices.DeleteFunc(strings.Split(path, "/"), func(s string) bool {
			return s == "" || s == "."
		})
		fmt.Println(len(dirs))
		return nil
	}
	var dirs1 []string = strings.Split(path, "/")
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
		var start, end = nums[0], nums[1]
		switch {
		case 0 < end && end < start:
			return fmt.Errorf(
				"invalid range: right side (%d) should be smaller than left side (%d)",
				start, end)
		case end < 0:
			positiveEnd := len(dirs1) + 1 + end
			if positiveEnd < start {
				return fmt.Errorf("index out of range %d: range should be -%d..%d", end, len(dirs1), len(dirs1))
			}
			end = positiveEnd
		case end == 0:
			end = len(dirs1)
		case len(dirs1) < end:
			return fmt.Errorf("index out of range %d: larger than length of given path (%d)", end, len(dirs1))
		}
		var tmp []int
		for i := start; i <= end; i++ {
			tmp = append(tmp, i)
		}
		nums = tmp
	}
	log.Printf("[DEBUG] %#v: nums is built", nums)
	convertToPositives := func() ([]int, error) {
		var indexes []int
		for _, num := range nums {
			switch {
			case num < 0:
				positiveNum := num + len(dirs1) + 1
				log.Printf("[DEBUG] convert %d to %d (count from backward)", num, positiveNum)
				if positiveNum < 0 {
					log.Printf("[ERROR] %d: converted but still negative", positiveNum)
					return []int{}, fmt.Errorf("index out of range %d: range should be -%d..%d", num, len(dirs1), len(dirs1))
				}
				indexes = append(indexes, positiveNum)
			case len(dirs1) < num:
				return []int{}, fmt.Errorf("index out of range %d: larger than length of given path (%d)", num, len(dirs1))
			default:
				indexes = append(indexes, num)
			}
		}
		pairs, err := zipmap(nums, indexes)
		if err == nil {
			for _, pair := range pairs {
				log.Printf("[TRACE] converted %d -> %d\n", pair.a, pair.b)
			}
		}
		return indexes, nil
	}
	nums, err := convertToPositives()
	if err != nil {
		return fmt.Errorf("cannot convert negative numbers: %w", err)
	}
	log.Printf("[DEBUG] %#v: nums is converted with positives", nums)
	for _, num := range nums {
		num -= 1 // convert 1-origin to 0-origin
		if dirs1[num] == "" {
			log.Printf("[DEBUG] removing a space from paths: %#v", dirs1)
			continue
		}
		dirs2 = append(dirs2, dirs1[num])
	}
	fmt.Fprintf(c.Stdout, "%s\n", strings.Join(dirs2, "/"))
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

type intTuple struct {
	a, b int
}

func zipmap(a, b []int) ([]intTuple, error) {
	if len(a) != len(b) {
		return nil, fmt.Errorf("zipmap: arguments must be of same length")
	}

	r := make([]intTuple, len(a), len(a))

	for i, e := range a {
		r[i] = intTuple{e, b[i]}
	}

	return r, nil
}
