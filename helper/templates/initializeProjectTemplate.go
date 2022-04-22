package templates

const GitIgnoreTemplate = `.idea
.env
coverage.*
.scannerwork/*
sonarqube-report/*
gometalinter-report.out
.DS_STORE`

const EnvTemplate = `
#Server
SERVICE_VERSION=v0.0.0
SERVICE_NAME={{.ProjectName}}
SERVER_PORT=8090
GRACEFULLY_DURATION=12s

INTERFACE_HTTP_ENABLE=true
INTERFACE_MESSAGING_ENABLE=false
INTERFACE_WORKER_ENABLE=false

#Log
LOG_LEVEL=INFO
LOG_FORMATTER=JSON

#New Relic
NEW_RELIC_ENABLE=false
NEW_RELIC_APP_NAME=
NEW_RELIC_LICENSE=

#MySQL
MYSQL_ENABLE=false
MYSQL_ADDRESS=localhost:3306
MYSQL_USERNAME=root
MYSQL_PASSWORD=
MYSQL_DB_NAME=db_name
MYSQL_MAX_IDLE_CONNECTION=15
MYSQL_MAX_OPEN_CONNECTION=15
MYSQL_MAX_LIFETIME_CONNECTION=60s
MYSQL_LOG_MODE=true

#Redis
REDIS_ENABLE=false
REDIS_HOST=localhost
REDIS_PASSWORD=
REDIS_DB=0
REDIS_POOL_SIZE=100
REDIS_MIN_IDLE_CONNS=10
REDIS_READ_ONLY=true
REDIS_DIAL_TIMEOUT=1s
REDIS_POOL_TIMEOUT=10s
REDIS_READ_TIMEOUT=1s
REDIS_WRITE_TIMEOUT=3s
REDIS_MAX_CONN_AGE=3m

#Memcache
MEMCACHE_ENABLE=false

#Kafka
KAFKA_ENABLE=false
KAFKA_HOST=localhost:9092
KAFKA_CLIENT_ID={{.ProjectName}}
KAFKA_CONSUMER_GROUP=default-{{.ProjectName}}-local
KAFKA_VERSION=2.6.0
KAFKA_CONSUMER_WORKER=100
KAFKA_CONSUMER_RETRY_MAX=10
KAFKA_STRATEGY=BalanceStrategyRange

#Kafka Slack Notification
KAFKA_SLACK_NOTIFICATION_ACTIVE=false
KAFKA_SLACK_NOTIFICATION_HOOK=https://hooks.slack.com/services/
KAFKA_SLACK_NOTIFICATION_TIMEOUT=12s

SENTRY_ENABLE=false
SENTRY_ENVIRONMENT=development
SENTRY_DSN=
SENTRY_APP_NAME={{.ProjectName}}
SENTRY_SAMPLE_RATE=1
SENTRY_TRACE_SAMPLE_RATE=1
`

