# Module: render

TODO: this is a draft, finalize

## What are its tasks
- drawing in stages, at different resolutions
- manage cameras per render stage
- manage rendering at different resolution while correcting for scale
    - camera zoom is managed, real zoom is hidden

## How does it work

Consists of:
- one render system with a screen resolution
- multiple render stages with their own resolution

each render stage draws at its resolution to its texture and is then flushed to
the render systems internal texture, when the next render stage starts or the
render system is flushed.

at the end, the render systems texture is drawn to the screen.

## Usage Example

Create a render system and any number of render stages:
```go
renderSystem := render.NewRenderSystem(rl.NewVector2(
    float32(settings.CurrentSettings().ScreenWidth),
    float32(settings.CurrentSettings().ScreenHeight)))

// renders at default resolution
defaultRenderStage := render.NewRenderStage(rl.NewVector2(
    float32(settings.CurrentSettings().RenderWidth),
    float32(settings.CurrentSettings().RenderHeight)),
    1)

// renders at double resolution
// we want this stage to behave as if it rendered at 1x resolution in terms of
// positioning/size, just with a sharper image.
doubleResRenderStage := render.NewRenderStage(rl.NewVector2(
    float32(settings.CurrentSettings().RenderWidth*2),
    float32(settings.CurrentSettings().RenderHeight*2)),
    2) // to correct back to 1x behaviour, we have to apply a correction factor.
```

Draw things into the render stages.
```go
for !rl.WindowShouldClose() {

    // inside the standard raylib draw block...
    rl.BeginDrawing()

    // first we clear the screen.
    rl.ClearBackground(rl.RayWhite)

    // then we enable one of the render stages we defined earlier.
    renderSystem.EnableRenderStage(defaultRenderStage)
    rl.ClearBackground(rl.Blank)
    rl.DrawCircleV(rl.NewVector2(100, 100), 50, rl.Red) // this circle will render at default resolution.

    // enabling another render stage automatically flushes the last one.
    renderSystem.EnableRenderStage(doubleResRenderStage)
    rl.ClearBackground(rl.Blank)
    rl.DrawCircleV(rl.NewVector2(100, 100), 50, rl.Green) // this circle will render at double resolution, so it will appear half as big!

    // Don't forget to flush the render system at the end, otherwise the last
    // render stage will be discarded!
    renderSystem.FlushRenderStage()

    // render the render systems internal texture to the screen.
    renderSystem.RenderToScreen()

    rl.EndDrawing()

}
```

## Controlling the Camera
Each render stage has its own camera. We can control it like so:
```go
if rl.IsKeyDown(rl.KeyJ) {
    curretCameraTarget := defaultRenderStage.GetCameraTarget()
    defaultRenderStage.SetCameraTarget(
        rl.Vector2Add(curretCameraTarget, rl.NewVector2(0, 10))
    )
}
```


## Debugging
To aid in debugging, we can draw a widget that visualizes all the stage
viewports like so:
```go
render.DebugDrawStageViewports(
    rl.NewVector2(10, 10), // position
    4, // scale of the debug widget, 4 means 1/4 of the original screen size
    renderSystem, // we pass our render system
    []*render.RenderStage{defaultRenderStage, doubleResRenderStage}, // we list all stages that need to be displayed
)
```
