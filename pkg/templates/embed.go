package templates

import (
	"embed"
	"html/template"
)

//go:embed *
var FS embed.FS

var (
	TemplateChat         = template.Must(template.ParseFS(FS, "chat.html", "components/*.html"))
	TemplateChatResponse = template.Must(template.ParseFS(FS, "chat_response.html", "components/chat.html"))
)
