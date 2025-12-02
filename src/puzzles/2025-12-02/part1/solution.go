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

	if len(str)%2 != 0 {
		if build.DEBUG {
			fmt.Println("    > NOT 2n chars; DONT count")
		}
		return false
	}

	lh := len(str) / 2

	if str[:lh] == str[lh:] {
		if build.DEBUG {
			fmt.Println("    > sequence; count")
		}

		return true
	}

	if build.DEBUG {
		fmt.Println("    > NOT sequence; DONT count it")
	}

	return false
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
