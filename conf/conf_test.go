package conf

import (
	"fmt"
	"testing"
)

func TestNewConfParse(t *testing.T) {

	testSmsCase := []string{
		`sign`,
		`templateid`,
		`appid`,
		`appkey`,
		`secretid`,
		`secretkey`,
		`phone_numbers`,
	}

	testLogCase :=[]string{
		`acc_path`,
	}

	var err error
	if err = NewConfParse(`/Users/osx/Documents/tsmsrv/tsmsrv.conf`); err != nil {
		t.Error(err.Error())
	}

	for _, lk := range testLogCase {
		fmt.Println(GlobConfig.LogSection().Key(lk).String())
	}
	for _, sk := range testSmsCase {
		fmt.Println(GlobConfig.SmsSection().Key(sk).String())

	}
}