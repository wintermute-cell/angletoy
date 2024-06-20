package physics

//
// This file contains helper functions internal to the physics system.
//

import rl "github.com/gen2brain/raylib-go/raylib"


// pixelToSimulationScaleV converts a vector in pixel scale to a vector in
// simulation scale.
func pixelToSimulationScaleV(vector rl.Vector2) rl.Vector2 {
    return rl.Vector2Scale(vector, float32(State.simulationScale))
}

// simulationToPixelScaleV converts a vector in simulation scale to a vector in
// pixel scale.
func simulationToPixelScaleV(vector rl.Vector2) rl.Vector2 {
    return rl.Vector2Scale(vector, float32(1.0 / State.simulationScale))
}

type number interface {
	int | int16 | int32 | uint | uint16 | uint32 | float32 | float64
}

// pixelToSimulationScale converts a number in pixel scale to a number in
// simulation scale.
func pixelToSimulationScale[T number](x T) T {
    return T(float64(x) * State.simulationScale)
}

// simulationToPixelScale converts a number in simulation scale to a number in
// pixel scale.
func simulationToPixelScale[T number](x T) T {
    return T(float64(x) / State.simulationScale)
}
