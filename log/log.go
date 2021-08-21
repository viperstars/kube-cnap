package log

import "go.uber.org/zap"

func NewLogger(path string) *zap.Logger {
    loggerConfig := &zap.Config{
        Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
        Development: false,
        Sampling: &zap.SamplingConfig{
            Initial:    100,
            Thereafter: 100,
        },
        Encoding:         "json",
        EncoderConfig:    zap.NewProductionEncoderConfig(),
        OutputPaths:      []string{path},
        ErrorOutputPaths: []string{path},
    }
    logger, err := loggerConfig.Build()
    if err != nil {
        panic(err)
    }
    return logger
}
