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
