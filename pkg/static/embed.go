package static

import (
	"embed"
	"html/template"
)

//go:embed *
var FS embed.FS

var (
	TemplateChat = template.Must(template.ParseFS(FS, "chat.gohtml"))
)
