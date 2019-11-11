package logger

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	stdLog "log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const severityChar = "DIWEF"
const bufferSize = 256 * 1024
const flushInterval = 30 * time.Second

var (
	pid  = os.Getpid()
	host = "unknownhost"
)

type Config struct {
	Console      bool   `yaml:"console"`
	File         bool   `yaml:"file"`
	ConsoleLevel string `yaml:"consoleLevel"`
	FileLevel    string `yaml:"fileLevel"`
	Path         string `yaml:"path"`
	Type         string `yaml:"type"`
	Strategy     string `yaml:"strategy"`
	Maxsize      int64  `yaml:"maxsize"`
}

// InitLogger 初始化日志模块
func InitLogger(fileFullPath string, consoleOut bool, level string) {
	h, err := os.Hostname()
	if err == nil {
		host = shortHostname(h)
	}
	conf, err := load()
	if err != nil {
		logging.consoleOut = ConsoleOut
		logging.consoleLevel = ConsoleLevel
		logging.fileOut = FileOut
		logging.fileLevel = FileLevel
		logging.logDir = LogDir
		logging.logType = FileType
		logging.maxsize = MaxSize
		logging.strategy = Strategy
	} else {
		// use input args
		conf.Path, _ = filepath.Abs(filepath.Dir(fileFullPath))
		conf.ConsoleLevel = level
		conf.FileLevel = level
		conf.Console = consoleOut

		logging.fileName = path.Base(fileFullPath)
		logging.consoleOut = ConsoleOut
		logging.consoleLevel = ConsoleLevel
		logging.fileOut = FileOut
		logging.fileLevel = FileLevel
		logging.logDir = LogDir
		logging.logType = FileType
		logging.maxsize = MaxSize
		logging.consoleOut = conf.Console
		v := strings.ToUpper(conf.ConsoleLevel)
		if v == "DEBUG" {
			logging.consoleLevel = 0
		} else if v == "INFO" {
			logging.consoleLevel = 1
		} else if v == "WARN" {
			logging.consoleLevel = 2
		} else if v == "ERROR" {
			logging.consoleLevel = 3
		} else if v == "FATAL" {
			logging.consoleLevel = 4
		} else {
			logging.consoleLevel = 1
		}

		logging.fileOut = conf.File
		v = strings.ToUpper(conf.FileLevel)
		if v == "DEBUG" {
			logging.fileLevel = 0
		} else if v == "INFO" {
			logging.consoleLevel = 1
		} else if v == "WARN" {
			logging.fileLevel = 2
		} else if v == "ERROR" {
			logging.fileLevel = 3
		} else if v == "FATAL" {
			logging.fileLevel = 4
		} else {
			logging.fileLevel = 1
		}
		logging.logDir = conf.Path
		v = strings.ToLower(conf.Type)
		if v == "size" || v == "date" {
			logging.logType = v
			if conf.Maxsize < 10 {
				conf.Maxsize = 10
			}
			logging.maxsize = conf.Maxsize * 1024 * 1024
		} else {
			v = "default"
			logging.maxsize = MaxSize
		}
		var format string = "2006-01-02"
		if logging.logType == "date" {
			switch conf.Strategy {
			case "yyyy-mm-dd hh:MM":
				format = "2006-01-02 15:04"
			case "yyyy-mm-dd hh":
				format = "2006--01-02 15"
			case "yyyy-mm-dd":
				format = "2006-01-02"
			case "yyyy-mm":
				format = "2006-01"
			case "yyyy":
				format = "2006"
			default:
				format = "2006-01-02"
			}
		} else {
			format = "2006-01-02"
		}
		logging.strategy = format
	}
	logging.createFiles()
	go logging.flushDaemon()
}

