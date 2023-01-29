package router

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
	Mode     string `json:"mode"`
	Config   string `json:"config"`
	Response string `json:"response"`
}

// RouterFactory
type RouterFactory interface {
	GetRouterMeta(string) (*RouterMeta, error)
}
