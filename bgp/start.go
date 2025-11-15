package bgp

import (
	"context"

	"github.com/osrg/gobgp/v4/pkg/server"
	"github.com/rs/zerolog"
)

func Start(ctx context.Context, s *server.BgpServer) {
	logger := zerolog.Ctx(ctx)

	logger.Debug().Msg("Starting BGP server")
	defer logger.Debug().Msg("Stopped BGP server")

	context.AfterFunc(ctx, func() {
		logger.Trace().Str("cause", context.Cause(ctx).Error()).Msg("Stopping BGP server")
		s.Stop()
	})

	s.Serve()
}
