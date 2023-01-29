package router

import "github.com/crackeer/gopkg/router/api"

// RouterExecuter
type RouterExecuter struct {
	env        string
	header     map[string]string
	input      map[string]interface{}
	apiFactory api.APIMetaFactory

	ErrorCode   int64
	ResponseRaw string
	Error       string
}

// NewRouterExecuter
//
//	@param meta
//	@param header
//	@return *RouterExecuter
func NewRouterExecuter(apiFactory api.APIMetaFactory) *RouterExecuter {
	return &RouterExecuter{
		apiFactory: apiFactory,
	}
}

// UseEnv
//
//	@receiver e
//	@param env
//	@return *RouterExecuter
func (e *RouterExecuter) UseEnv(env string) *RouterExecuter {
	e.env = env
	return e
}

// UseHeader //
//
//	@receiver e
//	@param env
//	@return *RouterExecuter
func (e *RouterExecuter) UseHeader(header map[string]string) *RouterExecuter {
	e.header = header
	return e
}

// UseInput
//
//	@receiver e
//	@param env
//	@return *RouterExecuter
func (e *RouterExecuter) UseInput(input map[string]interface{}) *RouterExecuter {
	e.input = input
	return e
}

func (e *RouterExecuter) Exec(routerMeta *RouterMeta) (map[string]interface{}, error) {
	if routerMeta.Mode == ModeRelay {
		return e.doRelay(routerMeta)
	}
	if routerMeta.Mode == ModeMesh {
		return e.doMesh(routerMeta)
	}
	return e.doStatic(routerMeta)
}

func (meta *RouterExecuter) doRelay(routerMeta *RouterMeta) (map[string]interface{}, error) {
	meta.apiFactory.GetAPIMeta(routerMeta.Config, meta.env)
	return nil, nil
}

func (meta *RouterExecuter) doMesh(routerMeta *RouterMeta) (map[string]interface{}, error) {
	return nil, nil
}

func (meta *RouterExecuter) doStatic(routerMeta *RouterMeta) (map[string]interface{}, error) {
	return nil, nil
}
