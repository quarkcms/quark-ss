package foundation

type Application struct {
	Container map[string]any
}

var application Application

// 简单绑定
func Bind[T any](abstract string, concrete T) {
	container := make(map[string]any)
	container[abstract] = concrete

	application.Container = container
}

// 单例绑定
func Singleton[T any](abstract string, concrete T) {
	if application.Container[abstract] == nil {
		container := make(map[string]any)
		container[abstract] = concrete

		application.Container = container
	}
}

// 获取实例
func Make(abstract string) any {
	return application.Container[abstract]
}
