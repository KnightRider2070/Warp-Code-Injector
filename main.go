package main

import (
	"fmt"
	"os"
	"time"
	"wci/cmd"
	"wci/config"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
	config.NoColor = os.Getenv("DEBUG") == ""

	// Configure zerolog without the "version" field for logs
	zerolog.SetGlobalLevel(logLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		NoColor:    config.NoColor,
	})

	// Log the selected log level
	log.Debug().Str("logLevel", logLevel.String()).Msg("Logging configured")
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
	log.Info().Msg("Application starting...")

	logo, err := loadLogo()
	if err != nil {
		log.Error().Msgf("Error loading logo: %v", err)
	} else {
		fmt.Printf("\n%s\n", logo)
	}
	fmt.Printf("---\n%s\nVersion: %s\nVersion Number: %s\n---\n", config.AppName, config.Version, config.VersionNumber)
}

// loadLogo returns the ASCII logo for the application (replace with your logo).
func loadLogo() (string, error) {

	var logo = `
 __    __                        ___            _       
/ / /\ \ \ __ _  _ __  _ __     / __\ ___    __| |  ___ 
\ \/  \/ // _' || '__|| '_ \   / /   / _ \  / _' | / _ \
 \  /\  /| (_| || |   | |_) | / /___| (_) || (_| ||  __/
  \/  \/  \__,_||_|   | .__/  \____/ \___/  \__,_| \___|
                      |_|                               
  _____          _              _                       
  \_   \ _ __   (_)  ___   ___ | |_  ___   _ __         
   / /\/| '_ \  | | / _ \ / __|| __|/ _ \ | '__|        
/\/ /_  | | | | | ||  __/| (__ | |_| (_) || |           
\____/  |_| |_|_/ | \___| \___| \__|\___/ |_|           
              |__/
`

	return logo, nil
}
