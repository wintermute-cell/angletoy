package gui

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// GUI BACKEND STATE
type GuiBackendState struct {
	fonts map[string]rl.Font
}

var Gbs GuiBackendState

// This function is automatically called on package import.
func InitBackend() {
	Gbs = GuiBackendState{}
	Gbs.fonts = make(map[string]rl.Font)
	Gbs.fonts["default"] = rl.GetFontDefault()
	Gbs.fonts["alagard"] = rl.LoadFont("fonts/alagard.png")
}

// BACKEND FUNCTIONS

func backend_label(label Label) {
	style := parseStyleDef(label.style_info)

	color := rl.Black
	if c, ok := style["color"]; ok && c != nil {
		color = c.(rl.Color)
	}

	font := Gbs.fonts["default"]
	if f, ok := style["font"]; ok && f != nil {
		font = Gbs.fonts[f.(string)]
	}

	font_scale := float32(1.0)
	if v, ok := style["font-scale"]; ok && v != nil {
		font_scale = v.(float32)
	}

	rl.DrawTextEx(
		font,
		label.text,
		label.position,
		float32(font.BaseSize)*font_scale,
		float32(font.BaseSize/10)*font_scale,
		color,
	)
}

func backend_label_finalize(label Label) {
	// nothing to do here
}

func backend_scroll_panel(scroll_panel ScrollPanel) {
	style := parseStyleDef(scroll_panel.style_info)

	draw_debug := false
	if v, ok := style["debug"]; ok && v != nil {
		draw_debug = v.(bool)
	}

	bg_color := rl.NewColor(0, 0, 0, 0)
	if v, ok := style["background"]; ok && v != nil {
		bg_color = v.(rl.Color)
	}

	if draw_debug {
		// this represents the full bounds shifted by the scroll position
		// (useful for visualizing the scroll concept)
		virt_fbounds := scroll_panel.full_bounds
		virt_fbounds.X += scroll_panel.state.scroll_position.X
		virt_fbounds.Y += scroll_panel.state.scroll_position.Y

		rl.DrawRectangleRec(virt_fbounds, rl.Blue)
	}

	rl.DrawRectangleRec(scroll_panel.visible_bounds, bg_color)

	rl.BeginScissorMode(
		int32(scroll_panel.visible_bounds.X),
		int32(scroll_panel.visible_bounds.Y),
		int32(scroll_panel.visible_bounds.Width),
		int32(scroll_panel.visible_bounds.Height),
	)
}

func backend_scroll_panel_finalize(scroll_panel ScrollPanel) {
	rl.EndScissorMode()
}

func backend_button(button Button) {
	style := parseStyleDef(button.style_info)

	color := rl.White
	if c, ok := style["color"]; ok && c != nil {
		color = c.(rl.Color)
	}

	font := Gbs.fonts["default"]
	if f, ok := style["font"]; ok && f != nil {
		font = Gbs.fonts[f.(string)]
	}

	font_scale := float32(1.0)
	if v, ok := style["font-scale"]; ok && v != nil {
		font_scale = v.(float32)
	}

	btn_color_normal := rl.Blue
	if v, ok := style["background"]; ok && v != nil {
		btn_color_normal = v.(rl.Color)
	}

	btn_color_hovered := rl.SkyBlue
	if v, ok := style["background-hovered"]; ok && v != nil {
		btn_color_hovered = v.(rl.Color)
	}

	btn_color_pressed := rl.DarkBlue
	if v, ok := style["background-pressed"]; ok && v != nil {
		btn_color_pressed = v.(rl.Color)
	}

	// determine appropriate colors based on current interaction state
	btn_color := btn_color_normal
	switch button.state {
	case ButtonStateHovered:
		btn_color = btn_color_hovered
	case ButtonStatePressed:
		btn_color = btn_color_pressed
	}

	bounds := rl.NewRectangle(button.position.X, button.position.Y, button.size.X, button.size.Y)
	rl.DrawRectangleRec(bounds, btn_color)
	rl.DrawTextEx(
		font,
		button.text,
		button.position,
		float32(font.BaseSize)*font_scale,
		float32(font.BaseSize/10)*font_scale,
		color,
	)
}

func backend_button_finalize(button Button) {
	// nothing to do here
}

func backend_slider(slider Slider) {
	style := parseStyleDef(slider.style_info)

	color := rl.White
	if c, ok := style["color"]; ok && c != nil {
		color = c.(rl.Color)
	}

	background := rl.Blue
	background_normal := rl.Blue
	if v, ok := style["background"]; ok && v != nil {
		background_normal = v.(rl.Color)
	}

	// TODO: hover and drag state for slider
	background = background_normal

	// draw slider background
	rl.DrawRectangle(int32(slider.position.X), int32(slider.position.Y), int32(slider.size.X), int32(slider.size.Y), background)

	// draw handle
	rl.DrawRectangle(
		int32(slider.handle_position.X),
		int32(slider.handle_position.Y),
		int32(slider.handle_size.X),
		int32(slider.handle_size.Y),
		color)
}

func backend_slider_finalize(slider Slider) {
	// nothing to do here
}
