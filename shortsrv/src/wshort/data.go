package wshort

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"path"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Short struct {
	LongURL    string    `json:"long_url"`
	ID         string    `json:"id"`
	Creation   time.Time `json:"creation"`
	LastAccess time.Time `json:"last_access"`
}

type Repository interface {
	Insert(context.Context, Short) error
	Get(context.Context, string) (Short, bool, error)
	Update(context.Context, Short) error
	List(context.Context) ([]Short, error)
	Close() error
}

type EtcdRepository struct {
	client *clientv3.Client
	prefix string
}

var (
	store  Repository
	w      logwriter
	Logger *log.Logger
)

type logwriter struct{}

func (logwriter) Write(p []byte) (n int, err error) {
	fmt.Printf("[%s]%s\n", time.Now().Format("2006-01-02T15:04:05"), string(p))
	return len(p), nil
}

func init() {
	Logger = log.New(w, "", log.Lshortfile)
}

func InitStore(endpoints []string, prefix string) error {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return err
	}
	store = &EtcdRepository{
		client: client,
		prefix: normalizePrefix(prefix),
	}
	return nil
}

func CloseStore() error {
	if store == nil {
		return nil
	}
	err := store.Close()
	store = nil
	return err
}

func insert(s Short) error {
	if store == nil {
		return errors.New("store not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return store.Insert(ctx, s)
}

func query(id string) (Short, bool, error) {
	if store == nil {
		return Short{}, false, errors.New("store not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return store.Get(ctx, id)
}

func update(s Short) error {
	if store == nil {
		return errors.New("store not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return store.Update(ctx, s)
}

func DumpData() {
	if store == nil {
		Logger.Println("store not initialized")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	shorts, err := store.List(ctx)
	if err != nil {
		Logger.Printf("dump error: %v", err)
		return
	}
	for _, v := range shorts {
		Logger.Printf("k=%s, v=%v\n", v.ID, v)
	}
	Logger.Println("Dumped")
}

func NewEtcdRepository(client *clientv3.Client, prefix string) *EtcdRepository {
	return &EtcdRepository{
		client: client,
		prefix: normalizePrefix(prefix),
	}
}

func (r *EtcdRepository) Insert(ctx context.Context, s Short) error {
	body, err := json.Marshal(s)
	if err != nil {
		return err
	}
	key := r.key(s.ID)
	resp, err := r.client.Txn(ctx).
		If(clientv3.Compare(clientv3.Version(key), "=", 0)).
		Then(clientv3.OpPut(key, string(body))).
		Commit()
	if err != nil {
		return err
	}
	if !resp.Succeeded {
		return errors.New("short id already exists")
	}
	return nil
}

func (r *EtcdRepository) Get(ctx context.Context, id string) (Short, bool, error) {
	resp, err := r.client.Get(ctx, r.key(id))
	if err != nil {
		return Short{}, false, err
	}
	if len(resp.Kvs) == 0 {
		return Short{}, false, nil
	}
	var s Short
	if err := json.Unmarshal(resp.Kvs[0].Value, &s); err != nil {
		return Short{}, false, err
	}
	return s, true, nil
}

func (r *EtcdRepository) Update(ctx context.Context, s Short) error {
	body, err := json.Marshal(s)
	if err != nil {
		return err
	}
	_, err = r.client.Put(ctx, r.key(s.ID), string(body))
	return err
}

func (r *EtcdRepository) List(ctx context.Context) ([]Short, error) {
	resp, err := r.client.Get(ctx, r.prefix+"/", clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	shorts := make([]Short, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		var s Short
		if err := json.Unmarshal(kv.Value, &s); err != nil {
			return nil, err
		}
		shorts = append(shorts, s)
	}
	return shorts, nil
}

func (r *EtcdRepository) Close() error {
	if r == nil || r.client == nil {
		return nil
	}
	return r.client.Close()
}

func (r *EtcdRepository) key(id string) string {
	return path.Join(r.prefix, id)
}

func normalizePrefix(prefix string) string {
	trimmed := strings.TrimSpace(prefix)
	trimmed = strings.Trim(trimmed, "/")
	if trimmed == "" {
		return "wshort"
	}
	return trimmed
}
