package control

import (
	"testing"
	"tsmsrv/conf"
)

func TestSendSmsSrv_SendSms(t *testing.T) {

	var err error
	if err = conf.NewConfParse(`/Users/osx/Documents/tsmsrv/tsmsrv.conf`); err != nil {

		t.Error(`t1`, err.Error())
	}

	smss := NewSendSmsSrv()
	smss.SendSms([]string{`13313331333`}, []string{`111111`, `222222`})

}