package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fpay/gopress"
	"github.com/fpay/lehuipay-hodor-go/middlewares/logging"
	"github.com/fpay/lehuipay-shorturl-go/internal/controllers"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var webCmd = &cobra.Command{
	Use:   "web",
	Short: "A shorturl Web Server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Web server called")
		opts := loadApplicationOptions()
		boot := NewBootstrap(opts)

		s := gopress.NewServer(opts.Server.Web)
		app := s.App()
		app.Logger.Logger = boot.Logger.Logger
		app.Use(logging.NewContextLoggerMiddleware(boot.Logger.NewEntry()))
		app.Use(logging.NewLoggingMiddleware(logging.LoggingMiddlewareConfig{
			Name: "shorturl.server",
			Skipper: func(ctx gopress.Context) bool {
				req := ctx.Request()
				return req.RequestURI == "/" && req.Method == "HEAD"
			},
		}))
		webController := controllers.NewWeb(boot.ShortURLService)
		app.GET("/:hash", webController.Index)

		go s.Start()

		// Graceful shutdown
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		<-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.Shutdown(ctx); err != nil {
			app.Logger.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(webCmd)
}
