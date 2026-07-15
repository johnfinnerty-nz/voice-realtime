package volcengine

const (
	defaultURL = "wss://openspeech.bytedance.com/api/v3/realtime/dialogue"

	protocolVersion = 0b0001

	clientFullRequest      = 0b0001
	clientAudioOnlyRequest = 0b0010
	serverFullResponse     = 0b1001
	serverACK              = 0b1011
	serverErrorResponse    = 0b1111

	msgWithEvent     = 0b0100
	noSerialization  = 0b0000
	jsonSerialization = 0b0001
	gzipCompression  = 0b0001

	eventStartConnection  = 1
	eventFinishConnection = 2
	eventStartSession     = 100
	eventFinishSession    = 102
	eventTaskRequest      = 200
)
