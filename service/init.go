package service

import "log"

func init() {
	if err := initUserDB(); err != nil {
		log.Printf("init user db error: %v\n", err)
	}
}
