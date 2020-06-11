package log

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)



// NewLogFile new log file
func NewLogFile(logPath string, isPrint bool) (d *log.Logger) {
	var (
		//f   *os.File
		rf  *rotatelogs.RotateLogs
		err error
	)
	d = log.New()

	// ensure director
	_dir := filepath.Dir(logPath)
	if _, _err := os.Stat(_dir); os.IsNotExist(_err) {
		if err = os.MkdirAll(_dir, os.ModePerm); err != nil {
			panic(err)
		}
	}

	// log file(s)
	if rf, err = rotatelogs.New(
		logPath+".%Y-%m-%d.log",
		//rotatelogs.WithLinkName(logPath),
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 文件最大保存时间
		rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
		// rotatelogs.WithRotationCount(RotationCount),
	); err == nil {
		if isPrint {
			d.Hooks.Add(lfshook.NewHook(
				lfshook.WriterMap{
					log.TraceLevel: rf,
					log.DebugLevel: rf,
					log.InfoLevel:  rf,
					log.WarnLevel:  rf,
					// log.ErrorLevel: rf,
					log.FatalLevel: rf,
					log.PanicLevel: rf,
				},
				&log.JSONFormatter{},
			))
		} else {
			d.Out = rf
		}

	} else {
		log.Warnln(err)
	}
	d.SetFormatter(&log.JSONFormatter{})
	// hook errors
	d.AddHook(&HookError{Filepath: logPath + ".error"})

	//if f, err = os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600); err == nil {
	//	d.Out = f
	//}

	return
}

// HookError hook error 错误钩子
type HookError struct {
	Filepath string
	Out      io.WriteCloser
}

// Levels level
func (h *HookError) Levels() []log.Level {
	return []log.Level{log.ErrorLevel}
}

func (h *HookError)SetFormat(logger *log.Logger,format log.Formatter)*HookError  {
	logger.SetFormatter(format)
	return h
}



// Fire fire
func (h *HookError) Fire(entry *log.Entry) error {
	if entry.Level == log.ErrorLevel && len(h.Filepath) > 0 {
		if _s, _err := entry.String(); _err == nil {
			if _, _err = os.Stat(h.Filepath); os.IsNotExist(_err) && h.Out != nil {
				h.Out.Close()
				h.Out = nil
			}
			if h.Out == nil {
				h.Out, _ = os.OpenFile(h.Filepath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			}
			if h.Out != nil {
				_, _ = h.Out.Write([]byte(_s))
			}
		} else {
			return _err
		}
	}
	return nil
}

// PanicRecover 统一处理panic
func PanicRecover(logger *log.Logger) {
	r := recover()
	if r == nil {
		return
	}
	if logger == nil {
		logger = log.New()
	}
	logger.Errorf(`[panic-recover] "%s" %v`, panicIdentify(), r)
}

// PanicRecoverError 统一处理panic, 并更新error
func PanicRecoverError(logger *log.Logger, err *error) {
	r := recover()
	if r == nil {
		return
	}
	if logger == nil {
		logger = log.New()
	}
	logger.Errorf(`[panic-recover] "%s" %v`, panicIdentify(), r)
	return
}

// panicIdentify  定位panic位置 参考自: https://gist.github.com/swdunlop/9629168
func panicIdentify() string {
	var (
		pc [16]uintptr
		n  = runtime.Callers(3, pc[:])
	)
	for _, _pc := range pc[:n] {
		fn := runtime.FuncForPC(_pc)
		if fn == nil {
			continue
		}
		_fnName := fn.Name()
		if strings.HasPrefix(_fnName, "runtime.") {
			continue
		}
		file, line := fn.FileLine(_pc)

		//
		var (
			_fnNameDir = strings.Split(_fnName, "/")
			_fnNameLis = strings.Split(_fnName, ".")
			_fnNameSrc string
		)
		if len(_fnNameDir) > 1 {
			_fnNameSrc = _fnNameDir[0] + "/" + _fnNameDir[1] + "/"
		} else {
			_fnNameSrc = _fnNameDir[0]
		}
		fnName := _fnNameLis[len(_fnNameLis)-1]

		// file
		_pcLis := strings.Split(file, _fnNameSrc)
		filePath := strings.Join(_pcLis, "")
		return fmt.Sprintf("%s:%d|%s", filePath, line, fnName)
	}
	return "unknown"
}
