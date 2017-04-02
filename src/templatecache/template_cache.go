package templatecache

import (
	"net/http"
	"html/template"

	"../state"
	"../config"
)

// Pointer to parsed struct of all templates
var templates *template.Template

func Init() {
	folder := config.Config.TemplateFolder
	// Parse all html files in template folder or panic
	templates = template.Must(template.ParseGlob(folder+"/*.html"))
}

func Render(w http.ResponseWriter, tmpl string, p *state.PlaylistInfo) {
	w.Header().Set("Content-type", "text/html")
	// Execute template by name from the parsed struct
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
