package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/lebauce/vlaunch/backend"
	"github.com/lebauce/vlaunch/config"
	"github.com/lebauce/vlaunch/vm"
	"github.com/spf13/cobra"
)

var (
	cfgFiles []string
	keepVM   bool
)

var RootCmd = &cobra.Command{
	Use: "vlaunch",
	Run: func(cmd *cobra.Command, args []string) {
		dataPath := config.GetConfig().GetString("data_path")
		logWriters := []io.Writer{os.Stdout}
		if logFile, err := os.Create(path.Join(dataPath, "vlaunch.log")); err == nil {
			logWriters = append(logWriters, logFile)
		} else if logFile, err := os.Create("/tmp/vlaunch.log"); err == nil {
			logWriters = append(logWriters, logFile)
		}

		multiLogger := io.MultiWriter(logWriters...)
		log.SetOutput(multiLogger)

		if os.Geteuid() != 0 {
			executable, err := os.Executable()
			if err != nil {
				log.Panic(fmt.Sprintf("Failed to determine executable: %s", err.Error()))
			}

			if err := backend.RunAsRoot(executable); err != nil {
				log.Panic(fmt.Sprintf("Failed to run as root: %s", err.Error()))
			}

			return
		}

		vm, err := vm.NewVM()
		if err != nil {
			log.Panic(fmt.Sprintf("Failed to create vm: %s", err.Error()))
		}
		defer func() {
			if !keepVM {
				if err := vm.Release(); err != nil {
					log.Panic(fmt.Sprintf("Failed to release vm: %s", err.Error()))
				}
			}
		}()

		if err := vm.Start(); err != nil {
			log.Panic(fmt.Sprintf("Failed to start vm: %s", err.Error()))
		}

		if err := vm.Run(); err != nil {
			log.Panic(fmt.Sprintf("Error during vm execution: %s", err.Error()))
		}

		if err := vm.Stop(); err != nil {
			log.Panic(fmt.Sprintf("Failed to stop vm: %s", err.Error()))
		}
	},
}

func initConfig() {
	if err := config.InitConfig(cfgFiles); err != nil {
		log.Panic(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringArrayVarP(&cfgFiles, "config", "c", []string{}, "location of Vlaunch configuration files")
	RootCmd.PersistentFlags().BoolVarP(&keepVM, "keep", "k", false, "do not destroy the VM when exiting")
}