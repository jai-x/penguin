package musicserver

import (
	"html/template"
	"net/http"
)

func adminHandler(w http.ResponseWriter, req *http.Request) {
	if ad.ValidSession(ip(req.RemoteAddr)) {
		tmpl, _ := template.ParseFiles("./templates/admin.html")
		tmpl.Execute(w, newPageInfo(req.RemoteAddr))
	} else {
		http.Redirect(w, req, url("/admin/login"), http.StatusFound)
	}
}

func adminLoginHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		tmpl, _ := template.ParseFiles("./templates/admin_login.html")
		tmpl.Execute(w, nil)
		return
	}

	pwd := req.PostFormValue("admin_pwd")
	if ad.ValidPassword(pwd) {
		ad.StartSession(ip(req.RemoteAddr))
		http.Redirect(w, req, url("/admin"), http.StatusSeeOther)
		return
	} else {
		tmpl, _ := template.ParseFiles("./templates/admin_bad_login.html")
		tmpl.Execute(w, nil)
		return
	}
}

func adminLogoutHandler(w http.ResponseWriter, req *http.Request) {
	ad.EndSession(ip(req.RemoteAddr))
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func adminRemoveHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost && ad.ValidSession(ip(req.RemoteAddr)) {
		remUUID := req.PostFormValue("video_id")
		pl.RemoveVideo(remUUID)
	}
	http.Redirect(w, req, url("/admin"), http.StatusSeeOther)
}
