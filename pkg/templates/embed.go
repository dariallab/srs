package templates

import (
	"embed"
	"html/template"
)

//go:embed *
var FS embed.FS

var (
	TemplateChat        = template.Must(template.ParseFS(FS, "chat.html", "components/*.html"))
	TemplateChatInput   = template.Must(template.ParseFS(FS, "chat_input.html", "components/chat.html"))
	TemplateChatMessage = template.Must(template.ParseFS(FS, "chat_message.html", "components/chat.html"))
)