const MainTemplate = `package main

import (
	"context"
	"{{.ProjectURL}}/{{.InboundPath}}"
	pkgDi "{{.ProjectURL}}/{{.DiPath}}"
	pkgResource "{{.ProjectURL}}/{{.ResourcePath}}"
	leafRunner "github.com/paulusrobin/leaf-utilities/appRunner"
	leafServer "github.com/paulusrobin/leaf-utilities/appRunner/server"
	leafLogger "github.com/paulusrobin/leaf-utilities/logger/logger"
	"github.com/paulusrobin/leaf-utilities/tracer/tracer/tracer"
	"go.uber.org/dig"
)

type (
	ContainerCall func() (*dig.Container, error)
	Invoke        func(container *dig.Container) error
	InvokeError   func(container *dig.Container, err error)
)

func Run(containerCall ContainerCall, invoke Invoke, onError InvokeError) error {
	container, err := containerCall()
	if err != nil {
		return err
	}
	if err := invoke(container); err != nil {
		onError(container, err)
	}
	return nil
}

func run(container *dig.Container) error {
	return container.Invoke(func(inbound inbound.Inbound, resource pkgResource.Resource) error {
		tracer.Start(resource.Tracer)

		leafRunner.New().
			With(leafServer.NewHttp(resource.ConfigApp.ServiceName, resource.ConfigApp.ServiceVersion,
				leafServer.WithHttpPort(resource.ConfigApp.ServerPort),
				leafServer.WithHttpHealthAccessKey(resource.ConfigApp.HealthCheckAccessKey),
				leafServer.WithHttpLogger(resource.Log),
				leafServer.WithHttpHealthCheck(inbound.Http.Health.Check),
				leafServer.WithHttpRegister(inbound.Http.Routes),
				leafServer.WithHttpValidator(resource.Validator),
				leafServer.WithHttpEnable(resource.ConfigApp.HttpEnable),
			)).
			Run()

		tracer.Stop()
		return nil
	})
}

func onError(container *dig.Container, err error) {
	_ = container.Invoke(func(logger leafLogger.Logger) error {
		if err != nil {
			logger.Info(leafLogger.BuildMessage(context.Background(), err.Error()))
		}
		return nil
	})
}

func main() {
	if err := Run(pkgDi.Container, run, onError); err != nil {
		panic(err)
	}
}
`

const GenerateMockTemplate = `# =================================================
# Example
# =================================================
# mockgen
# --source=internal/usecases/apiKey/usecase.go
# -package=mock_usecase_apiKey
# --destination=mocks/usecases/apiKey/mock_usecase.go

# - UseCase
mkdir -p mocks/usecases

# - Outbound
# - Outbound.Repositories
mkdir -p mocks/outbound/repositories

# - Outbound.Cache
mkdir -p mocks/outbound/cache

# - Outbound.Messaging
mkdir -p mocks/outbound/messaging

# - Outbound.Webservices
mkdir -p mocks/outbound/webservices

# ==========================================
# ================= RESOURCES ==============
# ==========================================
mkdir -p mocks/resource
`

const GoModTemplate = `module {{.ProjectURL}}

go 1.18
`

const InboundDITemplate = `package inbound

import (
	"{{.ProjectURL}}/{{.InboundPath}}/http"
	"go.uber.org/dig"
)

type Inbound struct {
	dig.In

	Http http.Inbound
}`

const OutboundDITemplate = `package outbound

import (
	"go.uber.org/dig"
)

type Outbound struct {
	dig.In
}

func Register(container *dig.Container) error {
	return nil
}`

const UseCasesDITemplate = `package usecases

import (
	"go.uber.org/dig"
)

type UseCase struct {
	dig.In
}

func Register(container *dig.Container) error {
	return nil
}`

const HttpRoutesTemplate = `package http

import (
	"{{.ProjectURL}}/{{.InboundPath}}/http/health"
	pkgResource "{{.ProjectURL}}/{{.ResourcePath}}"
	"github.com/labstack/echo/v4"
	leafHttpMiddleware "github.com/paulusrobin/leaf-utilities/appRunner/middleware/http"
	leafPrivilege "github.com/paulusrobin/leaf-utilities/common/constants/privilege"
	leafLogger "github.com/paulusrobin/leaf-utilities/logger/logger"
	"go.uber.org/dig"
)

type Inbound struct {
	dig.In

	Health   health.Controller
	Resource pkgResource.Resource
}

func (i Inbound) Routes(ec *echo.Echo, logger leafLogger.Logger) {
	ec.Use(leafHttpMiddleware.AppContextWithLogger(logger), leafHttpMiddleware.Tracer())
	ec.GET("/healthz/routes", i.Health.Routes).Name = leafPrivilege.Trusted
}`

