package bgp

type Configuration struct {
	Listen         string `description:"GRPC service listen address"`
	RecvMaxMsgSize int    `description:"Maximum message size received by the gRPC server"`
	SendMaxMsgSize int    `description:"Maximum message size sent by the gRPC server"`
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Listen:         ":50051",
		RecvMaxMsgSize: 256 << 20,
		SendMaxMsgSize: 256 << 20,
	}
}
