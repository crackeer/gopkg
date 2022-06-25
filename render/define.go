package render

const (
	_defaultPlaceholderJSData = "{{-JSData-}}"

	_defaultPlaceholderContent = "{{-Content-}}"
)

// Option
type Option struct {
	PlaceholderJSData  string
	PlaceholderContent string
	InjectData         interface{}
}

// DefaultOption
//  @return *Option
func DefaultOption() *Option {
	return &Option{
		PlaceholderJSData:  _defaultPlaceholderJSData,
		PlaceholderContent: _defaultPlaceholderContent,
	}
}
