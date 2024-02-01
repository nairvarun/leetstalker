package color

import "fmt"

const (
	NONE = iota
	RED
	GREEN
	YELLOW
	BLUE
	PURPLE
	BOLD
	ITALIC
)

const escape = "\x1b"

func color(c int) string {
	switch c {
	case NONE:
		return fmt.Sprintf("%s[%dm", escape, c)
	case BOLD:
		return fmt.Sprintf("%s[1m", escape)
	case ITALIC:
		return fmt.Sprintf("%s[3m", escape)
	default:
		return fmt.Sprintf("%s[3%dm", escape, c)
	}
}

func Format(c int, text string) string {
	return color(c) + text + color(NONE)
}
