/*
Star Nepgear BOT
===
High-Speed Group Protective BOT for LINE.

Copyright(c) 2020 Star Inc. All Rights Reserved.
The software is licensed under Apache License 2.0.
*/
package main

import (
	"fmt"
	api "github.com/star-inc/NepCoreO"
	core "github.com/star-inc/olsb_cores/libs/NepCoreO"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

type configInterface struct {
	LINE *struct {
		Server *struct {
			ClientPath string `yaml:"Command_Path"`
			ListenPath string `yaml:"LongPoll_path"`
		} `yaml:"Server"`
		Account *struct {
			AuthToken string `yaml:"X-Line-Access"`
		} `yaml:"Account"`
	} `yaml:"LINE"`
}

const (
	commandPrefix         = "$"
	configPath            = "config.yaml"
	routineTimeout        = 100
	fetchOpeartionsLength = 50
)

var (
	config *configInterface
	client *api.ClientInterface
	listen *api.ClientInterface
)

func main() {
	declare()
	readConfig()
	connect()
	ctx, _ := api.SetRoutine(routineTimeout)
	revision, _ := client.TalkServiceClient.GetLastOpRevision(ctx)
	for {
		ctx, _ = api.SetRoutine(routineTimeout)
		ops, _ := listen.TalkServiceClient.FetchOperations(ctx, revision, fetchOpeartionsLength)
		go func(ops []*core.Operation) {
			for _, task := range ops {
				switch task.Type {
				case core.OpType_RECEIVE_MESSAGE:
					go messageHandle(task)
					break
				}
			}
		}(ops)
		if len(ops) > 1 {
			revision = func(ops []*core.Operation) int64 {
				if ops[len(ops)-1].Revision > ops[len(ops)-2].Revision {
					return ops[len(ops)-1].Revision
				}
				return ops[len(ops)-2].Revision
			}(ops)
		}
	}
}

func declare() {
	fmt.Println("")
	fmt.Println("Star Nepgear BOT")
	fmt.Println("===")
	fmt.Println("High-Speed Group Protective BOT for LINE.")
	fmt.Println("\nCopyright(c) 2020 Star Inc. All Rights Reserved.")
	fmt.Println("The software is licensed under Apache License 2.0.")
	fmt.Println("")
	fmt.Println("")
}

func readConfig() {
	yamlFile, _ := os.Open(configPath)
	defer yamlFile.Close()
	srcYAML, _ := ioutil.ReadAll(yamlFile)
	_ = yaml.Unmarshal(srcYAML, &config)
}

func sendToWho(op *core.Operation) string {
	switch op.Message.ToType {
	case core.MIDType_USER:
		return op.Message.From_
	case core.MIDType_ROOM:
	case core.MIDType_GROUP:
		return op.Message.To
	}
	return ""
}

func connect() {
	client = api.NewClientInterface(config.LINE.Server.ClientPath)
	listen = api.NewClientInterface(config.LINE.Server.ListenPath)
	client.Authorize(config.LINE.Account.AuthToken)
	listen.Authorize(config.LINE.Account.AuthToken)
}

func messageHandle(op *core.Operation) {
	switch op.Message.ContentType {
	case core.ContentType_NONE:
		text(op)
		break
	}
}

func text(op *core.Operation) {
	msg := strings.Split(op.Message.Text, " ")
	switch msg[0] {
	case commandPrefix + "help":
		client.SendText(sendToWho(op), "Star Nepgear BOT")
		break
	}
}
