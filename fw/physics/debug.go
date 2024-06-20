package physics

import (
	"github.com/ByteArena/box2d"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Deprecated: does not work with new cameras / rendering.
// Instead use myCollder.GetVectices() and draw them yourself.
func DrawColliders(draw_shapes, draw_joints, draw_bounding bool) {
	if draw_shapes { // shape
		for b := State.physicsWorld.GetBodyList(); b != nil; b = b.GetNext() {
			xf := b.GetTransform()
			for f := b.GetFixtureList(); f != nil; f = f.GetNext() {
				drawShape(f, xf)
			}
		}
	}

	//if draw_joints { // joint
	//	for j := State.physicsWorld.GetJointList(); j != nil; j = j.GetNext() {
	//		drawJoint(j)
	//	}
	//}

	if draw_bounding { // bounding-box
		bp := &State.physicsWorld.M_contactManager.M_broadPhase
		for b := State.physicsWorld.GetBodyList(); b != nil; b = b.GetNext() {
			if !b.IsActive() {
				continue
			}
			for f := b.GetFixtureList(); f != nil; f = f.GetNext() {
				for i := 0; i < f.M_proxyCount; i++ {
					proxy := f.M_proxies[i]
					aabb := bp.GetFatAABB(proxy.ProxyId)
					vs := [4]rl.Vector2{}
					vs[0] = simulationToPixelScaleV(rl.Vector2{X: float32(aabb.LowerBound.X), Y: float32(aabb.LowerBound.Y)})
					vs[1] = simulationToPixelScaleV(rl.Vector2{X: float32(aabb.UpperBound.X), Y: float32(aabb.LowerBound.Y)})
					vs[2] = simulationToPixelScaleV(rl.Vector2{X: float32(aabb.UpperBound.X), Y: float32(aabb.UpperBound.Y)})
					vs[3] = simulationToPixelScaleV(rl.Vector2{X: float32(aabb.LowerBound.X), Y: float32(aabb.UpperBound.Y)})
					drawPolygon(vs[:])
				}
			}
		}
	}
}

func drawShape(fixture *box2d.B2Fixture, transform box2d.B2Transform) {
	if fixture.GetType() == box2d.B2Shape_Type.E_circle {
		pos := rl.NewVector2(float32(transform.P.X), float32(transform.P.Y))
		pos = simulationToPixelScaleV(pos)
		radius := fixture.GetShape().GetRadius()
		radius = simulationToPixelScale(radius)
		rl.DrawCircleV(pos, float32(radius), rl.Red)
	} else if fixture.GetType() == box2d.B2Shape_Type.E_polygon {
		polygonShape := fixture.GetShape().(*box2d.B2PolygonShape) // Cast to specific shape type
		vertexCount := polygonShape.M_count

		// Prepare a slice to hold our transformed vertices
		vs := make([]rl.Vector2, vertexCount)

		for i := 0; i < vertexCount; i++ {
			// Transform the vertex
			// TODO: what is even going on here
			transformedVertex := polygonShape.M_vertices[i] //box2d.B2TransformVec2Mul(transform, polygonShape.M_vertices[i])

			// Convert to raylib Vector2 and scale
			vpos := rl.NewVector2(float32(transformedVertex.X), float32(transformedVertex.Y))
			pos := rl.NewVector2(float32(transform.P.X), float32(transform.P.Y))
			vs[i] = simulationToPixelScaleV(rl.Vector2Add(pos, vpos))
		}

		// Draw the polygon using raylib
		drawPolygon(vs)
	}
}

// Draw a polygon in pixel coordinates
func drawPolygon(points []rl.Vector2) {
	for i, point := range points {
		if i < len(points)-1 {
			rl.DrawLineV(point, points[i+1], rl.Color{R: 255, G: 0, B: 0, A: 100})
		} else {
			rl.DrawLineV(point, points[0], rl.Color{R: 255, G: 0, B: 0, A: 100})
		}
	}
}
