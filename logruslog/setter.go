package logruslog

import (
	"github.com/sirupsen/logrus"
)

type Setter func(config *Config)

func SetLevel(level logrus.Level) Setter {
	return func(config *Config) {
		config.level = level
	}
}

func SetFormatter(formatter logrus.Formatter) Setter {
	return func(config *Config) {
		config.formatter = formatter
	}
}

func SetProjectName(projectName string) Setter {
	return func(config *Config) {
		config.projectName = projectName
	}
}

func SetIncludeReportCaller() Setter {
	return func(config *Config) {
		config.includeReportCaller = true
	}
}

func AddHook(hooks ...logrus.Hook) Setter {
	return func(config *Config) {
		config.hooks = append(config.hooks, hooks...)
	}
}
