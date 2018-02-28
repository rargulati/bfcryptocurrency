package gossip

// TODO: sync.atomic
type update int64

type updateTracker struct {
	current update
	seen    map[update]bool
}

func (t *updateTracker) See(u update) bool {
	if u < t.current || t.seen[u] {
		return false
	}

	if t.seen == nil {
		t.seen = make(map[update]bool)
	}
	t.seen[u] = true
	for t.seen[t.current] {
		delete(t.seen, t.current)
		t.current++
	}
	return true
}
