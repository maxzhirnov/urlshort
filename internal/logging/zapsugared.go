package logging

import "go.uber.org/zap"

type ZapSugaredAdapter struct {
	*zap.SugaredLogger
}

func NewZapSugared() (*ZapSugaredAdapter, error) {
	logger, err := zap.NewProduction(zap.AddCallerSkip(1))

	if err != nil {
		return nil, err
	}
	sugared := logger.Sugar()
	return &ZapSugaredAdapter{
		sugared,
	}, nil
}

func (z *ZapSugaredAdapter) Debug(msg string, keysAndValues ...interface{}) {
	z.SugaredLogger.Debugln(msg, keysAndValues)
}

func (z *ZapSugaredAdapter) Info(msg string, keysAndValues ...interface{}) {
	z.SugaredLogger.Infoln(msg, keysAndValues)
}

func (z *ZapSugaredAdapter) Error(msg string, keysAndValues ...interface{}) {
	z.SugaredLogger.Errorln(msg, keysAndValues)
}

func (z *ZapSugaredAdapter) Fatal(msg string, keysAndValues ...interface{}) {
	z.SugaredLogger.Fatalln(msg, keysAndValues)
}

func (z *ZapSugaredAdapter) Warn(msg string, keysAndValues ...interface{}) {
	z.SugaredLogger.Warnln(msg, keysAndValues)
}
