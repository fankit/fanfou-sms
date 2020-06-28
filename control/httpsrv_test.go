package control

import (
	"bytes"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"testing"
	"tsmsrv/conf"
	"tsmsrv/utils"
)

func TestInitRouters(t *testing.T) {
	var err error

	if err = conf.NewConfParse(`/Users/osx/Documents/tsmsrv/tsmsrv.conf`); err != nil {
		t.Error(`t1`, err.Error())
	}
	if err = utils.NewLoggerMgr(`info`); err != nil {
		t.Error(err.Error())
	}


	if err = NewRouters(true);err != nil {
		t.Error(err.Error())
	}
	utils.Logger.Log.Info(`service`, zap.String(`start`, `sucessful`))
	Routers.SetApiGroup()
	if err = Routers.Run("127.0.0.1", "8080"); err != nil {
		t.Error(err.Error())
	}
}

func TestRequestRouter(t *testing.T)  {
	var err error
	var fd []byte
	var r *http.Request
	var rs *http.Response
	var bs []byte
	if fd, err = ioutil.ReadFile("/Users/osx/Documents/tsmsrv/request.json"); err != nil {
		t.Error(err.Error())
	}
	fmt.Println(string(fd))
	if r, err = http.NewRequest(http.MethodPost, "http://127.0.0.1:8080/v1/sms", bytes.NewBuffer(fd)); err != nil {
		t.Error(err.Error())
	}
	client := http.Client{}
	if rs ,  err = client.Do(r); err != nil {
		t.Error(err.Error())
	}
	defer rs.Body.Close()
	if bs , err = ioutil.ReadAll(rs.Body); err != nil {
		t.Error(err.Error())
	}
	fmt.Println(string(bs))

}