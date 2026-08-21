package shared

const NetworkBufferSize int = 4 * 1024 * 1024

const (
	AuthTypePassword AuthType = iota
	AuthTypeDiscord
)

const (
	DiscordHistoryEventLinked DiscordHistoryEvent = iota
	DiscordHistoryEventUnlinked
)

const (
	MsgTypePleaseDiscard uint8 = iota
	MsgTypeGameCutoffConst
	GameMsgTypeWorldStateUpdate
	GameMsgTypeCameraDelta
	GameMsgTypeSpawn
	GameMsgTypeInput
)

const NameLengthLimit int = 16
