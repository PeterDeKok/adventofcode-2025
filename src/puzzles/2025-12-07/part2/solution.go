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
	var timelines []int

	for lineNr, line := range input.LineReader(rd) {
		if lineNr == 0 {
			timelines = make([]int, len(line))

			for i, r := range line {
				if r == 'S' {
					timelines[i] = 1
					if build.DEBUG {
						fmt.Printf("%-2d", timelines[i])
					}
				} else if build.DEBUG {
					fmt.Printf(". ")
				}
			}

			if build.DEBUG {
				fmt.Printf("\n")
			}

			continue
		}

		for i, r := range line {
			if r == '^' && timelines[i] > 0 {
				timelines[i-1] += timelines[i]
				timelines[i+1] += timelines[i]
				timelines[i] = 0
				if build.DEBUG {
					fmt.Printf("\033[D\033[D%-2d^ ", timelines[i-1])
				}
			} else if build.DEBUG {
				if timelines[i] > 0 {
					fmt.Printf("%-2d", timelines[i])
				} else {
					fmt.Printf(". ")
				}
			}
		}

		if build.DEBUG {
			fmt.Printf("\n")
		}
	}

	for _, b := range timelines {
		sum += b
	}

	if build.DEBUG {
		fmt.Printf("\n\n")
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
