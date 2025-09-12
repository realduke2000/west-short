package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	ws "shortsrv/wshort"

	"github.com/gorilla/mux"
)

func main() {
	genurl := func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		lurl := bytes.NewBuffer(b).String()
		surl, err := ws.GenerateShortURL("http://localhost", lurl)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(surl))
	}
	router := mux.NewRouter()
	router.HandleFunc("/surl", genurl).Methods(http.MethodPut)
	router.HandleFunc("/debug", func(w http.ResponseWriter, r *http.Request) {
		ws.DumpData()
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)
	router.HandleFunc("/s/{surlid}", func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		s := params["surlid"]
		l, err := ws.GetURL(s)
		if err != nil {
			fmt.Printf("redirect short url error with: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Add("Location", l)
		w.WriteHeader(http.StatusMovedPermanently)
		w.Write([]byte("301 Moved"))
	}).Methods(http.MethodGet)

	s := &http.Server{
		Addr:    "127.0.0.1:80",
		Handler: router,
	}
	err := s.ListenAndServe()
	if err != nil {
		fmt.Printf("server started error: %s", err.Error())
		os.Exit(1)
	}
}
