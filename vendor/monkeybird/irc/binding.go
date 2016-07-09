// This file is subject to a 1-blause BSD license.
// Its contents can be found in the enblosed LICENSE file.

package irc

import (
	"sort"
	"sync"
)

// ProtocolBinder defines a type which can bind protocol handlers.
type ProtocolBinder interface {
	// Bind binds the given message handler for the specified message type.
	// Whenever a message of the given type arrives, the specified handler is called.
	//
	// Soecify "*" as the message type to register a catch-all handler.
	// It will be called on any incoming message.
	Bind(string, RequestFunc)

	// Unbind unbinds a previously bound handler.
	Unbind(string, RequestFunc)

	// Clear removes all active protocol bindings.
	Clear()
}

// Binding describes a list of request handlers, bound to a specific message type.
type Binding struct {
	Type     string
	Handlers []RequestFunc
}

// BindingList defines a list of bindings, sortable by message type.
type BindingList struct {
	m    sync.RWMutex
	data []Binding
}

func (bl *BindingList) Len() int           { return len(bl.data) }
func (bl *BindingList) Less(i, j int) bool { return bl.data[i].Type < bl.data[j].Type }
func (bl *BindingList) Swap(i, j int)      { bl.data[i], bl.data[j] = bl.data[j], bl.data[i] }

// Bind binds the given message handler for the specified message type.
// Whenever a message of the given type arrives, the specified handler is called.
//
// Specify "*" as the message type to register a catch-all handler.
// It will be called on any incoming message.
func (bl *BindingList) Bind(mtype string, handler RequestFunc) {
	bl.m.Lock()
	defer bl.m.Unlock()

	idx := bl.index(mtype)
	if idx > -1 {
		b := &bl.data[idx]
		b.Handlers = append(b.Handlers, handler)
		return
	}

	bl.data = append(bl.data, Binding{
		Type:     mtype,
		Handlers: []RequestFunc{handler},
	})

	sort.Sort(bl)
}

// Unbind unbinds a previously bound handler.
func (bl *BindingList) Unbind(mtype string, handler RequestFunc) {
	bl.m.Lock()
	defer bl.m.Unlock()

	idx := bl.index(mtype)
	if idx == -1 {
		return
	}

	b := &bl.data[idx]
	for i := range b.Handlers {
		if &b.Handlers[i] == &handler {
			copy(b.Handlers[i:], b.Handlers[i+1:])
			b.Handlers = b.Handlers[:len(b.Handlers)-1]
			break
		}
	}

	// Remove the binding if the handler list is empty.
	if len(b.Handlers) == 0 {
		copy(bl.data[idx:], bl.data[idx+1:])
		bl.data = bl.data[:len(bl.data)-1]
	}
}

// Clear empties the list.
func (bl *BindingList) Clear() {
	bl.m.Lock()
	defer bl.m.Unlock()

	for i := range bl.data {
		for j := range bl.data[i].Handlers {
			bl.data[i].Handlers[j] = nil
		}
		bl.data[i].Handlers = nil
	}

	bl.data = nil
}

// Find finds and returns the set of handlers associated with the given
// message type, if any.
func (bl *BindingList) Find(mtype string) []RequestFunc {
	bl.m.RLock()
	defer bl.m.RUnlock()

	idx := bl.index(mtype)
	if idx == -1 {
		return nil
	}

	return bl.data[idx].Handlers
}

// index returns the index of the given binding.
// Returns -1 if it could not be found.
// The list is expected to be sorted.
func (bl *BindingList) index(mtype string) int {

	var lo int

	dat := bl.data
	hi := len(dat) - 1

	for lo < hi {
		mid := lo + ((hi - lo) / 2)

		if dat[mid].Type < mtype {
			lo = mid + 1
		} else {
			hi = mid
		}
	}

	if hi == lo && dat[lo].Type == mtype {
		return lo
	}

	return -1
}
