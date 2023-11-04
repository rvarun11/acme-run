package amqp

import config "github.com/CAS735-F23/macrun-teamvsl/challenge_manager/config"

var conf *rabbitmqConfig

const (
	exchangeKind       = "direct"
	exchangeDurable    = true
	exchangeAutoDelete = false
	exchangeInternal   = false
	exchangeNoWait     = false

	queueDurable    = true
	queueAutoDelete = false
	queueExclusive  = false
	queueNoWait     = false

	// publishMandatory = false
	// publishImmediate = false

	prefetchCount  = 1
	prefetchSize   = 0
	prefetchGlobal = false

	consumeAutoAck   = false
	consumeExclusive = false
	consumeNoLocal   = false
	consumeNoWait    = false
)

type rabbitmqConfig struct {
	host     string
	port     string
	user     string
	password string
}

func init() {
	conf = newConfig()
}

func newConfig() *rabbitmqConfig {
	return &rabbitmqConfig{
		host:     config.GetEnv("RABBITMQ_HOSTNAME", "localhost"),
		port:     config.GetEnv("RABBITMQ_PORT", "5432"),
		user:     config.GetEnv("RABBITMQ_USER", "guest"),
		password: config.GetEnv("RABBITMQ_PASSWORD", "guest"),
	}
}
