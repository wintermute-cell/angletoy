package main

import (
	"fmt"
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func parsePolygons(polyStr string) [][]point {
	var polygons [][]point
	// Extract individual polygons
	polyStrs := strings.Split(polyStr, "][")
	for _, polyStr := range polyStrs {
		polyStr = strings.Trim(polyStr, "[] ")
		pts := strings.Split(polyStr, "}")
		var polygon []point
		for _, ptStr := range pts {
			if ptStr == "" {
				continue
			}
			xy := strings.Fields(strings.Trim(ptStr, "{} "))
			x, _ := strconv.Atoi(xy[0])
			y, _ := strconv.Atoi(xy[1])
			polygon = append(polygon, point{int32(x), int32(y)})
		}
		polygons = append(polygons, polygon)
	}

	return polygons
}

func drawPolygons(imagePath string, polyStr string) {
	polygons := parsePolygons(polyStr)

	image := rl.LoadImage(imagePath)
	rl.InitWindow(1280, 720, "Polygon Preview")
	defer rl.CloseWindow()

	texture := rl.LoadTextureFromImage(image)

	rl.SetTargetFPS(60)

	camera := rl.NewCamera2D(rl.NewVector2(0, 0), rl.NewVector2(0, 0), 0, 1.0)

	for !rl.WindowShouldClose() {
		// Handle zooming
		if rl.GetMouseWheelMoveV().Y > 0 {
			camera.Zoom += 0.25
		} else if rl.GetMouseWheelMoveV().Y < 0 {
			camera.Zoom -= 0.25
			if camera.Zoom < 0.1 {
				camera.Zoom = 0.1
			}
		}

		// Handle panning
		if rl.IsMouseButtonDown(rl.MouseMiddleButton) {
			camera.Target.X -= float32(rl.GetMouseDelta().X) / camera.Zoom
			camera.Target.Y -= float32(rl.GetMouseDelta().Y) / camera.Zoom
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode2D(camera)
		rl.DrawTexture(texture, 0, 0, rl.White)

		for _, polygon := range polygons {
			// Draw polygon outline
			for i := 0; i < len(polygon)-1; i++ {
				rl.DrawCircle(polygon[i].X, polygon[i].Y, 0.1, rl.Green)
				rl.DrawLine(
					int32(polygon[i].X),
					int32(polygon[i].Y),
					int32(polygon[i+1].X),
					int32(polygon[i+1].Y),
					rl.Red,
				)
			}
			// Close the polygon by drawing a line between the last and first points
			rl.DrawCircle(
				polygon[len(polygon)-1].X,
				polygon[len(polygon)-1].Y,
				0.2, rl.Green)
			rl.DrawLine(
				int32(polygon[len(polygon)-1].X),
				int32(polygon[len(polygon)-1].Y),
				int32(polygon[0].X),
				int32(polygon[0].Y),
				rl.Red,
			)
		}

		rl.EndMode2D()
		rl.DrawText("Use Mouse Wheel to Zoom in/out and Middle Mouse Button to Pan", 10, 10, 10, rl.Gray)
		rl.EndDrawing()
	}

	rl.UnloadImage(image)     // Once image is in GPU texture memory, we can unload RAM
	rl.UnloadTexture(texture) // Unload texture from VRAM
}

func PreviewPolygonsFromString(args []string) {
	if len(args) != 2 {
		fmt.Println("Usage: preview-string <image-path> <polygon-string>")
		return
	}

	imagePath := args[0]
	polyStr := args[1]

	drawPolygons(imagePath, polyStr)
}
