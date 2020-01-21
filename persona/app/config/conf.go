package config

import (
	"fmt"
	"os"

	"github.com/Unknwon/goconfig"
)

var Cfg *goconfig.ConfigFile

const (
	// DebugMode indicates gin mode is debug.
	CfgDebug = "debug"
	CfgAlpha = "alpha"
	// ReleaseMode indicates gin mode is release.
	CfgRelease = "release"
	// TestMode indicates gin mode is test.
	TestMode = "test"
)

func InitCfg(evn *string) (err error) {
	fmt.Println("cfg InitCfg evn", *evn)
	var replacePath string = ""
	if *evn != CfgDebug {
		replacePath = "app/config/config." + *evn + ".ini"
		fileInfo, _ := os.Stat(replacePath)
		if fileInfo == nil {
			return fmt.Errorf("err InitCfg : %s not have", replacePath)
		}
	}
	if replacePath == "" {
		Cfg, err = goconfig.LoadConfigFile("app/config/config.default.ini")
	} else {
		Cfg, err = goconfig.LoadConfigFile("app/config/config.default.ini", replacePath)
	}
	if err != nil {
		return fmt.Errorf("lew InitCfg  LoadConfigFile have err", err.Error())
	} else {
		return nil
	}

}
