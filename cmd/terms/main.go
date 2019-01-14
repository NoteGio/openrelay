package main

import (
	dbModule "github.com/notegio/openrelay/db"
	"github.com/notegio/openrelay/terms"
	"github.com/rs/cors"
	"net/http"
	"log"
	"os"
	"regexp"
	"strconv"
)

type route struct {
    pattern *regexp.Regexp
    handler http.Handler
}

type regexpHandler struct {
    routes []*route
}

func (h *regexpHandler) Handler(pattern *regexp.Regexp, handler http.Handler) {
    h.routes = append(h.routes, &route{pattern, handler})
}

func (h *regexpHandler) HandleFunc(pattern *regexp.Regexp, handler func(http.ResponseWriter, *http.Request)) {
    h.routes = append(h.routes, &route{pattern, http.HandlerFunc(handler)})
}

func (h *regexpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    for _, route := range h.routes {
        if route.pattern.MatchString(r.URL.Path) {
            route.handler.ServeHTTP(w, r)
            return
        }
    }
    // no pattern matched; send 404 response
    http.NotFound(w, r)
}


func main() {
	db, err := dbModule.GetDB(os.Args[1], os.Args[2])
	if err != nil {
		log.Fatalf("Could not open database connection: %v", err.Error())
	}
	port := "8080"
	for _, arg := range os.Args[3:] {
		if _, err := strconv.Atoi(arg); err == nil {
			// If the argument is castable as an integer,
			port = arg
		}
	}
	mux := &regexpHandler{[]*route{}}
	mux.HandleFunc(regexp.MustCompile("^(/[^/]+)?/v2/_tos/"), terms.TermsCheckHandler(db))
	mux.HandleFunc(regexp.MustCompile("^(/[^/]+)?/v2/_tos"), terms.TermsHandler(db))
	mux.HandleFunc(regexp.MustCompile("^/_hc$"), terms.HealthCheckHandler(db))
	corsHandler := cors.Default().Handler(mux)
	log.Printf("ToS Serving on :%v", port)
	http.ListenAndServe(":"+port, corsHandler)
}