func load() (c *Config, err error) {
	// yamlFile, err := ioutil.ReadFile("etc/log.yaml")
	// if err != nil {
	// 	//stdLog.Println("read log config file error:", err)
	// 	return nil, err
	// }
	// err = yaml.Unmarshal(yamlFile, &c)
	// if err != nil {
	// 	stdLog.Println("format log.yaml error:", err)
	// 	return nil, err
	// }
	// return
	c = &Config{
		Console:      false,
		File:         true,
		ConsoleLevel: "WARN",
		FileLevel:    "WARN",
		Path:         "./log",
		Type:         "size",
		Strategy:     "yyyy-mm-dd hh:MM",
		Maxsize:      50,
	}
	return
}

const (
	debugLog int = iota
	infoLog
	warningLog
	errorLog
	fatalLog
)

type flushSyncWriter interface {
	Flush() error
	Sync() error
	io.Writer
}

func createLogDir() {
	if logging.logDir != "" {
		err := os.MkdirAll(logging.logDir, 0777)
		if err != nil {
			fmt.Println("日志目录失败")
		}
	} else {
		logging.logDir = LogDir
		err := os.MkdirAll(logging.logDir, 0777)
		if err != nil {
			fmt.Println("日志目录失败")
		}
	}
}

func shortHostname(hostname string) string {
	if i := strings.Index(hostname, "."); i >= 0 {
		return hostname[:i]
	}
	return hostname
}

var onceLogDirs sync.Once

func create(fname string) (f *os.File, filename string, n int, err error) {
	onceLogDirs.Do(createLogDir)
	if _, err := os.Stat(fname); err != nil {
		f, err = os.Create(fname)
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "Log file created at: %s\n", time.Now().Format("2006/01/02 15:04:05"))
		fmt.Fprintf(&buf, "Running on machine: %s\n", host)
		fmt.Fprintf(&buf, "Binary: Built with %s %s for %s/%s\n", runtime.Compiler, runtime.Version(), runtime.GOOS, runtime.GOARCH)
		fmt.Fprintf(&buf, "Log line format: [IWEF]mmdd hh:mm:ss.uuuuuu threadid file:line] msg\n")
		n, _ = f.Write(buf.Bytes())
		return f, filename, n, err
	} else {
		f, err = os.OpenFile(fname, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		t, _ := f.Seek(0, os.SEEK_END)
		return f, filename, int(t), err
	}
}

const (
	ConsoleOut   bool   = true
	FileOut      bool   = true
	ConsoleLevel int    = 0
	FileLevel    int    = 0
	MaxSize      int64  = 1024 * 1024 * 100
	FileType     string = "size"
	LogDir       string = "./logs"
	Strategy     string = "2006-01-02"
)

func (self *loggingT) formatLogName() (fileName string) {
	fileName = self.fileName + time.Now().Format(self.strategy) + ".log"
	fileName = filepath.Join(self.logDir, fileName)
	if fileName == self.logName {
		logging.files = append(logging.files, self.logName)
		l := len(self.files)
		for index, name := range logging.files {
			os.Rename(name, fileName+strconv.FormatInt(int64(l-index), 10))
		}
	} else {
		logging.files = make([]string, 0)
	}
	logging.logName = fileName
	return
}

func Flush() {
	logging.lockAndFlushAll()
}

type loggingT struct {
	consoleOut   bool
	fileOut      bool
	consoleLevel int
	fileLevel    int
	fileName     string
	logDir       string
	logType      string
	logName      string
	strategy     string
	namestyle    string
	maxsize      int64
	freeList     *buffer
	files        []string
	freeListMu   sync.Mutex
	mu           sync.Mutex
	file         flushSyncWriter
	filterLength int32
}

type buffer struct {
	bytes.Buffer
	tmp  [64]byte
	next *buffer
}

var logging loggingT

func (l *loggingT) getBuffer() *buffer {
	l.freeListMu.Lock()
	b := l.freeList
	if b != nil {
		l.freeList = b.next
	}
	l.freeListMu.Unlock()
	if b == nil {
		b = new(buffer)
	} else {
		b.next = nil
		b.Reset()
	}
	return b
}

func (l *loggingT) putBuffer(b *buffer) {
	if b.Len() >= 256 {
		return
	}
	l.freeListMu.Lock()
	b.next = l.freeList
	l.freeList = b
	l.freeListMu.Unlock()
}

