package main

import (
	"github.com/ElliottCepin/go-sqlite-pastebin/internal/server"
	"github.com/ElliottCepin/go-sqlite-pastebin/internal/store"
	"log/slog"
	"net/http"
	"flag"
	"errors"
)
var storeType string
var storeLocation string
func init() {
	flag.StringVar(&storeType, "STORE", "memory", "Choose between in-memory storage and a persistent sqlite store")
	flag.StringVar(&storeLocation, "SQLITE_PATH", "default.db", "Select the path for your sqlite db")
	flag.Parse()

}
func main() {
	var st server.Store
	var err error
	if (storeType == "memory") {
		st = store.NewMemoryStore()
	} else if (storeType == "sqlite") {
		st, err = store.NewSQLiteStore(storeLocation)
		if (err != nil) {
			panic(err)
		}
	} else {
		panic(errors.New("could not identify storage type"))
	}
	emp := server.NewServer(st, slog.Default())
	s := &emp
	srv := &http.Server{Addr: ":8080", Handler:s.Routes()}
	srv.ListenAndServe()
}
