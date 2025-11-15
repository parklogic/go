package bgp

import (
	"context"
	"log/slog"

	"github.com/osrg/gobgp/v4/pkg/server"
	"github.com/parklogic/go/log"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

func NewServer(ctx context.Context, cfg *Configuration) *server.BgpServer {
	logger := zerolog.Ctx(ctx).With().Logger()

	slogLevel := &slog.LevelVar{}
	slogLevel.Set(log.GetSlogLevel(logger.GetLevel()))

	return server.NewBgpServer(
		server.GrpcListenAddress(cfg.Listen),
		server.GrpcOption([]grpc.ServerOption{grpc.MaxRecvMsgSize(cfg.RecvMaxMsgSize), grpc.MaxSendMsgSize(cfg.SendMaxMsgSize)}),
		server.LoggerOption(log.GetSlogger(&logger), slogLevel),
	)

}
