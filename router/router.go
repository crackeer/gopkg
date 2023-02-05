package router

import "github.com/crackeer/gopkg/router/api"

const (

	// ModeRelay
	ModeRelay = "relay"

	// ModeMesh
	ModeMesh = "mesh"

	// ModeStatic
	ModeStatic = "static"
)

// RouterMeta
type RouterMeta struct {
	Mode       string               `json:"mode"`
	RelayAPI   string               `json:"relay_api"`
	MeshConfig [][]*api.RequestItem `json:"mesh_config"`
	Response   interface{}          `json:"response"`
}

// RouterFactory
type RouterFactory interface {
	Get(string) *RouterMeta
	LoadAll() error
}
