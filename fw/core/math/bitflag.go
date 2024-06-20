package math

import "fmt"

// BitFlag represents a 64-bit flag.
// It is used to store multiple boolean values in a single int64.
type BitFlag int64

// Set sets the given bits in the BitFlag.
func (b *BitFlag) Set(flags BitFlag) {
	*b |= flags
}

// Clear clears the given bits in the BitFlag.
func (b *BitFlag) Clear(flags BitFlag) {
	*b &= ^flags
}

// Toggle toggles the given bits in the BitFlag.
func (b *BitFlag) Toggle(flags BitFlag) {
	*b ^= flags
}

// Has checks if the given bits are set in the BitFlag.
func (b BitFlag) Has(flags BitFlag) bool {
	return b&flags == flags
}

// IsAny checks if any of the given bits are set in the BitFlag.
func (b BitFlag) IsAny(flags BitFlag) bool {
	return b&flags != 0
}

// CreateBitFlag creates a new BitFlag with the given bits set.
func CreateBitFlag(flags ...int) BitFlag {
	var result BitFlag
	for _, flag := range flags {
		result |= 1 << flag
	}
	return result
}

// FromInt64 creates a BitFlag from an int64 value.
func FromInt64(value int64) BitFlag {
	return BitFlag(value)
}

// ToInt64 converts the BitFlag to an int64 value.
func (b BitFlag) ToInt64() int64 {
	return int64(b)
}

// String returns the binary representation of the BitFlag.
func (b BitFlag) String() string {
	return fmt.Sprintf("%064b", b)
}

// Example flags.
// May be used directly, renamed/aliased or just as an example.
const (
	Flag0 BitFlag = 1 << iota
	Flag1
	Flag2
	Flag3
	Flag4
	Flag5
	Flag6
	Flag7
	Flag8
	Flag9
	Flag10
	Flag11
	Flag12
	Flag13
	Flag14
	Flag15
	Flag16
	Flag17
	Flag18
	Flag19
	Flag20
	Flag21
	Flag22
	Flag23
	Flag24
	Flag25
	Flag26
	Flag27
	Flag28
	Flag29
	Flag30
	Flag31
	Flag32
	Flag33
	Flag34
	Flag35
	Flag36
	Flag37
	Flag38
	Flag39
	Flag40
	Flag41
	Flag42
	Flag43
	Flag44
	Flag45
	Flag46
	Flag47
	Flag48
	Flag49
	Flag50
	Flag51
	Flag52
	Flag53
	Flag54
	Flag55
	Flag56
	Flag57
	Flag58
	Flag59
	Flag60
	Flag61
	Flag62
)
