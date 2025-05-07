package types

type NabrunnerTasks string

func (t NabrunnerTasks) String() string {
	return string(t)
}

const (
	Dataworker          NabrunnerTasks = "run:dataworker"
	CancelDataWorkerJob NabrunnerTasks = "cancel:dataworker"
	MetadataTask        NabrunnerTasks = "sync:datasource:metadata"
	AIAgentTask         NabrunnerTasks = "run:aiagent:task"
	NabhikTask          NabrunnerTasks = "run:nabhik:task"
	CancelTasks         NabrunnerTasks = "cancel:task"
)
