package utils

import (
	"go.uber.org/zap"
	"testing"
	"tsmsrv/conf"
)

func TestNewLoggerMgr(t *testing.T) {
	var err error
	if err = conf.NewConfParse(`/Users/osx/Documents/tsmsrv/tsmsrv.conf`); err != nil {

		t.Error(`t1`, err.Error())
	}
	if err = NewLoggerMgr(`info`); err != nil {
		t.Error(`t2`, err.Error())
	}
	Logger.Log.Info(`sssssssss`, zap.Int(`ssss`, 1))
	Logger.Log.Error(`a`)
	Logger.Log.Fatal(`eeeeee`)
	Logger.Log.Debug(`sssss`)
}
