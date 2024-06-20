package game

import (
	"gorl/fw/core/store"
	"gorl/fw/modules/scenes"
	gscenes "gorl/game/scenes"
)

type ControlState struct {
	SliderVal float32
}

func Init() {
	cs := ControlState{}
	store.Add(cs)

	scenes.RegisterScene("VFH", &gscenes.VfhScene{})
	scenes.RegisterScene("Angles", &gscenes.AnglesScene{})

	scenes.EnableScene("Angles")
}
