package argsflag

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
)

// FlagConfig 通用Flag的解析结果
type FlagConfig struct {
	LogPath     string
	LogFullPath string
	LogConsole  bool
	LogLevel    string
	Addr        string
	DBPath      string
}

func (f *FlagConfig) checkFlags() bool {
	if "" == f.LogPath || "" == f.Addr || "" == f.DBPath || "" == f.LogLevel {
		return false
	}
	return true
}

// PasteFlag 包装一些通用的Flag处理方法
func PasteFlag(BuildVersion string, userPasteFuc func()) *FlagConfig {
	f := &FlagConfig{}
	version := flag.Bool("v", false, "(default false)")
	logPath := flag.String("logpath", "", "type dir")
	logConsole := flag.Bool("logconsole", false, "(default false) print log to stdout")
	logLevel := flag.String("loglevel", "info", "{ debug | info | warn | error }")
	addr := flag.String("addr", "", "type ip:port")
	dbPath := flag.String("dbpath", "", "type dir")
	if nil != userPasteFuc {
		userPasteFuc()
	}
	flag.Parse()

	f.LogPath = *logPath
	f.LogFullPath = fmt.Sprintf("%v/%v.log", f.LogPath, path.Base(os.Args[0]))
	f.LogConsole = *logConsole
	f.LogLevel = *logLevel
	f.Addr = *addr
	f.DBPath = *dbPath

	isVersion(*version, BuildVersion)
	if false == f.checkFlags() {
		flag.Usage()
		os.Exit(0)
		return f
	}
	return f
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
