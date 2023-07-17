package common

const (
	EncryptKey              = "yeyongping202210"
	ConsolePipeName         = "console"
	ConsoleListenerPipeName = "consoleListener"
	UploadFilePipeName      = "uploadFile"
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
