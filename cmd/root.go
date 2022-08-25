package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	version string
	commit  string
	verbose bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{Use: "photo-copier", Version: fmt.Sprintf("%s / %s", version, commit)}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	setupCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func setupCmd() {
	cobra.OnInitialize(initConfig)

	copyInit()
	resizeInit()
	watermarkInit()

	cmdCopyPhoto.PersistentFlags().BoolVarP(&verbose, "verbose", "", false, "verbose output")
	cmdCopyPhoto.PersistentFlags().StringVarP(&cfgFile, "config", "", "", "override configuration file")
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	if verbose {
		log.SetLevel(log.DebugLevel)
		err := os.Setenv("LOGLEVEL", "DEBUG")
		if err != nil {
			log.Fatalf("Cannot add LOGLEVEL variable into Environment")
		}
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find config directory.
		cfg, err := os.UserConfigDir()
		if err != nil {
			log.Debugf("Error getting the config dir: %s", err)
		} else {
			viper.AddConfigPath(cfg)
		}
		home, err := os.UserHomeDir()
		if err != nil {
			log.Debugf("Error getting the home dir: %s", err)
		} else {
			viper.AddConfigPath(home)
		}

		viper.AddConfigPath(".")
		viper.SetConfigName("photo-copier")
	}

	setConfigDefaults()
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok { // nolint:errorlint
			log.Debugln("No config file found")
		} else {
			log.Infoln("Error reading config file:", err)
		}
	} else {
		log.Debugln("Using config file:", viper.ConfigFileUsed())
	}
}

func setConfigDefaults() {}
