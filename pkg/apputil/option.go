package apputil

import (
	"github.com/gobuffalo/packr/v2"
	"go.uber.org/zap"
)

type Option func(*Syncer) error

func Box(b *packr.Box) Option {
	return func(s *Syncer) error {
		s.box = b
		return nil
	}
}

func Logger(l *zap.SugaredLogger) Option {
	return func(s *Syncer) error {
		s.logger = l
		return nil
	}
}

func Assets(d string) Option {
	return func(s *Syncer) error {
		s.assetsDir = d
		return nil
	}
}

