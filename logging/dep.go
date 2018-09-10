package logging

import (
	_ "github.com/apex/log"
	_ "github.com/francoispqt/onelog"
	_ "github.com/inconshreveable/log15"
	_ "github.com/rs/zerolog"
	_ "github.com/rs/zerolog/diode"
	_ "github.com/rs/zerolog/hlog"
	_ "github.com/rs/zerolog/journald"
	_ "github.com/rs/zerolog/log"
	_ "github.com/sirupsen/logrus"
	_ "github.com/sirupsen/logrus/hooks/syslog"
	_ "github.com/sirupsen/logrus/hooks/test"
	_ "go.uber.org/zap"
	_ "go.uber.org/zap/zapgrpc"
)
