package user

import (
	"fmt"
	"github.com/mughub/mughub/db"
	"github.com/mughub/ssr/routes"
	"github.com/mughub/ssr/template"
	"net/http"
	"strings"
)

// GetLoginTmpl handles GET requests for the register page
func GetLoginTmpl(w http.ResponseWriter, req *http.Request) {
	// Get template context data from request
	tmplCtx, _ := routes.GetTmplCtx(req, routes.DL())

	// Execute login page template
	err := template.ExecuteTemplate(w, "user/login", tmplCtx)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
	}
}

// PostLoginTmpl handles POST requests for the register page
func PostLoginTmpl(w http.ResponseWriter, req *http.Request) {
	// Parse form
	err := req.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Extract credentials
	name := req.Form.Get("username") // Note: username is a catch all for username and email
	if name == "" {
		http.Error(w, "username or email must be provided", http.StatusBadRequest)
		return
	}

	passwd := req.Form.Get("password")
	if passwd == "" {
		http.Error(w, "password must be provided", http.StatusBadRequest)
		return
	}

	// Build request
	var reqStr strings.Builder
	reqStr.WriteString("{ login(")

	var key string
	if strings.ContainsRune(name, '@') {
		reqStr.WriteString("email: $email, ")
		key = "$email"
	} else {
		reqStr.WriteString("username: $username, ")
		key = "$username"
	}
	reqStr.WriteString("password: $password) { id } }")

	vars := map[string]interface{}{
		key:         name,
		"$password": passwd,
	}

	// Perform database lookup
	res := db.Do(req.Context(), reqStr.String(), vars)
	if len(res.Errors) > 0 {
		// TODO: Inspect and handle/report errors
		return
	}

	// Return macaron for successful login

}
