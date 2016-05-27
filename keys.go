package engo

import (
	"sync"
)

//TODO: Need better names for this stuff mabye?

type KeyAction int

const (
	KeyUp       = KeyAction(0)
	KeyDown     = KeyAction(1)
	KeyJustDown = KeyAction(2)
	KeyJustUp   = KeyAction(3)
)

// KeyManager tracks which keys are pressed and released at the current point of time.
type KeyManager struct {
	mapper map[Key]KeyState
}

// Set is used for updating whether or not a key is held down, or not held down.
func (km *KeyManager) Set(k Key, state bool) {
	ks := km.mapper[k]
	ks.set(state)
	km.mapper[k] = ks
}

// Get retrieves a keys state.
func (km *KeyManager) Get(k Key) KeyState {
	ks, ok := km.mapper[k]
	if !ok {
		return KeyState{lastState: false, currentState: false}
	}

	return ks
}

func (km *KeyManager) update() {
	// Set all keys to their current states
	for key, state := range km.mapper {
		state.set(state.currentState)
		km.mapper[key] = state
	}
}

// KeyState is used for detecting the state of a key press.
type KeyState struct {
	lastState    bool
	currentState bool

	mutex sync.RWMutex
}

func (key *KeyState) set(state bool) {
	key.mutex.Lock()

	key.lastState = key.currentState
	key.currentState = state

	key.mutex.Unlock()
}

// State returns the raw state of a key.
func (key KeyState) State() KeyAction {
	key.mutex.RLock()
	defer key.mutex.RUnlock()

	if !key.lastState && key.currentState {
		return KeyJustDown
	} else if key.lastState && !key.currentState {
		return KeyJustUp
	} else if key.lastState && key.currentState {
		return KeyDown
	} else if !key.lastState && !key.currentState {
		return KeyUp
	}

	return KeyUp
}

func (key KeyState) JustPressed() bool {
	return key.State() == KeyJustDown
}

func (key KeyState) JustReleased() bool {
	return key.State() == KeyJustUp
}

func (key KeyState) Up() bool {
	return key.State() == KeyUp
}

func (key KeyState) Down() bool {
	return key.State() == KeyDown
}
