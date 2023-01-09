package utils

import (
	"github.com/sirupsen/logrus"
)

func Die(err error) {
	if err != nil {
		logrus.Fatalf("%+v", err)
	}
}

func DoOrDie[A any](out A, err error) A {
	Die(err)
	return out
}
