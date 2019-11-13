package cfgargs

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// InitSrvConfig 加载服务配置信息
func InitSrvConfig(BuildVersion string, usrFlagParse func()) (*SrvConfig, error) {
	srvCfg := newSrvConfig()

	// 通过本地文件加载默认配置
	srvCfg.syncLocal()

	// 通过命令行加载配置
	version := flag.Bool("v", false, "(default false)")
	cfgaddr := flag.String("cfgaddr", "", "type dir")
	if nil != usrFlagParse {
		usrFlagParse()
	}
	flag.Parse()
	isVersion(*version, BuildVersion)
	srvCfg.CfgCenter.Addr = *cfgaddr

	// 通过远程配置中心同步配置
	err := srvCfg.syncConfig()
	if nil != err {
		return nil, err
	}

	fmt.Println(srvCfg.Print())
	return srvCfg, nil
}

func isVersion(v bool, BuildVersion string) {
	if v {
		versions := strings.Split(BuildVersion, "*")
		fmt.Printf("VERSION    : %s\n", strings.Replace(versions[0], "_", " ", -1))
		fmt.Printf("BUILD BY   : %s\n", strings.Replace(versions[1], "_", " ", -1))
		fmt.Printf("BUILD TIME : %s\n", strings.Replace(versions[2], "_", " ", -1))
		fmt.Printf("ON MACHINE : %s\n", strings.Replace(versions[3], "_", " ", -1))
		os.Exit(0)
	}
}
