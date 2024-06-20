package gem

import (
	"gorl/fw/core/entities"
	"gorl/fw/core/logging"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// gem represents the Global Entity Manager graph.
type gem struct {
	root    *gemNode
	nodeMap map[entities.IEntity]*gemNode
}

// gemNode represents a node in the Gem graph.
// It wraps a entities.IEntity and keeps track of its parent and children.
type gemNode struct {
	entity   entities.IEntity
	parent   *gemNode
	children []*gemNode
}

const DefaultLayer = 0

// gemInstance is the global Gem graph.
var gemInstance *gem

// Init initializes the global Gem graph.
// This should be called once at the start of the program.
func Init() {
	rootEntity := entities.NewEntity("root", rl.Vector2Zero(), 0, rl.Vector2One())
	gemInstance = &gem{
		root: &gemNode{
			entity:   rootEntity,
			parent:   nil,
			children: make([]*gemNode, 0),
		},
		nodeMap: make(map[entities.IEntity]*gemNode),
	}
	// self-map the root entity
	gemInstance.nodeMap[gemInstance.root.entity] = gemInstance.root
}

// GetRoot returns the root entity of the Gem graph.
func GetRoot() entities.IEntity {
	return gemInstance.root.entity
}

// Deinit deinitializes the global Gem graph, calling Deinit on all entities.
func Deinit() {
	// FIXME: causes segfault, fix
	Remove(gemInstance.root.entity)
}

// Append adds a entities.IEntity to the Gem graph, as a child of the given parent.
func Append(parent, entity entities.IEntity) {
	parentNode, ok := gemInstance.nodeMap[parent]
	if !ok {
		logging.Error("Parent not found in graph, can't add child")
		return
	}

	node := &gemNode{
		entity:   entity,
		parent:   parentNode,
		children: make([]*gemNode, 0),
	}
	parentNode.children = append(parentNode.children, node)
	gemInstance.nodeMap[entity] = node

	entity.Init()
}

// Remove removes a entities.IEntity from the graph.
// All children of the removed entity are also removed.
func Remove(entity entities.IEntity) {
	node, ok := gemInstance.nodeMap[entity]
	if !ok {
		logging.Error("entity not found in graph, can't remove")
		return
	}

	// recursively remove all children. This will call Deinit on all children.
	for _, child := range node.children {
		Remove(child.entity)
	}

	// remove the node from the parent's children
	parent := node.parent
	for i, child := range parent.children {
		if child == node {
			parent.children = append(parent.children[:i], parent.children[i+1:]...)
			break
		}
	}

	// remove the node from the node map
	delete(gemInstance.nodeMap, entity)

	entity.Deinit()
}

// ReParent changes the parent of a entities.IEntity.
func ReParent(entity, newParent entities.IEntity) {
	node, ok := gemInstance.nodeMap[entity]
	if !ok {
		logging.Error("entity not found in graph, can't reparent")
		return
	}

	newParentNode, ok := gemInstance.nodeMap[newParent]
	if !ok {
		logging.Error("New parent not found in graph, can't reparent")
		return
	}

	// remove the node from the old parent's children
	oldParent := node.parent
	for i, child := range oldParent.children {
		if child == node {
			oldParent.children = append(oldParent.children[:i], oldParent.children[i+1:]...)
			break
		}
	}

	// add the node to the new parent's children
	newParentNode.children = append(newParentNode.children, node)
	node.parent = newParentNode
}

// GetChildren returns the children of a entities.IEntity.
func GetChildren(entity entities.IEntity) []entities.IEntity {
	node, ok := gemInstance.nodeMap[entity]
	if !ok {
		logging.Error("entity not found in graph, can't get children")
		return nil
	}

	children := make([]entities.IEntity, len(node.children))
	for i, child := range node.children {
		children[i] = child.entity
	}
	return children
}

// GetParent returns the parent of a entities.IEntity.
func GetParent(entity entities.IEntity) entities.IEntity {
	node, ok := gemInstance.nodeMap[entity]
	if !ok {
		logging.Error("entity not found in graph, can't get parent")
		return nil
	}

	return node.parent.entity
}
