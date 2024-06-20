package tool

import (
	"fmt"
	"image"
	_ "image/png"
	"io"
	"math"
	"os"
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
	tri "github.com/osuushi/triangulate"
	"github.com/spf13/cobra"
)

// EXAMPLE USAGE:
// go run cmd/tool/main.go triangulate_image ./vfh_map.png 0.1 | go run cmd/tool/main.go triangulate_image -p ./vfh_map.png .
// (The dot at the end is necessary since this script is crap and should be refactored)

// triangulate_imageCmd represents the triangulate_image command
var triangulate_imageCmd = &cobra.Command{
	Use:   "triangulate_image <image_path> <simplification_range>",
	Short: "Triangulate an image",
	Long:  `Triangulate the white blobs of a black and white image and print the resulting triangles to stdout.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("preview") {
			imgPath := args[0]
			simplificationRange, err := strconv.ParseFloat(args[1], 32) // recommended: 0.1
			if err != nil {
				fmt.Println("Please provide a valid simplification range, defaulting to 0.1")
				simplificationRange = 0.1
			}
			GenerateCollidersFromImage(imgPath, float32(simplificationRange))
		} else {
			PreviewPolygonsFromString(args)
		}
	},
}

func init() {
	rootCmd.AddCommand(triangulate_imageCmd)

	triangulate_imageCmd.Flags().BoolP("preview", "p", false, "Preview a list of triangles from stdin. Ignores the <simplification_range> argument.")
}

type point struct {
	X, Y int32
}

type polygon []point

// GenerateCollidersFromImage takes in a black and white image, and builds a
// set of trianges, covering all white shapes in the image.
// The results are then printed to stdout, in the following format:
//
//	Triangle 1          Triangle 2         ...
//
// [{X Y} {X Y} {X Y}] [{X Y} {X Y} {X Y}] ...
func GenerateCollidersFromImage(image_path string, simplification_range float32) {
	img, err := loadImage(image_path)
	if err != nil {
		panic(err)
	}

	binImg := thresholding(img, 128)
	blobs := findBlobs(binImg)
	for _, blob := range blobs {
		polygon := extractContour(blob)
		simplificationOne := mergeStraightLines(polygon)
		merge_iterations := 30
		simplificationTwo := simplificationOne
		for i := 0; i < merge_iterations; i++ {
			simplificationTwo = mergeCloseSubsequentPoints(simplificationTwo, simplification_range)
		}

		tri_points := []*tri.Point{}
		for _, point := range simplificationTwo {
			tri_points = append(tri_points, &tri.Point{
				X: float64(point.X),
				Y: float64(point.Y),
			})
		}
		//triangles, _ := tri.Triangulate(tri_points)
		//for _, triangle := range triangles {
		//	as_slice := []point{
		//		{X: int32(triangle.A.X), Y: int32(triangle.A.Y)},
		//		{X: int32(triangle.B.X), Y: int32(triangle.B.Y)},
		//		{X: int32(triangle.C.X), Y: int32(triangle.C.Y)},
		//	}
		//	fmt.Printf("%v", as_slice)
		//}
		fmt.Printf("%v", simplificationTwo)
	}
}

func loadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	return img, err
}

func thresholding(img image.Image, threshold uint32) [][]bool {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	binImg := make([][]bool, height)
	for y := 0; y < height; y++ {
		binImg[y] = make([]bool, width)
		for x := 0; x < width; x++ {
			_, _, b, _ := img.At(x, y).RGBA()
			binImg[y][x] = b >= threshold
		}
	}

	return binImg
}

func findBlobs(binImg [][]bool) [][]point {
	height := len(binImg)
	width := len(binImg[0])

	// visited keeps track of visited pixels.
	visited := make([][]bool, height)
	for i := range visited {
		visited[i] = make([]bool, width)
	}

	var blobs [][]point

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if binImg[y][x] && !visited[y][x] {
				var blob []point
				blob = floodFill(binImg, int32(x), int32(y), visited, blob)
				blobs = append(blobs, blob)
			}
		}
	}
	return blobs
}

func floodFill(binImg [][]bool, x, y int32, visited [][]bool, blob []point) []point {
	// Base case: if x, y is out of bound or the pixel is not white or it's already visited, return the blob.
	if x < 0 || y < 0 || int(x) >= len(binImg[0]) || int(y) >= len(binImg) || visited[y][x] || !binImg[y][x] {
		return blob
	}

	// Mark this pixel as visited.
	visited[y][x] = true

	// Add this pixel to the blob.
	blob = append(blob, point{X: x, Y: y})

	// Check its eight neighbors.
	blob = floodFill(binImg, x-1, y-1, visited, blob)
	blob = floodFill(binImg, x-1, y, visited, blob)
	blob = floodFill(binImg, x-1, y+1, visited, blob)
	blob = floodFill(binImg, x, y-1, visited, blob)
	blob = floodFill(binImg, x, y+1, visited, blob)
	blob = floodFill(binImg, x+1, y-1, visited, blob)
	blob = floodFill(binImg, x+1, y, visited, blob)
	blob = floodFill(binImg, x+1, y+1, visited, blob)

	return blob
}

// Stack implementation from
// https://gist.github.com/bemasher/1777766

type stack struct {
	top  *stack_element
	size int
}

type stack_element struct {
	value interface{} // All types satisfy the empty interface, so we can store anything here.
	next  *stack_element
}

// Return the stack's length
func (s *stack) Len() int {
	return s.size
}

// push a new element onto the stack
func (s *stack) push(value interface{}) {
	s.top = &stack_element{value, s.top}
	s.size++
}

// Remove the top element from the stack and return it's value
// If the stack is empty, return nil
func (s *stack) pop() (value interface{}) {
	if s.size > 0 {
		value, s.top = s.top.value, s.top.next
		s.size--
		return
	}
	return nil
}

func findLowestPoint(p []point) {
	m := 0
	for i := 1; i < len(p); i++ {
		//If lowest points are on the same line, take the rightmost point
		if (p[i].Y < p[m].Y) || ((p[i].Y == p[m].Y) && p[i].X > p[m].X) {
			m = i
		}
	}
	p[0], p[m] = p[m], p[0]
}

// findContour follows the edge of a blob and returns the contour.
func extractContour(blob []point) polygon {
	var contour polygon
	minX, minY, maxX, maxY := blob[0].X, blob[0].Y, blob[0].X, blob[0].Y

	// Get bounding box of the blob.
	for _, p := range blob {
		if p.X < minX {
			minX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}

	// Loop through the bounding box to find the starting point on the perimeter of the blob.
	start := point{X: -1, Y: -1}
OuterLoop:
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if isInBlob(point{X: x, Y: y}, blob) {
				start = point{X: x, Y: y}
				break OuterLoop
			}
		}
	}

	// Define directions (Up, Right, Down, Left) for marching.
	directions := []point{
		{X: 0, Y: -1}, // Up
		{X: 1, Y: 0},  // Right
		{X: 0, Y: 1},  // Down
		{X: -1, Y: 0}, // Left
	}

	// Start marching.
	// TODO: this breaks when there is a one pixel wide gap in the blob
	if start.X != -1 && start.Y != -1 { // Ensure a starting point was found.
		dir := 0 // Initial direction: Up
		cur := start

		// Continue marching until returning to the starting point.
		for {
			for i, point := range contour {
				if point == cur {
					contour = contour[:i]
					break
				}
			}
			contour = append(contour, cur) // Add current point to contour.

			// Check directions in CW order.
			for i := 0; i < 4; i++ {
				dir = (dir + 1) % 4 // Modulo to wrap around.
				nextPoint := point{
					X: cur.X + directions[dir].X,
					Y: cur.Y + directions[dir].Y,
				}

				if isInBlob(nextPoint, blob) {
					// Move to the next white pixel.
					cur = nextPoint
					break
				}
			}

			// here we subtract 1 to look "left" relative to the last direction
			// and 1 more, since we increment directly at the start of the loop
			dir = (dir - 2 + 4) % 4

			// Stop if returned to the starting point.
			if cur == start {
				break
			}
		}
	}

	return contour
}

func isInBlob(p point, blob []point) bool {
	for _, bp := range blob {
		if bp == p {
			return true
		}
	}
	return false
}

func mergeStraightLines(poly polygon) polygon {
	// If the polygon has less than 3 points, return as-is.
	if len(poly) < 3 {
		return poly
	}

	var simplifiedPoly polygon

	// Always keep the first point.
	simplifiedPoly = append(simplifiedPoly, poly[0])

	// Check for collinear points and skip them.
	for i := 1; i < len(poly)-1; i++ {
		if !isCollinear(poly[i-1], poly[i], poly[i+1]) {
			simplifiedPoly = append(simplifiedPoly, poly[i])
		}
	}

	// Always keep the last point.
	simplifiedPoly = append(simplifiedPoly, poly[len(poly)-1])

	return simplifiedPoly
}

// isCollinear checks if three points are collinear by using the cross product.
func isCollinear(p1, p2, p3 point) bool {
	ax, ay := p2.X-p1.X, p2.Y-p1.Y
	bx, by := p3.X-p2.X, p3.Y-p2.Y

	// Cross product = 0 for collinear points.
	return ax*by == ay*bx
}

func mergeCloseSubsequentPoints(poly polygon, distanceThreshold float32) polygon {
	if len(poly) < 2 {
		// No points to merge.
		return poly
	}

	mergedPoly := make(polygon, 0, len(poly))
	mergedPoly = append(mergedPoly, poly[0]) // Add the first point

	i := 1
	for i < len(poly)-1 {
		// If the two points are close enough to merge, create a midpoint and add it to mergedPoly.
		if shouldMerge(poly[i], poly[i+1], distanceThreshold) {
			mergedPoint := midpoint(poly[i], poly[i+1])
			mergedPoly = append(mergedPoly, mergedPoint)
			i++ // skip the next point because we've just merged it
		} else {
			// If should not merge, add the current point to mergedPoly.
			mergedPoly = append(mergedPoly, poly[i])
		}
		i++
	}

	// Append the last point if it wasn't part of a merge.
	if !shouldMerge(poly[len(poly)-2], poly[len(poly)-1], distanceThreshold) {
		mergedPoly = append(mergedPoly, poly[len(poly)-1])
	}

	return mergedPoly
}

func shouldMerge(p1, p2 point, threshold float32) bool {
	return euclideanDistance(p1, p2) <= threshold
}

func euclideanDistance(p1, p2 point) float32 {
	dx := float64(p2.X - p1.X)
	dy := float64(p2.Y - p1.Y)
	return float32(math.Sqrt(dx*dx + dy*dy))
}

func midpoint(p1, p2 point) point {
	return point{
		X: (p1.X + p2.X) / 2,
		Y: (p1.Y + p2.Y) / 2,
	}
}

// ============================================================================
// PREVIEW
// ============================================================================

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
				rl.DrawCircle(polygon[i].X, polygon[i].Y, 2, rl.Green)
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
				2, rl.Green)
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
	polyStr, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	drawPolygons(imagePath, string(polyStr))
}
