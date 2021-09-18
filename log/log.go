package log

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func SetLogger(fileRootPath string, logLife time.Duration, logSpliteTime time.Duration, stdOut bool, isDebug bool) * logrus.Logger{
	logger := logrus.New()
	//显示log来源：哪个文件，多少行，哪个函数
	logger.SetReportCaller(true)
	//是否启动debug
	if isDebug{
		logger.SetLevel(logrus.DebugLevel)
	}else{
		logger.SetLevel(logrus.InfoLevel)
	}
	rotatelog, err := rotatelogs.New(
		fileRootPath + ".%Y%m%d%H%M.log",
		rotatelogs.WithLinkName(fileRootPath),
		// 设置保存时间     其他类别 : WithRotationCount 设置文件最多保存个数
		rotatelogs.WithMaxAge(logLife),
		// 日志分割时间
		rotatelogs.WithRotationTime(logSpliteTime),
	)
	if err != nil{
		logger.Errorf("config rotetelog failed err : %v", err)
	}
	//时间格式
	logger.Formatter = &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	}
	//选择日志打印位置（log文件或者终端）
	if !stdOut{
		logger.SetOutput(rotatelog)
	}else{
		logger.SetOutput(os.Stdout)
	}
	return logger
}
