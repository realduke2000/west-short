package wshort

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"math/rand"
	"time"
)

func generateId() string {
	n := time.Now().UTC()
	var buf bytes.Buffer
	// 1 byte for random
	// 1 byte in 1/256 second
	// seconds in every 136 years, 4 bytes
	r := byte(rand.Intn(256))
	binary.Write(&buf, binary.BigEndian, r)
	ms := byte(n.Nanosecond() % 256)
	binary.Write(&buf, binary.BigEndian, ms)
	sInSec := (n.UnixMicro() / 1000000) % 0xFFFFFFFF
	binary.Write(&buf, binary.BigEndian, []byte{
		byte(sInSec >> 24),
		byte(sInSec >> 16),
		byte(sInSec >> 8),
		byte(sInSec),
	})
	bs := buf.Bytes()
	logger.Printf("random=%X, ms=%X, sinSec=%d, unixmicro=%d\n%v\n", r, ms, sInSec, n.UnixMicro(), bs)
	return base64.RawURLEncoding.EncodeToString(bs)
}

func GenerateShortURL(prefix, lurl string) (surl string, err error) {
	sid := generateId()
	surl = prefix + "/s/" + sid
	insert(sid, lurl)
	return surl, nil
}

func GetURL(surlid string) (lurl string, err error) {
	s, ok := query(surlid)
	if !ok {
		return "", errors.New("unregistered url")
	}
	s.LastAccess = time.Now().Unix()
	return s.LongURL, nil
}
