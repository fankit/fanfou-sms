package control

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tsmsrv/conf"
	"tsmsrv/utils"
)

type RouterMgr struct {
	engine *gin.Engine
}

var Routers *RouterMgr

func NewRouters(d bool) (err error) {

	gin.SetMode(gin.ReleaseMode)

	if d {
		gin.SetMode(gin.DebugMode)
	}


	r := gin.New()
	r.NoRoute(SourceNotFound)
	r.Use(LoggerMidderWare(), gin.Recovery())

	Routers = &RouterMgr{engine: r}
	return
}

func (r *RouterMgr) SetApiGroup() {
	var v1 *gin.RouterGroup
	r.engine.Handle("GET", `/ping`, r.ping)
	v1 = r.engine.Group(`/v1`)
	v1.Use()
	{
		v1.GET(`/ping`, r.v1ping)
		v1.POST(`/sms`, r.WebHook)
	}
}

func (router *RouterMgr) Run(addr, port string) (err error) {

	var (
		httpAddress string
		httpServer  *http.Server
	)
	httpAddress = fmt.Sprintf(`%s:%s`, addr, port)
	httpServer = &http.Server{
		Addr:           httpAddress,
		Handler:        router.engine,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			os.Exit(-1)
		}
	}()

	if err = router.exitHttpServer(httpServer); err != nil {
		return
	}
	return
}

func (router *RouterMgr) exitHttpServer(httpserver *http.Server) (err error) {
	var (
		sigch   chan os.Signal
		cancel  context.CancelFunc
		cxt     context.Context
	)

	sigch = make(chan os.Signal)
	signal.Notify(sigch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	sig := <-sigch
	cxt, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = httpserver.Shutdown(cxt); err != nil {
		return
	}
	fmt.Println(`Signal:`, sig)
	return
}

func (r *RouterMgr) ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, fmt.Sprintf(`ok`))
}

func (r *RouterMgr) v1ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, `ok`)
}

func (r *RouterMgr) WebHook(ctx *gin.Context) {
	var b Body
	ctx.BindJSON(&b)
	utils.Logger.Log.Info(`webhook`, zap.String("Title", b.Title),
		zap.Int("RuleID", b.RuleID),
		zap.String("RuleName", b.RuleName),
		zap.String("RuleURL", b.RuleURL),
		zap.String("State", string(b.State)),
		zap.String("ImageURL", b.ImageURL),
		zap.String("Message", b.Message),
		zap.Any("EvalMatches", b.EvalMatches),
	)
	phns := conf.GlobConfig.SmsSection().Key(`phone_numbers`).Strings(`,`)

	var alertContext string

	switch st := string(b.State); st {
	case `alerting`:
		alertContext = b.RuleName + "\n" + `告警详情:` + b.Message  + "\n" + `[告警]`

	case `ok`:
		alertContext = b.RuleName + "\n" + `告警详情:` + b.Message + "\n" + `[已恢复]`
	}

	wcontext := []string{alertContext}

	if err := NewSendSmsSrv().SendSms(phns, wcontext); err != nil {
		utils.Logger.Log.Error(`SendSms`, zap.String(`sendsms`, err.Error()))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "succ",
	})

	/*	ctx.JSON(http.StatusOK, gin.H{
			"Title":       b.Title,
			"RuleID":      b.RuleID,
			"RuleName":    b.RuleName,
			"RuleURL":     b.RuleURL,
			"State":       b.State,
			"ImageURL":    b.ImageURL,
			"Message":     b.Message,
			"EvalMatches": b.EvalMatches,
		})
	*/
}
