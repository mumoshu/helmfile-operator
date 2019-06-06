package helmfile_applier

type Option func(runner *Runner) error

func AssetDir(d string) Option {
	return func(r *Runner) error {
		r.assetsDir = d
		return nil
	}
}

func Source(s string) Option {
	return func(r *Runner) error {
		r.source = s
		return nil
	}
}

func Once(b bool) Option {
	return func(r *Runner) error {
		r.once = b
		return nil
	}
}
