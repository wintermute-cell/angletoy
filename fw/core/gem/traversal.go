package gem

import (
	"gorl/fw/core/datastructures"
	"gorl/fw/core/entities"
	input "gorl/fw/core/input/input_handling"
	"gorl/fw/core/logging"
	"gorl/fw/core/math"
	"gorl/fw/core/render"
)

// GetAbsoluteTransform returns the absolute transform of the entity.
func GetAbsoluteTransform(entity entities.IEntity) math.Transform2D {
	entityNode, ok := gemInstance.nodeMap[entity]
	if !ok {
		logging.Error("Tried to get absolute transform for entity not existent in gem.")
		return math.Transform2DZero()
	}

	// step through parents until we arrive at the root.
	transformMat3 := entity.GetTransform().GenerateMatrix()
	for entityNode.parent != gemInstance.root && entityNode.parent != nil { // we do a nil check just for good measure, normally it should stop at the root.
		parentMat3 := entityNode.parent.entity.GetTransform().GenerateMatrix()
		transformMat3 = transformMat3.Multiply(parentMat3)
		entityNode = entityNode.parent
	}

	return math.NewTransform2DFromMatrix3(transformMat3)
}

// Traverse traverses through the entity graph, updating the entities.
// In the process, it produces a list of DrawableEntity objects.
func Traverse(withFixedUpdate bool) ([]render.Drawable, []input.InputReceiver) {

	root := gemInstance.root

	nodeStack := datastructures.NewStack[*gemNode](len(gemInstance.nodeMap))
	nodeStack.Push(root)

	transformStack := datastructures.NewStack[math.Matrix3](len(gemInstance.nodeMap))
	transformStack.Push(math.Matrix3Identity())

	drawables := make([]render.Drawable, 0, len(gemInstance.nodeMap)/2)
	inputReceivers := make([]input.InputReceiver, 0, len(gemInstance.nodeMap)/2)

	for !nodeStack.IsEmpty() {

		node, _ := nodeStack.Pop()
		tMat3, _ := transformStack.Pop()

		// if the entity is not enabled, skip it and its children
		if !node.entity.IsEnabled() {
			continue
		}

		// add the entity to the input receivers
		if node.entity.IsEnabled() {
			inputReceivers = append(inputReceivers, node.entity)
		}

		// Update the entity
		node.entity.Update()
		if withFixedUpdate {
			node.entity.FixedUpdate()
		}

		drawables = append(drawables, WrappedEntity{
			IEntity:      node.entity,
			absTransform: math.NewTransform2DFromMatrix3(tMat3),
		})

		for _, child := range node.children {
			nodeStack.Push(child)
			transformStack.Push( // we push M_child * M_stack
				child.entity.
					GetTransform().
					GenerateMatrix().
					Multiply(tMat3))
		}
	}

	return drawables, inputReceivers
}