var timeNow = time.Now

func (l *loggingT) header(s int, depth int) (*buffer, string, int) {
	_, file, line, ok := runtime.Caller(3 + depth)
	if !ok {
		file = "???"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return l.formatHeader(s, file, line), file, line
}

func (l *loggingT) formatHeader(s int, file string, line int) *buffer {
	now := timeNow()
	if line < 0 {
		line = 0
	}
	if s > fatalLog {
		s = infoLog
	}
	buf := l.getBuffer()
	year, month, day := now.Date()
	hour, minute, second := now.Clock()
	// Lmmdd hh:mm:ss.uuuuuu threadid file:line]
	buf.tmp[0] = severityChar[s]
	buf.tmp[1] = ' '
	buf.nDigits(4, 2, year, ' ')
	buf.twoDigits(6, int(month))
	buf.twoDigits(8, day)
	buf.tmp[10] = ' '
	buf.twoDigits(11, hour)
	buf.tmp[13] = ':'
	buf.twoDigits(14, minute)
	buf.tmp[16] = ':'
	buf.twoDigits(17, second)
	buf.tmp[19] = '.'
	buf.nDigits(6, 20, now.Nanosecond()/1000, '0')
	buf.tmp[26] = ' '
	buf.nDigits(6, 27, pid, ' ')
	buf.tmp[33] = ' '
	buf.Write(buf.tmp[:34])
	buf.WriteString(file)
	buf.tmp[0] = ':'
	n := buf.someDigits(1, line)
	buf.tmp[n+1] = ']'
	buf.tmp[n+2] = ' '
	buf.Write(buf.tmp[:n+3])
	return buf
}

const digits = "0123456789"

func (buf *buffer) twoDigits(i, d int) {
	buf.tmp[i+1] = digits[d%10]
	d /= 10
	buf.tmp[i] = digits[d%10]
}

func (buf *buffer) nDigits(n, i, d int, pad byte) {
	j := n - 1
	for ; j >= 0 && d > 0; j-- {
		buf.tmp[i+j] = digits[d%10]
		d /= 10
	}
	for ; j >= 0; j-- {
		buf.tmp[i+j] = pad
	}
}

func (buf *buffer) someDigits(i, d int) int {
	j := len(buf.tmp)
	for {
		j--
		buf.tmp[j] = digits[d%10]
		d /= 10
		if d == 0 {
			break
		}
	}
	return copy(buf.tmp[i:], buf.tmp[j:])
}

func (l *loggingT) println(s int, args ...interface{}) {
	buf, file, line := l.header(s, 0)
	fmt.Fprintln(buf, args...)
	l.output(s, buf, file, line, false)
}

func (l *loggingT) print(s int, args ...interface{}) {
	l.printDepth(s, 1, args...)
}

func (l *loggingT) printDepth(s int, depth int, args ...interface{}) {
	buf, file, line := l.header(s, depth)
	fmt.Fprint(buf, args...)
	if buf.Bytes()[buf.Len()-1] != '\n' {
		buf.WriteByte('\n')
	}
	l.output(s, buf, file, line, false)
}

func (l *loggingT) printf(s int, format string, args ...interface{}) {
	buf, file, line := l.header(s, 0)
	fmt.Fprintf(buf, format, args...)
	if buf.Bytes()[buf.Len()-1] != '\n' {
		buf.WriteByte('\n')
	}
	l.output(s, buf, file, line, false)
}

func (l *loggingT) output(s int, buf *buffer, file string, line int, alsoToStderr bool) {
	l.mu.Lock()
	data := buf.Bytes()
	if alsoToStderr || (l.consoleOut && s >= l.consoleLevel) {
		os.Stderr.Write(data)
	}
	if l.fileOut && s >= l.fileLevel {
		l.file.Write(data)
	}
	if s == fatalLog {
		l.mu.Unlock()
		timeoutFlush(10 * time.Second)
		os.Exit(255)
	}
	l.putBuffer(buf)
	l.mu.Unlock()
}

func timeoutFlush(timeout time.Duration) {
	done := make(chan bool, 1)
	go func() {
		Flush()
		done <- true
	}()
	select {
	case <-done:
	case <-time.After(timeout):
		fmt.Fprintln(os.Stderr, "glog: Flush took longer than", timeout)
	}
}

type syncBuffer struct {
	logger *loggingT
	*bufio.Writer
	file   *os.File
	sev    int
	nbytes int64
}

func (sb *syncBuffer) Sync() error {
	return sb.file.Sync()
}

func (sb *syncBuffer) Write(p []byte) (n int, err error) {
	if sb.nbytes+int64(len(p)) >= logging.maxsize {
		if err := sb.rotateFile(logging.formatLogName()); err != nil {
			stdLog.Println("create log file error:", err)
		}
		sb.nbytes = int64(n)
	} else {
		sb.nbytes += int64(len(p))
	}
	n, err = sb.Writer.Write(p)
	if err != nil {
		stdLog.Println("create log file error:", err)
	}
	return
}

func (sb *syncBuffer) rotateFile(name string) error {
	if sb.file != nil {
		sb.Flush()
		sb.file.Close()
	}
	var err error
	var n int
	sb.file, _, n, err = create(name)
	if err != nil {
		return err
	}
	sb.Writer = bufio.NewWriterSize(sb.file, bufferSize)
	sb.nbytes += int64(n)
	return err
}

func (l *loggingT) createFiles() error {
	if l.fileOut {
		sb := &syncBuffer{logger: l}
		if err := sb.rotateFile(logging.formatLogName()); err != nil {
			return err
		}
		l.file = sb
	}
	return nil
}

func (l *loggingT) flushDaemon() {
	for _ = range time.NewTicker(flushInterval).C {
		l.lockAndFlushAll()
	}
}

func (l *loggingT) lockAndFlushAll() {
	l.mu.Lock()
	l.flushAll()
	l.mu.Unlock()
}

func (l *loggingT) flushAll() {
	if l.file != nil {
		l.file.Flush()
		l.file.Sync()
	}
}

func Info(args ...interface{}) {
	logging.print(infoLog, args...)
}

func InfoDepth(depth int, args ...interface{}) {
	logging.printDepth(infoLog, depth, args...)
}

func Infoln(args ...interface{}) {
	logging.println(infoLog, args...)
}

func Infof(format string, args ...interface{}) {
	logging.printf(infoLog, format, args...)
}

func Debug(args ...interface{}) {
	logging.print(debugLog, args...)
}

func DebugDepth(depth int, args ...interface{}) {
	logging.printDepth(debugLog, depth, args...)
}

func Debugln(args ...interface{}) {
	logging.println(debugLog, args...)
}

func Debugf(format string, args ...interface{}) {
	logging.printf(debugLog, format, args...)
}

func Warn(args ...interface{}) {
	logging.print(warningLog, args...)
}

func WarnDepth(depth int, args ...interface{}) {
	logging.printDepth(warningLog, depth, args...)
}

func Warnln(args ...interface{}) {
	logging.println(warningLog, args...)
}

func Warnf(format string, args ...interface{}) {
	logging.printf(warningLog, format, args...)
}

func Error(args ...interface{}) {
	logging.print(errorLog, args...)
}

func ErrorDepth(depth int, args ...interface{}) {
	logging.printDepth(errorLog, depth, args...)
}

func Errorln(args ...interface{}) {
	logging.println(errorLog, args...)
}

func Errorf(format string, args ...interface{}) {
	logging.printf(errorLog, format, args...)
}

func Fatal(args ...interface{}) {
	logging.print(fatalLog, args...)
}

func FatalDepth(depth int, args ...interface{}) {
	logging.printDepth(fatalLog, depth, args...)
}

func Fatalln(args ...interface{}) {
	logging.println(fatalLog, args...)
}

func Fatalf(format string, args ...interface{}) {
	logging.printf(fatalLog, format, args...)
}
