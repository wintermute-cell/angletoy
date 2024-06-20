package collision

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func DistancePointToLine(point, lineStart, lineEnd rl.Vector2) (float32, rl.Vector2) {
	AB := rl.NewVector2(lineEnd.X-lineStart.X, lineEnd.Y-lineStart.Y)
	AP := rl.NewVector2(point.X-lineStart.X, point.Y-lineStart.Y)

	dotProduct := AP.X*AB.X + AP.Y*AB.Y
	lengthSquared := AB.X*AB.X + AB.Y*AB.Y

	t := dotProduct / lengthSquared
	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}

	closest := rl.NewVector2(
		lineStart.X+t*AB.X,
		lineStart.Y+t*AB.Y,
	)

	dx := point.X - closest.X
	dy := point.Y - closest.Y

	return float32(math.Sqrt(float64(dx*dx + dy*dy))), closest
}

func CheckCollisionLineCircle(line_start, line_end, circle_origin rl.Vector2, circle_radius float32) (bool, []rl.Vector2) {
	distance, closest := DistancePointToLine(circle_origin, line_start, line_end)
	intersects := distance <= circle_radius

	if !intersects {
		return false, nil
	}

	h := math.Sqrt(float64(circle_radius*circle_radius) - float64(distance*distance))

	direction := rl.NewVector2(line_end.X-line_start.X, line_end.Y-line_start.Y)
	direction = rl.Vector2DivideV(
		direction,
		rl.NewVector2(
			float32(math.Sqrt(float64(direction.X*direction.X+direction.Y*direction.Y))),
			float32(math.Sqrt(float64(direction.X*direction.X+direction.Y*direction.Y)))))

	intersection1 := rl.NewVector2(
		closest.X+float32(h)*direction.X,
		closest.Y+float32(h)*direction.Y,
	)

	intersection2 := rl.NewVector2(
		closest.X-float32(h)*direction.X,
		closest.Y-float32(h)*direction.Y,
	)

	return true, []rl.Vector2{intersection1, intersection2}
}

func CheckCollisionLineRectangle(line_start, line_end rl.Vector2, rect rl.Rectangle) (bool, []rl.Vector2) {
	// Define the four edges of the rectangle
	edges := [4][2]rl.Vector2{
		{{X: rect.X, Y: rect.Y}, {X: rect.X + rect.Width, Y: rect.Y}},                             // Top
		{{X: rect.X + rect.Width, Y: rect.Y}, {X: rect.X + rect.Width, Y: rect.Y + rect.Height}},  // Right
		{{X: rect.X, Y: rect.Y + rect.Height}, {X: rect.X + rect.Width, Y: rect.Y + rect.Height}}, // Bottom
		{{X: rect.X, Y: rect.Y}, {X: rect.X, Y: rect.Y + rect.Height}},                            // Left
	}

	intersectionPoints := []rl.Vector2{}

	for _, edge := range edges {
		intersect, point := LineIntersection(line_start, line_end, edge[0], edge[1])
		if intersect {
			intersectionPoints = append(intersectionPoints, point)
		}
	}

	if len(intersectionPoints) == 0 {
		return false, nil
	}

	return true, intersectionPoints
}

// Helper function to calculate line intersection
func LineIntersection(p1, p2, q1, q2 rl.Vector2) (bool, rl.Vector2) {
	denom := (p1.X-p2.X)*(q1.Y-q2.Y) - (p1.Y-p2.Y)*(q1.X-q2.X)
	if denom == 0 {
		// Lines are parallel
		return false, rl.Vector2{}
	}

	t := ((p1.X-q1.X)*(q1.Y-q2.Y) - (p1.Y-q1.Y)*(q1.X-q2.X)) / denom
	u := -((p1.X-p2.X)*(p1.Y-q1.Y) - (p1.Y-p2.Y)*(p1.X-q1.X)) / denom

	if t >= 0 && t <= 1 && u >= 0 && u <= 1 {
		return true, rl.NewVector2(p1.X+t*(p2.X-p1.X), p1.Y+t*(p2.Y-p1.Y))
	}

	return false, rl.Vector2{}
}
