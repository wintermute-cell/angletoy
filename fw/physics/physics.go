package physics

import (
	"gorl/fw/core/logging"
	"gorl/fw/util"

	"github.com/ByteArena/box2d"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// ------------
//  PHYSICS
// ------------

type PhysicsState struct {
	timestep           float64
	velocityIterations int
	positionIterations int
	updateTimer        util.Timer

	physicsWorld     box2d.B2World
	destructionQueue []*box2d.B2Body

	// The physics world needs a factor to calculate between pixels and meters.
	// If your player is 32 pixels high and should be ~2m tall, the
	// simulationScale should be (1/16).
	simulationScale float64
}

var State PhysicsState

// ----------------
//  MAIN FUNCTIONS
// ----------------

// InitPhysics initializes the physics state
func InitPhysics(timestep float32, gravity rl.Vector2, simulationScale float32) {

	if simulationScale == 0 {
		logging.Error("Provided simulation scale is zero!")
	}
	if timestep == 0 {
		logging.Error("Provided timestep is zero!")
	}

	State = PhysicsState{
		timestep:           float64(timestep),
		velocityIterations: 8,
		positionIterations: 3,
		updateTimer:        *util.NewTimer(timestep),
		physicsWorld:       box2d.MakeB2World(box2d.MakeB2Vec2(float64(gravity.X), float64(gravity.Y))),
		simulationScale:    float64(simulationScale),
	}

	State.physicsWorld.SetContactListener(&ContactListener{})
}

// DeinitPhysics deinitializes the physics state
func DeinitPhysics() {
	State.physicsWorld.Destroy()
}

// Update the physics world. This must be called every frame, the fixed
// timestep is managed internally.
// Returns true if the physics world was updated, false otherwise.
func Update() bool {
	if !State.updateTimer.Check() {
		return false
	}

	State.physicsWorld.Step(State.timestep, State.velocityIterations, State.positionIterations)

	// remove all bodies queued for destruction. Destroying an object while the
	// physics world is updating (for example in a collision callback) causes a
	// crash, so we delay the destruction until the update is finished.
	State.destructionQueue = util.SliceRemoveDuplicate(State.destructionQueue)
	for _, body := range State.destructionQueue {
		State.physicsWorld.DestroyBody(body)
	}
	State.destructionQueue = []*box2d.B2Body{}
	return true
}

// ------------------
//  CONFIG FUNCTIONS
// ------------------

// SetGravity sets the gravity of the physics world
func SetGravity(gravity rl.Vector2) {
	State.physicsWorld.SetGravity(box2d.MakeB2Vec2(float64(gravity.X), float64(gravity.Y)))
}

// ------------------
//  GETTER FUNCTIONS
// ------------------

// GetTimestep returns the timestep of the physics world in seconds
func GetTimestep() float32 {
	if State.timestep == 0 {
		logging.Error("Tried to get the physics timestep before it was set!")
	}
	return float32(State.timestep)
}

// ------------------
//  CONTACT LISTENER
// ------------------
//
// This custom ContactListener allows us to call callbacks attached to the
// collider struct whenever a collision occurs.

var _ box2d.B2ContactListenerInterface = (*ContactListener)(nil)

type ContactListener struct{}

func (*ContactListener) BeginContact(contact box2d.B2ContactInterface) {
	fA := contact.GetFixtureA()
	fB := contact.GetFixtureB()
	colA := fA.GetBody().GetUserData().(*Collider)
	colB := fB.GetBody().GetUserData().(*Collider)
	if colA == nil || colB == nil {
		logging.Error("Missing collider in body userdata!")
		return
	}

	// If the collider has a callback registered for the category of the other
	// collider, we call that callback.
	for category, callbackFunc := range colA.callbacks {
		if uint16(category)&fB.GetFilterData().CategoryBits != 0 {
			callbackFunc()
		}
	}

	for category, callbackFunc := range colB.callbacks {
		if uint16(category)&fA.GetFilterData().CategoryBits != 0 {
			callbackFunc()
		}
	}
}
func (*ContactListener) EndContact(contact box2d.B2ContactInterface) {
	// Nothing to do here
}
func (*ContactListener) PreSolve(contact box2d.B2ContactInterface, oldManifold box2d.B2Manifold) {
	// Nothing to do here
}
func (*ContactListener) PostSolve(contact box2d.B2ContactInterface, impulse *box2d.B2ContactImpulse) {
	// Nothing to do here
}