const HealthRoutesTemplate = `package health

import (
	"github.com/labstack/echo/v4"
	leafHttpResponse "github.com/paulusrobin/leaf-utilities/appRunner/response/http"
	leafModel "github.com/paulusrobin/leaf-utilities/common/model"
	"net/http"
	"sort"
)

func getListEndpoints(e *echo.Echo) leafModel.HealthEndpoints {
	var (
		routes    []*echo.Route             = e.Routes()
		routeLen  int                       = len(routes)
		endpoints leafModel.HealthEndpoints = make(leafModel.HealthEndpoints, routeLen)
	)

	for i := 0; i < routeLen; i++ {
		endpoints[i] = leafModel.HealthEndpoint{
			Method: routes[i].Method,
			Path:   routes[i].Path,
			Verify: routes[i].Name,
		}
	}

	sort.Slice(endpoints, func(i, j int) bool {
		if endpoints[i].Path < endpoints[j].Path {
			return true
		}

		if endpoints[i].Path > endpoints[j].Path {
			return false
		}

		return endpoints[i].Method < endpoints[j].Method
	})

	return endpoints
}

func (c *Controller) Routes(eCtx echo.Context) error {
	endpoints := getListEndpoints(eCtx.Echo())
	response := leafModel.HealthRoutesResponse{Routes: endpoints.String()}
	return leafHttpResponse.NewJSONResponse(eCtx, http.StatusOK, response)
}`

const HealthControllerTemplate = `package health

import (
	pkgResource "{{.ProjectURL}}/{{.ResourcePath}}"
	"go.uber.org/dig"
)

type Controller struct {
	dig.In

	Resource pkgResource.Resource
}`

const HealthCheckTemplate = `package health

import (
	"context"
	"net/http"
)

const (
	toggle = "toggle"
	status = "status"
)

func (c *Controller) Check(ctx context.Context) (int, map[string]map[string]interface{}) {
	var (
		httpStatus   = http.StatusOK
		dependencies = make(map[string]map[string]interface{})
	)
	dependencies = map[string]map[string]interface{}{
		// pkgResource.Memcache:   {toggle: c.Resource.ConfigMemcache.MemcacheEnable},
		// pkgResource.Redis:      {toggle: c.Resource.ConfigRedis.RedisEnable},
		// pkgResource.Kafka:      {toggle: c.Resource.ConfigKafka.KafkaEnable},
		// pkgResource.MySQL:      {toggle: c.Resource.ConfigMySQL.MySqlEnable},
		// pkgResource.PostgreSQL: {toggle: c.Resource.ConfigMySQL.PostgreSQLEnable},
		// pkgResource.MongoDB:    {toggle: c.Resource.ConfigMySQL.MongoDBEnable},
	}
	// httpStatus, dependencies[pkgResource.Memcache] 	= ping(ctx, httpStatus, dependencies[pkgResource.Memcache], c.Resource.Cache.Memcache.Ping)
	// httpStatus, dependencies[pkgResource.Redis] 		= ping(ctx, httpStatus, dependencies[pkgResource.Redis], c.Resource.Cache.Redis.Ping)
	// httpStatus, dependencies[pkgResource.Kafka] 		= ping(ctx, httpStatus, dependencies[pkgResource.Kafka], c.Resource.MQ.Kafka.Ping)
	// httpStatus, dependencies[pkgResource.MySQL] 		= ping(ctx, httpStatus, dependencies[pkgResource.MySQL], c.Resource.DatabaseSQL.MySQL.Ping)
	// httpStatus, dependencies[pkgResource.PostgreSQL] = ping(ctx, httpStatus, dependencies[pkgResource.PostgreSQL], c.Resource.DatabaseSQL.PostgreSQL.Ping)
	// httpStatus, dependencies[pkgResource.MongoDB]    = ping(ctx, httpStatus, dependencies[pkgResource.MongoDB], c.Resource.DatabaseNoSQL.MongoDB.Ping)
	return httpStatus, dependencies
}

func ping(ctx context.Context, httpStatus int, data map[string]interface{}, fn func(ctx context.Context) error) (int, map[string]interface{}) {
	if activeValue, found := data[toggle]; found {
		if !activeValue.(bool) {
			data[status] = "DISABLED"
			return httpStatus, data
		}
	}

	if err := fn(ctx); err != nil {
		data[status] = "DOWN"
		return http.StatusInternalServerError, data
	}
	data[status] = "UP"
	return http.StatusOK, data
}`

