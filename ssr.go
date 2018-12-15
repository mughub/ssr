// Package ssr contains a Server Side Rendered (SSR) GoHub UI implementation.
package ssr

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mughub/mughub/bare"
	"github.com/mughub/ssr/routes"
	"github.com/mughub/ssr/routes/user"
	"github.com/mughub/ssr/template"
	"github.com/spf13/viper"
	"net/http"
)

// UI represents a server side rendered website for GoHub.
type UI struct{}

// Init initializes the UI to the provided Router.
// The provided is assumed to be a sub viper rooted at gohub.ui.
//
func (_ *UI) Init(base bare.Router, cfg *viper.Viper) {
	domain := cfg.GetString("domain")
	if domain == "" {
		panic("gohub.ui.domain must be specified")
	}

	// Parse UI templates
	template.ParseTmpls(cfg)

	// Set domain
	base.Host(domain)

	// Handle assests
	assestDir := cfg.GetString("assests")
	base.PathPrefix("/assests/").Handler(http.StripPrefix("/assests/", http.FileServer(http.Dir(assestDir))))

	// TODO: Add Middlewares
	// TODO: Middlewares to add - Auth, CSRF
	// TODO: Optional middlewares - Logging, metrics,

	/***** Base Routes: START ******/
	base.Path("/").Methods("GET").HandlerFunc(routes.GetHome)
	/***** Base Routes: END ******/

	/***** User Routes: START ******/
	usr := base.PathPrefix("/{user}/").Subrouter()

	// Handle user base
	usr.Path("/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		fmt.Fprintf(w, "Welcome %s to your GoHub page!", vars["user"])
	})

	// TODO: Handle register and login
	usr.Path("/register").Methods("GET", "POST")

	usr.Path("/login").Methods("GET", "POST").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodGet {
			user.GetLoginTmpl(w, req)
			return
		}
		user.PostLoginTmpl(w, req)
	})

	/***** User Routes: END ******/

	/***** Repo Routes: START ******/
	repo := usr.PathPrefix("/{repo}/").Subrouter()

	repo.Path("/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		fmt.Fprintf(w, "Welcome %s to your repo: %s", vars["user"], vars["repo"])
	})

	repo.PathPrefix("/tree/{branch}")
	/***** Repo Routes: END ******/
}
