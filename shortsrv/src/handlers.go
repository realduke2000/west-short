package main

import (
	"encoding/base64"
	"net/http"
	ws "shortsrv/wshort"

	"github.com/gorilla/mux"
)

// ShortUrlPutHandler handles PUT /s.
// URL: /s?urlinbase64=<base64url-encoded-long-url>
// Parameters: query parameter "urlinbase64", encoded with base64.RawURLEncoding.
// Returns: 200 with the generated short ID as plain text.
// Errors: 400 "invalid url" when the query value cannot be decoded, 500 "invalid url" when storage fails.
// Behavior: decodes the long URL, stores it in etcd, and returns the generated short ID.
func ShortUrlPutHandler(w http.ResponseWriter, r *http.Request) {
	urlb64 := r.URL.Query().Get(URLInBase64Param)
	lurl, err := base64.RawURLEncoding.DecodeString(urlb64)
	if err != nil {
		ws.Logger.Printf("Invalid URL %s, err: %s", urlb64, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid url"))
		return
	}

	s, err := ws.CreateShort(string(lurl))
	if err != nil {
		ws.Logger.Printf("Failed to create short due to: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("invalid url"))
		return
	}
	w.Write([]byte(s.ID))
}

// ShortUrlGetHandler handles GET /s/{surlid}.
// URL: /s/{surlid}
// Parameters: path parameter "surlid", the short ID returned by PUT /s.
// Returns: 301 with a Location header pointing to the original URL and body "301 Moved".
// Errors: 500 when the short ID cannot be loaded from etcd.
// Behavior: fetches the stored long URL, updates LastAccess in etcd, and redirects the client.
func ShortUrlGetHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	sid := params[ShortIDParam]
	short, err := ws.GetShort(sid)
	if err != nil {
		ws.Logger.Printf("redirect short url (%s) error: %v", sid, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := ws.UpdateShortAccess(short); err != nil {
		ws.Logger.Printf("update short access (%s) error: %v", sid, err)
	}
	w.Header().Add("Location", short.LongURL)
	w.WriteHeader(http.StatusMovedPermanently)
	w.Write([]byte("301 Moved"))
}

// DebugHandler handles GET /debug.
// URL: /debug
// Parameters: none.
// Returns: 200 with an empty body.
// Behavior: dumps all stored short-link records to the server log for inspection.
func DebugHandler(w http.ResponseWriter, r *http.Request) {
	ws.DumpData()
	w.WriteHeader(http.StatusOK)
}
