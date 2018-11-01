package conf

import (
	"github.com/Waitfantasy/unicorn/util"
	"github.com/Waitfantasy/unicorn/util/logger"
	"strings"
)

type LogConf struct {
	log        *logger.Log
	Level      string `yaml:"level"`
	Output     string `yaml:"output"`
	Split      bool   `yaml:"split"`
	FilePath   string `yaml:"filePath"`
	FilePrefix string `yaml:"filePrefix"`
	FileSuffix string `yaml:"fileSuffix"`
}

func (c *LogConf) Init() error {
	if c.Level == "" {
		if v, err := util.GetEnv("UNICORN_LOG_LEVEL", "string"); err != nil {
			c.Level = "info|warning|error|debug"
		} else {
			c.Level = v.(string)
		}
	}


	if c.Output == "" {
		if v, err := util.GetEnv("UNICORN_LOG_OUTPUT", "string"); err != nil {
			c.Output = "console"
		} else {
			c.Output = v.(string)
		}
	}

	if c.FilePath == "" {
		if v, err := util.GetEnv("UNICORN_LOG_FILE_PATH", "string"); err != nil {
			c.FilePath = "/var/log/unicorn"
		} else {
			c.FilePath = v.(string)
		}
	}

	if c.FilePrefix == "" {
		if v, err := util.GetEnv("UNICORN_LOG_FILE_PREFIX", "string"); err != nil {
			c.FilePrefix = "unicorn"
		} else {
			c.FilePrefix = v.(string)
		}
	}

	if c.FileSuffix == "" {
		if v, err := util.GetEnv("UNICORN_LOG_FILE_SUFFIX", "string"); err != nil {
			c.FileSuffix = "log"
		} else {
			c.FileSuffix = v.(string)
		}
	}

	if c.Split == false {
		if v, err := util.GetEnv("UNICORN_LOG_SPLIT", "bool"); err != nil {
			c.Split = false
		} else {
			c.Split = v.(bool)
		}
	}

	c.log = &logger.Log{}
	c.setLogLevel()
	c.setLogOutput()
	if err := c.log.InitLogger(); err != nil {
		return err
	}
	return nil
}

func (c *LogConf) setLogLevel() {
	levels := strings.Split(c.Level, "|")
	for _, level := range levels {
		switch strings.ToUpper(level) {
		case "DEBUG":
			c.log.SetDebugLevel()
		case "INFO":
			c.log.SetInfoLevel()
		case "WARN":
			c.log.SetWarnLevel()
		case "ERROR":
			c.log.SetErrLevel()
		}
	}
}

func (c *LogConf) setLogOutput() {
	outputs := strings.Split(c.Output, "|")
	for _, output := range outputs {
		switch strings.ToUpper(output) {
		case "CONSOLE":
			c.log.SetConsoleOut()
		case "FILE":
			c.log.SetFileOut(c.FilePath)
			c.log.SetFilePrefix(c.FilePrefix)
			c.log.SetFileSuffix(c.FileSuffix)
			c.log.SetFileSplit(c.Split)
		}
	}
}

func (c *LogConf) GetLogger() *logger.Log {
	return c.log
}