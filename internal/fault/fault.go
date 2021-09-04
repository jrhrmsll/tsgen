package fault

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

const (
	upperBound    = 1.0
	maxRate       = 0.999999
	defaultIndent = "  "
)

func key(path string, code int) string {
	return fmt.Sprintf("%s:%d", path, code)
}

func split(k string) (string, int) {
	parts := strings.Split(k, ":")

	code, _ := strconv.Atoi(parts[1])

	return parts[0], code
}

type kv struct {
	entries map[string]float32
	mu      sync.RWMutex
}

var store = kv{
	entries: map[string]float32{},
}

func (kv *kv) has(k string) bool {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	_, ok := kv.entries[k]
	return ok
}

func (kv *kv) set(k string, v float32) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	if v >= upperBound {
		v = maxRate
	}

	kv.entries[k] = v
}

func (kv *kv) get(k string) float32 {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	if v, ok := kv.entries[k]; ok {
		return v
	}

	return 0
}

func (kv *kv) faults() interface{} {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	type Fault struct {
		Path       string  `json:"path"`
		Code       int     `json:"code"`
		StatusText string  `json:"status_text"`
		Rate       float32 `json:"rate"`
	}

	faults := []Fault{}

	for k, v := range store.entries {
		path, code := split(k)

		faults = append(faults, Fault{
			Path:       path,
			Code:       code,
			StatusText: http.StatusText(code),
			Rate:       v,
		})

	}

	return faults
}
