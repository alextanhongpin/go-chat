package server

type Mapper struct {
	data map[string]map[string]struct{}
}

func NewMapper() *Mapper {
	return &Mapper{
		data: make(map[string]map[string]struct{}),
	}
}

func (m *Mapper) Add(a, b string) {
	// a is online, add b as a's peer.
	if peers, found := m.data[a]; !found {
		m.data[a] = make(map[string]struct{})
		m.data[a][b] = struct{}{}
	} else {
		peers[b] = struct{}{}
	}
	// If b is also online, add a as b's peer.
	if peers, found := m.data[b]; found {
		peers[a] = struct{}{}
	}
}

func (m *Mapper) Delete(a string) {
	if peers, found := m.data[a]; found {
		// For each of a's peers, delete a from their list.
		for p := range peers {
			m.DeleteChild(p, a)
		}
		// Set a to nil.
		m.data[a] = nil

		// Delete a.
		delete(m.data, a)
	}
}

func (m *Mapper) DeleteChild(a, b string) {
	if peers, found := m.data[a]; found {
		delete(peers, b)
	}
}

func (m *Mapper) Get(a string) map[string]struct{} {
	return m.data[a]
}

func (m *Mapper) Has(a string) bool {
	_, found := m.data[a]
	return found
}

func (m *Mapper) HasChild(a, b string) bool {
	if _, found := m.data[a]; found {
		_, hasChild := m.data[a][b]
		return hasChild
	}
	return false
}

func (m *Mapper) Data() map[string]map[string]struct{} {
	return m.data
}
