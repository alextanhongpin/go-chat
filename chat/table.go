package chat

import "sync"

type Table struct {
	mu   sync.RWMutex
	data map[interface{}]map[interface{}]struct{}
}

func NewTableInMemory() *Table {
	return &Table{
		data: make(map[interface{}]map[interface{}]struct{}),
	}
}

// If room id is provided, return an array of users in the room,
// if user id is provided, return an array of rooms of the user.

// Get returns slice string of given item.
func (t *Table) Get(id interface{}) []string {
	t.mu.RLock()
	items := t.data[id]
	t.mu.RUnlock()

	result := make([]string, len(items))
	var i int
	for item := range items {
		result[i] = item.(string)
		i++
	}
	return result
}

func (t *Table) Add(a, b interface{}) error {
	// Add the user in the room.
	t.add(a, b)
	// Keep track of the rooms the user is in.
	t.add(b, a)
	return nil
}

func (t *Table) add(a, b interface{}) {
	t.mu.Lock()
	if _, exist := t.data[a]; !exist {
		t.data[a] = make(map[interface{}]struct{})
	}
	t.data[a][b] = struct{}{}
	t.mu.Unlock()
}

// To delete a, delete instance of a in every b first. If b is empty, delete b.
// Then delete all bs in a.
func (t *Table) Delete(a interface{}, fn func(string)) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	// bs == rooms
	if bs, exist := t.data[a]; exist {
		// Remove the user from each room.
		for b := range bs {
			delete(t.data[b], a)
			fn(b.(string))
			// Room is empty, delete.
			if len(t.data[b]) == 0 {
				delete(t.data, b)
			}
		}

	}
	// Delete all user rooms.
	delete(t.data, a)
	return nil
}
