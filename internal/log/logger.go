package log

import (
	"fmt"
	"goshield/internal/gbl"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

func SetupLogger() {
	logFilePath := getBaseFilePath() + getFileName()
	log.Printf("Log file path: %+v\n", logFilePath)

	writerSyncer := getZapLogWriter(logFilePath)
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writerSyncer, zapcore.DebugLevel)
	zapLogger := zap.New(core, zap.AddCaller())
	Logger = zapLogger.Sugar()
}

func GetAccessLogWriter() io.Writer {
	logPath := getBaseFilePath() + getAccessLogFileName()
	return getLogWriter(logPath)
}

func getZapLogWriter(logPath string) zapcore.WriteSyncer {
	return zapcore.AddSync(getLogWriter(logPath))
}

func getLogWriter(logPath string) io.Writer {
	logFile := logPath
	rotator, err := rotatelogs.New(
		logFile,
		rotatelogs.WithRotationTime(time.Hour))
	if err != nil {
		fmt.Println("[logger.go][getLogWriter]:Error occurred initializing log files", err)
	}

	var multiWriter io.Writer
	if runtime.GOOS == "windows" {
		multiWriter = io.MultiWriter(os.Stderr, rotator)
	} else {
		multiWriter = io.MultiWriter(rotator)
	}
	return multiWriter
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(gbl.DateTimeDispFmt)
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = PathEncodeCaller
	encoderConfig.ConsoleSeparator = " "
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func PathEncodeCaller(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	fileName := filepath.Base(caller.FullPath())
	enc.AppendString("[" + fileName + "]")
}

func getBaseFilePath() string {
	wd, err := os.Getwd()
	if err != nil {
		log.Printf("Error while reading base path: %+v\n", err)
	}

	return filepath.Join(filepath.Dir(wd), "log") + "\\"
}

func getFileName() string {
	return "GoShieldServer_%Y%m%d_%H.log"
}

func getAccessLogFileName() string {
	return "GoSheildServiceAccess_%Y%m%d_%H.log"
}
