package chat

import (
	"sync"
)

type Table struct {
	mu   sync.RWMutex
	data map[interface{}]map[string]struct{}
}

func NewTableInMemory() *Table {
	return &Table{
		data: make(map[interface{}]map[string]struct{}),
	}
}

// Get returns slice string of given item.
func (t *Table) GetSessions(user string) []string {
	items := t.get(UserID(user))
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

func (t *Table) GetUsers(room string) []string {
	items := t.get(SessionID(room))
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

func (t *Table) Add(user, room string) {
	// Add the user in the room.
	t.add(UserID(user), room)
	// Keep track of the rooms the user is in.
	t.add(SessionID(room), user)
}

func (t *Table) add(a interface{}, b string) {
	t.mu.Lock()
	if _, exist := t.data[a]; !exist {
		t.data[a] = make(map[string]struct{})
	}
	t.data[a][b] = struct{}{}
	t.mu.Unlock()
}

// To delete a, delete instance of a in every b first. If b is empty, delete b.
// Then delete all bs in a.
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
