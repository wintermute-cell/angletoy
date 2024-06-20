package physics

import (
    "github.com/ByteArena/box2d"
	rl "github.com/gen2brain/raylib-go/raylib"
)

/*
*
* This file contains wrapper functions around box2d functions. These wrappers
* calculate the given positions between pixel and simulation scale. They also
* provide a more convenient interface for the user, only taking and returning
* usually used types in this framework, instead of box2d types.
*
 */

// SetDensity sets the density of every fixture attacked to the given collider.
func SetDensity(collider *Collider, density float32) {
    for f := collider.GetB2Body().GetFixtureList(); f != nil; f = f.GetNext() {
        f.SetDensity(float64(density))
    }
}

// ApplyForce applies a force to the given collider at the given point.
func (col *Collider) ApplyForce(force, point rl.Vector2) {
    force = pixelToSimulationScaleV(force)
    point = pixelToSimulationScaleV(point)
    b2f := box2d.MakeB2Vec2(float64(force.X), float64(force.Y))
    b2p := box2d.MakeB2Vec2(float64(point.X), float64(point.Y))
    col.GetB2Body().ApplyForce(b2f, b2p, true)
}

// ApplyForceToCenter applies a force to the given collider at the center of
// mass of the given collider.
func (col *Collider) ApplyForceToCenter(force rl.Vector2) {
    force = pixelToSimulationScaleV(force)
    b2f := box2d.MakeB2Vec2(float64(force.X), float64(force.Y))
    col.GetB2Body().ApplyForceToCenter(b2f, true)
}

// ApplyLinearImpulse applies an impulse to the given collider at the given
// point.
func (col *Collider) ApplyLinearImpulse(impulse, point rl.Vector2) {
    impulse = pixelToSimulationScaleV(impulse)
    point = pixelToSimulationScaleV(point)
    b2i := box2d.MakeB2Vec2(float64(impulse.X), float64(impulse.Y))
    b2p := box2d.MakeB2Vec2(float64(point.X), float64(point.Y))
    col.GetB2Body().ApplyLinearImpulse(b2i, b2p, true)
}

// ApplyLinearImpulseToCenter applies an impulse to the given collider at the
// center of mass of the given collider.
func (col *Collider) ApplyLinearImpulseToCenter(impulse rl.Vector2) {
    impulse = pixelToSimulationScaleV(impulse)
    b2i := box2d.MakeB2Vec2(float64(impulse.X), float64(impulse.Y))
    col.GetB2Body().ApplyLinearImpulseToCenter(b2i, true)
}

// ApplyTorque applies a torque to the given collider.
func (col *Collider) ApplyTorque(torque float32) {
    col.GetB2Body().ApplyTorque(float64(torque), true)
}

// ApplyAngularImpulse applies an angular impulse to the given collider.
func (col *Collider) ApplyAngularImpulse(impulse float32) {
    col.GetB2Body().ApplyAngularImpulse(float64(impulse), true)
}

// SetLinearVelocity sets the linear velocity of the given collider.
func (col *Collider) SetLinearVelocity(velocity rl.Vector2) {
    velocity = pixelToSimulationScaleV(velocity)
    b2v := box2d.MakeB2Vec2(float64(velocity.X), float64(velocity.Y))
    col.GetB2Body().SetLinearVelocity(b2v)
}

// SetAngularVelocity sets the angular velocity of the given collider.
func (col *Collider) SetAngularVelocity(velocity float32) {
    velocity = pixelToSimulationScale(velocity)
    col.GetB2Body().SetAngularVelocity(float64(velocity))
}

// GetLinearVelocity returns the linear velocity of the given collider.
func (col *Collider) GetLinearVelocity() rl.Vector2 {
    b2v := col.GetB2Body().GetLinearVelocity()
    return simulationToPixelScaleV(rl.NewVector2(float32(b2v.X), float32(b2v.Y)))
}

// GetAngularVelocity returns the angular velocity of the given collider.
func (col *Collider) GetAngularVelocity() float32 {
    return float32(col.GetB2Body().GetAngularVelocity())
}

