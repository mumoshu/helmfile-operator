package helmfile_applier

import (
	"github.com/roboll/helmfile/pkg/app"
	"go.uber.org/zap"
)

// TODO Generalize this and upstream to the helmfile project
type Config struct {
	args, env, fileOrDir, kubeContext, ns, helmBinary string

	sels, vals []string
	set        map[string]interface{}
	logger     *zap.SugaredLogger
}

func (c *Config) HelmBinary() string {
	return c.helmBinary
}

func (c *Config) KubeContext() string {
	return c.kubeContext
}

func (c *Config) Namespace() string {
	return c.ns
}

func (c *Config) Selectors() []string {
	return c.sels
}

func (c *Config) Set() map[string]interface{} {
	return c.set
}

func (c *Config) ValuesFiles() []string {
	return c.vals
}

func (c *Config) Logger() *zap.SugaredLogger {
	return c.logger
}

func (c *Config) Args() string {
	return c.args
}

func (c *Config) Env() string {
	return c.env
}

func (c *Config) FileOrDir() string {
	return c.fileOrDir
}

type DiffConfig struct {
	args string
	vals []string
	concurrency int
	skipDeps, detailedExitcode, suppressSecrets bool
}

func (c DiffConfig) Args() string {
	return c.args
}

func (c DiffConfig) Values() []string {
	return c.vals
}

func (c DiffConfig) SkipDeps() bool {
	return c.skipDeps
}

func (c DiffConfig) SuppressSecrets() bool {
	return c.suppressSecrets
}

func (c DiffConfig) DetailedExitcode() bool {
	return c.detailedExitcode
}

func (c DiffConfig) Concurrency() int {
	return c.concurrency
}

type ApplyConfig struct {
	args string
	skipDeps, suppressSecrets, interactive bool
	concurrency int
	vals []string
	logger *zap.SugaredLogger
}

func (c ApplyConfig) Args() string {
	return c.args
}

func (c ApplyConfig) Values() []string {
	return c.vals
}

func (c ApplyConfig) SkipDeps() bool {
	return c.skipDeps
}

func (c ApplyConfig) SuppressSecrets() bool {
	return c.suppressSecrets
}

func (c ApplyConfig) Concurrency() int {
	return c.concurrency
}

func (c ApplyConfig) Interactive() bool {
	return c.interactive
}

func (c ApplyConfig) Logger() *zap.SugaredLogger {
	return c.logger
}

var _ app.ConfigProvider = &Config{}

var _ app.DiffConfigProvider = &DiffConfig{}

var _ app.ApplyConfigProvider = &ApplyConfig{}
