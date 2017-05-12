package server

import (
	"net/http"
	"html/template"

	"../admin"
	"../playlist"
)

func adminHandler(w http.ResponseWriter, req *http.Request) {
	if admin.ValidSession(req.RemoteAddr) {
		tmpl, _ := template.ParseFiles("templates/admin.html")
		tmpl.Execute(w, fetchInfo(req.RemoteAddr))
	} else {
		http.Redirect(w, req, "/admin/login", http.StatusFound)
	}
}

func adminLoginHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		tmpl, _ := template.ParseFiles("templates/admin_login.html")
		tmpl.Execute(w, nil)
		return
	}

	pwd := req.PostFormValue("admin_pwd")
	if admin.VerifyPassword(pwd) {
		admin.StartSession(req.RemoteAddr)
		http.Redirect(w, req, "/admin", http.StatusSeeOther)
		return
	} else {
		tmpl, _ := template.ParseFiles("templates/admin_bad_login.html")
		tmpl.Execute(w, nil)
		return
	}
}

func adminLogoutHandler(w http.ResponseWriter, req *http.Request) {
	admin.EndSession(req.RemoteAddr)
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func adminRemoveHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost && admin.ValidSession(req.RemoteAddr) {
		remUUID := req.PostFormValue("video_id")
		playlist.RemoveVideo(remUUID)
	}
	http.Redirect(w, req, "/admin", http.StatusSeeOther)
}
