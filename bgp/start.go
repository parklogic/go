package bgp

import (
	"context"
	"fmt"
	"sync"

	"github.com/osrg/gobgp/v4/api"
	"github.com/osrg/gobgp/v4/pkg/server"
	"github.com/rs/zerolog"
)

func Start(ctx context.Context, s *server.BgpServer, g *api.Global) error {
	logger := zerolog.Ctx(ctx)

	logger.Debug().Msg("Starting BGP server")
	defer logger.Debug().Msg("Stopped BGP server")

	wg := &sync.WaitGroup{}
	wg.Go(func() {
		s.Serve()
	})

	context.AfterFunc(ctx, func() {
		logger.Trace().Str("cause", context.Cause(ctx).Error()).Msg("Stopping BGP server")
		s.Stop()
	})

	if err := s.StartBgp(ctx, &api.StartBgpRequest{
		Global: g,
	}); err != nil {
		return fmt.Errorf("error starting BGP: %w", err)
	}

	wg.Wait()

	return nil
}
