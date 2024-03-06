package templates

import (
	"embed"
	"html/template"
)

//go:embed *
var FS embed.FS

var (
	TemplateLogin        = template.Must(template.ParseFS(FS, "login.html", "components/*.html"))
	TemplateChat         = template.Must(template.ParseFS(FS, "chat.html", "components/*.html"))
	TemplateChatInput    = template.Must(template.ParseFS(FS, "chat_input.html", "components/chat_input.html"))
	TemplateChatMessages = template.Must(template.ParseFS(FS, "chat_messages.html", "components/chat_messages.html"))
)
