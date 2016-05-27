package engo

const (
	// AxisMax is the maximum value a joystick or keypress axis will reach
	AxisMax float32 = 1
	// AxisMin is the value an axis returns if there has been to state change.
	AxisNeutral float32 = 0
	// AxisMin is the minimum value a joystick or keypress axis will reach
	AxisMin float32 = -1
)

// NewInputManager holds onto anything input related for engo
func NewInputManager() *InputManager {
	return &InputManager{
		axes:    make(map[string]Axis),
		buttons: make(map[string]Button),
		keys:    &KeyManager{mapper: make(map[Key]KeyState)},
		mouse:   &MouseManager{mapper: make(map[MouseButton]KeyState)},
	}
}

// InputManager contains information about all forms of input.
type InputManager struct {
	// Mouse is InputManager's reference to the mouse. It is recommended to use the
	// Axis and Button system if at all possible.
	axes    map[string]Axis
	buttons map[string]Button
	keys    *KeyManager
	mouse   *MouseManager
}

func (im *InputManager) update() {
	im.keys.update()
	im.mouse.update()
}

// RegisterAxis registers a new axis which can be used to retrieve inputs which are spectrums.
func (im *InputManager) RegisterAxis(name string, pairs ...AxisPair) {
	im.axes[name] = Axis{
		Name:  name,
		Pairs: pairs,
	}
}

// RegisterButton registers a new button input.
func (im *InputManager) RegisterButton(name string, keys ...Key) {
	im.buttons[name] = Button{
		Triggers: keys,
		Name:     name,
	}
}

// Axis retrieves an Axis with a specified name.
func (im *InputManager) Axis(name string) Axis {
	return im.axes[name]
}

// Button retrieves a Button with a specified name.
func (im *InputManager) Button(name string) Button {
	return im.buttons[name]
}

func (im *InputManager) Key(key Key) KeyState {
	return im.keys.Get(key)
}

func (im *InputManager) MouseButton(mb MouseButton) KeyState {
	return im.mouse.Button(mb)
}

func (im *InputManager) MouseModifierKey(mod ModifierKey) ModifierState {
	return im.mouse.Modifier(mod)
}

func (im *InputManager) MouseScroll() (float32, float32) {
	return im.mouse.Scroll()
}

func (im *InputManager) MousePosition() (float32, float32) {
	return im.mouse.Position()
}
