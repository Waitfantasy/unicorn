package conf

import (
	"github.com/Waitfantasy/unicorn/util/logger"
	"strings"
)

type LogConf struct {
	Enable     bool   `yaml:"enable"`
	Level      string `yaml:"level"`
	Output     string `yaml:"output"`
	Split      bool   `yaml:"split"`
	Filepath   string `yaml:"filePath"`
	FilePrefix string `yaml:"filePrefix"`
	FileSuffix string `yaml:"fileSuffix"`
}

func (c *LogConf) InitLog() (*logger.Log, error) {
	log := &logger.Log{}

	// init level
	c.initLogLevel(log)

	// int log output
	c.initLogOutput(log)

	// init log writer
	if err := log.InitLogger(); err != nil {
		return nil, err
	}

	return log, nil
}

func (c *LogConf) initLogLevel(log *logger.Log) {
	levels := strings.Split(c.Level, "|")
	for _, level := range levels {
		switch strings.ToUpper(level) {
		case "DEBUG":
			log.SetDebugLevel()
		case "INFO":
			log.SetInfoLevel()
		case "WARN":
			log.SetWarnLevel()
		case "ERR":
			log.SetErrLevel()
		}
	}
}

func (c *LogConf) initLogOutput(log *logger.Log) {
	outputs := strings.Split(c.Output, "|")
	for _, output := range outputs {
		switch strings.ToUpper(output) {
		case "CONSOLE":
			log.SetConsoleOut()
		case "FILE":
			log.SetFileOut(c.Filepath)
			log.SetFilePrefix(c.FilePrefix)
			log.SetFileSuffix(c.FileSuffix)
			log.SetFileSplit(c.Split)
		}
	}
}
