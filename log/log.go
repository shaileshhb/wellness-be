package log

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

// init will create instance of logger
func InitializeLogger() zerolog.Logger {
	return logFormatter()
}

func logFormatter() zerolog.Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}

	log := zerolog.New(output).Level(zerolog.InfoLevel).With().Timestamp().Caller().Logger()
	return log
}
