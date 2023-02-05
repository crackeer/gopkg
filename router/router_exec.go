package router

import (
	"github.com/crackeer/gopkg/mapbuilder"
	"github.com/crackeer/gopkg/router/api"
	"github.com/sirupsen/logrus"
)

// RouterExecuter
type RouterExecuter struct {
	logrusLogger *logrus.Logger
	env          string
	header       map[string]string
	input        map[string]interface{}
	apiFactory   api.APIFactory

	meshAPIPrefix    string
	meshAPISeperator string

	Respone map[string]*api.APIResponse
	Error   map[string]string
	Data    map[string]interface{}
}

// NewRouterExecuter
//
//	@param meta
//	@param header
//	@return *RouterExecuter
func NewRouterExecuter(apiFactory api.APIFactory) *RouterExecuter {
	return &RouterExecuter{
		apiFactory: apiFactory,
		Respone:    map[string]*api.APIResponse{},
		Error:      map[string]string{},
		Data:       map[string]interface{}{},

		meshAPIPrefix:    mapbuilder.DefaultPrefix,
		meshAPISeperator: mapbuilder.DefaultSeperator,
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

// UseLogrusLogger func
//
//	@receiver e
//	@param logger
//	@return *RouterExecuter
func (e *RouterExecuter) UseLogrusLogger(logger *logrus.Logger) *RouterExecuter {
	e.logrusLogger = logger
	return e
}

// UseHeader //
//
//	@receiver e
//	@param env
//	@return *RouterExecuter
func (e *RouterExecuter) UseHeader(header map[string]string) *RouterExecuter {
	e.header = header
	e.Data[api.HeaderKey] = header
	return e
}

// SetMeshAPIConfig
//
//	@receiver e
//	@param prefix
//	@param seperator
func (e *RouterExecuter) SetMeshAPIConfig(prefix, seperator string) {
	e.meshAPIPrefix = prefix
	e.meshAPISeperator = seperator
}

// UseInput
//
//	@receiver e
//	@param env
//	@return *RouterExecuter
func (e *RouterExecuter) UseInput(input map[string]interface{}) *RouterExecuter {
	e.input = input
	e.Data[api.InputKey] = input
	return e
}

// Exec
//
//	@receiver executor
//	@param routerMeta
//	@return error
func (executor *RouterExecuter) Exec(routerMeta *RouterMeta) error {
	if routerMeta.Mode == ModeRelay {
		return executor.Relay(routerMeta)
	}

	if routerMeta.Mode == ModeMesh {
		return executor.Mesh(routerMeta)
	}

	if routerMeta.Mode == ModeStatic {
		return executor.Static(routerMeta)
	}

	return nil
}

// Relay
//
//	@receiver executor
//	@param routerMeta
//	@return error
func (executor *RouterExecuter) Relay(routerMeta *RouterMeta) error {
	client := api.NewRequestClient(executor.apiFactory).UseEnv(executor.env).UseLogrusLogger(executor.logrusLogger)
	result, err := client.Request(routerMeta.RelayAPI, executor.input, executor.header)
	executor.Respone[result.Name] = result
	executor.Data[result.Name] = result.Data
	if err != nil {
		executor.Error[result.Name] = err.Error()
	}
	return err
}

// Mesh
//
//	@receiver executor
//	@param routerMeta
//	@return error
func (executor *RouterExecuter) Mesh(routerMeta *RouterMeta) error {
	client := api.NewRequestClient(executor.apiFactory).UseEnv(executor.env).UseLogrusLogger(executor.logrusLogger)
	result, errMap, err := client.Mesh(routerMeta.MeshConfig, executor.input, executor.header)
	for key, value := range result {
		executor.Respone[key] = value
		executor.Data[key] = value.Data
	}
	for key, value := range errMap {
		executor.Error[key] = value
	}
	return err
}

// Static
//
//	@receiver executor
//	@param routerMeta
//	@return error
func (executor *RouterExecuter) Static(routerMeta *RouterMeta) error {
	executor.Data["static"] = map[string]interface{}{
		"header": executor.header,
		"input":  executor.input,
		"env":    executor.env,
	}
	return nil
}

// BuildResponse
//
//	@receiver executor
//	@param routerMeta
//	@return interface{}
func (executor *RouterExecuter) BuildResponse(routerMeta *RouterMeta) interface{} {
	if routerMeta.Response == nil {
		if routerMeta.Mode == ModeRelay {
			for _, value := range executor.Data {
				return value
			}
		}
		if routerMeta.Mode == ModeMesh {
			return executor.Data
		}
		return nil
	}
	builder, _ := mapbuilder.MapBuilderFrom(executor.Data)
	return builder.Build(routerMeta.Response)
}
