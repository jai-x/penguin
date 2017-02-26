package musicserver

import (
	"fmt"
	"net/http"
	"html/template"

	"../admin"
)

func adminHandler(w http.ResponseWriter, req *http.Request) {
	token, sessExists := A.GetSession(req.RemoteAddr)

	if !sessExists || admin.Expired(token) {
		// Not logged in or session expired
		http.Redirect(w, req, "/admin/login", http.StatusSeeOther)
	}

	plInfo := Q.GetPlaylistInfo(req.RemoteAddr)
	adminTemplate, _ := template.ParseFiles("templates/admin.html")
	adminTemplate.Execute(w, plInfo)
}

func adminLoginHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		req.ParseForm()
		correct := A.VerifyPassword(req.Form["admin_pwd"][0])

		if !correct {
			http.Redirect(w, req, "/admin/login?info=incorrect", http.StatusSeeOther)
		} else {
			A.StartSession(req.RemoteAddr)
			http.Redirect(w, req, "/admin", http.StatusSeeOther)
		}

	} else {
		loginTemplate, _ := template.ParseFiles("templates/admin_login.html")
		loginTemplate.Execute(w, nil)
	}
}

func adminLogoutHandler(w http.ResponseWriter, req *http.Request) {
	A.EndSession(req.RemoteAddr)
	http.Redirect(w, req, "/home", http.StatusSeeOther)
}

func adminKillHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "hi fam")
}