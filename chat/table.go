package chat

import (
	"sync"
)

// Table represents a mapping of a one-to-many relationship.
type Table struct {
	mu   sync.RWMutex
	data map[interface{}]map[string]struct{}
}

// NewTableInMemory returns a new one-to-many table that persist the data
// in-memory.
func NewTableInMemory() *Table {
	return &Table{
		data: make(map[interface{}]map[string]struct{}),
	}
}

// Get returns slice string of given item.
func (t *Table) Get(key interface{}) []string {
	items := t.get(key)
	if items == nil {
		return []string{}
	}
	result := make([]string, len(items))
	var i int
	for item := range items {
		result[i] = item
		i++
	}
	return result
}

func (t *Table) get(id interface{}) map[string]struct{} {
	t.mu.RLock()
	items, _ := t.data[id]
	t.mu.RUnlock()
	return items
}

func (t *Table) Add(user, session string) {
	// Add the session it the user's session list.
	t.add(UserID(user), session)
	// Map the current session to the user.
	t.add(SessionID(session), user)
}

func (t *Table) add(a interface{}, b string) {
	t.mu.Lock()
	if _, exist := t.data[a]; !exist {
		t.data[a] = make(map[string]struct{})
	}
	t.data[a][b] = struct{}{}
	t.mu.Unlock()
}

// Delete clears the many-to-many relationship, first by clearing the sessions
// the user belong to, and then the sessions the user has.
func (t *Table) Delete(sid SessionID) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if users, exist := t.data[sid]; exist {
		for user := range users {
			delete(t.data[UserID(user)], sid.String())
			if len(t.data[user]) == 0 {
				delete(t.data, user)
			}
		}

	}
	delete(t.data, sid)
	return nil
}
