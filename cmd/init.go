package main

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()


func init() {
	err := godotenv.Load()
	
	if err != nil {
		log.WithFields(logrus.Fields{
			"event": "dot_env_fault",
		}).Warning("Файл .env не был загружен")
	} else {
		log.WithFields(logrus.Fields{
			"event": "dot_env_success",
		}).Info("Файл .env был загружен")
	}
}