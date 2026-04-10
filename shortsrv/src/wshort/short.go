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
	Logger.Printf("random=%X, ms=%X, sinSec=%d, unixmicro=%d, %v", r, ms, sInSec, n.UnixMicro(), bs)
	return base64.RawURLEncoding.EncodeToString(bs)
}

func CreateShort(lurl string) (short Short, err error) {
	short.Creation = time.Now()
	short.LongURL = lurl
	short.ID = generateId()
	short.LastAccess = time.Now()
	err = insert(short)
	for i := 0; i < 5 && err != nil; i++ {
		short.ID = generateId()
		err = insert(short)
	}
	return short, err
}

func GetShort(id string) (s Short, err error) {
	s, ok, err := query(id)
	if err != nil {
		return s, err
	}
	if !ok {
		return s, errors.New("unregistered url")
	}
	return s, err
}

func UpdateShortAccess(s Short) (err error) {
	s.LastAccess = time.Now()
	return update(s)
}
