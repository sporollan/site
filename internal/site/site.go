package site

type Site struct {
	InputDir  string
	OutputDir string
	StaticDir string
}

type Page struct {
	Path  string
	Title string
	Body  string
}

func New(input string, output string, static string) *Site {
	return &Site{
		InputDir:  input,
		OutputDir: output,
		StaticDir: static,
	}
}
