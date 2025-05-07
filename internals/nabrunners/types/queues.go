package types

var NabRunnerQueues = map[string]int{
	DataworkerQueue.String():   4,
	MetadataSyncQueue.String(): 3,
	AIAgentQueue.String():      3,
}

type NabRunnerQueue string

func (q NabRunnerQueue) String() string {
	return string(q)
}

const (
	DataworkerQueue   NabRunnerQueue = "dataworker"
	MetadataSyncQueue NabRunnerQueue = "metadataSync"
	AIAgentQueue      NabRunnerQueue = "aiagent"
)
