package apputil

import (
	"github.com/davidovich/summon/pkg/summon"
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/packr/v2/file"
	"go.uber.org/zap"
)

type Syncer struct {
	box    *packr.Box
	logger *zap.SugaredLogger
	assetsDir string

	s      *summon.Summoner
	synced bool
}

func New(opts ...Option) (*Syncer, error) {
	s := &Syncer{}
	for _, o := range opts {
		if err := o(s); err != nil {
			return nil, err
		}
	}
	return s, nil
}

func (r *Syncer) SyncOnce() error {
	if !r.synced && r.box != nil {
		assetsDir := r.assetsDir

		if err := r.box.Walk(func(path string, info file.File) error {
			println(path)
			return nil
		}); err != nil {
			return err
		}

		if r.s == nil {
			var err error
			r.s, err = summon.New(r.box)
			if err != nil {
				return err
			}
		}

		str, err := r.s.Summon(
			summon.All(true),
			summon.Raw(true),
			summon.Dest(assetsDir),
		)
		if err != nil {
			return err
		}

		r.logger.Debugf("summon: %s", str)

		r.synced = true
	}
	return nil
}