const ConfigAppTemplate = `package pkgConfig

import (
	leafConfig "github.com/paulusrobin/leaf-utilities/config"
	"time"
)

type (
	AppConfig struct {
		// - Server
		ServiceVersion       string        ` + "`envconfig:\"SERVICE_VERSION\" required:\"true\"`" + `
		ServiceName          string        ` + "`envconfig:\"SERVICE_NAME\" required:\"true\"`" + `
		ServerPort           int           ` + "`envconfig:\"SERVER_PORT\" default:\"8090\" required:\"true\"`" + `
		GracefullyDuration   time.Duration ` + "`envconfig:\"GRACEFULLY_DURATION\" default:\"10s\"`" + `
		HealthCheckAccessKey string        ` + "`envconfig:\"HEALTH_CHECK_ACCESS_KEY\"`" + `

		// - Interface Setting
		HttpEnable      bool ` + "`envconfig:\"INTERFACE_HTTP_ENABLE\" default:\"true\"`" + `
		MessagingEnable bool ` + "`envconfig:\"INTERFACE_MESSAGING_ENABLE\" default:\"true\"`" + `
		WorkerEnable    bool ` + "`envconfig:\"INTERFACE_WORKER_ENABLE\" default:\"true\"`" + `

		// - Log
		LogLevel     string ` + "`envconfig:\"LOG_LEVEL\" default:\"INFO\"`" + `
		LogFilePath  string ` + "`envconfig:\"LOG_FILE_NAME\"`" + `
		LogFormatter string ` + "`envconfig:\"LOG_FORMATTER\" default:\"JSON\"`" + `
	}
)

func NewAppConfig() (AppConfig, error) {
	configuration := AppConfig{}
	if err := leafConfig.NewFromEnv(&configuration); err != nil {
		return AppConfig{}, err
	}
	return configuration, nil
}
`

const ConfigSentryTemplate = `package pkgConfig

import (
	"fmt"
	leafConfig "github.com/paulusrobin/leaf-utilities/config"
)

type (
	SentryConfig struct {
		// - Sentry
		SentryEnable          bool    ` + "`envconfig:\"SENTRY_ENABLE\" default:\"false\"`" + `
		SentryEnvironment     string  ` + "`envconfig:\"SENTRY_ENVIRONMENT\" default:\"development\"`" + `
		SentryDSN             string  ` + "`envconfig:\"SENTRY_DSN\" required:\"true\"`" + `
		SentryAppName         string  ` + "`envconfig:\"SENTRY_APP_NAME\" default:\"{{.ProjectName}}\"`" + `
		SentrySampleRate      float64 ` + "`envconfig:\"SENTRY_SAMPLE_RATE\"`" + `
		SentryTraceSampleRate float64 ` + "`envconfig:\"SENTRY_TRACE_SAMPLE_RATE\"`" + `
	}
)

func NewSentryConfig() (SentryConfig, error) {
	configuration := SentryConfig{}
	if err := leafConfig.NewFromEnv(&configuration); err != nil {
		return SentryConfig{}, err
	}

	if configuration.SentryEnable {
		if configuration.SentryDSN == "" {
			return SentryConfig{}, fmt.Errorf("setry dsn is required")
		}
	}

	return configuration, nil
}`

