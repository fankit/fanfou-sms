package conf

import "gopkg.in/ini.v1"

type ConfParse struct {
	Reader *ini.File
}

var GlobConfig *ConfParse

func NewConfParse(c string) (err error) {
	var r *ini.File
	if r, err = ini.Load(c); err != nil {
		return
	}
	GlobConfig = &ConfParse{Reader: r}
	return
}

func (c *ConfParse) LogSection() *ini.Section {
	return c.Reader.Section(`log`)
}

func (c *ConfParse) SmsSection() *ini.Section {
	return c.Reader.Section(`sms`)
}

