package main

import (
	"flag"
	"fmt"
	"github.com/bassbeaver/eventhouse/controller"
	"github.com/bassbeaver/eventhouse/service/auth"
	"github.com/bassbeaver/eventhouse/service/clickhouse"
	"github.com/bassbeaver/eventhouse/service/grpc"
	"github.com/bassbeaver/eventhouse/service/logger"
	"github.com/bassbeaver/eventhouse/service/recovery"
	"github.com/bassbeaver/eventhouse/service/request_id_setter"
	"github.com/bassbeaver/eventhouse/storage"
	"github.com/bassbeaver/gioc"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

func main() {
	config := readConfigFiles(determineConfigPath())

	container := gioc.NewContainer()

	// Setting parameters to container
	if !config.IsSet("config") {
		panic("Invalid config files format")
	}
	container.SetParameters(config.GetStringMapString("config"))

	registerServicesToContainer(container)

	if noCycles, cycledService := container.CheckCycles(); !noCycles {
		panic("Failed to start application, errors in DI container: service " + cycledService + " has circular dependencies")
	}

	container.GetByAlias(grpc.GrpcServerServiceAlias).(*grpc.GrpcServer).Serve()

	fmt.Println("Application terminated")
}

func determineConfigPath() string {
	flags := flag.NewFlagSet("flags", flag.PanicOnError)
	configPathFlag := flags.String("config", "", "Path to application config folder")
	flagsErr := flags.Parse(os.Args[1:])
	if nil != flagsErr {
		panic(flagsErr)
	}

	var configPath string
	if "" == *configPathFlag {
		curBinDir, curBinDirError := os.Getwd()
		if curBinDirError != nil {
			panic("Failed to determine path to workdir:" + curBinDirError.Error())
		}
		configPath = curBinDir + "/config"
	} else {
		configPath = *configPathFlag
	}

	return configPath
}

func readConfigFiles(configPath string) *viper.Viper {
	config := viper.New()

	workdir, workdirError := os.Getwd()
	if nil != workdirError {
		panic("failed to determine working directory, error: " + workdirError.Error())
	}
	config.Set("workdir", workdir)

	var configDir string
	configPathStat, configPathStatError := os.Stat(configPath)
	if nil != configPathStatError {
		panic("failed to read configs: " + configPathStatError.Error())
	}
	if configPathStat.IsDir() {
		configDir = configPath
	} else {
		configDir = filepath.Dir(configPath)
	}

	firstConfigFile := true
	pathWalkError := filepath.Walk(
		configDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				panic("failed to read config file " + path + ", error: " + err.Error())
			}

			if info.IsDir() {
				return nil
			}

			configFilePath := filepath.Dir(path)
			configFileExt := filepath.Ext(info.Name())
			// if extension is not allowed - take next file
			if !stringInSlice(configFileExt[1:], viper.SupportedExts) {
				return nil
			}

			configFileName := info.Name()[0 : len(info.Name())-len(configFileExt)]

			config.AddConfigPath(configFilePath)
			config.SetConfigName(configFileName)

			if firstConfigFile {
				if configError := config.ReadInConfig(); nil != configError {
					return configError
				}

				firstConfigFile = false
			} else {
				if configError := config.MergeInConfig(); nil != configError {
					return configError
				}
			}

			return nil
		},
	)
	if nil != pathWalkError {
		panic("failed to read configs: " + pathWalkError.Error())
	}

	return config
}

func registerServicesToContainer(container *gioc.Container) {
	container.RegisterServiceFactoryByAlias(
		clickhouse.ClickhouseConnectServiceAlias,
		gioc.Factory{
			Create:    clickhouse.NewClickhouseConnect,
			Arguments: []string{"#clickhouse.dsn"},
		},
		true,
	)

	container.RegisterServiceFactoryByAlias(
		storage.EventRepositoryAlias,
		gioc.Factory{
			Create: storage.NewBatchSavingClickhouseRepository,
			Arguments: []string{
				"#event_repository.max_entities_in_batch",
				"#event_repository.batch_lifetime_ms",
				"@" + clickhouse.ClickhouseConnectServiceAlias,
			},
		},
		true,
	)

	container.RegisterServiceFactoryByAlias(
		controller.EventControllerServiceAlias,
		gioc.Factory{
			Create:    controller.NewEventController,
			Arguments: []string{"@" + storage.EventRepositoryAlias},
		},
		true,
	)

	container.RegisterServiceFactoryByAlias(
		logger.LoggerFactoryServiceAlias,
		gioc.Factory{
			Create:    logger.NewLoggerFactory,
			Arguments: []string{"#logs.path"},
		},
		true,
	)

	container.RegisterServiceFactoryByAlias(
		request_id_setter.RequestIdSetterServiceAlias,
		gioc.Factory{
			Create:    request_id_setter.NewRequestIdSetter,
			Arguments: []string{},
		},
		true,
	)

	container.RegisterServiceFactoryByAlias(
		recovery.RecoveryServiceAlias,
		gioc.Factory{
			Create:    recovery.NewRecoveryService,
			Arguments: []string{"@" + logger.LoggerFactoryServiceAlias},
		},
		true,
	)

	container.RegisterServiceFactoryByAlias(
		logger.RequestContextLoggerSetterServiceAlias,
		gioc.Factory{
			Create:    logger.NewRequestContextLoggerInterceptor,
			Arguments: []string{"@" + logger.LoggerFactoryServiceAlias},
		},
		true,
	)

	container.RegisterServiceFactoryByAlias(
		auth.AuthServiceAlias,
		gioc.Factory{
			Create:    auth.NewAuthService,
			Arguments: []string{"@" + clickhouse.ClickhouseConnectServiceAlias},
		},
		true,
	)

	container.RegisterServiceFactoryByAlias(
		grpc.GrpcServerServiceAlias,
		gioc.Factory{
			Create: grpc.NewGrpcServer,
			Arguments: []string{
				"#port",
				"@" + recovery.RecoveryServiceAlias,
				"@" + request_id_setter.RequestIdSetterServiceAlias,
				"@" + logger.RequestContextLoggerSetterServiceAlias,
				"@" + auth.AuthServiceAlias,
				"@" + controller.EventControllerServiceAlias,
			},
		},
		true,
	)
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}

	return false
}
