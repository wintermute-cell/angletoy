package gui

import (
	"gorl/fw/core/logging"
	"gorl/fw/util"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Widget interface {
	SetPosition(p rl.Vector2)
	GetPosition() rl.Vector2
	SetSize(s rl.Vector2)
	GetSize() rl.Vector2
	Bounds() rl.Rectangle
}

// ----------------
//   BASE WIDGET  |
// ----------------

type BaseWidget struct {
	position rl.Vector2
	size     rl.Vector2
}

func (w *BaseWidget) SetPosition(p rl.Vector2) {
	w.position = p
}

func (w *BaseWidget) GetPosition() rl.Vector2 {
	return w.position
}

func (w *BaseWidget) SetSize(s rl.Vector2) {
	w.size = s
}

func (w *BaseWidget) GetSize() rl.Vector2 {
	return w.size
}

func (w *BaseWidget) Bounds() rl.Rectangle {
	return rl.NewRectangle(w.position.X, w.position.Y, w.size.X, w.size.Y)
}

// ----------------
//    CONTAINER   |
// ----------------

type Container struct {
	BaseWidget
	Children []Widget
}

func NewContainer() *Container {
	return &Container{}
}

func (c *Container) AddChild(child Widget) {
	c.Children = append(c.Children, child)
}

func (c *Container) RemoveChild(target Widget) {
	for i, child := range c.Children {
		if child == target {
			c.Children = util.SliceDelete(c.Children, i, i+1)
			return
		}
	}
}

// ----------------
//      LABEL     |
// ----------------

// widget definition
type Label struct {
	BaseWidget
	text                   string
	state                  *LabelState
	style_info             string
	watched_string         *string
	watched_int32          *int32
	watched_int32_format   string
	watched_float32        *float32
	watched_float32_format string
	watch_flag             int32 // 0 = none, 1 = string, 2 = int32, 3 = float32
	// we have a state pointer inside the widget itself, since each widget must
	// have the ability to update its own state, and do these updates based on
	// it's state. (see ScrollPanel for example, the scroll position is not
	// ephemeral)
}

// state info
type LabelState struct {
	// label does not need state info yet
}

// update function
func (label *Label) update_label() {
	// maybe add a hover tooltip here, or whatever else
	switch label.watch_flag {
	case 1:
		label.text = *label.watched_string
	case 2:
		label.text = fmt.Sprintf(label.watched_int32_format, *label.watched_int32)
	case 3:
		label.text = fmt.Sprintf(label.watched_float32_format, *label.watched_float32)
	}
}

// constructor
func NewLabel(text string, position rl.Vector2, style_info string) *Label {
	l := Label{text: text, style_info: style_info}
	l.position = position
	l.state = &LabelState{}
	return &l
}

// Set the labels text
func (label *Label) SetText(new_text string) {
	label.text = new_text
}

// Begin watching a string and change the label text whenever the watched
// string changes.
func (label *Label) WatchString(watch_target *string) {
	label.watch_flag = 1
	label.watched_string = watch_target
}

// Begin watching an int32 and change the label text whenever the watched
// int32 changes.
func (label *Label) WatchInt32(watch_target *int32, format string) {
	label.watch_flag = 2
	label.watched_int32 = watch_target
	label.watched_int32_format = format
}

// Begin watching a float32 and change the label text whenever the watched
// float32 changes.
func (label *Label) WatchFloat32(watch_target *float32, format string) {
	label.watch_flag = 3
	label.watched_float32 = watch_target
	label.watched_float32_format = format
}

// Stop watching whatever the label is watching.
func (label *Label) StopWatching() {
	label.watch_flag = 0
}

// ----------------
//     BUTTON     |
// ----------------

// widget definiton
type Button struct {
	BaseWidget
	text       string
	callback   func(ButtonState)
	state      ButtonState
	style_info string
}

// state info
type ButtonState int32

const (
	ButtonStateNone ButtonState = iota
	ButtonStateHovered
	ButtonStatePressed
	ButtonStateReleased
)

// update function
func (button *Button) update_button() {
	bounds := rl.NewRectangle(button.position.X, button.position.Y, button.size.X, button.size.Y)
	button.state = ButtonStateNone
	if rl.CheckCollisionPointRec(rl.GetMousePosition(), bounds) {
		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			// mouse down, button is pressed down
			button.state = ButtonStatePressed
		} else {
			// mouse is hovering, no click occured
			button.state = ButtonStateHovered
		}
		if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
			// button released, do the button action
			button.state = ButtonStateReleased
		}
		button.callback(button.state)
	}
}

