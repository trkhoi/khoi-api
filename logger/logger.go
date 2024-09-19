// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package logger configures the Logrus logging library.
package logger

import (
	"strings"

	stackdriver "github.com/TV4/logrus-stackdriver-formatter"
	"github.com/sirupsen/logrus"
)

func New(serviceName string) *logrus.Entry {
	l := logrus.WithFields(logrus.Fields{
		"app": serviceName,
	})
	l.Logger.SetFormatter(&logrus.JSONFormatter{})
	l.Level = logrus.InfoLevel
	return l
}

func newFormatter(formatter string) logrus.Formatter {
	switch strings.ToLower(formatter) {
	case "stackdriver":
		return stackdriver.NewFormatter()
	case "json":
		return &logrus.JSONFormatter{}
	}
	return &logrus.TextFormatter{}
}

func toLevel(level string) logrus.Level {
	switch strings.ToLower(level) {
	case "trace":
		return logrus.TraceLevel
	case "debug":
		return logrus.DebugLevel
	case "warn":
		fallthrough
	case "warning":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	}
	return logrus.InfoLevel
}

func isDebugLevel(level logrus.Level) bool {
	switch level {
	case logrus.TraceLevel:
		fallthrough
	case logrus.DebugLevel:
		return true
	}
	return false
}
