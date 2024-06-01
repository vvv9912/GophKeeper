package main

import "GophKeeper/pkg/logger"

var flagLogLevel string

//func init() {
//	//
//}

func Run() error {
	flagLogLevel = "debug"
	if err := logger.Initialize(flagLogLevel); err != nil {
		return err
	}
	return nil
}