// constructor
func NewButton(text string, position, size rl.Vector2, callback func(ButtonState), style_info string) *Button {
	// NOTE: The commented code below won't work out, since data only flows
	// logic -> rendering, and the font is a part of the rendering.
	//
	//if size == rl.Vector2Zero() {
	//	// NOTE: this might have to be changed if we use custom spacings (using MeasureTextEx)
	//	text_size := rl.MeasureText(text, font.BaseSize*int32(font_scale))
	//	size = rl.NewVector2(float32(text_size), float32(font.BaseSize)*font_scale)
	//}
	new_button := &Button{text: text, callback: callback, state: ButtonStateNone, style_info: style_info}
	new_button.size = size
	new_button.position = position
	return new_button
}

// ----------------
//  SCROLL PANEL  |
// ----------------

// widget definiton
type ScrollPanel struct {
	BaseWidget
	visible_bounds      rl.Rectangle
	full_bounds         rl.Rectangle
	state               *ScrollPanelState
	container           *Container
	reference_positions map[Widget]rl.Vector2 // stores the original positions of children
	style_info          string
}

// state info
type ScrollPanelState struct {
	scroll_position rl.Vector2
}

// update function
func (scroll_panel *ScrollPanel) update_scroll_panel() {
	// check if the mouse overlaps the visible bounds
	if rl.CheckCollisionPointRec(rl.GetMousePosition(), scroll_panel.visible_bounds) {
		wheel_move := rl.GetMouseWheelMoveV()

		// Compute the new scroll position
		new_scroll_position := rl.Vector2Add(
			scroll_panel.state.scroll_position,
			rl.Vector2Scale(wheel_move, rl.GetFrameTime()*1800), // scroll speed
		)

		maxXScroll := scroll_panel.visible_bounds.Width - scroll_panel.full_bounds.Width
		maxYScroll := scroll_panel.visible_bounds.Height - scroll_panel.full_bounds.Height

		// Limit the scroll_position based on full_bounds.
		if new_scroll_position.X > 0 {
			new_scroll_position.X = 0
		}
		if new_scroll_position.Y > 0 {
			new_scroll_position.Y = 0
		}
		if new_scroll_position.X < maxXScroll {
			new_scroll_position.X = maxXScroll
		}
		if new_scroll_position.Y < maxYScroll {
			new_scroll_position.Y = maxYScroll
		}

		// Apply the new scroll position
		scroll_panel.state.scroll_position = new_scroll_position

		for _, child_widget := range scroll_panel.container.Children {
			new_pos := rl.Vector2Add(
				scroll_panel.reference_positions[child_widget],
				scroll_panel.state.scroll_position)
			child_widget.SetPosition(new_pos)
		}
	}
}

// constructor
func NewScrollPanel(visible_bounds, full_bounds rl.Rectangle, style_info string) *ScrollPanel {
	return &ScrollPanel{
		visible_bounds:      visible_bounds,
		full_bounds:         full_bounds,
		state:               &ScrollPanelState{},
		container:           NewContainer(),
		reference_positions: make(map[Widget]rl.Vector2),
		style_info:          style_info,
	}
}

func (scroll_panel *ScrollPanel) AddChild(child Widget) {
	scroll_panel.container.Children = append(scroll_panel.container.Children, child)
	scroll_panel.reference_positions[child] = child.GetPosition()
}

func (scroll_panel *ScrollPanel) RemoveChild(target Widget) {
	for i, child := range scroll_panel.container.Children {
		if child == target {
			// Remove child while preserving order
			scroll_panel.container.Children = util.SliceDelete(scroll_panel.container.Children, i, i+1)
			return
		}
	}
}

// ----------------
//     SLIDER     |
// ----------------

// widget definiton
type Slider struct {
	BaseWidget
	max_value              float32
	min_value              float32
	value_interval         float32
	current_value          float32
	last_value             float32
	handle_size            rl.Vector2
	handle_position        rl.Vector2
	value_changed_callback func(new_value float32)
	is_dragging            bool
	style_info             string
}

