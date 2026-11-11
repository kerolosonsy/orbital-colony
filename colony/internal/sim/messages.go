package sim

import (
	"fmt"
	"sync"

	"orbital/colony/i18n"
)

// A colony stream can have subscribers in both languages. Carry the exact
// message pairs with each event; formatted values never build up in memory.
var pendingMessages = struct {
	sync.Mutex
	pairs map[[2]string]struct{}
}{pairs: make(map[[2]string]struct{})}

func say(ar, en string) string {
	pendingMessages.Lock()
	pendingMessages.pairs[[2]string{ar, en}] = struct{}{}
	pendingMessages.Unlock()
	return i18n.T(ar, en)
}

func sayf(ar, en string, args ...any) string {
	return say(fmt.Sprintf(ar, args...), fmt.Sprintf(en, args...))
}

func takeMessages() [][2]string {
	pendingMessages.Lock()
	defer pendingMessages.Unlock()
	pairs := make([][2]string, 0, len(pendingMessages.pairs))
	for pair := range pendingMessages.pairs {
		pairs = append(pairs, pair)
	}
	clear(pendingMessages.pairs)
	return pairs
}
