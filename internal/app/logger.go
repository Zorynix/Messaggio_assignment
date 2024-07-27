package app

import (
	"os"

	"github.com/sirupsen/logrus"
)

func SetLogrus(logLevel string) {
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.Fatalf("Invalid log level: %v", err)
	}

	logrus.SetOutput(os.Stdout)
	// file, err := os.OpenFile("logfile.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// if err != nil {
	// 	logrus.Fatalf("Failed to open log file: %v", err)
	// }

	// multiWriter := io.MultiWriter(os.Stdout, file)
	// logrus.SetOutput(multiWriter)

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(level)
}
