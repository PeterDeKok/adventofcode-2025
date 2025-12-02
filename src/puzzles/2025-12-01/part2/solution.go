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
)

func Solution(_ context.Context, il *logger.IterationLogger, rd io.Reader, w io.Writer) error {
	var sum int

	position := 50

	for lineNr, line := range input.LineReader(rd) {
		if len(line) < 2 {
			return fmt.Errorf("line %d is too short", lineNr)
		}

		i, err := strconv.Atoi(line[1:])
		if err != nil {
			return fmt.Errorf("line %d; failed to parse number. %s", lineNr, line)
		}

		sum += i / 100
		remainder := i % 100

		// Going Right hits zero AFTER going over the limit
		if line[0] == 'R' && (position+remainder) > 99 {
			// overflow
			sum++
		}

		// Going Left hits zero BEFORE going over the limit
		if line[0] == 'L' {
			if (remainder > position && position > 0) || (remainder == position && remainder > 0) {
				// overflow while not starting at zero or
				// hits 0 exactly AND is not a no-op (nor exactly n full-turns)
				sum++
			}

			remainder = -remainder
		}

		// Adding 100 is simpler then (and equivalent to) modulo using euclidean division.
		position = (position + remainder + 100) % 100

		if build.DEBUG {
			fmt.Printf("%6s: %4d  --  %6d\n", line, position, sum)
		}
	}

	if _, err := w.Write([]byte(strconv.Itoa(sum))); err != nil {
		return err
	}

	return nil
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
