package Server

import "GophKeeper/pkg/logger"

var flagLogLevel string

func init() {
	//
}

func Run() error {
	flagLogLevel = "debug"
	if err := logger.Initialize(flagLogLevel); err != nil {
		return err
	}
	return nil
}

/*
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, err
	}
	_ = conn
	return nil, err
*/
