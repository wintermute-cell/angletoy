package lighting

import (
	"gorl/fw/core/render"
	"gorl/fw/core/logging"
	"gorl/fw/util"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// the lighting system keeps track of all the existing lights
type lighting_system struct {
	lights                        []*Light2D
	occluders                     []Occluder2D
	lighting_target               rl.RenderTexture2D
	normal_collection_tex         rl.RenderTexture2D
	lightmap_blur_shader          rl.Shader
	blur_shader_ambient_level_loc int32
	ambient_light_level           float32
	is_lighting_enabled           bool
	lightmap_extend               int32 // a number by which to extend the lightmap beyond the screen resolution
}

var ls lighting_system

func Enable() {
	ls.is_lighting_enabled = true
	logging.Info("Enabled lighting")
}

func Disable() {
	ls.is_lighting_enabled = false
	logging.Info("Disabled lighting")
}

// SetAmbientLight sets the ambient light level. This is a value between 0 and 1.
func SetAmbientLight(light_level float32) {
	l := util.Clamp(light_level, 0.0, 1.0)
	ls.ambient_light_level = l
}

// GetAmbientLight returns the ambient light level. This is a value between 0 and 1.
func GetAmbientLight() float32 {
	return ls.ambient_light_level
}

func UnloadOccluder(occluder Occluder2D) {
	for i, o := range ls.occluders {
		if occluder == o {
			ls.occluders = util.SliceDelete(ls.occluders, i, i+1)
			return
		}
	}
}

func UnloadLight(light *Light2D) {
	for i, l := range ls.lights {
		if light == l {
			ls.lights = util.SliceDelete(ls.lights, i, i+1)
			return
		}
	}
}

type Occluder2D interface {
	Draw()
	DrawNormal()
	GetPosition() rl.Vector2
	GetOrigin() rl.Vector2
	GetRotation() float32
	GetSize() rl.Vector2
	getTexCoordOrigin() rl.Vector2
}

func InitLighting() {
	ls = lighting_system{
		lightmap_blur_shader: rl.LoadShader("", "shaders/lightmap-blur.glsl"),
		lightmap_extend:      400,
		ambient_light_level:  0.13,
	}
	ls.blur_shader_ambient_level_loc = rl.GetShaderLocation(ls.lightmap_blur_shader, "ambient_light_level")
	ls.lighting_target = rl.LoadRenderTexture(
		int32(render.Rs.RenderResolution.X)+ls.lightmap_extend, int32(render.Rs.RenderResolution.Y)+ls.lightmap_extend)
	ls.normal_collection_tex = rl.LoadRenderTexture(
		int32(render.Rs.RenderResolution.X)+ls.lightmap_extend, int32(render.Rs.RenderResolution.Y)+ls.lightmap_extend)
}

func DeinitLighting() {
	rl.UnloadRenderTexture(ls.lighting_target)
	for _, l := range ls.lights {
		rl.UnloadRenderTexture(l.occlusion_map)
		rl.UnloadRenderTexture(l.polar_shadowmap)
		rl.UnloadShader(l.polar_transform_shader)
		rl.UnloadShader(l.light_render_shader)
		rl.UnloadShader(l.normal_lighting_shader)
	}
}

// ----------------------------
//
//    OCCLUDER DEFINITIONS
//
// ----------------------------

// ------------------
//
//	OCCLUDER SPRITE
//
// ------------------
type OccluderSprite2D struct {
	sprite   rl.Texture2D
	position rl.Vector2
	size     rl.Vector2
	origin   rl.Vector2
	rotation float32
}

// An OccluderSprite2D is a simple light blocking sprite. Passing a size of (0, 0) will use the sprite's size.
func NewOccluderSprite2D(sprite rl.Texture2D, position rl.Vector2, size rl.Vector2, origin rl.Vector2, rotation float32) *OccluderSprite2D {
	if size == rl.Vector2Zero() {
		size = rl.NewVector2(float32(sprite.Width), float32(sprite.Height))
	}
	new_occluder := OccluderSprite2D{
		sprite:   sprite,
		position: position,
		size:     size,
		origin:   origin,
		rotation: rotation,
	}
	ls.occluders = append(ls.occluders, &new_occluder)
	return &new_occluder
}

// Like NewOccluderSprite2D, but with a position, origin and rotation 0 and a
// size of the sprite's size.
func NewOccluderSprite2DZ(
	sprite rl.Texture2D,
) *OccluderSprite2D {
	size := rl.NewVector2(float32(sprite.Width), float32(sprite.Height))
	new_occluder := OccluderSprite2D{
		sprite:   sprite,
		position: rl.Vector2Zero(),
		size:     size,
		origin:   rl.Vector2Zero(),
		rotation: 0.0,
	}
	ls.occluders = append(ls.occluders, &new_occluder)
	return &new_occluder
}

func (o *OccluderSprite2D) Draw() {
	rl.DrawTexturePro(
		o.sprite,
		rl.NewRectangle(0, 0, float32(o.sprite.Width), float32(o.sprite.Height)),
		rl.NewRectangle(o.position.X, o.position.Y, o.size.X, o.size.Y),
		o.origin,
		o.rotation,
		rl.White,
	)
}

func (o *OccluderSprite2D) DrawNormal() {
	// do nothing
}

func (o *OccluderSprite2D) GetRotation() float32 {
	return o.rotation
}

func (o *OccluderSprite2D) GetSize() rl.Vector2 {
	return o.size
}

func (o *OccluderSprite2D) GetPosition() rl.Vector2 {
	return o.position
}

func (o *OccluderSprite2D) GetOrigin() rl.Vector2 {
	return o.origin
}

func (o *OccluderSprite2D) getTexCoordOrigin() rl.Vector2 {
	tex_coord_origin := rl.NewVector2(
		o.origin.X/float32(o.sprite.Width),
		o.origin.Y/float32(o.sprite.Height),
	)
	return tex_coord_origin
}

// Update the occluder's position, size, origin and rotation. If size is (0,
// 0), the sprite's existing size will be used.
func (o *OccluderSprite2D) Update(position, size, origin rl.Vector2, rotation float32) {
	if size == rl.Vector2Zero() {
		size = o.size
	}
	o.position = position
	o.size = size
	o.origin = origin
	o.rotation = rotation
}

// ------------------
//
//	OCCLUDER NORMAL
//
// ------------------
type OccluderNormal2D struct {
	sprite        rl.Texture2D
	normal_sprite rl.Texture2D
	position      rl.Vector2
	size          rl.Vector2
	origin        rl.Vector2
	rotation      float32
}

// An OccluderNormal2D blocks light like a sprite, but also has a normal map for more detailed lighting on the sprite.
// Passing a size of (0, 0) will use the sprite's size.
func NewOccluderNormal2D(
	sprite rl.Texture2D,
	normal_sprite rl.Texture2D,
	position rl.Vector2,
	size rl.Vector2,
	origin rl.Vector2,
	rotation float32,
) *OccluderNormal2D {
	if size == rl.Vector2Zero() {
		size = rl.NewVector2(float32(sprite.Width), float32(sprite.Height))
	}
	new_occluder := OccluderNormal2D{
		sprite:        sprite,
		normal_sprite: normal_sprite,
		position:      position,
		size:          size,
		origin:        origin,
		rotation:      rotation,
	}
	ls.occluders = append(ls.occluders, &new_occluder)
	return &new_occluder
}

// Like NewOccluderNormal2D, but with a position, origin and rotation 0 and a
// size of the sprite's size.
func NewOccluderNormal2DZ(
	sprite rl.Texture2D,
	normal_sprite rl.Texture2D,
) *OccluderNormal2D {
	size := rl.NewVector2(float32(sprite.Width), float32(sprite.Height))
	new_occluder := OccluderNormal2D{
		sprite:        sprite,
		normal_sprite: normal_sprite,
		position:      rl.Vector2Zero(),
		size:          size,
		origin:        rl.Vector2Zero(),
		rotation:      0.0,
	}
	ls.occluders = append(ls.occluders, &new_occluder)
	return &new_occluder
}

func (o *OccluderNormal2D) Draw() {
	rl.DrawTexturePro(
		o.sprite,
		rl.NewRectangle(0, 0, float32(o.sprite.Width), float32(o.sprite.Height)),
		rl.NewRectangle(o.position.X, o.position.Y, o.size.X, o.size.Y),
		o.origin,
		o.rotation,
		rl.White,
	)
}

func (o *OccluderNormal2D) DrawNormal() {
	rl.DrawTexturePro(
		o.normal_sprite,
		rl.NewRectangle(0, 0, float32(o.normal_sprite.Width), float32(o.normal_sprite.Height)),
		rl.NewRectangle(o.position.X, o.position.Y, o.size.X, o.size.Y),
		o.origin,
		o.rotation,
		rl.White,
	)
}

func (o *OccluderNormal2D) GetRotation() float32 {
	return o.rotation
}

func (o *OccluderNormal2D) GetSize() rl.Vector2 {
	return o.size
}

func (o *OccluderNormal2D) GetPosition() rl.Vector2 {
	return o.position
}

func (o *OccluderNormal2D) GetOrigin() rl.Vector2 {
	return o.origin
}

func (o *OccluderNormal2D) getTexCoordOrigin() rl.Vector2 {
	tex_coord_origin := rl.NewVector2(
		o.origin.X/float32(o.sprite.Width),
		o.origin.Y/float32(o.sprite.Height),
	)
	return tex_coord_origin
}

// Update the occluder's position, size, origin and rotation. If size is (0,
// 0), the sprite's existing size will be used.
func (o *OccluderNormal2D) Update(position, size, origin rl.Vector2, rotation float32) {
	if size == rl.Vector2Zero() {
		size = o.size
	}
	o.position = position
	o.size = size
	o.rotation = rotation
	o.origin = origin
}

// ------------------
//
//	OCCLUDER CIRCLE
//
// ------------------
type OccluderCircle2D struct {
	radius   float32
	position rl.Vector2
}

// An OccluderCircle2D is a simple light blocking circle.
func NewOccluderCircle2D(position rl.Vector2, radius float32) *OccluderCircle2D {
	new_occluder := OccluderCircle2D{
		radius:   radius,
		position: position,
	}
	ls.occluders = append(ls.occluders, &new_occluder)
	return &new_occluder
}

func (o *OccluderCircle2D) Draw() {
	rl.DrawCircleV(o.position, o.radius, rl.White)
}

func (o *OccluderCircle2D) DrawNormal() {
	// do nothing
}

func (o *OccluderCircle2D) GetRotation() float32 {
	return 0.0
}

func (o *OccluderCircle2D) GetSize() rl.Vector2 {
	return rl.NewVector2(o.radius, o.radius)
}

func (o *OccluderCircle2D) GetPosition() rl.Vector2 {
	return o.position
}

func (o *OccluderCircle2D) GetOrigin() rl.Vector2 {
	return rl.Vector2Zero()
}

func (o *OccluderCircle2D) getTexCoordOrigin() rl.Vector2 {
	return rl.Vector2Zero()
}

func (o *OccluderCircle2D) Update(position rl.Vector2, radius float32) {
	o.position = position
	o.radius = radius
}

type Light2D struct {
	// Custom Fields
	occlusion_map   rl.RenderTexture2D
	polar_shadowmap rl.RenderTexture2D

	polar_transform_shader         rl.Shader
	polar_transform_resolution_loc int32

	light_render_shader          rl.Shader
	light_render_resolution_loc  int32
	light_render_light_color_loc int32

	normal_lighting_shader      rl.Shader
	normal_light_pos_loc        int32
	normal_light_color_loc      int32
	normal_light_rot_loc        int32
	normal_light_origin_loc     int32
	normal_light_sprite_pos_loc int32
	normal_light_range_loc      int32

	occlusion_shift_cam rl.Camera2D

	position   rl.Vector2
	size       rl.Vector2
	resolution float32

	color         rl.Color
	mask          rl.Texture2D
	mask_rotation float32

	is_enabled bool
}

/*
NewLight2D creates a new light and adds it to the lighting system

size:

	the size of the light in pixels

resolution:

	the resolution of the light's shadowmap. Could be thought of like the
	amount of "rays" shooting out from the lightsource to determine shadows.
*/
func NewLight2D(size rl.Vector2, resolution float32, color rl.Color) *Light2D {
	new_light := Light2D{
		size:                   size,
		resolution:             resolution,
		color:                  color,
		occlusion_map:          rl.LoadRenderTexture(int32(size.X), int32(size.Y)),
		polar_shadowmap:        rl.LoadRenderTexture(int32(resolution), 1),
		polar_transform_shader: rl.LoadShader("", "shaders/polar-lightmap.glsl"),
		light_render_shader:    rl.LoadShader("", "shaders/light-render.glsl"),
		normal_lighting_shader: rl.LoadShader("", "shaders/normal-lighting.glsl"),
		is_enabled:             true,
	}
	new_light.polar_transform_resolution_loc = rl.GetShaderLocation(new_light.polar_transform_shader, "resolution")
	new_light.light_render_resolution_loc = rl.GetShaderLocation(new_light.light_render_shader, "resolution")
	new_light.light_render_light_color_loc = rl.GetShaderLocation(new_light.light_render_shader, "light_color")
	new_light.normal_light_pos_loc = rl.GetShaderLocation(new_light.normal_lighting_shader, "light_position")
	new_light.normal_light_color_loc = rl.GetShaderLocation(new_light.normal_lighting_shader, "light_color")
	new_light.normal_light_rot_loc = rl.GetShaderLocation(new_light.normal_lighting_shader, "rotation_angle")
	new_light.normal_light_origin_loc = rl.GetShaderLocation(new_light.normal_lighting_shader, "self_origin")
	new_light.normal_light_sprite_pos_loc = rl.GetShaderLocation(new_light.normal_lighting_shader, "self_position")
	new_light.normal_light_range_loc = rl.GetShaderLocation(new_light.normal_lighting_shader, "light_range")
	new_light.occlusion_shift_cam = rl.NewCamera2D(rl.Vector2Zero(), rl.NewVector2(size.X/4, size.Y/4), 0.0, 1.0)
	ls.lights = append(ls.lights, &new_light)
	return &new_light
}

func (l *Light2D) SetPosition(position rl.Vector2) {
	l.position = position
}

func (l *Light2D) SetMask(mask rl.Texture2D) {
	l.mask = mask
}

func (l *Light2D) SetMaskRotation(angle float32) {
	l.mask_rotation = angle
}

func (l *Light2D) DrawMask() {
	if l.mask != (rl.Texture2D{}) {
		// we move with the center offset to account for things like screenshake
		co := render.GetWSCameraCenterOffset()

		rl.DrawTexturePro(
			l.mask,
			rl.NewRectangle(0, 0, float32(l.mask.Width), float32(l.mask.Height)),
			rl.NewRectangle(
				l.position.X+float32(ls.lightmap_extend)/2-co.X,
				l.position.Y+float32(ls.lightmap_extend)/2-co.Y,
				float32(l.mask.Width), float32(l.mask.Height)),
			rl.NewVector2(float32(l.mask.Width)/2, float32(l.mask.Height)/2),
			l.mask_rotation, rl.White,
		)
	}
}

func (l *Light2D) Enable() {
	l.is_enabled = true
}

func (l *Light2D) Disable() {
	l.is_enabled = false
}

func DrawLight() {
	if ls.is_lighting_enabled {
		render.PauseTargetTex()

		// clean last frame's data from the lighting target
		rl.BeginTextureMode(ls.lighting_target)
		rl.ClearBackground(rl.Black)
		rl.EndTextureMode()

		rl.BeginTextureMode(ls.normal_collection_tex)
		rl.ClearBackground(rl.Blank)
		rl.EndTextureMode()

		for _, l := range ls.lights {
			if !l.is_enabled {
				continue
			}
			ws_light_pos := l.position
			//ot := render.Rs.CurrentStage.Camera.Target
			l.occlusion_shift_cam.Target = rl.Vector2Subtract(ws_light_pos, rl.Vector2Scale(l.size, 0.5))

			// draw to the occlusion texture
			rl.BeginTextureMode(l.occlusion_map)
			rl.BeginMode2D(l.occlusion_shift_cam)
			rl.ClearBackground(rl.Blank)

			// draw all the occluders
			for _, occ := range ls.occluders {
				occ.Draw()
			}

			rl.EndMode2D()
			rl.EndTextureMode()

			// do the polar transformation
			rl.BeginTextureMode(l.polar_shadowmap)
			rl.ClearBackground(rl.Blank)
			rl.BeginShaderMode(l.polar_transform_shader)
			rl.SetShaderValueV(l.polar_transform_shader, l.polar_transform_resolution_loc, []float32{l.resolution, l.resolution}, rl.ShaderUniformVec2, 1)
			rl.DrawTexturePro(
				l.occlusion_map.Texture,
				rl.NewRectangle(0, 0, l.size.X, -l.size.Y),
				rl.NewRectangle(0, 0, l.resolution, float32(l.occlusion_map.Texture.Height)),
				rl.NewVector2(0, 0),
				0.0,
				rl.White,
			)
			rl.EndShaderMode()
			rl.EndTextureMode()

			// add the shadows to the lighting target
			rl.BeginTextureMode(ls.lighting_target)
			rl.BeginShaderMode(l.light_render_shader)
			rl.BeginBlendMode(rl.BlendAdditive)
			rl.BeginMode2D(*render.Rs.CurrentStage.WSCamera)
			rl.SetShaderValueV(l.light_render_shader, l.light_render_resolution_loc, []float32{l.resolution, l.resolution}, rl.ShaderUniformVec2, 1)
			rl.SetShaderValueV(
				l.light_render_shader,
				l.light_render_light_color_loc,
				[]float32{
					float32(l.color.R) / 255,
					float32(l.color.G) / 255,
					float32(l.color.B) / 255,
				},
				rl.ShaderUniformVec3, 1)

			// we move with the center offset to account for things like screenshake
			co := render.GetWSCameraCenterOffset()

			rl.DrawTexturePro(
				l.polar_shadowmap.Texture,
				rl.NewRectangle(0, 0, l.resolution, 1),
				rl.NewRectangle(
					l.position.X+float32(ls.lightmap_extend)/2-co.X,
					l.position.Y+float32(ls.lightmap_extend)/2-co.Y,
					l.size.X, l.size.Y),
				rl.NewVector2(l.size.X/2, l.size.Y/2),
				0.0,
				rl.White,
			)
			rl.EndMode2D()
			rl.EndBlendMode()
			rl.EndShaderMode()
			rl.EndTextureMode()

			// NORMAL MAPS

			rl.BeginTextureMode(ls.normal_collection_tex)
			rl.BeginMode2D(*render.Rs.CurrentStage.WSCamera)
			rl.SetShaderValueV(
				l.normal_lighting_shader,
				l.normal_light_color_loc,
				[]float32{
					float32(l.color.R) / 255,
					float32(l.color.G) / 255,
					float32(l.color.B) / 255,
				},
				rl.ShaderUniformVec3, 1)
			// draw all the occluders
			for _, occ := range ls.occluders {
				rl.BeginShaderMode(l.normal_lighting_shader)
				light_ws_pos := l.position
				rl.SetShaderValueV(
					l.normal_lighting_shader,
					l.normal_light_pos_loc,
					[]float32{
						light_ws_pos.X,
						light_ws_pos.Y,
					},
					rl.ShaderUniformVec2, 1)
				rl.SetShaderValue(
					l.normal_lighting_shader,
					l.normal_light_rot_loc,
					[]float32{occ.GetRotation() * rl.Deg2rad},
					rl.ShaderUniformFloat)
				rl.SetShaderValueV(
					l.normal_lighting_shader,
					l.normal_light_origin_loc,
					[]float32{occ.GetOrigin().X, occ.GetOrigin().Y},
					rl.ShaderUniformVec2, 1)
				rl.SetShaderValueV(
					l.normal_lighting_shader,
					l.normal_light_sprite_pos_loc,
					[]float32{occ.GetPosition().X, occ.GetPosition().Y},
					rl.ShaderUniformVec2, 1)
				rl.SetShaderValueV(
					l.normal_lighting_shader,
					l.normal_light_range_loc,
					[]float32{l.size.X * 0.4, l.size.Y * 0.4, 100},
					rl.ShaderUniformVec3, 1)
				occ.DrawNormal()
				rl.EndShaderMode()
			}
			rl.EndMode2D()
			rl.EndTextureMode()

			rl.BeginTextureMode(ls.lighting_target)
			rl.DrawTexturePro(
				ls.normal_collection_tex.Texture,
				rl.NewRectangle(0, 0, float32(ls.normal_collection_tex.Texture.Width), -float32(ls.normal_collection_tex.Texture.Height)),
				rl.NewRectangle(float32(ls.lightmap_extend)/2-co.X, float32(ls.lightmap_extend)/2-co.Y, float32(ls.lighting_target.Texture.Width), float32(ls.lighting_target.Texture.Height)),
				rl.Vector2Zero(), 0.0, rl.White,
			)

			// multiply the mask over this light layer
			rl.BeginBlendMode(rl.BlendMultiplied)
			rl.BeginMode2D(*render.Rs.CurrentStage.WSCamera)
			l.DrawMask()
			rl.EndMode2D()
			rl.EndBlendMode()
			rl.EndTextureMode()
		}

		// draw the shadows to the normal rendertex
		// draw all the occluders
		render.ContinueTargetTex()
		rl.BeginBlendMode(rl.BlendMultiplied)
		rl.BeginShaderMode(ls.lightmap_blur_shader)
		rl.SetShaderValue(
			ls.lightmap_blur_shader,
			ls.blur_shader_ambient_level_loc,
			[]float32{ls.ambient_light_level},
			rl.ShaderUniformFloat,
		)
		// we need to pass worldspace coords, but want the light texture to stick to the screen
		p := render.ScreenToWorldPoint(rl.NewVector2(-float32(ls.lightmap_extend)/2, -float32(ls.lightmap_extend)/2))
		rl.DrawTexturePro(
			ls.lighting_target.Texture,
			rl.NewRectangle(0, 0, float32(ls.lighting_target.Texture.Width), -float32(ls.lighting_target.Texture.Height)),
			rl.NewRectangle(p.X, p.Y, float32(ls.lighting_target.Texture.Width), float32(ls.lighting_target.Texture.Height)),
			rl.Vector2Zero(),
			0.0,
			rl.White,
		)
		rl.EndShaderMode()
		rl.EndBlendMode()

		// UNCOMMENT to draw the occluders to the screen for debugging
		//for _, occ := range ls.occluders {
		//    occ.Draw()
		//}
	}
}
