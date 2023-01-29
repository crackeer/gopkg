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

func (e *RouterExecuter) Exec(routerMeta *RouterMeta) (interface{}, error) {
	if routerMeta.Mode == ModeRelay {
		_, err := e.doRelay(routerMeta)
		if err != nil {
			return nil, err
		}
		//apiResponse.Data
	}
	if routerMeta.Mode == ModeMesh {
		_, _, err := e.doMesh(routerMeta)
		if err != nil {
			return nil, err
		}
	}
	return e.doStatic(routerMeta)
}

func (meta *RouterExecuter) doRelay(routerMeta *RouterMeta) (*api.APIResponse, error) {
	client := api.NewRequestClient(meta.apiFactory)
	client.UseEnv(meta.env)
	return client.Request(routerMeta.Config, meta.input, meta.header)
}

func (meta *RouterExecuter) doMesh(routerMeta *RouterMeta) (map[string]*api.APIResponse, map[string]string, error) {
	client := api.NewRequestClient(meta.apiFactory)
	client.UseEnv(meta.env)
	list, _ := api.ParseMeshConfig(routerMeta.Config)
	return client.Mesh(list, meta.input, meta.header)
}

func (meta *RouterExecuter) doStatic(routerMeta *RouterMeta) (map[string]interface{}, error) {
	return nil, nil
}
