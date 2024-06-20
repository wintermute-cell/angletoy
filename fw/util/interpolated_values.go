package util

import rl "github.com/gen2brain/raylib-go/raylib"

type SmoothVector2 struct {
    value      rl.Vector2
    history    []rl.Vector2
    index      int32  // Add an index for the circular buffer
    smoothness int32
}

func NewSmoothVector2(value rl.Vector2, smoothness int32) SmoothVector2 {
    return SmoothVector2{
        value: value,
        history: make([]rl.Vector2, smoothness),  // Pre-allocate the buffer to its max size
        index: 0,
        smoothness: smoothness,
    }
}

func (i *SmoothVector2) SetValue(value rl.Vector2) {
    // Update the current index with the new value
    i.history[i.index] = value
    
    // Update the index, making sure it wraps around when reaching the buffer limit
    i.index = (i.index + 1) % int32(i.smoothness)

    // Calculate the average of the history buffer
    sum := rl.Vector2Zero()
    count := Min(int32(len(i.history)), i.smoothness)  // count will ensure we only consider the values we've added so far
    for _, val := range i.history {
        sum = rl.Vector2Add(sum, val)
    }
    i.value = rl.Vector2Scale(sum, 1.0/float32(count))
}

func (i *SmoothVector2) GetValue() rl.Vector2 {
    return i.value
}
