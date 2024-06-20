package scenes

import (
	"gorl/fw/core/gem"
	"gorl/fw/core/logging"
	"gorl/fw/util"
)

type sceneManager struct {
	scenes         map[string]IScene
	enabled_scenes map[string]bool
}

// Create a new SceneManager. A SceneManager will automatically take care of
// your Scenes (calling their Init(), Deinit(), Draw(), DrawGUI() functions).
func newSceneManager() *sceneManager {
	return &sceneManager{
		scenes:         make(map[string]IScene),
		enabled_scenes: make(map[string]bool),
	}
}

// The global instance of the SceneManager
var sm *sceneManager = newSceneManager()

// Register a scene with the SceneManager for automatic control
func RegisterScene(name string, scene IScene) {
	if _, exists := sm.scenes[name]; exists {
		logging.Fatal("A scene with name \"%v\" is already registered.", name)
	}
	sm.scenes[name] = scene
}

// Enable the Scene. The Scenes Init() function will be called.
func EnableScene(name string) {
	scene, exists := sm.scenes[name]
	if !exists {
		logging.Fatal("Scene with name \"%v\" not found.", name)
	}

	// Initialize the scene if it's not already enabled
	if !sm.enabled_scenes[name] {
		gem.Append(gem.GetRoot(), scene.GetRoot())
		scene.Init()
		sm.enabled_scenes[name] = true
	}
}

// Disable the Scene. The Scenes Deinit() function will be called.
func DisableScene(name string) {
	scene, exists := sm.scenes[name]
	if !exists {
		logging.Fatal("Scene with name \"%v\" not found.", name)
	}

	// De-initialize the scene if it's currently enabled
	if sm.enabled_scenes[name] {
		scene.Deinit()
		gem.Remove(scene.GetRoot())
		sm.enabled_scenes[name] = false
	}
}

// Disable all Scenes that are currently enabled.
func DisableAllScenes() {
	for name, _ := range sm.scenes {
		if sm.enabled_scenes[name] {
			sm.scenes[name].Deinit()
			gem.Remove(sm.scenes[name].GetRoot())
			sm.enabled_scenes[name] = false
		}
	}
}

// Disable all Scenes that are currently enabled, except for the ones specified
// by name in the `exception_slice` parameter.
func DisableAllScenesExcept(exception_slice []string) {
	for name, _ := range sm.scenes {
		if sm.enabled_scenes[name] && !util.SliceContains(exception_slice, name) {
			sm.scenes[name].Deinit()
			sm.enabled_scenes[name] = false
		}
	}
}
