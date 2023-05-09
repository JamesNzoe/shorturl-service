// Copyright © 2020 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/fpay/lehuipay-shorturl-go/api"
	"github.com/fpay/lehuipay-shorturl-go/internal/controllers"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
)

var (
	serverDefaultStackSize = 4 << 10
)

type GRPCOptions struct {
	Port int `yaml:"port" mapstructure:"port"`
}

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "grpc",
	Short: "A shorturl gRPC Server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gRPC server called")
		opts := loadApplicationOptions()
		boot := NewBootstrap(opts)
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", opts.Server.Grpc.Port))
		handleInitError("grpc", err)

		logger := boot.Logger.WithScope("grpc").Entry
		logger.Debugf("Listen on port: %d", opts.Server.Grpc.Port)
		gs := grpc.NewServer(
			grpc.KeepaliveParams(keepalive.ServerParameters{
				Time: 10 * time.Minute,
			}),
			grpc_middleware.WithUnaryServerChain(
				grpc_logrus.UnaryServerInterceptor(logger),
				grpc_logrus.PayloadUnaryServerInterceptor(
					logger, func(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
						return true
					},
				),
				grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(func(p interface{}) error {
					err, ok := p.(error)
					if !ok {
						err = fmt.Errorf("%v", p)
					}
					stack := make([]byte, serverDefaultStackSize)
					length := runtime.Stack(stack, false)
					logger.WithError(err).Errorf("recovered from panic: %v [stack]: %s", err, stack[:length])
					return status.Errorf(codes.Internal, "Unexcepted internal server error")
				})),
			),
		)

		gServer := controllers.NewGrpcServer(boot.ShortURLService)
		api.RegisterShortURLServiceServer(gs, gServer)
		go gs.Serve(lis)

		// Graceful shutdown
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		<-quit

		gs.GracefulStop()
	},
}

func init() {
	RootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
