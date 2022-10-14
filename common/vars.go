package common

const (
	EncryptKey              = "xxxxxxxxxxxxxxxx"
	ConsolePipeName         = "console"
	ConsoleListenerPipeName = "consoleListener"
	HoldMaxLimit            = 5 * 60 // 5分钟
	RootPrefix              = "root"
)

const (
	ApiTypeFireForget ApiType = iota
	ApiTypeRequestResponse
	ApiTypeRequestStream
)

const (
	EventTypeRegister EventType = iota
	EventTypeUnRegister
	EventSoundOut
)
