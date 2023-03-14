package contract

type IDistrbution struct {
	dynamic any
}

type ITask struct{}

type ITaskResult struct{}

type IPluginContract interface {
	Name() string
	//master
	PlanRead() []IDistrbution
	//worker method
	Distrbute(dist IDistrbution, task ITask) ITaskResult

	Configs(map[string]string) IPluginContract
}
