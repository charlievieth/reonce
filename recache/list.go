package recache

// list represents a doubly linked list.
type list struct {
	root entry // sentinel list Element, only &root, root.prev, and root.next are used
	len  int   // current list length excluding (this) sentinel Element
}

// Init initializes or clears list l.
func (l *list) Init() *list {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	return l
}

func newList() *list { return new(list).Init() }

// Len returns the number of Elements of list l.
// The complexity is O(1).
func (l *list) Len() int { return l.len }

// Back returns the last Element of list l or nil if the list is empty.
func (l *list) Back() *entry {
	if l.len != 0 {
		return l.root.prev
	}
	return nil
}

// insert inserts e after at, increments l.len, and returns e.
func (l *list) insert(e, at *entry) *entry {
	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
	l.len++
	return e
}

// insertValue is a convenience wrapper for insert(&Element{Value: v}, at).
func (l *list) insertValue(v *entry, at *entry) *entry {
	return l.insert(v, at)
}

// remove removes e from its list, decrements l.len, and returns e.
func (l *list) remove(e *entry) {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil // avoid memory leaks
	e.prev = nil // avoid memory leaks
	l.len--
}

// move moves e to next to at and returns e.
func (l *list) move(e, at *entry) {
	if e == at {
		return
	}
	e.prev.next = e.next
	e.next.prev = e.prev

	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
}

// Remove removes e from l if e is an Element of list l.
// It returns the Element value e.Value.
// The Element must not be nil.
func (l *list) Remove(e *entry) {
	l.remove(e)
}

// PushFront inserts a new Element e with value v at the front of list l and returns e.
func (l *list) PushFront(v *entry) *entry {
	return l.insertValue(v, &l.root)
}

// MoveToFront moves Element e to the front of list l.
// If e is not an Element of l, the list is not modified.
// The Element must not be nil.
func (l *list) MoveToFront(e *entry) {
	if l.root.next == e {
		return
	}
	l.move(e, &l.root)
}
