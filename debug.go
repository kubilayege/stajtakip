// +build debug

package main

import "github.com/sirupsen/logrus"

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}