// update function
func (slider *Slider) update() {
	// Define slider boundary
	slider_bounds := rl.NewRectangle(slider.position.X, slider.position.Y, slider.size.X, slider.size.Y)

	// Calculate the position of the handle based on the current value
	percentage := (slider.current_value - slider.min_value) / (slider.max_value - slider.min_value)
	handlePosX := slider.position.X + (slider.size.X-slider.handle_size.X)*percentage
	handlePosY := slider.position.Y + (slider.size.Y-slider.handle_size.Y)/2
	slider.handle_position = rl.NewVector2(handlePosX, handlePosY)
	handle_bounds := rl.NewRectangle(
		handlePosX,
		handlePosY,
		slider.handle_size.X,
		slider.handle_size.Y,
	)

	// Check for collisions with mouse position
	slider_collision := rl.CheckCollisionPointRec(rl.GetMousePosition(), slider_bounds)
	handle_collision := rl.CheckCollisionPointRec(rl.GetMousePosition(), handle_bounds)

	is_mouse_down := rl.IsMouseButtonDown(rl.MouseLeftButton)
	mouse_x := float32(rl.GetMouseX())

	// Clicked on the slider but not on the handle
	if slider_collision && !handle_collision && is_mouse_down {
		// Update the current value based on clicked position
		click_percentage := (mouse_x - slider.position.X) / slider.size.X
		slider.current_value = slider.min_value + click_percentage*(slider.max_value-slider.min_value)

		// Clamp the value to available intervals
		slider.current_value = util.Round(slider.current_value/slider.value_interval) * slider.value_interval
		slider.current_value = util.Clamp(slider.current_value, slider.min_value, slider.max_value)
	}

	// Handle dragging logic
	if handle_collision && is_mouse_down {
		// Calculate the new value based on the mouse's X position
		drag_percentage := (mouse_x - slider.position.X) / slider.size.X
		slider.current_value = slider.min_value + drag_percentage*(slider.max_value-slider.min_value)

		// Again, clamp to the available intervals
		slider.current_value = util.Round(slider.current_value/slider.value_interval) * slider.value_interval
		slider.current_value = util.Clamp(slider.current_value, slider.min_value, slider.max_value)
	}

	// call value changed callback if value was changed
	if slider.last_value != slider.current_value && slider.value_changed_callback != nil {
		slider.value_changed_callback(slider.current_value)
	}
	slider.last_value = slider.current_value
}

// constructor
func NewSlider(min, max, starting_value, interval float32, position, size, handle_size rl.Vector2, style_info string) *Slider {
	s := &Slider{
		min_value:              min,
		max_value:              max,
		value_interval:         interval,
		current_value:          starting_value,
		last_value:             starting_value,
		handle_size:            handle_size,
		is_dragging:            false,
		value_changed_callback: nil,
		style_info:             style_info,
	}
	s.position = position
	s.size = size
	return s
}

// Register a callback that is called whenever the Sliders value changes.
// The new value is passed to that callback.
func (slider *Slider) SetValueChangedCallback(callback func(new_value float32)) {
	slider.value_changed_callback = callback
}

// Get the sliders current value
func (slider *Slider) GetCurrentValue() float32 {
	return slider.current_value
}

// Get a pointer to the sliders current value (useful for hooking up with label)
func (slider *Slider) GetCurrentValuePointer() *float32 {
	return &slider.current_value
}

// Set the sliders current value
func (slider *Slider) SetCurrentValue(new_value float32) {
	slider.current_value = new_value
}

// ----------------
//       GUI      |
// ----------------

//type Gui struct {
//    // Widget is an interface, thus no need to make this a slice of pointers.
//    // (as long as we make sure we only feed pointers into it lol)
//    widgets []Widget
//}

type Gui struct {
	container Container
}

func NewGui() *Gui {
	return &Gui{}
}

func (gui *Gui) AddWidget(widget Widget) {
	gui.container.AddChild(widget)
}

func (gui *Gui) RemoveWidget(widget Widget) {
	gui.container.RemoveChild(widget)
}

func (gui *Gui) Draw() {
	doRecursiveDraw(gui.container)
}

func doRecursiveDraw(container Container) {
	for _, widget := range container.Children {
		switch w := any(widget).(type) {
		case *Label:
			w.update_label()
			backend_label(*w)
			backend_label_finalize(*w)
		case *Button:
			w.update_button()
			backend_button(*w)
			backend_button_finalize(*w)
		case *ScrollPanel:
			w.update_scroll_panel()
			backend_scroll_panel(*w)
			doRecursiveDraw(*w.container) // draw the panels children
			backend_scroll_panel_finalize(*w)
		case *Slider:
			w.update()
			backend_slider(*w)
			backend_slider_finalize(*w)
		default:
			logging.Error("Attempted to draw GUI widget type with missing draw case: %v", w)
		}
	}
}