const ConfigNewRelicTemplate = `package pkgConfig

import (
	"fmt"
	leafConfig "github.com/paulusrobin/leaf-utilities/config"
)

type (
	NewRelicConfig struct {
		// - NewRelic
		NewRelicEnable  bool   ` + "`envconfig:\"NEW_RELIC_ENABLE\" default:\"false\"`" + `
		NewRelicAppName string ` + "`envconfig:\"NEW_RELIC_APP_NAME\"`" + `
		NewRelicLicense string ` + "`envconfig:\"NEW_RELIC_LICENSE\"`" + `
	}
)

func InitNewRelicConfig() (NewRelicConfig, error) {
	configuration := NewRelicConfig{}
	if err := leafConfig.NewFromEnv(&configuration); err != nil {
		return NewRelicConfig{}, err
	}

	if configuration.NewRelicEnable {
		if configuration.NewRelicAppName == "" || configuration.NewRelicLicense == "" {
			return NewRelicConfig{}, fmt.Errorf("newrelic app name / license is required")
		}
	}

	return configuration, nil
}`

const DITemplate = `package pkgDi

import (
	"{{.ProjectURL}}/internal/outbound"
	"{{.ProjectURL}}/internal/usecases"
	pkgConfig "{{.ProjectURL}}/{{.ConfigPath}}"
	pkgResource "{{.ProjectURL}}/{{.ResourcePath}}"
	"go.uber.org/dig"
	"sync"
)

var (
	container *dig.Container
	once      sync.Once
)

func Container() (*dig.Container, error) {
	var outer error

	once.Do(func() {
		container = dig.New()

		if err := pkgConfig.Register(container); err != nil {
			outer = err
			return
		}

		if err := pkgResource.Register(container); err != nil {
			outer = err
			return
		}

		if err := outbound.Register(container); err != nil {
			outer = err
			return
		}
		
		if err := usecases.Register(container); err != nil {
			outer = err
			return
		}
	})

	if outer != nil {
		return nil, outer
	}

	return container, nil
}
`

const LoggerTemplate = `package injection

import (
	pkgConfig "{{.ProjectURL}}/{{.ConfigPath}}"
	leafLogrus "github.com/paulusrobin/leaf-utilities/logger/integrations/logrus"
	leafLogger "github.com/paulusrobin/leaf-utilities/logger/logger"
)

func NewLogger(config pkgConfig.AppConfig) (leafLogger.Logger, error) {
	formatter, err := leafLogrus.GetLoggerFormatter(config.LogFormatter)
	if err != nil {
		return nil, err
	}
	return leafLogrus.New(
		leafLogrus.WithLevel(leafLogger.GetLoggerLevel(config.LogLevel)),
		leafLogrus.WithFormatter(formatter),
		leafLogrus.WithLogFilePath(config.LogFilePath),
	)
}
`

const TracerTemplate = `package injection

import (
	pkgConfig "{{.ProjectURL}}/{{.ConfigPath}}"
	"github.com/getsentry/sentry-go"
	leafNewRelicTracer "github.com/paulusrobin/leaf-utilities/tracer/integrations/newRelic"
	leafSentryTracer "github.com/paulusrobin/leaf-utilities/tracer/integrations/sentry"
	leafTracer "github.com/paulusrobin/leaf-utilities/tracer/tracer"
	"github.com/paulusrobin/leaf-utilities/tracer/tracer/tracer"
	"time"
)

func InitTracer(appConfig pkgConfig.AppConfig, nrConfig pkgConfig.NewRelicConfig, sentryConfig pkgConfig.SentryConfig) (leafTracer.Tracer, error) {
	if nrConfig.NewRelicEnable {
		return leafNewRelicTracer.InitTracing(
			leafNewRelicTracer.WithAppName(nrConfig.NewRelicAppName),
			leafNewRelicTracer.WithLicense(nrConfig.NewRelicLicense),
		)
	}

	if sentryConfig.SentryEnable {
		return leafSentryTracer.InitTracing(
			leafSentryTracer.WithSentryOptions(sentry.ClientOptions{
				Dsn:              sentryConfig.SentryDSN,
				ServerName:       sentryConfig.SentryAppName,
				Release:          appConfig.ServiceVersion,
				Environment:      sentryConfig.SentryEnvironment,
				SampleRate:       sentryConfig.SentrySampleRate,
				TracesSampleRate: sentryConfig.SentryTraceSampleRate,
			}),
			leafSentryTracer.WithDeferFlushDuration(12*time.Second),
		)
	}

	sentry.CaptureMessage("Capture Initialization")

	return tracer.NoopTracer(), nil
}
`

