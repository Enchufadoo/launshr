package app_style

type Color struct {
	Light string
	Dark  string
}

func BackColor() Color {
	return Color{Light: "#D9DCCF", Dark: "#383838"}
}

func SecondaryColor() Color {
	return Color{Light: "#D9DCCF", Dark: "#c7c7c7"}
}
