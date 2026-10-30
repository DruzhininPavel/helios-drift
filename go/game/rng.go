package game

import "math/rand"

// The game keeps its own random source so a run can be reproduced exactly:
// the same seed and the same inputs must produce the same world.
var (
	rng           = rand.New(rand.NewSource(1))
	rngSeed int64 = 1
)

// SetSeed reseeds the world generator. Every run that must be reproducible
// calls this before the first frame.
func SetSeed(seed int64) {
	rngSeed = seed
	rng = rand.New(rand.NewSource(seed))
}

// Seed reports the seed in use, for logs and replays.
func Seed() int64 { return rngSeed }

var rngCalls int

func randFloat() float64             { rngCalls++; return rng.Float64() }
func randInt(n int) int              { rngCalls++; return rng.Intn(n) }
func randRange(a, b float64) float64 { rngCalls++; return a + rng.Float64()*(b-a) }

// RngCalls reports how many numbers the world has drawn. Two runs that draw a
// different amount have already diverged, even when the visible state matches.
func RngCalls() int { return rngCalls }

// RandProbe returns the next number without moving the world counter; tests use
// it to compare two supposedly identical runs.
func RandProbe() float64 { return rng.Float64() }
