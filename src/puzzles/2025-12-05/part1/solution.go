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

func Solution(_ context.Context, _ *logger.IterationLogger, rd io.Reader, w io.Writer) error {
	var sum int

	ranges := make([]func(id int) bool, 0, 1024)
	procLine := procRangeLine

	for lineNr, line := range input.LineReader(rd) {
		if len(line) == 0 {
			procLine = procIngredientLine
			continue
		}

		if err := procLine(lineNr, line, &ranges, &sum); err != nil {
			return err
		}
	}

	if _, err := w.Write([]byte(strconv.Itoa(sum))); err != nil {
		return err
	}

	return nil
}

func procRangeLine(lineNr int, line string, ranges *[]func(id int) bool, _ *int) error {
	// Possible optimisation, consolidate overlapping ranges
	var l, r int

	_, err := fmt.Sscanf(line, "%d-%d", &l, &r)
	if err != nil {
		return err
	}

	*ranges = append(*ranges, func(id int) bool {
		return id >= l && id <= r
	})

	if build.DEBUG {
		fmt.Printf("range (line %d): %d - %d\n", lineNr, l, r)
	}

	return nil
}

func procIngredientLine(lineNr int, line string, ranges *[]func(id int) bool, sum *int) error {
	if build.DEBUG {
		fmt.Printf("ingredient (line %d): %s\n", lineNr, line)
	}

	id, err := strconv.Atoi(line)
	if err != nil {
		return err
	}

	for rangeIndex, rangeFn := range *ranges {
		if rangeFn(id) {
			*sum++
			if build.DEBUG {
				fmt.Printf("  > range %d: FRESH\n", rangeIndex)
			}

			return nil
		} else {
			if build.DEBUG {
				fmt.Printf("  > range %d: next\n", rangeIndex)
			}
		}
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
