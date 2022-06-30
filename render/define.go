package render

const (
	_defaultPlaceholderJSData = "{{-JSData-}}"

	_defaultPlaceholderContent = "{{-Content-}}"

	_defaultPlaceholderTitle = "{{-Title-}}"
)

// Option
type Option struct {
	PlaceholderJSData  string
	PlaceholderContent string
	PlaceholderTitle   string
	InjectData         interface{}
	Title              string
}

// DefaultOption
//  @return *Option
func DefaultOption() *Option {
	return &Option{
		PlaceholderJSData:  _defaultPlaceholderJSData,
		PlaceholderContent: _defaultPlaceholderContent,
		PlaceholderTitle:   _defaultPlaceholderTitle,
	}
}
