package plugins

/*
import (
    "flag"
    "log"
    "os"
    "github.com/devopsfaith/krakend/config/viper"
    "github.com/devopsfaith/krakend/logging/gologging"
    "github.com/devopsfaith/krakend/proxy"
    "github.com/devopsfaith/krakend/router/gin"
)

func Krakend() {
    port := flag.Int("p", 0, "Port of the service")
    logLevel := flag.String("l", "ERROR", "Logging level")
    debug := flag.Bool("d", false, "Enable the debug")
    configFile := flag.String("c", "/etc/krakend/configuration.json", "Path to the configuration filename")
    flag.Parse()

    parser := viper.New()
    serviceConfig, err := parser.Parse(*configFile)
    if err != nil {
        log.Fatal("ERROR:", err.Error())
    }
    serviceConfig.Debug = serviceConfig.Debug || *debug
    if *port != 0 {
        serviceConfig.Port = *port
    }

    logger := gologging.NewLogger(*logLevel, os.Stdout, "[KRAKEND]")

    routerFactory := gin.DefaultFactory(proxy.DefaultFactory(logger), logger)

    routerFactory.New().Run(serviceConfig)
}
*/