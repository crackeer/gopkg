package router

// RouterExecuter
type RouterExecuter struct {
	*RouterMeta
	Env         string
	Header      map[string]string
	ErrorCode   int64
	ResponseRaw string
	Error       string
}

// NewRouterExecuter
//
//	@param meta
//	@param header
//	@return *RouterExecuter
func NewRouterExecuter(meta *RouterMeta, env string, header map[string]string) *RouterExecuter {
	return &RouterExecuter{
		Env:        env,
		RouterMeta: meta,
		Header:     header,
	}
}

func (e *RouterExecuter) Exec(params map[string]interface{}) (map[string]interface{}, error) {
	if e.Mode == ModeRelay {
		return e.doRelay(params)
	}
	if e.Mode == ModeMesh {
		return e.doMesh(params)
	}
	return e.doStatic(params)
}

func (meta *RouterExecuter) doRelay(params map[string]interface{}) (map[string]interface{}, error) {
	return nil, nil
}

func (meta *RouterExecuter) doMesh(params map[string]interface{}) (map[string]interface{}, error) {
	return nil, nil
}

func (meta *RouterExecuter) doStatic(params map[string]interface{}) (map[string]interface{}, error) {
	return nil, nil
}
