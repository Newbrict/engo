package common

import (
	"github.com/luxengine/math"

	"engo.io/ecs"
	"engo.io/engo"
)

type Event interface {
	Type() string
	Notify(*SpaceComponent, *CameraSystem) bool
}

type KeyEvent struct {
	engo.Key
	engo.KeyAction
}

func (ke KeyEvent) Type() string {
	return "Key Event"
}

func (ke *KeyEvent) Notify(*SpaceComponent, *CameraSystem) bool {
	if engo.Input.Key(ke.Key).State() == ke.KeyAction {
		return true
	}

	return false
}

type MouseButtonEvent struct {
	engo.MouseButton
	engo.KeyAction
}

func (mbe MouseButtonEvent) Type() string {
	return "MouseButton Event"
}

func (mbe *MouseButtonEvent) Notify(*SpaceComponent, *CameraSystem) bool {
	if engo.Input.MouseButton(mbe.MouseButton).State() == mbe.KeyAction {
		return true
	}
	return false
}

type MouseEventType interface {
	State() bool
}

type hoverEvent struct {
	triggered bool
}

func (he hoverEvent) State() bool {
	return he.triggered
}

type dragEvent struct {
	triggered bool
}

func (de dragEvent) State() bool {
	return de.triggered
}

var (
	Hover hoverEvent = hoverEvent{false}
	Drag  dragEvent  = dragEvent{false}
)

type MouseAction int

const (
	JustHovered MouseAction = MouseAction(0)
	Hovering    MouseAction = MouseAction(1)
	JustExited  MouseAction = MouseAction(2)
	NotHovering MouseAction = MouseAction(3)

	JustDragged  MouseAction = MouseAction(4)
	Dragging     MouseAction = MouseAction(5)
	JustReleased MouseAction = MouseAction(6)
	NotDragging  MouseAction = MouseAction(7)
)

type MouseEvent struct {
	MouseEventType
	MouseAction
}

func (me MouseEvent) Type() string {
	return "MouseEvent"
}

func (me *MouseEvent) Notify(space *SpaceComponent, camera *CameraSystem) bool {

	if camera != nil {
		x, y := engo.Input.MousePosition()
		gameWidth, gameHeight := engo.GameWidth(), engo.GameHeight()
		canvasWidth, canvasHeight := engo.CanvasWidth(), engo.CanvasHeight()

		// Translate Mouse.X and Mouse.Y into "game coordinates"
		x = x*camera.z*(gameWidth/canvasWidth) + camera.x - (gameWidth/2)*camera.z
		y = y*camera.z*(gameHeight/canvasHeight) + camera.y - (gameHeight/2)*camera.z

		// Rotate if needed
		if camera.angle != 0 {
			sin, cos := math.Sincos(camera.angle * math.Pi / 180)
			x, y = x*cos+y*sin, y*cos-x*sin
		}

		within := space.Within(engo.Point{x, y})
		switch mouseEvent := me.MouseEventType.(type) {
		case hoverEvent:
			if within {
				if !mouseEvent.State() {
					me.MouseEventType = hoverEvent{true}
					if me.MouseAction == JustHovered {
						return true
					}
				} else {
					if me.MouseAction == Hovering {
						return true
					}
				}
			} else {
				if mouseEvent.State() {
					me.MouseEventType = hoverEvent{false}
					if me.MouseAction == JustExited {
						return true
					}
				} else {
					if me.MouseAction == NotHovering {
						return true
					}
				}
			}

		case dragEvent:
			mouseBtnLeftState := engo.Input.MouseButton(engo.MouseButtonLeft)
			if within && (mouseBtnLeftState.JustPressed() || mouseBtnLeftState.Down()) {
				if !mouseEvent.State() {
					me.MouseEventType = dragEvent{true}
					if me.MouseAction == JustDragged {
						return true
					}
				} else {
					if me.MouseAction == Dragging {
						return true
					}
				}
			} else {
				if mouseEvent.State() {
					me.MouseEventType = dragEvent{false}
					if me.MouseAction == JustReleased {
						return true
					}
				} else {
					if me.MouseAction == NotDragging {
						return true
					}
				}
			}
		}
	}

	return false
}

type EventHandler func(Event)

type EventComponent struct {
	*SpaceComponent
	handlers map[Event][]EventHandler
}

func (ec *EventComponent) Bind(event Event, handlers ...EventHandler) {
	if ec.handlers == nil {
		println("Event System not added.  Unable to bind event.")
		return
	}

	for _, handler := range handlers {
		ec.handlers[event] = append(ec.handlers[event], handler)
	}
}

type eventEntity struct {
	*ecs.BasicEntity
	*EventComponent
}

type EventSystem struct {
	entities []eventEntity
	camera   *CameraSystem
}

func (es *EventSystem) New(w *ecs.World) {
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *CameraSystem:
			es.camera = sys
		}
	}
}

func (es *EventSystem) Add(basic *ecs.BasicEntity, event *EventComponent, space *SpaceComponent) {
	event.SpaceComponent = space
	event.handlers = make(map[Event][]EventHandler)

	es.entities = append(es.entities, eventEntity{basic, event})
}

func (es *EventSystem) Remove(basic ecs.BasicEntity) {
	for index, entity := range es.entities {
		if entity.ID() == basic.ID() {
			es.entities = append(es.entities[:index], es.entities[index+1:]...)
			break
		}
	}
}

func (es *EventSystem) Update(dt float32) {
	for _, entity := range es.entities {
		for event, eventHandlers := range entity.handlers {
			if event.Notify(entity.SpaceComponent, es.camera) {
				for _, handler := range eventHandlers {
					go handler(event)
				}
			}
		}
	}
}
