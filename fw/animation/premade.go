package animation

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type premade_anim_state struct {
	camera_shake_offset    rl.Vector2
	camera_bounce_offset   rl.Vector2
	camera_shake_anim      *Animation[float32]
	camera_bounce_anim     *Animation[float32]
	ws_camera              *rl.Camera2D
	original_ws_cam_offset rl.Vector2
}

var pas premade_anim_state

func InitPremades(worldspace_cam *rl.Camera2D, orig_cam_offset rl.Vector2) {
	pas = premade_anim_state{
		ws_camera:              worldspace_cam,
		original_ws_cam_offset: orig_cam_offset,
	}
}

func UpdatePremades() {
	shakeX, shakeY, bounceX, bounceY := float32(0), float32(0), float32(0), float32(0)

	if pas.camera_shake_anim != nil {
		pas.camera_shake_anim.Update()
		shakeX = pas.camera_shake_offset.X
		shakeY = pas.camera_shake_offset.Y
	}

	if pas.camera_bounce_anim != nil {
		pas.camera_bounce_anim.Update()
		bounceX = pas.camera_bounce_offset.X
		bounceY = pas.camera_bounce_offset.Y
	}

	// Sum the computed offsets from both animations and apply to camera
	pas.ws_camera.Offset.X = pas.original_ws_cam_offset.X + shakeX + bounceX
	pas.ws_camera.Offset.Y = pas.original_ws_cam_offset.Y + shakeY + bounceY
}

// CameraShake will shake the camera with the given intensity
func CameraShake(intensity float32) {
	pas.camera_shake_anim = CreateAnimation[float32](0.3)

	pas.camera_shake_anim.AddKeyframe(&pas.camera_shake_offset.X, 0.0, 0)
	pas.camera_shake_anim.AddKeyframe(&pas.camera_shake_offset.Y, 0.0, 0)

	pas.camera_shake_anim.AddKeyframe(&pas.camera_shake_offset.X, 0.05, 0+intensity*0.7)
	pas.camera_shake_anim.AddKeyframe(&pas.camera_shake_offset.Y, 0.05, 0+intensity*0.7)

	pas.camera_shake_anim.AddKeyframe(&pas.camera_shake_offset.X, 0.1, 0-intensity*0.7)
	pas.camera_shake_anim.AddKeyframe(&pas.camera_shake_offset.Y, 0.1, 0-intensity*0.7)

	pas.camera_shake_anim.AddKeyframe(&pas.camera_shake_offset.X, 0.15, 0+intensity*0.5)
	pas.camera_shake_anim.AddKeyframe(&pas.camera_shake_offset.Y, 0.15, 0-intensity*0.5)

	pas.camera_shake_anim.AddKeyframe(&pas.camera_shake_offset.X, 0.2, 0-intensity*0.5)
	pas.camera_shake_anim.AddKeyframe(&pas.camera_shake_offset.Y, 0.2, 0+intensity*0.5)

	pas.camera_shake_anim.AddKeyframe(&pas.camera_shake_offset.X, 0.28, 0)
	pas.camera_shake_anim.AddKeyframe(&pas.camera_shake_offset.Y, 0.28, 0)
	pas.camera_shake_anim.AddKeyframe(&pas.camera_shake_offset.X, 0.3, 0)
	pas.camera_shake_anim.AddKeyframe(&pas.camera_shake_offset.Y, 0.3, 0)

	pas.camera_shake_anim.Play(false, false)
}

func CameraBounce(direction rl.Vector2, intensity float32) {
	pas.camera_bounce_anim = CreateAnimation[float32](0.30)

	pas.camera_bounce_anim.AddKeyframe(&pas.camera_bounce_offset.X, 0.0, 0)
	pas.camera_bounce_anim.AddKeyframe(&pas.camera_bounce_offset.X, 0.02, 0+direction.X*intensity*10.0)
	pas.camera_bounce_anim.AddKeyframe(&pas.camera_bounce_offset.X, 0.05, 0+direction.X*intensity*8.0)
	pas.camera_bounce_anim.AddKeyframe(&pas.camera_bounce_offset.X, 0.3, 0)

	pas.camera_bounce_anim.AddKeyframe(&pas.camera_bounce_offset.Y, 0.0, 0)
	pas.camera_bounce_anim.AddKeyframe(&pas.camera_bounce_offset.Y, 0.02, 0+direction.Y*intensity*10.0)
	pas.camera_bounce_anim.AddKeyframe(&pas.camera_bounce_offset.Y, 0.05, 0+direction.Y*intensity*8.0)
	pas.camera_bounce_anim.AddKeyframe(&pas.camera_bounce_offset.Y, 0.3, 0)

	pas.camera_bounce_anim.Play(false, false)
}
