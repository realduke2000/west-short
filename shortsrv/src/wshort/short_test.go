package wshort

import (
	"context"
	"encoding/base64"
	"errors"
	"testing"
	"time"
)

type fakeRepository struct {
	data map[string]Short
}

func newFakeRepository() *fakeRepository {
	return &fakeRepository{data: make(map[string]Short)}
}

func (f *fakeRepository) Insert(_ context.Context, s Short) error {
	if _, ok := f.data[s.ID]; ok {
		return errors.New("duplicate")
	}
	f.data[s.ID] = s
	return nil
}

func (f *fakeRepository) Get(_ context.Context, id string) (Short, bool, error) {
	s, ok := f.data[id]
	return s, ok, nil
}

func (f *fakeRepository) Update(_ context.Context, s Short) error {
	f.data[s.ID] = s
	return nil
}

func (f *fakeRepository) List(_ context.Context) ([]Short, error) {
	shorts := make([]Short, 0, len(f.data))
	for _, s := range f.data {
		shorts = append(shorts, s)
	}
	return shorts, nil
}

func (f *fakeRepository) Close() error {
	return nil
}

func TestGenerateIDProducesURLSafeToken(t *testing.T) {
	id := generateId()
	if len(id) != 8 {
		t.Fatalf("unexpected id length: %d", len(id))
	}
	if _, err := base64.RawURLEncoding.DecodeString(id); err != nil {
		t.Fatalf("expected base64url token, got error: %v", err)
	}
}

func TestCreateAndGetShort(t *testing.T) {
	prevStore := store
	store = newFakeRepository()
	defer func() { store = prevStore }()

	short, err := CreateShort("https://example.com")
	if err != nil {
		t.Fatalf("CreateShort returned error: %v", err)
	}
	if short.ID == "" {
		t.Fatal("expected generated short id")
	}

	got, err := GetShort(short.ID)
	if err != nil {
		t.Fatalf("GetShort returned error: %v", err)
	}
	if got.LongURL != "https://example.com" {
		t.Fatalf("unexpected long url: %s", got.LongURL)
	}
}

func TestUpdateShortAccessPersistsTimestamp(t *testing.T) {
	prevStore := store
	store = newFakeRepository()
	defer func() { store = prevStore }()

	short := Short{
		ID:         "abc12345",
		LongURL:    "https://example.com",
		Creation:   time.Now().Add(-time.Hour),
		LastAccess: time.Now().Add(-time.Hour),
	}
	if err := insert(short); err != nil {
		t.Fatalf("insert returned error: %v", err)
	}

	before := short.LastAccess
	if err := UpdateShortAccess(short); err != nil {
		t.Fatalf("UpdateShortAccess returned error: %v", err)
	}
	updated, _, err := query(short.ID)
	if err != nil {
		t.Fatalf("query returned error: %v", err)
	}
	if !updated.LastAccess.After(before) {
		t.Fatalf("expected LastAccess to advance, before=%v after=%v", before, updated.LastAccess)
	}
}

func TestNormalizePrefix(t *testing.T) {
	if got := normalizePrefix(" /custom/prefix/ "); got != "custom/prefix" {
		t.Fatalf("unexpected prefix: %q", got)
	}
	if got := normalizePrefix(" "); got != "wshort" {
		t.Fatalf("expected default prefix, got %q", got)
	}
}
