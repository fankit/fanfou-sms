package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"os"
	"time"
	"tsmsrv/conf"
	"tsmsrv/control"
	"tsmsrv/utils"
)

var app *cli.App

func init() {
	var (
		auths []*cli.Author
		auth  *cli.Author
	)

	auth = &cli.Author{Name: "seraphico", Email: "osx1260@163.com"}
	auths = append(auths, auth)

	app = cli.NewApp()

	app.Name = `tsmsrv`
	app.Usage = `腾讯短信接口`
	app.Authors = auths
	app.Version = `beta`
	app.Copyright = `©2007~2020 fanfou.com 粤公网安备 44030502004681号 深圳市中经饭否科技有限公司 版权所有`
}

func main() {

	app.Commands = []*cli.Command{
		{
			Name:        "start",
			Usage:       "指定启动参数",
			Description: "启动服务",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "conf",
					Aliases: []string{`c`},
					Usage:   "指定配置文件",
					Value:   "tsmsrv.conf",
				},
				&cli.BoolFlag{
					Name:    "debug",
					Aliases: []string{`D`},
					Usage:   "是否开启DEBUG",
					Value:   false,
				},
			},
			Action: func(ctx *cli.Context) (err error) {
				if err = conf.NewConfParse(ctx.String(`c`)); err != nil {
					return
				}

				if err = utils.NewLoggerMgr(conf.GlobConfig.LogSection().Key(`level`).String()); err != nil {
					return
				}

				if err = control.NewRouters(ctx.Bool(`D`)); err != nil {
					return
				}
				addr := conf.GlobConfig.Reader.Section(``).Key(`bind`).String()
				port := conf.GlobConfig.Reader.Section(``).Key(`port`).String()
				httpaddrs := fmt.Sprintf(`%s:%s`, addr, port)

				utils.Logger.Log.Info(`service`,
					zap.String(`start`, `succ`),
					zap.String(`bind`, httpaddrs),
					zap.String(`time`, time.Now().Format(`2006-01-02 15:04:05`)),
				)
				control.Routers.SetApiGroup()
				if err = control.Routers.Run(addr, port); err != nil {
					cli.Exit(err.Error(), 1)
					return
				}
				return
			},
		},
	}
	app.Run(os.Args)
}
