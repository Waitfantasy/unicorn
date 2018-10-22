package conf

import (
	"github.com/Waitfantasy/unicorn/util"
	"strings"
		"os"
)


const UnixDefaultLogFilepath = "/var/run/log/unicorn"

type LogConf struct {
	Enable         bool   `yaml:"enable"`
	Level          string `yaml:"level"`
	Output         string `yaml:"output"`
	Split          bool   `yaml:"split"`
	SplitDimension bool   `yaml:"splitDimension"`
	Filepath       string `yaml:"filepath"`
	FilePrefix     string `yaml:"filePrefix"`
	FileSuffix     string `yaml:"fileSuffix"`
}

func (c *LogConf) InitLog() (*util.Log, error) {
	log := &util.Log{}

	// init level
	c.initLogLevel(log)

	// int log output
	if err := c.initLogOutput(log); err != nil {
		return nil, err
	}

	return log, nil
}

func (c *LogConf) initLogLevel(log *util.Log)  {
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

func (c *LogConf) initLogOutput(log *util.Log)  error {
	outputs := strings.Split(c.Output, "|")
	for _, output := range outputs {
		switch strings.ToUpper(output) {
		case "CONSOLE":
			log.SetConsoleOut()
		case "FILE":
			if path, err := c.createLogPath(); err != nil {
				return err
			} else {
				log.SetFileOut(path)
			}
		}
	}
	return nil
}

func (c *LogConf) createLogPath() (string, error) {
	var path string

	if c.Filepath == "" {
		path = UnixDefaultLogFilepath
	} else {
		path = c.Filepath
	}

	_, err := os.Stat(path)
	if err == nil {
		return path, nil
	}

	if os.IsNotExist(err) {
		// create log path dir
		if err = os.Mkdir(path, 0644); err != nil {
			return path, err
		}
	}

	return path ,err
}