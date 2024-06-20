// animation.go

package animation

import (
	"gorl/fw/core/logging"
	"math"
	"math/rand"
	"sort"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type number interface {
	~int32 | ~float32
}

type animatable interface {
	number | bool | string
}

type Keyframe[T animatable] struct {
	Time  float32
	Value T
}

type Animation[T animatable] struct {
	variables  map[*T][]Keyframe[T]
	duration   float32
	isLooping  bool
	playTime   float32
	isPlaying  bool
	isReversed bool
}

// Create a new animation object
func CreateAnimation[T animatable](duration float32) *Animation[T] {
	return &Animation[T]{
		variables: map[*T][]Keyframe[T]{},
		duration:  duration,
	}
}

// Add a keyframe to the animation
func (a *Animation[T]) AddKeyframe(variable *T, time float32, value T) {
	// check if variable already exists...
	if _, exists := a.variables[variable]; !exists {
		// ... if not, create empty keyframe slice
		a.variables[variable] = []Keyframe[T]{}
	}
	a.variables[variable] = append(a.variables[variable], Keyframe[T]{Time: time, Value: value})
	// ensure the new keyframe is in the correct time location
	sort.Slice(a.variables[variable], func(i, j int) bool {
		return a.variables[variable][i].Time < a.variables[variable][j].Time
	})
}

// Play the animation, optionally with a random time offset (uselful for
// multiple entities with the same animation)
func (a *Animation[T]) Play(loop bool, random_time_offset bool) {
	a.isLooping = loop
	a.playTime = 0
	a.isPlaying = true

	if random_time_offset {
		offs := rand.Float32() * a.duration
		a.playTime += offs
	}
}

// ToggleReverse inverts the playback direction
func (a *Animation[T]) ToggleReverse() {
	a.isReversed = !a.isReversed
}

// SetReverse sets the playback direction. True = reverse, False = forward
func (a *Animation[T]) SetReverse(reverse bool) {
	a.isReversed = reverse
}

// Return whether the animation is currently playing
func (a *Animation[T]) IsPlaying() bool {
	return a.isPlaying
}

// Return the current time of the animation
func (a *Animation[T]) GetTime() float32 {
	return a.playTime
}

// Pause the animation
func (a *Animation[T]) Pause() {
	a.isPlaying = false
}

// Stop the animation, resetting its playtime
func (a *Animation[T]) Stop() {
	a.isPlaying = false
	a.playTime = 0
}

// Continue the animation after it has been paused
func (a *Animation[T]) Continue() {
	a.isPlaying = true
}

// Set the current time of the animation
func (a *Animation[T]) SetTime(time float32) {
	a.playTime = time
}

func (a *Animation[T]) GetLength() float32 {
	return a.duration
}

// Update the animation, should be called every frame
func (a *Animation[T]) Update() {
	if !a.isPlaying {
		return
	}

	if a.isReversed {
		a.playTime -= rl.GetFrameTime()
		if a.playTime < 0 {
			if a.isLooping {
				a.playTime = a.duration
			} else {
				a.isPlaying = false
				a.playTime = 0
			}
		}
	} else {
		a.playTime += rl.GetFrameTime()
		if a.playTime > a.duration {
			if a.isLooping {
				a.playTime = 0
			} else {
				a.isPlaying = false
				a.playTime = a.duration
			}
		}
	}

	// interpolate animated variables
	for variable, keyframes := range a.variables {
		for i := 0; i < len(keyframes)-1; i++ {
			// we iterate over all the keyframes until we find the pair of
			// keyframes that the animation time is currently between.
			// In the following example we will stop with i = 1, keyframes[i] = B
			//          | A------B----------C |
			//                        ^
			//                    a.playTime
			if a.playTime >= keyframes[i].Time && a.playTime <= keyframes[i+1].Time {
				// we divide the time passed since B by the time from B to C,
				// to find out the interpolation ratio between B and C.
				r := (a.playTime - keyframes[i].Time) / (keyframes[i+1].Time - keyframes[i].Time)

				// Holy shit this is ugly, but I'm pretty damn certain there is
				// no better way...
				switch any(variable).(type) {
				case *int32:
					rhs := int32(any(keyframes[i+1].Value).(int32))
					lhs := int32(any(keyframes[i].Value).(int32))
					// bit ugly with all the casting, but were basically just
					// rounding the rhs-lhs difference up or down, around the middle
					var v any = (lhs + int32(math.Round(float64(r*float32(rhs-lhs)))))
					*variable = v.(T)
				case *float32:
					rhs := float32(any(keyframes[i+1].Value).(float32))
					lhs := float32(any(keyframes[i].Value).(float32))
					var v any = (lhs + float32(r*float32(rhs-lhs)))
					*variable = v.(T)
				// NOTE: bool and string might need some tweaking, so that the
				// value is switched in the middle, and not just the lhs is taken
				case *bool:
					lhs := bool(any(keyframes[i].Value).(bool))
					var v any = lhs
					*variable = v.(T)
				case *string:
					lhs := string(any(keyframes[i].Value).(string))
					var v any = lhs
					*variable = v.(T)
				default:
					logging.Warning("Failed to match animation type during interpolation!")
				}
				break
			}
		}
	}
}
