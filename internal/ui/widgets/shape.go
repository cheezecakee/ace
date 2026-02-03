package widgets

import "fmt"

type ShapeType int

type Shapes []Shape

const (
	Grid ShapeType = iota
	Bar
	Column
)

type Shape struct {
	Type ShapeType
	Cols int
	Rows int
}

type ShapeRule interface {
	Type() ShapeType

	AllowMultiple() bool

	CanAttach(other ShapeType) bool

	Normalize(rows, cols int) (int, int)
}

type GridRule struct{}

func (GridRule) Type() ShapeType { return Grid }

func (GridRule) AllowMultiple() bool { return false }

func (GridRule) CanAttach(other ShapeType) bool { return false }

func (GridRule) Normalize(rows, cols int) (int, int) {
	return rows, cols
}

type BarRule struct{}

func (BarRule) Type() ShapeType { return Bar }

func (BarRule) AllowMultiple() bool { return true }

func (BarRule) CanAttach(other ShapeType) bool {
	return other == Bar
}

func (BarRule) Normalize(rows, cols int) (int, int) {
	return 1, cols
}

type ColumnRule struct{}

func (ColumnRule) Type() ShapeType { return Column }

func (ColumnRule) AllowMultiple() bool { return true }

func (ColumnRule) CanAttach(other ShapeType) bool {
	return other == Column
}

func (ColumnRule) Normalize(rows, cols int) (int, int) {
	return rows, 1
}

func NewShape(rule ShapeRule, rows, cols int) Shape {
	if rows == 0 || cols == 0 {
		fmt.Println("invalid shape size 0")
		return Shape{}
	}
	r, c := rule.Normalize(rows, cols)

	return Shape{
		Type: rule.Type(),
		Rows: r,
		Cols: c,
	}
}
