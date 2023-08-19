package gdk_helpers

import "github.com/sirupsen/logrus"

func FailOnError(err error, message string) {
	if err != nil {
		log := logrus.New()
		log.Errorf("%s: %s", message, err.Error())
	}
}

func FatalOnError(err error, message string) {
	if err != nil {
		log := logrus.New()
		log.Fatalf("%s: %s", message, err.Error())
	}
}
