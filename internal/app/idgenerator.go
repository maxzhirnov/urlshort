package app

import (
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type RandIDGenerator struct {
	IDLen int
}

func NewRandIDGenerator(idLen int) *RandIDGenerator {
	return &RandIDGenerator{IDLen: idLen}
}

func (g *RandIDGenerator) Generate() string {
	if g.IDLen < 4 {
		g.IDLen = 4
	}

	seededRand := rand.New(
		rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, g.IDLen)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
