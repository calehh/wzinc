package rpc

import (
	"context"

	"net/http"

	"strings"
	"sync"
	"time"
	"wzinc/log"

	"github.com/gin-gonic/gin"
)

const InternalError = "internal server error"

const (
	Success            = 200
	ErrorCodeUnknow    = -500
	ErrorCodeReadReq   = -501
	ErrorCodeParseReq  = -502
	ErrorCodeUnmarshal = -503
)

const (
	HealthCheckUrl = "/health"
	QuestionUrl    = "/api"
)

var SessionCookieName = "session_id"

var Host = "127.0.0.1"

const (
	Avalible = 1
	InUse    = 2
	Down     = 0
)

const (
	CheckRelayHealthInterval = time.Second * 10
	WaitForAnswer            = time.Minute * 3
	MaxRetry                 = 0
)

var once sync.Once

var client *http.Client

var RpcServer *Service

type Service struct {
	port     string
	zincUrl  string
	username string
	password string
}

func InitRpcService(url, port, username, password string) {
	once.Do(func() {
		client = &http.Client{Timeout: time.Minute * 3}
		RpcServer = &Service{
			port:     port,
			zincUrl:  url,
			username: username,
			password: password,
		}
	})
}

type LoggerMy struct {
}

func (*LoggerMy) Write(p []byte) (n int, err error) {
	msg := strings.TrimSpace(string(p))
	if strings.Index(msg, `"/healthcheck"`) > 0 {
		return
	}
	log.Debug(msg)
	return
}

func (c *Service) Start(ctx context.Context) error {
	//start gin
	gin.DefaultWriter = &LoggerMy{}
	r := gin.Default()

	//cors middleware
	r.SetTrustedProxies(nil)
	r.GET("/healthcheck", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	r.POST("/api/input", c.HandleInput)
	r.POST("/api/delete", c.HandleDelete)
	r.POST("/api/query", c.HandleQuery)
	address := "0.0.0.0:" + c.port

	r.Run(address)
	log.Info("start rpc on port:" + c.port)
	return nil
}

type Resp struct {
	ResultCode int         `json:"ret"`
	ResultMsg  string      `json:"msg"`
	ResultBody interface{} `json:"data"`
}

func (s *Service) HandleInput(c *gin.Context) {
	rep := Resp{
		ResultCode: ErrorCodeUnknow,
		ResultMsg:  "",
		ResultBody: "",
	}
	defer func() {
		if rep.ResultCode == Success {
			c.JSON(http.StatusOK, rep)
		} else {
			c.JSON(http.StatusInternalServerError, rep)
		}
	}()

	//get session state
	// msg := c.PostForm("message")
	// if msg == "" {
	// 	log.Debug("no message")
	// 	return
	// }
	// header, err := c.FormFile("document")
	// if err != nil {
	// 	log.Error("get form file error", err)
	// 	return
	// }
	// file, err := header.Open()
	// if err != nil {
	// 	log.Error("open file err", err)
	// 	return
	// }
	// defer file.Close()
	// docB, err := ioutil.ReadAll(file)
	// if err != nil {
	// 	log.Error("read file", err)
	// 	return
	// }

	return
}

func (s *Service) HandleDelete(c *gin.Context) {
	rep := Resp{
		ResultCode: ErrorCodeUnknow,
		ResultMsg:  "",
		ResultBody: "",
	}
	defer func() {
		if rep.ResultCode == Success {
			c.JSON(http.StatusOK, rep)
		} else {
			c.JSON(http.StatusInternalServerError, rep)
		}
	}()

	//get session state

	return
}

func (s *Service) HandleQuery(c *gin.Context) {
	rep := Resp{
		ResultCode: ErrorCodeUnknow,
		ResultMsg:  "",
		ResultBody: "",
	}
	defer func() {
		if rep.ResultCode == Success {
			c.JSON(http.StatusOK, rep)
		} else {
			c.JSON(http.StatusInternalServerError, rep)
		}
	}()

	//get session state
	return
}
