package metrics

import "sync"

type Registry struct {
	mu     sync.RWMutex
	values map[string]int64
}

func New() *Registry                             { return &Registry{values: make(map[string]int64)} }
func (r *Registry) Inc(name string)              { r.mu.Lock(); r.values[name]++; r.mu.Unlock() }
func (r *Registry) Add(name string, value int64) { r.mu.Lock(); r.values[name] += value; r.mu.Unlock() }
func (r *Registry) Value(name string) int64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.values[name]
}
func (r *Registry) Snapshot() map[string]int64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[string]int64, len(r.values))
	for key, value := range r.values {
		out[key] = value
	}
	return out
}
