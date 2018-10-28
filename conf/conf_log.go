package conf

import (
	"github.com/Waitfantasy/unicorn/util/logger"
	"strings"
)

type LogConf struct {
	log        *logger.Log
	Enable     bool   `yaml:"enable"`
	Level      string `yaml:"level"`
	Output     string `yaml:"output"`
	Split      bool   `yaml:"split"`
	Filepath   string `yaml:"filePath"`
	FilePrefix string `yaml:"filePrefix"`
	FileSuffix string `yaml:"fileSuffix"`
}

func (c *LogConf) Init() error {
	c.log = &logger.Log{}
	c.setLogLevel()
	c.setLogOutput()
	if err := c.log.InitLogger(); err != nil {
		return err
	}
	return nil
}

func (c *LogConf) GetLogger() *logger.Log {
	return c.log
}

func (c *LogConf) InitLog() (*logger.Log, error) {
	log := &logger.Log{}

	// init log writer
	if err := log.InitLogger(); err != nil {
		return nil, err
	}

	return log, nil
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
			c.log.SetFileOut(c.Filepath)
			c.log.SetFilePrefix(c.FilePrefix)
			c.log.SetFileSuffix(c.FileSuffix)
			c.log.SetFileSplit(c.Split)
		}
	}
}
