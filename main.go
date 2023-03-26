package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"proxy/db"
	"proxy/log"
	"proxy/rpc"
	"proxy/trie"
	"strings"
	"syscall"
	"time"

	cli "gopkg.in/urfave/cli.v1"
	yaml "gopkg.in/yaml.v2"
)

var (
	OriginCommandHelpTemplate = `{{.Name}}{{if .Subcommands}} command{{end}}{{if .Flags}} [command options]{{end}} {{.ArgsUsage}}
{{if .Description}}{{.Description}}
{{end}}{{if .Subcommands}}
SUBCOMMANDS:
  {{range .Subcommands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
  {{end}}{{end}}{{if .Flags}}
OPTIONS:
{{range $.Flags}}   {{.}}
{{end}}
{{end}}`
)
var app *cli.App

func init() {
	app = cli.NewApp()
	app.Version = "v1.0.0"
	app.Commands = []cli.Command{
		commandStart,
	}

	cli.CommandHelpTemplate = OriginCommandHelpTemplate
}

var (
	configPathFlag = cli.StringFlag{
		Name:  "config",
		Usage: "config path",
		Value: "./config.yml",
	}
)

var commandStart = cli.Command{
	Name:  "start",
	Usage: "start loading contract gas fee",
	Flags: []cli.Flag{
		configPathFlag,
	},
	Action: Start,
}

type ProxyConfig struct {
	Port             string            `yaml:"port"`
	OpenAiKey        []string          `yaml:"openai_key"`
	MaxPendingLength int               `yaml:"max_pending"`
	Host             string            `yaml:"host"`
	ModelConfig      map[string]string `yaml:"bs_model"`
	MongoURI         string            `yaml:"mongo_uri"`
	Sensitive        string            `yaml:"sensitive"`
}

func Start(ctx *cli.Context) {
	conf := loadConfig(ctx)
	if conf.Host != "" {
		rpc.Host = conf.Host
	}
	rpc.InitRpcService(conf.Port, conf.OpenAiKey, conf.MaxPendingLength, conf.ModelConfig)
	contx := context.Background()
	err = rpc.RpcServer.Start(contx)
	if err != nil {
		log.Fatal(err)
	}
	waitToExit()
}

func loadConfig(ctx *cli.Context) ProxyConfig {
	var proxyConfig ProxyConfig
	if ctx.IsSet(configPathFlag.Name) {
		configPath := ctx.String(configPathFlag.Name)
		b, err := ioutil.ReadFile(configPath)
		if err != nil {
			log.Fatal("read config error", err)
		}
		err = yaml.Unmarshal(b, &proxyConfig)
		if err != nil {
			log.Fatal(err)
		}
	}
	return proxyConfig
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func waitToExit() {
	exit := make(chan bool, 0)
	sc := make(chan os.Signal, 1)
	if !signal.Ignored(syscall.SIGHUP) {
		signal.Notify(sc, syscall.SIGHUP)
	}
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for sig := range sc {
			fmt.Printf("received exit signal:%v", sig.String())
			close(exit)
			break
		}
	}()
	<-exit
}
