package wshort

import (
	"fmt"
	"log"
	"time"
)

type shorturl struct {
	LongURL    string
	ShortURLID string
	Creation   int64
	LastAccess int64
}

var (
	db     map[string]shorturl
	w      logwriter
	logger *log.Logger
)

type logwriter struct {
}

func (logwriter) Write(p []byte) (n int, err error) {
	fmt.Printf("[%s]%s\n", time.Now().Format("2006-01-02T15:04:05"), string(p))
	return len(p), nil
}

func init() {
	db = make(map[string]shorturl)
	logger = log.New(w, "", log.Lshortfile)
}

func query(surlid string) (s *shorturl, ok bool) {
	_s, ok := db[surlid]
	return &_s, ok
}

func insert(id, lurl string) error {
	db[id] = shorturl{
		LongURL:    lurl,
		ShortURLID: id,
		Creation:   time.Now().Unix(),
		LastAccess: -1,
	}
	return nil
}

func DumpData() {
	for k, v := range db {
		logger.Printf("k=%s, v=%v\n", k, v)
	}
	logger.Println("Dumped")
}
