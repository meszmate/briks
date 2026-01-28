package game

import "math/rand"

// Bag implements the 7-bag randomizer.
// Each bag contains one of each piece type, shuffled randomly.
type Bag struct {
	pieces []PieceType
}

// NewBag creates a new 7-bag randomizer.
func NewBag() *Bag {
	b := &Bag{}
	b.refill()
	return b
}

// Next returns the next piece type from the bag.
func (b *Bag) Next() PieceType {
	if len(b.pieces) == 0 {
		b.refill()
	}
	p := b.pieces[0]
	b.pieces = b.pieces[1:]
	return p
}

// Preview returns the next n piece types without consuming them.
func (b *Bag) Preview(n int) []PieceType {
	// Ensure we have enough pieces.
	for len(b.pieces) < n {
		b.refill()
	}
	result := make([]PieceType, n)
	copy(result, b.pieces[:n])
	return result
}

func (b *Bag) refill() {
	newBag := make([]PieceType, len(AllPieceTypes))
	copy(newBag, AllPieceTypes)
	rand.Shuffle(len(newBag), func(i, j int) {
		newBag[i], newBag[j] = newBag[j], newBag[i]
	})
	b.pieces = append(b.pieces, newBag...)
}
