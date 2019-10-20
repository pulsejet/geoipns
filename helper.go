package main

type rowStack []*DatabaseRow

func (s rowStack) Push(v *DatabaseRow) rowStack {
	return append(s, v)
}

func (s rowStack) Peek() *DatabaseRow {
	if len(s) == 0 {
		return nil
	}
	return s[len(s)-1]
}

func (s rowStack) Pop() (rowStack, *DatabaseRow) {
	if len(s) == 0 {
		return s, nil
	}
	l := len(s)
	return s[:l-1], s[l-1]
}
