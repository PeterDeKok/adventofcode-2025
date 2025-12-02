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

		if line[0] == 'L' {
			i = -i
		}

		position = (position + i + 100) % 100

		if position == 0 {
			sum++
		}

		if build.DEBUG {
			fmt.Printf("%6s: %d\n", line, position)
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
