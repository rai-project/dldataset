package dldataset

import (
	"github.com/k0kubun/pp"
	"github.com/rai-project/config"
	"github.com/rai-project/vipertags"
)

type dldatasetConfig struct {
	WorkingDirectory string        `json:"working_directory" config:"dldataset.working_directory" default:""`
	done             chan struct{} `json:"-" config:"-"`
}

var (
	// Config holds the data read by rai-project/config
	Config = &dldatasetConfig{
		done: make(chan struct{}),
	}
)

func (dldatasetConfig) ConfigName() string {
	return "DLDataset"
}

func (c *dldatasetConfig) SetDefaults() {
	vipertags.SetDefaults(c)
}

func (c *dldatasetConfig) Read() {
	defer close(c.done)
	config.App.Wait()
	vipertags.Fill(c)
	if c.WorkingDirectory == "" || c.WorkingDirectory == "default" {
		c.WorkingDirectory = config.App.TempDir
	}
}

func (c dldatasetConfig) Wait() {
	<-c.done
}

func (c dldatasetConfig) String() string {
	return pp.Sprintln(c)
}

func (c dldatasetConfig) Debug() {
	log.Debug("DLDataset Config = ", c)
}

func init() {
	config.Register(Config)
}
