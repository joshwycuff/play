package model

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Exit                      key.Binding
	Enter                     key.Binding
	Tab                       key.Binding
	NavigateToParent          key.Binding
	NavigateToPreviousSibling key.Binding
	NavigateToNextSibling     key.Binding
	NavigateToLatestChild     key.Binding
	ToggleInputVisibility     key.Binding
	ToggleOutputVisibility    key.Binding
}

func GetDefaultKeyMap() KeyMap {
	return KeyMap{
		Exit:                      key.NewBinding(key.WithKeys("ctrl+c", "esc")),
		Enter:                     key.NewBinding(key.WithKeys("enter")),
		Tab:                       key.NewBinding(key.WithKeys("tab")),
		NavigateToParent:          key.NewBinding(key.WithKeys("alt+h", "alt+,")),
		NavigateToNextSibling:     key.NewBinding(key.WithKeys("alt+j", "down")),
		NavigateToPreviousSibling: key.NewBinding(key.WithKeys("alt+k", "up")),
		NavigateToLatestChild:     key.NewBinding(key.WithKeys("alt+l", "alt+.")),
		ToggleInputVisibility:     key.NewBinding(key.WithKeys("alt+i")),
		ToggleOutputVisibility:    key.NewBinding(key.WithKeys("alt+o")),
	}
}
