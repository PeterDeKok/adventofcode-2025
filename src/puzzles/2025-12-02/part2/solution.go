package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/build"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/input"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/logger"
	"strconv"
	"strings"
)

func Solution(_ context.Context, _ *logger.IterationLogger, rd io.Reader, w io.Writer) error {
	var sum int

	for lineNr, line := range input.LineReader(rd) {
		for i, pair := range strings.Split(line, ",") {
			if build.DEBUG {
				fmt.Printf("%s\n", pair)
			}

			if len(pair) == 0 {
				continue
			}

			var l, r int

			_, err := fmt.Sscanf(pair, "%d-%d", &l, &r)
			if err != nil {
				return fmt.Errorf("failed to parse line %d, pair %d: %w", lineNr, i, err)
			}

			for id := l; id <= r; id++ {
				if build.DEBUG {
					fmt.Printf("  %d\n", id)
				}

				if procId(id) {
					sum += id
				}
			}
		}
	}

	if _, err := w.Write([]byte(strconv.Itoa(sum))); err != nil {
		return err
	}

	return nil
}

func procId(id int) bool {
	str := strconv.Itoa(id)

	// 11111111 12121212 123123123 12341234
	for l := 1; l+l <= len(str); l++ {
		if len(str)%l != 0 {
			if build.DEBUG {
				fmt.Printf("    > partlength %d not modulo fit\n", l)
			}

			continue
		}

		if build.DEBUG {
			fmt.Printf("    > partlength %d\n", l)
		}

		if checkSequences(str, l) {
			return true
		}
	}

	return false
}

func checkSequences(str string, partLen int) bool {
	for i := 0; i < len(str)-partLen-partLen+1; i += partLen {
		if str[i:i+partLen] != str[i+partLen:i+2*partLen] {
			if build.DEBUG {
				fmt.Printf("      > NOT sequence; DONT count it\n")
			}

			return false
		}
	}

	if build.DEBUG {
		fmt.Printf("        > sequence; count\n")
	}

	return true
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}

	if err := Solution(context.Background(), logger.CreateIterationLogger(context.Background()), f, os.Stdout); err != nil {
		panic(err)
	}
}
