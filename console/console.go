package console

import "fmt"

const (
	TextBlack = iota + 30
	TextRed
	TextGreen
	TextYellow
	TextBlue
	TextMagenta
	TextCyan
	TextWhite
)

func ColorfulText(color int, text string) string {
	return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", color, text)
}
