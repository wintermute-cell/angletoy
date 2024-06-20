# Module: Scenes

A scene provides an environment where a set of entities are created, updated
and destroyed together. Multiple scenes may be active at the same time.

For example, you could have a scene for a `MainMenu`, one for the `GameView`
and one for the `DebugConsole`. You enable the debug console scene at the start
and never disable it. You then load the main menu scene. When the player
pressen "play", disable the main menu scene and enable the game view scene.

## Usage

First, you must define a new scene. The easiest way to do this is to use the gorl tool to create a new scene file:

```bash
go run cmd/tool/main.go new_scene <scene_name>
```

TODO: write this out proper, implement tool.
This just copies bla from bla and replaces bla...

In the new scene file, you should start adding entities:

```go
func (scn *MyScene) Init() {
	testEntity := entities.NewTestEntity("My Test Entity")
	gem.AddEntity(scn.GetRoot(), testEntity, gem.DefaultLayer)
}
```

Note that each scene has a root entity. This is the parent node to which we can
directly attach our entities. As long as an entity exists somewhere in the
subtree of an active scene, it will be Updated automatically every frame.

Drawing is not done by scenes. 
TODO: explain how drawing is done or link to the correct place.

Once you have created your scene definition, you must register the scene with
the scene manager. By doing that you create a single instance of you scene,
which can then be enabled and disabled at will. All of this can be done
directly from the main.go file, before the game loop begins.

```go
scenes.RegisterScene("some_name", &uscenes.MyScene{})
scenes.EnableScene("some_name")
scenes.DisableScene("some_name")
```

The main game loop should already be set up to call `scenes.UpdateScenes()` and
`scenes.FixedUpdateScenes()` for you each frame. Should you have modified the
main loop, make sure these two functions are properly called.

TODO: explain what functions can be overwritten like Update() and why and how



