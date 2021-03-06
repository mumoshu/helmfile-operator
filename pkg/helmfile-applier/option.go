package helmfile_applier

import (
	"fmt"
	"os/user"
	"path/filepath"
)

type Option func(runner *Runner) error

func AssetDir(d string) Option {
	return func(r *Runner) error {
		r.assetsDir = d
		return nil
	}
}

func Source(s string) Option {
	return func(r *Runner) error {
		r.config.fileOrDir = s
		return nil
	}
}

func Once(b bool) Option {
	return func(r *Runner) error {
		r.once = b
		return nil
	}
}

func HelmX(b bool) Option {
	return func(r *Runner) error {
		if b {
			usr, err := user.Current()
			if err != nil {
				return fmt.Errorf("enabling helm-x integration: %v", err)
			}
			r.config.helmBinary = filepath.Join(usr.HomeDir, ".helm/plugins/helm-x/bin/helm-x")
		}
		return nil
	}
}

func Environment(e string) Option {
	return func(r *Runner) error {
		r.config.env = e
		return nil
	}
}

func Values(m map[string]interface{}) Option {
	return func(r *Runner) error {
		r.config.set = m
		return nil
	}
}

func ValuesFiles(f []string) Option {
	return func(r *Runner) error {
		r.config.vals = f
		return nil
	}
}
