package types

type NabhikMessageAssetType string

func (a NabhikMessageAssetType) String() string {
	return string(a)
}

const (
	NabhikChat = "nabhik"
)

type NabhikAgentType string

func (a NabhikAgentType) String() string {
	return string(a)
}

const (
	NabhikRecurringAgent NabhikAgentType = "recurring"
	NabhikWebhookAgent   NabhikAgentType = "webhook"
	NabhikEventAgent     NabhikAgentType = "event"
)

type AssetOwner string

func (a AssetOwner) String() string {
	return string(a)
}

const (
	AssetOwnerSystem AssetOwner = "system"
	AssetOwnerUser   AssetOwner = "user"
)

type ActionStatus string

func (a ActionStatus) String() string {
	return string(a)
}

const (
	ActionStatusPending    ActionStatus = "pending"
	ActionStatusCompleted  ActionStatus = "completed"
	ActionStatusInProgress ActionStatus = "in-progress"
	ActionStatusFailed     ActionStatus = "failed"
)

const (
	NabhikTaskPackageName = "vapusdata-python-env"
)

type NabhikChatMessageType string

func (a NabhikChatMessageType) String() string {
	return string(a)
}

const (
	NabhikTaskMessage  NabhikChatMessageType = "task"
	NabhikAgentMessage NabhikChatMessageType = "agent"
	NabhikChatMessage  NabhikChatMessageType = "chat"
	NabhikErrorMessage NabhikChatMessageType = "error"
)

type NabhikTaskQueueType string

func (a NabhikTaskQueueType) String() string {
	return string(a)
}

const (
	NabhikTaskBasicQueue NabhikTaskQueueType = "nbtask::"
)
