package templatecache

import (
	"html/template"
	"net/http"
)

type TmplCache struct {
	cache *template.Template
}

func NewTemplateCache(dir, domain string) TmplCache {
	funcMap := template.FuncMap {
		"serverDomain": func() string { return domain },
	}
	c := template.Must(template.New("main").Funcs(funcMap).ParseGlob(dir + "/*.html"))
	return TmplCache{c}
}

func (t TmplCache) Render(w http.ResponseWriter, tmpl string, d interface{}) {
	w.Header().Set("Content-type", "text/html")
	err := t.cache.ExecuteTemplate(w, tmpl + ".html", d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
