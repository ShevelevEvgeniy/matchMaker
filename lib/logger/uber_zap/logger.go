package uber_zap

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger() (logger *zap.Logger, stop func(), err error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, nil, err
	}

	switch cfg.EnvType {
	case "dev":
		log, stop := InitDevLogger()
		return log, stop, nil
	case "prod":
		log, stop := InitProdLogger()
		return log, stop, nil
	default:
		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		return logger, func() {
			_ = logger.Sync()
		}, nil
	}
}

func InitDevLogger() (logger *zap.Logger, stop func()) {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stderr), zap.DebugLevel)
	stop = func() {
		_ = core.Sync()
	}
	logger = zap.New(core, zap.AddStacktrace(zap.FatalLevel), zap.AddCaller())

	return
}

func InitProdLogger() (logger *zap.Logger, stop func()) {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zap.InfoLevel)

	stop = func() {
		_ = core.Sync()
	}
	logger = zap.New(core, zap.AddStacktrace(zap.FatalLevel), zap.AddCaller())

	return
}
