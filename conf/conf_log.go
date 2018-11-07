package conf

import (
	"github.com/Waitfantasy/unicorn/util"
)

type LogConfig struct {
	Level      string `yaml:"level"`
	Output     string `yaml:"output"`
	Split      bool   `yaml:"split"`
	FilePath   string `yaml:"filePath"`
	FilePrefix string `yaml:"filePrefix"`
	FileSuffix string `yaml:"fileSuffix"`
}

func (c *LogConfig) fromEnvInitConfig() {
	if c.Level == "" {
		if v, err := util.Getenv("UNICORN_LOG_LEVEL", "string"); err == nil {
			c.Level = v.(string)
		}
	}

	if c.Output == "" {
		if v, err := util.Getenv("UNICORN_LOG_OUTPUT", "string"); err == nil {
			c.Output = v.(string)
		}
	}

	if c.FilePath == "" {
		if v, err := util.Getenv("UNICORN_LOG_FILE_PATH", "string"); err == nil {
			c.FilePath = v.(string)
		}
	}

	if c.FilePrefix == "" {
		if v, err := util.Getenv("UNICORN_LOG_FILE_PREFIX", "string"); err == nil {
			c.FilePrefix = v.(string)
		}
	}

	if c.FileSuffix == "" {
		if v, err := util.Getenv("UNICORN_LOG_FILE_SUFFIX", "string"); err == nil {
			c.FileSuffix = v.(string)
		}
	}

	if c.Split == false {
		if v, err := util.Getenv("UNICORN_LOG_SPLIT", "bool"); err == nil {
			c.Split = v.(bool)
		}
	}
}