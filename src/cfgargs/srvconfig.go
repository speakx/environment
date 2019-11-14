package cfgargs

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"

	"gopkg.in/yaml.v2"
)

// SrvConfig 通过默认配置+cfgaddr同步后得到的配置信息
type SrvConfig struct {
	Addr string `yaml:"addr"`
	Name string `yaml:"name"`
	Log  struct {
		Path    string `yaml:"path"`    // default: ./log/(appname).log
		Console bool   `yaml:"console"` // default: false
		Level   string `yaml:"level"`   // default: info
	}
	CfgCenter struct {
		Addr string `yaml:"addr"`
	}
	Cache struct {
		Path     string `yaml:"path"`     // default: ./cache
		MMapSize int    `yaml:"mmapsize"` // default: 1024*1024*1
		DataSize int    `yaml:"datasize"` // default: 1024*8
		PreAlloc int    `yaml:"prealloc"` // default: 100
	}
	Dump struct {
		Interval int    `yaml:"interval"` // default: 5
		Addr     string `yaml:"addr"`     // default:
	}
	RemoteCfg map[string]string
}

func newSrvConfig() *SrvConfig {
	return &SrvConfig{
		RemoteCfg: make(map[string]string),
	}
}

func (s *SrvConfig) syncLocal() error {
	binPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	appPath, _ := filepath.Split(binPath)
	fmt.Println(appPath)
	yamlFile := filepath.Join(appPath, "/etc/local.yaml")

	data, err := ioutil.ReadFile(yamlFile)
	if nil != err {
		fmt.Printf("load local yaml err:%v path:%v\n", err, yamlFile)
		return err
	}

	err = yaml.Unmarshal([]byte(data), s)
	if nil != err {
		fmt.Printf("unmarshal local yaml err:%v path:%v\n", err, yamlFile)
		return err
	}

	if "" == s.Log.Level {
		s.Log.Level = "INFO"
	}
	if "" == s.Log.Path {
		s.Log.Path = filepath.Join(appPath, fmt.Sprintf("/log/%v.log", path.Base(os.Args[0])))
	}

	if "" == s.Cache.Path {
		s.Cache.Path = filepath.Join(appPath, fmt.Sprintf("/cache"))
	}
	if 0 == s.Cache.MMapSize {
		s.Cache.MMapSize = 1024 * 1024 * 1
	}
	if 0 == s.Cache.DataSize {
		s.Cache.DataSize = 1024 * 8
	}
	if 0 == s.Cache.PreAlloc {
		s.Cache.PreAlloc = 100
	}

	if 0 == s.Dump.Interval {
		s.Dump.Interval = 5
	}
	return nil
}

func (s *SrvConfig) syncConfig() error {
	return nil
}

// Print 把SrvConfig格式化为字符串信息
func (s *SrvConfig) Print() string {
	var buf bytes.Buffer
	v := reflect.ValueOf(*s)
	t := reflect.TypeOf(*s)
	buf.WriteString("Config --------------------------------\r\n")
	s.printField(&buf, 1, "", v, t)
	buf.WriteString("Config --------------------------------")
	return buf.String()
}

func (s *SrvConfig) printField(buf *bytes.Buffer, depth int, parentField string, v reflect.Value, t reflect.Type) {
	nameFmt := fmt.Sprintf("%%- %vv: %%v\r\n", 16)
	for i := 0; i < v.NumField(); i++ {
		switch v.Field(i).Kind() {
		case reflect.Struct:
			fieldType := t.Field(i)
			v1 := reflect.ValueOf(v.Field(i).Interface())
			t1 := reflect.TypeOf(v.Field(i).Interface())
			s.printField(buf, depth+1, fmt.Sprintf("%v%v.", parentField, fieldType.Name), v1, t1)
		default:
			fieldType := t.Field(i)
			buf.WriteString(fmt.Sprintf(nameFmt, fmt.Sprintf("%v%v", parentField, fieldType.Name), v.Field(i).Interface()))
		}
	}
}
