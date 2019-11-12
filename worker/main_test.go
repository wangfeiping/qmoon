// Copyright 2018 The QOS Authors

package worker

import (
	"testing"

	"github.com/QOSGroup/qmoon/config"
	"github.com/QOSGroup/qmoon/lib/qstarscli"
	"github.com/QOSGroup/qmoon/lib/tmcli"
	"github.com/QOSGroup/qmoon/models"
)

func TestMain(m *testing.M) {
	dbTest, err := models.NewTestEngine(config.TestDBConfig())
	if err != nil {
		panic(err)
	}
	defer dbTest.Close()

	//err := service.CreateNode("qstars", "http://192.168.1.223:26657", "", nil)
	//if err != nil {
	//	panic(err)
	//}

	tq := qstarscli.NewTestQstarsServer()
	defer tq.Close()

	tts := tmcli.NewTestTmServer()
	defer tts.Close()

	m.Run()
}
