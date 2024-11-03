package cmd

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/lab42/ingress-backend/templates"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ingress-backend",
	Short: "Default ingress backend",
	Run: func(cmd *cobra.Command, args []string) {
		serverPort := viper.GetString(keyPort)
		filePath := viper.GetString(keyPath)

		e := echo.New()
		e.HideBanner = true

		// Add middleware
		e.Use(
			echoprometheus.NewMiddleware("not_found"),
			middleware.Logger(),
			middleware.RequestID(),
			middleware.Recover(),
			middleware.Timeout(),
			middleware.GzipWithConfig(middleware.GzipConfig{}),
		)

		// Health check route
		e.GET("/health", func(c echo.Context) error {
			return c.String(http.StatusOK, "OK")
		})

		// Determine which file source to use
		var handler http.Handler

		if filePath != "" {
			// Check if the provided path exists
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				e.Logger.Printf("Configured file path does not exist: %s. Falling back to embedded files.", filePath)
				handler = http.FileServer(http.FS(templates.FS))
			} else {
				// Use the configured file path
				handler = http.FileServer(http.Dir(filePath))
				e.Logger.Printf("Serving files from configured path: %s", filePath)
			}
		} else {
			// Use embedded files if no path is configured
			handler = http.FileServer(http.FS(templates.FS))
		}

		// Metrics server
		go func() {
			// Create a separate Echo instance for metrics
			eProm := echo.New()
			eProm.HideBanner = true

			// Add middleware
			eProm.Use(
				middleware.Logger(),
				middleware.RequestID(),
				middleware.Recover(),
				middleware.Timeout(),
				middleware.GzipWithConfig(middleware.GzipConfig{}),
			)

			// Custom metrics handler that uses promhttp
			eProm.GET("/metrics", echoprometheus.NewHandler())

			eProm.Logger.Info("Starting metrics server on :8080/metrics")

			// Start metrics server
			if err := eProm.Start(":8080"); err != nil {
				e.Logger.Errorf("Metrics server failed: %v", err)
			}
		}()

		e.GET("/*", echo.WrapHandler(handler))
		e.Logger.Printf("Starting server on port %s", serverPort)
		e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", serverPort)))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ingress-backend.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Local flags for the root command
	rootCmd.Flags().Int(keyPort, defaultPort, "HTTP server port")
	rootCmd.Flags().String(keyPath, defaultPath, "Root file path")

	// Bind flags to viper
	viper.BindPFlag(keyPort, rootCmd.Flags().Lookup(keyPort))
	viper.BindPFlag(keyPath, rootCmd.Flags().Lookup(keyPath))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// 1. Set defaults first
	viper.SetDefault(keyPort, defaultPort)
	viper.SetDefault(keyPath, defaultPath)

	// 2. Read config file
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".ingress-backend")
	}

	// Read config file if it exists
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	// 3. Set up environment variables
	viper.SetEnvPrefix("ingress-backend")
	viper.AutomaticEnv()

	// Configure environment variable mapping
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
}
