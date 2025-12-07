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
	var sum int64
	var cols []int32

	for lineNr, line := range input.LineReader(rd) {
		runes := []rune(line)

		if runes[0] == '+' || runes[0] == '*' {
			if build.DEBUG {
				fmt.Printf("cols: %v\n", cols)
			}

			op := &Op{r: runes[0], sum: &sum}

			for i, r := range runes {
				if r != ' ' {
					op.NewOp(r)
				}

				op.RunValue(cols[i])
			}

			op.Finish()

			break
		}

		if lineNr == 0 {
			cols = make([]int32, len(runes))
		}

		for i, r := range runes {
			if r == ' ' {
				continue
			}

			cols[i] *= 10
			cols[i] += r - '0'
		}
	}

	if _, err := w.Write([]byte(strconv.Itoa(int(sum)))); err != nil {
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

type Op struct {
	r   rune
	sum *int64
	val int64
}

func (op *Op) NewOp(r rune) {
	if build.DEBUG {
		fmt.Printf("New op: %s | Previous value: %d\n", string(r), op.val)
	}

	*op.sum += op.val
	if r == '*' {
		op.val = 1
	} else {
		op.val = 0
	}
	op.r = r
}

func (op *Op) Finish() {
	if build.DEBUG {
		fmt.Printf("Finish: %d\n", op.val)
	}

	*op.sum += op.val
}

func (op *Op) RunValue(v int32) {
	if v == 0 {
		return
	}

	switch op.r {
	case '+':
		op.val += int64(v)
	case '*':
		op.val *= int64(v)
	}
}
