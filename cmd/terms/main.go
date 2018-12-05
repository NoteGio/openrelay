package main

import (
	dbModule "github.com/notegio/openrelay/db"
	"github.com/notegio/openrelay/terms"
	"github.com/rs/cors"
	"net/http"
	"log"
	"os"
	"strconv"
)


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
	mux := http.NewServeMux()
	mux.HandleFunc("/v2/_tos", terms.TermsHandler(db))
	mux.HandleFunc("/_hc", terms.HealthCheckHandler(db))
	corsHandler := cors.Default().Handler(mux)
	log.Printf("ToS Serving on :%v", port)
	http.ListenAndServe(":"+port, corsHandler)
}
