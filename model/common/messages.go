package common

type ResizeMsg struct {
	Height int
	Width  int
}

func NewResizeMsg(height int, width int) ResizeMsg {
	return ResizeMsg{Height: height, Width: width}
}
