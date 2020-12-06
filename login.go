package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/securecookie"
	"gopkg.in/ini.v1"
)

// cookie handling

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func getUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("sessionid"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		}
	}
	return userName
}

func setSession(userName string, pass string, response http.ResponseWriter) {
	value := map[string]string{
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode("sessionus", value); err == nil {
		cookie := &http.Cookie{
			Name:  "sessionus",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}
	valuepw := map[string]string{
		"name": pass,
	}
	if encoded, err := cookieHandler.Encode("sessionpw", valuepw); err == nil {
		cookie := &http.Cookie{
			Name:  "sessionpw",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}
}

func clearSession(response http.ResponseWriter) {
	cookieus := &http.Cookie{
		Name:   "sessionus",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookieus)

	cookiepw := &http.Cookie{
		Name:   "sessionpw",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookiepw)
}

// login handler

func loginHandler(response http.ResponseWriter, request *http.Request) {
	name := request.FormValue("name")
	pass := request.FormValue("password")
	redirectTarget := "/main"
	if name != "" && pass != "" {
		cfgUser, cfgPass := getUserData()

		if cfgUser != "" && cfgPass != "" && cfgUser == name && cfgPass == pass {
			setSession(name, pass, response)
			redirectTarget = "/"
		}
	}
	http.Redirect(response, request, redirectTarget, 302)
}

// logout handler

func logoutHandler(response http.ResponseWriter, request *http.Request) {
	clearSession(response)
	http.Redirect(response, request, "/", 302)
}

// index page

const indexPage = `
<h1>Login</h1>
<form method="post" action="/login">
    <label for="name">User name</label>
    <input type="text" id="name" name="name">
    <label for="password">Password</label>
    <input type="password" id="password" name="password">
    <button type="submit">Login</button>
</form>
`

func indexPageHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, indexPage)
}

func loginCheck(w http.ResponseWriter, r *http.Request) bool {

	var name string
	var pass string

	if cookie, err := r.Cookie("sessionus"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("sessionus", cookie.Value, &cookieValue); err == nil {
			name = cookieValue["name"]
		}
	}

	if cookie, err := r.Cookie("sessionpw"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("sessionpw", cookie.Value, &cookieValue); err == nil {
			pass = cookieValue["name"]
		}
	}

	cfgUser, cfgPass := getUserData()
	if name != cfgUser || pass != cfgPass {
		http.Redirect(w, r, "/main", 302)
		return false
	}
	return true
}

func getUserData() (string, string) {
	cfg, err := ini.Load("ImageCloud.conf")
	if err != nil {
		return "", ""
	}
	confUsername := cfg.Section("account").Key("username").String()
	confPasswd := cfg.Section("account").Key("passwd").String()

	return confUsername, confPasswd
}
