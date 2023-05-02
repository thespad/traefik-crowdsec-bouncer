package main

import (
  "os"
  "strings"

  "github.com/gin-contrib/logger"
  "github.com/gin-gonic/gin"
  "github.com/rs/zerolog"
  "github.com/rs/zerolog/log"
  "github.com/thespad/traefik-crowdsec-bouncer/config"
  "github.com/thespad/traefik-crowdsec-bouncer/controller"
)

const (
  PingPath        = "/api/v1/ping"
  HealthzPath     = "/api/v1/healthz"
  ForwardAuthPath = "/api/v1/forwardAuth"
  MetricsPath     = "/api/v1/metrics"
)

func main() {
  router, err := setupRouter()
  if err != nil {
    log.Fatal().Err(err).Msgf("An error occurred while starting webserver")
    return
  }

  err = router.Run()
  if err != nil {
    log.Fatal().Err(err).Msgf("An error occurred while starting bouncer")
    return
  }

}

func setupRouter() (*gin.Engine, error) {
  // logger framework
  logOutput := zerolog.ConsoleWriter{Out: os.Stderr, NoColor: !gin.IsDebugging(), TimeFormat: zerolog.TimeFieldFormat}
  log.Logger = log.Output(logOutput)

  logLevel := config.OptionalEnv("CROWDSEC_BOUNCER_LOG_LEVEL", "1")
  level, err := zerolog.ParseLevel(logLevel)
  if err != nil {
    return nil, err
  }
  zerolog.SetGlobalLevel(level)

  // Web framework
  router := gin.New()
  if err := router.SetTrustedProxies(strings.Split(config.OptionalEnv("TRUSTED_PROXIES", "0.0.0.0/0"), ",")); err != nil {
    return nil, err
  }

  router.Use(logger.SetLogger(logger.WithSkipPath([]string{PingPath, HealthzPath})))
  router.GET(PingPath, controller.Ping)
  router.GET(HealthzPath, controller.Healthz)
  router.GET(ForwardAuthPath, controller.ForwardAuth)
  router.GET(MetricsPath, controller.Metrics)

  return router, nil
}
