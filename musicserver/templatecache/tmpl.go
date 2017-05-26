package templatecache

import (
	"html/template"
	"net/http"
)

var (
	domain  string

	funcMap = template.FuncMap {
		"serverDomain": func() string { return domain },
	}
)

func SetDomain(dom string) {
	domain = dom
}

type TmplCache struct {
	cache    *template.Template
	dir      string
	useCache bool
}

func NewTemplateCache(dir string, useCache bool) TmplCache {
	c := template.Must(template.New("").Funcs(funcMap).ParseGlob(dir + "/*.html"))
	return TmplCache{c, dir, useCache}
}

func (t TmplCache) Render(w http.ResponseWriter, name string, d interface{}) {
	w.Header().Set("Content-type", "text/html")

	// Optionally bypass cache and reparse template at every request
	// Good for debugging
	if !t.useCache {
		t.noCacheRender(w, name, d)
		return
	}

	err := t.cache.ExecuteTemplate(w, name + ".html", d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (t TmplCache) noCacheRender(w http.ResponseWriter, name string, d interface{}) {
	out, err := template.New("").Funcs(funcMap).ParseFiles(t.dir + "/" + name + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = out.ExecuteTemplate(w, name + ".html", d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