const TranslatorTemplate = `package injection

import (
	ut "github.com/go-playground/universal-translator"
	leafTranslatorValidator "github.com/paulusrobin/leaf-utilities/translator/validator"
)

type Translator struct {
	Validator *ut.UniversalTranslator
}

func NewValidatorTranslator() Translator {
	return Translator{
		Validator: leafTranslatorValidator.GetTranslator(),
	}
}`

const ValidatorTemplate = `package injection

import (
	leafValidatorV10 "github.com/paulusrobin/leaf-utilities/validator/integrations/v10"
	leafValidator "github.com/paulusrobin/leaf-utilities/validator/validator"
)

func NewValidator(translator Translator) (leafValidator.Validator, error) {
	return leafValidatorV10.New(leafValidatorV10.WithTranslator(translator.Validator))
}
`

const ResourceTemplate = `package pkgResource

import (
	pkgConfig "{{.ProjectURL}}/{{.ConfigPath}}"
	leafLogger "github.com/paulusrobin/leaf-utilities/logger/logger"
	leafTracer "github.com/paulusrobin/leaf-utilities/tracer/tracer"
	leafValidator "github.com/paulusrobin/leaf-utilities/validator/validator"
	"go.uber.org/dig"
)

//const (
//	Memcache   = "memcache"
//	Redis      = "redis"
//	MySQL      = "mysql"
//	PostgreSQL = "postgresql"
//	Kafka      = "kafka"
//	MongoDB    = "mongodb"
//)

type (
	Resource struct {
		dig.In

		// Config for application
		ConfigApp               pkgConfig.AppConfig

		// Log for logger
		Log leafLogger.Logger

		// Tracer is utilities tracer
		Tracer leafTracer.Tracer

		// Validator
		Validator leafValidator.Validator
	}
)
`

const ResourceDiTemplate = `package pkgResource

import (
	"fmt"
	"{{.ProjectURL}}/{{.ResourcePath}}/injection"
	"go.uber.org/dig"
)

func Register(container *dig.Container) error {
	if err := container.Provide(injection.NewLogger); err != nil {
		return fmt.Errorf("[DI] can not initialize logger: %+v", err)
	}
	if err := container.Provide(injection.InitTracer); err != nil {
		return fmt.Errorf("[DI] can not initialize tracer: %+v", err)
	}
	if err := container.Provide(injection.NewValidatorTranslator); err != nil {
		return fmt.Errorf("[DI] can not initialize validator translator: %+v", err)
	}
	if err := container.Provide(injection.NewValidator); err != nil {
		return fmt.Errorf("[DI] can not initialize validator: %+v", err)
	}
	return nil
}
`

const ConfigDITemplate = `package pkgConfig

import (
	"fmt"

	"go.uber.org/dig"
)

func Register(container *dig.Container) error {
	if err := container.Provide(NewAppConfig); err != nil {
		return fmt.Errorf("[DI] can not initialize application config: %+v", err)
	}
	if err := container.Provide(InitNewRelicConfig); err != nil {
		return fmt.Errorf("[DI] can not initialize new relic config: %+v", err)
	}
	if err := container.Provide(NewSentryConfig); err != nil {
		return fmt.Errorf("[DI] can not initialize sentry config: %+v", err)
	}
	return nil
}`
