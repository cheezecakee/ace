package widgets

import "errors"

var ErrInvalidItem = errors.New("invalid item type")

type Action interface {
	Exec() any
}

type Items []Item

type Item struct {
	Label  string
	Action Action // optional
}

type Link struct {
	Target string
}

func (a Link) Exec() any {
	return a.Target
}

func NewLinkItem(label, target string) Item {
	return Item{
		Label:  label,
		Action: Link{Target: target},
	}
}

type Button struct {
	OnPress func() any
}

func (a Button) Exec() any {
	if a.OnPress != nil {
		return a.OnPress()
	}
	return nil
}

func NewButtonItem(label string, fn func() any) Item {
	return Item{
		Label:  label,
		Action: Button{OnPress: fn},
	}
}

type Text struct{}

func (a Text) Exec() any {
	return nil // no-op
}

func NewTextItem(label string) Item {
	return Item{
		Label: label,
	}
}

type Check struct{}

func (a Check) Exec() any {
	return nil // no-op
}

func NewCheckItem(label string) Item {
	return Item{
		Label:  label,
		Action: nil, // or optional toggle action later
	}
}
