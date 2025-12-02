package result

type Rows interface {
	Col(row, col int) Col
	AddRow(row Row) Rows
}

type defaultRows[R Row] struct {
	rows []Row
}

var _ Rows = &defaultRows[Row]{}

func (r *defaultRows[R]) Col(row, col int) Col {
	if row < 0 || row >= len(r.rows) {
		return &defaultCol{}
	}

	return r.rows[row].Get(col)
}

func (r *defaultRows[R]) AddRow(row Row) Rows {
	r.rows = append(r.rows, row)

	return r
}
