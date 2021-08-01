package main

import (
	"database/sql"
	"flag"
	"fmt"
	apiClient "github.com/bassbeaver/eventhouse/projector/service/api_client"
	"github.com/bassbeaver/eventhouse/projector/service/event_reader_reducer"
	"github.com/bassbeaver/eventhouse/projector/service/mysql"
	projectorOpentracing "github.com/bassbeaver/eventhouse/projector/service/opentracing"
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

	client := container.GetByAlias(apiClient.ApiClientServiceAlias).(*apiClient.ApiClient)
	opentracingBridge := container.GetByAlias(projectorOpentracing.OpentracingBridgeServiceAlias).(*projectorOpentracing.Bridge)
	readerReducer := container.GetByAlias(event_reader_reducer.EventReaderReducerServiceAlias).(*event_reader_reducer.EventReaderReducer)

	// Start event stream processing
	readerReducer.SubscribeReduceGlobalStream()

	opentracingBridge.CloseTracer()
	if err := client.CloseConnection(); nil != err {
		fmt.Printf("Failed to close gRPC connection: %s \n", err.Error())
	}
	if err := container.GetByAlias(mysql.MysqlDbServiceAlias).(*sql.DB).Close(); nil != err {
		fmt.Printf("Failed to close MySQL connection: %s \n", err.Error())
	}

	fmt.Println("Projector stopped")
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
		projectorOpentracing.OpentracingBridgeServiceAlias,
		gioc.Factory{
			Create: projectorOpentracing.NewBridge,
			Arguments: []string{
				"#jaegger.host",
				"#jaegger.port",
				"#jaegger.service_name",
			},
		},
		true,
	)

	container.RegisterServiceFactoryByAlias(
		apiClient.ApiClientServiceAlias,
		gioc.Factory{
			Create: apiClient.NewApiClient,
			Arguments: []string{
				"#eventhouse.host",
				"#eventhouse.port",
				"#eventhouse.auth_token",
				"@" + projectorOpentracing.OpentracingBridgeServiceAlias,
			},
		},
		true,
	)

	container.RegisterServiceFactoryByAlias(
		event_reader_reducer.EventReducerServiceAlias,
		gioc.Factory{
			Create: event_reader_reducer.NewEventReducer,
			Arguments: []string{
				"@" + mysql.MysqlDbServiceAlias,
			},
		},
		true,
	)

	container.RegisterServiceFactoryByAlias(
		event_reader_reducer.EventReaderReducerServiceAlias,
		gioc.Factory{
			Create: event_reader_reducer.NewEventReader,
			Arguments: []string{
				"@" + apiClient.ApiClientServiceAlias,
				"@" + projectorOpentracing.OpentracingBridgeServiceAlias,
				"@" + event_reader_reducer.EventReducerServiceAlias,
			},
		},
		true,
	)

	container.RegisterServiceFactoryByAlias(
		mysql.MysqlDbServiceAlias,
		gioc.Factory{
			Create: mysql.NewDB,
			Arguments: []string{
				"#mysql.dsn",
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
