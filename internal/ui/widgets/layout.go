package widgets

import "fmt"

type Layout struct {
	Header *Container
	Body   *Container
	Footer *Container
}

type Block struct {
	Items Items // This is []items
	Rows  int
	Cols  int
}

// Connect connects the sections to each other for
// graph navigation creation
func (l *Layout) Connect() ([]Block, error) {
	var blocks []Block

	containers := []*Container{
		l.Header,
		l.Body,
		l.Footer,
	}

	for _, c := range containers {
		if c == nil {
			continue
		}

		rows, cols, err := c.Dimensions()
		if err != nil {
			return nil, err
		}

		blocks = append(blocks, Block{
			Rows: rows,
			Cols: cols,
		})
	}

	if len(blocks) == 0 {
		return nil, fmt.Errorf("layout has no containers")
	}

	return blocks, nil
}

type Container struct {
	Shapes Shapes
}

// Attach attaches shapes inside the Container
// with certain restrictions
func (c *Container) Attach(rule ShapeRule, rows, cols int) error {
	// Rule: only one grid
	if !rule.AllowMultiple() && len(c.Shapes) > 0 {
		return fmt.Errorf("shape %v cannot be combined", rule.Type())
	}

	// Rule: compatible with existing shapes
	for _, s := range c.Shapes {
		if !rule.CanAttach(s.Type) {
			return fmt.Errorf("cannot attach %v to %v", rule.Type(), s.Type)
		}
	}

	c.Shapes = append(c.Shapes, NewShape(rule, rows, cols))
	return nil
}

func (c *Container) Dimensions() (rows int, cols int, err error) {
	if len(c.Shapes) == 0 {
		return 0, 0, fmt.Errorf("container has no shapes")
	}

	switch c.Shapes[0].Type {
	case Grid:
		// grid is guaranteed to be single
		s := c.Shapes[0]
		return s.Rows, s.Cols, nil

	case Bar:
		// bars stack horizontally
		rows = len(c.Shapes)
		cols = c.Shapes[0].Cols
		return rows, cols, nil

	case Column:
		// columns stack vertically
		rows = 0
		for _, s := range c.Shapes {
			rows += s.Rows
		}
		return rows, 1, nil

	default:
		return 0, 0, fmt.Errorf("unknown shape type")
	}
}
