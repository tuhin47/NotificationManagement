package cmd

import (
	"NotificationManagement/config"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the notification management server",
	Long:  `Start the notification management server with the specified configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		e := echo.New()

		// Health check route
		e.GET("/health", func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
		})

		port := config.App().Port
		if port == 0 {
			port = 8080
		}
		addr := fmt.Sprintf(":%d", port)

		fmt.Printf("Starting %s server on port %d\n", config.App().Name, port)
		fmt.Println("Server is running in", config.App().Env, "mode")

		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "Failed to start server: %v\n", err)
			os.Exit(1)
		}
	},
}
