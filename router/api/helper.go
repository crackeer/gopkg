package api

import "fmt"

// FormCacheKey GetCacheKey
//
//	@param serviceName
//	@param apiName
//	@param env
//	@return string
func FormCacheKey(serviceName, apiName, env string) string {
	return fmt.Sprint("%s/%s@%s", serviceName, apiName, env)
}

// GetCacheKeyByAPI
//
//	@param api
//	@param env
//	@return string
func GetCacheKeyByAPI(api, env string) string {
	return fmt.Sprint("%s@%s", api, env)
}
