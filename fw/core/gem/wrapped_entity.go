package gem

import (
	"gorl/fw/core/entities"
	input "gorl/fw/core/input/input_handling"
	"gorl/fw/core/math"
	"gorl/fw/core/render"
)

var _ render.Drawable = &WrappedEntity{}

type WrappedEntity struct {
	entities.IEntity
	absTransform math.Transform2D
}

// ShouldDraw checks if the entity should be drawn based on its layer flags,
// enabled and visible properties.
func (d WrappedEntity) ShouldDraw(layerFlags math.BitFlag) bool {
	e := d.IEntity.IsEnabled()
	v := d.IEntity.IsVisible()
	f := d.IEntity.GetLayerFlags().IsAny(layerFlags)
	return e && v && f
}

// Draw draws the entity.
func (d WrappedEntity) Draw() {
	oldTransform := *d.IEntity.GetTransform() // save the entity's old *local* transform
	d.IEntity.SetTransform(d.absTransform)    // set the entity's transform to the new transform matrix
	d.IEntity.Draw()                          // draw the entity
	d.IEntity.SetTransform(oldTransform)      // restore the entity's old *local* transform
}

// GetEntity retrieves the wrapped entity.
func (d WrappedEntity) GetEntity() entities.IEntity {
	return d.IEntity
}

// AsDrawable returns the entity as a drawable.
func (d WrappedEntity) AsDrawable() render.Drawable {
	return d
}

// AsInputReceiver returns the entity as an input receiver.
func (d WrappedEntity) AsInputReceiver() input.InputReceiver {
	return d
}
