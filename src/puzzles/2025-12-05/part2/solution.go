package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/build"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/input"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/logger"
	"sort"
	"strconv"
)

func Solution(_ context.Context, _ *logger.IterationLogger, rd io.Reader, w io.Writer) error {
	var sum int

	ids := make([]*struct{ l, r int }, 0, 1024)

	for lineNr, line := range input.LineReader(rd) {
		if len(line) == 0 {
			break
		}

		var l, r int

		_, err := fmt.Sscanf(line, "%d-%d", &l, &r)
		if err != nil {
			return err
		}

		if build.DEBUG {
			fmt.Printf("range (line %d): %d - %d\n", lineNr, l, r)
		}

		shouldAppend := true

		var i int
		for ; i < len(ids); i++ {
			pair := ids[i]

			// -1 & +1, so adjacency is also seen as overlap
			if r < (pair.l - 1) {
				// ids is sorted, r stays the same, pair.l is strictly increasing.
				// This condition will always be true from now on.
				i = len(ids) // Also skip the 'next' overlap check
				break
			} else if l > (pair.r + 1) {
				continue
			}

			if build.DEBUG {
				fmt.Printf(" > overlaps with range %d (%d - %d)\n", i, pair.l, pair.r)
			}

			shouldAppend = false

			if l >= pair.l && r <= pair.r {
				// Fully overlapping; skip checking for other overlaps
				i = len(ids)
				break
			}

			var diff int

			if l < pair.l {
				diff += pair.l - l
				pair.l = l
			}

			if r > pair.r {
				diff += r - pair.r
				pair.r = r
			}

			if build.DEBUG {
				fmt.Printf("   > + %d = %d\n", diff, sum)
			}

			sum += diff

			// Problem is, there could be more overlaps...
			// Good thing is, its only with the ids AFTER this one.
			// Bad thing, we can't continue with the 'remainder' only, as it might be on 2 sides.
			// So we continue iterating in the next for-loop, reducing overlaps we find, nil-ing zero length remainders.
			// However, it still leaves adjacent pairs if any were reduced, but not removed.
			break
		}

		for i++; i < len(ids); i++ {
			pair := ids[i]

			if pair == nil {
				// Previously removed; need to investigate if actually removing (re-slicing & copy) is better?
				continue
			}

			if l > pair.r {
				// No overlap
				continue
			} else if r < pair.l {
				// ids is sorted, r stays the same, pair.l is strictly increasing.
				// This condition will always be true from now on.
				break
			}

			if build.DEBUG {
				fmt.Printf(" > overlaps again with range %d (%d - %d)\n", i, pair.l, pair.r)
			}

			// Note; fully contained overlap (l > pair.l && r < pair.r) is not possible.
			// This loop would not trigger if there was no other overlap before
			// and any other overlap would then automatically overlap this one and been merged or negated.
			// An exact overlap is still possible though.

			pre := pair.r - pair.l

			if l <= pair.l {
				pair.l = r + 1
			} else if r >= pair.r {
				pair.r = l - 1
			}

			if pair.l > pair.r {
				// And... invalid, remove
				copy(ids[i:], ids[i+1:])
				ids = ids[:len(ids)-1]
				i--
				sum -= pre + 1
			} else {
				sum -= pre - pair.r + pair.l
			}
		}

		if shouldAppend {
			sum += r - l + 1

			index := sort.Search(len(ids), func(j int) bool { return ids[j].l >= l && (ids[j].l != l || ids[j].r >= r) })
			ids = append(ids, nil)
			copy(ids[index+1:], ids[index:])
			ids[index] = &struct{ l, r int }{l, r}
		}
	}

	if build.DEBUG {
		fmt.Printf("ranges:\n")
		for i, pair := range ids {
			if pair == nil {
				fmt.Printf(" %3d: <nil>\n", i)
			} else {
				fmt.Printf(" %3d: %20d - %-20d\n", i, pair.l, pair.r)
			}
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
