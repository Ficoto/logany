package slog

import "log/slog"

type Setter func(c *Config)

func SetLevel(level slog.Level) Setter {
	return func(c *Config) {
		c.level = level
		c.replaceAttrList = append([]func(groups []string, a slog.Attr) slog.Attr{LevelReplaceAttr}, c.replaceAttrList...)
	}
}

func SetFormatterText() Setter {
	return func(c *Config) {
		c.formatter = FormatterText
	}
}

func SetFormatterJson() Setter {
	return func(c *Config) {
		c.formatter = FormatterJson
	}
}

func SetFormatter(formatter Formatter) Setter {
	return func(c *Config) {
		c.formatter = formatter
	}
}

func SetProjectName(projectName string) Setter {
	return func(c *Config) {
		c.projectName = projectName
	}
}

func SetAddSource() Setter {
	return func(c *Config) {
		c.addSource = true
		c.replaceAttrList = append([]func(groups []string, a slog.Attr) slog.Attr{AddSourceReplaceAttr}, c.replaceAttrList...)
	}
}

func SetReplaceAttr(ra func(groups []string, a slog.Attr) slog.Attr) Setter {
	return func(c *Config) {
		c.replaceAttrList = append(c.replaceAttrList, ra)
	}
}
