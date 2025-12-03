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

	for _, line := range input.LineReader(rd) {
		joltage := make([]int, 12)

		if build.DEBUG {
			fmt.Printf("bank: %s\n", line)
		}

		for i := 0; i < len(line)-11; i++ {
			force := false

			for j := 0; j < 12; j++ {
				c := int(line[i+j]) - '0'
				if force || c > joltage[j] {
					force = true
					joltage[j] = c
				}
			}
		}

		for i, d := len(joltage)-1, 1; i >= 0; i-- {
			sum += d * joltage[i]
			d *= 10
		}
	}

	if _, err := w.Write([]byte(strconv.Itoa(sum))); err != nil {
		return err
	}

	return nil
}

func main() {
	f, err := os.Open("sample-input-1.txt")
	if err != nil {
		panic(err)
	}

	if err := Solution(context.Background(), logger.CreateIterationLogger(context.Background()), f, os.Stdout); err != nil {
		panic(err)
	}
}
