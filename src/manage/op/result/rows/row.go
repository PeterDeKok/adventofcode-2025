package result

type Row interface {
	Get(col int) Col
	Add(col Col) Row
}

type defaultRow struct {
	cols []Col
}

var _ Row = &defaultRow{}

func (r *defaultRow) Get(col int) Col {
	if col < 0 || col >= len(r.cols) {
		return &defaultCol{}
	}

	return r.cols[col]
}

func (r *defaultRow) Add(col Col) Row {
	r.cols = append(r.cols, col)

	return r
}
