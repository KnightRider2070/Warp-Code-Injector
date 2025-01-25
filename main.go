package main

import (
	"fmt"
	"os"
	"time"
	"wci/cmd"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	Version = "dev"
	noColor = false
	appName = "Warp Code Injector"
)

func main() {
	// Configure zerolog duration format
	zerolog.DurationFieldUnit = time.Second

	// Set up global log level
	configureLogging()

	// Display startup information
	displayStartupInfo()

	// Execute the application commands
	if err := cmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("Application encountered an error")
		os.Exit(1)
	}
}

// configureLogging sets the global log level and console writer based on environment variables.
func configureLogging() {
	// Set the global log level based on the LOGLEVEL environment variable
	logLevel := parseLogLevel(os.Getenv("LOGLEVEL"))

	// Enable or disable colored output based on the DEBUG environment variable
	if os.Getenv("DEBUG") != "" {
		noColor = true
	}

	zerolog.SetGlobalLevel(logLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		NoColor:    noColor,
	})
}

// parseLogLevel converts a string log level to a zerolog.Level.
func parseLogLevel(level string) zerolog.Level {
	switch level {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	case "info":
		fallthrough // Default case falls through to info level
	default:
		return zerolog.InfoLevel
	}
}

// displayStartupInfo prints the application banner and version information.
func displayStartupInfo() {
	logo, err := loadLogo()
	if err != nil {
		log.Error().Msgf("Error loading logo: %v", err)
	} else {
		fmt.Printf("\n%s\n", logo)
	}
	fmt.Printf("---\n%s Version: %s\n---\n", appName, Version)
}
