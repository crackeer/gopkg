package mapbuilder

// ResponseBuildClient response struct build client
type MapBuilder struct {
	Struct map[string]interface{}
}

const (
	nilString = "<nil>"
)

const (
	_sep = "/"
	_pre = "@"
)

const (
	_defaultMap    = "@map"
	_defaultString = "@string"
	_defaultArray  = "@array"
)
