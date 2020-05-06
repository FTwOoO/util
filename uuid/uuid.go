package uuid

import "github.com/rogpeppe/fastuuid"

var (
	g *fastuuid.Generator
)

func init() {
	g = fastuuid.MustNewGenerator()
}

func NewUUID() string {
	return g.Hex128()
}
