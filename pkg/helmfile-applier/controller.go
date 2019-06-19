package helmfile_applier

import (
	"fmt"
	"github.com/gobuffalo/packr/v2"
	"github.com/mumoshu/helmfile-operator/pkg/apputil"
	"github.com/roboll/helmfile/pkg/app"
	"github.com/stefanprodan/k8s-podinfo/pkg/signals"
	"go.uber.org/zap"
	"os"
	"strings"
	"time"
)

type Runner struct {
	logger *zap.SugaredLogger

	syncer *apputil.Syncer

	config    *Config
	diffConf  DiffConfig
	applyConf ApplyConfig

	assetsDir string
	interval  time.Duration
	once      bool
	synced    bool
}

func New(box *packr.Box, opts ...Option) (*Runner, error) {
	l := apputil.NewLogger(os.Stderr, "debug")

	r := &Runner{
		interval: 10 * time.Second,
		diffConf: DiffConfig{
			detailedExitcode: true,
		},
		applyConf: ApplyConfig{
			logger: l,
		},
		config: &Config{
			logger: l,
			env:    "default",
		},
		assetsDir: "assets",
	}

	for i := range opts {
		if err := opts[i](r); err != nil {
			return nil, err
		}
	}

	if r.config.fileOrDir == "" {
		r.config.fileOrDir = fmt.Sprintf("%s/helmfile.yaml", r.assetsDir)
	}

	syncer, err := apputil.New(
		apputil.Box(box),
		apputil.Logger(l),
		apputil.Assets(r.assetsDir),
	)
	if err != nil {
		return nil, err
	}

	r.syncer = syncer
	r.logger = l

	return r, nil
}

var DefaultAssetsDir = "assets"

func (r *Runner) RunOnce() error {
	if err := r.syncer.SyncOnce(); err != nil {
		return err
	}

	logger := r.logger

	helmfile := app.New(r.config)
	if err := helmfile.Diff(r.diffConf); err != nil {
		switch e := err.(type) {
		case *app.Error:
			if e.Code() == 2 {
				logger.Info("Changes detected. Applying...")

				if err2 := helmfile.Apply(r.applyConf); err2 != nil {
					return err2
				}

				logger.Infof("Changes applied.")
			} else {
				return err
			}
		default:
			r.logger.Errorf("Error: %v", err)
			if strings.HasSuffix(err.Error(), "no state file found") {
				return err
			}
			return nil
		}
	} else {
		logger.Infof("No changes detected.")
	}

	return nil
}

func (r *Runner) Run() error {
	stopSig := signals.SetupSignalHandler()

	if r.once {
		return r.RunOnce()
	}

	stop := make(chan struct{}, 0)
	errs := make(chan error, 0)

	go func() {
		for {
			if err := r.RunOnce(); err != nil {
				errs <- err
				return
			}

			r.logger.Infof("Waiting for %-8v", r.interval)
			nextTime := time.Now()
			nextTime = nextTime.Add(r.interval)
			time.Sleep(time.Until(nextTime))

			select {
			case <-stop:
				r.logger.Info("Gracefully stopped the run loop")
				return
			default:
			}
		}
	}()

	select {
	case <-stopSig:
		// TODO Immediately cancel the RunOnce call on SIGTERM
		r.logger.Info("Stopping the run loop")
		stop <- struct{}{}
		return nil
	case err := <-errs:
		return err
	}

	return nil
}
