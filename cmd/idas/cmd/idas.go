/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	lightstep "github.com/lightstep/lightstep-tracer-go"
	"github.com/oklog/oklog/pkg/group"
	stdopentracing "github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"
	"sourcegraph.com/sourcegraph/appdash"
	appdashot "sourcegraph.com/sourcegraph/appdash/opentracing"

	"idas/config"
	"idas/pkg/endpoint"
	"idas/pkg/logs"
	"idas/pkg/logs/flag"
	"idas/pkg/service"
	"idas/pkg/transport"
	"idas/pkg/utils/signals"
)

var (
	cfgFile        string
	logConfig      logs.Config
	debugAddr      string
	httpAddr       string
	zipkinURL      string
	zipkinBridge   bool
	lightstepToken string
	appdashAddr    string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gateway",
	Short: "The idas gateway server.",
	Long:  `The idas gateway server.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return Run(context.Background(), logs.GetRootLogger(), signals.SetupSignalHandler(logs.GetRootLogger()))
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

func Run(ctx context.Context, logger log.Logger, stopCh *signals.StopChan) (err error) {
	var zipkinTracer *zipkin.Tracer
	{
		if zipkinURL != "" {
			var (
				err         error
				hostPort    = "localhost:80"
				serviceName = "addsvc"
				reporter    = zipkinhttp.NewReporter(zipkinURL)
			)
			defer reporter.Close()
			zEP, _ := zipkin.NewEndpoint(serviceName, hostPort)
			zipkinTracer, err = zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(zEP))
			if err != nil {
				logger.Log("err", err)
				os.Exit(1)
			}
			if !(zipkinBridge) {
				logger.Log("tracer", "Zipkin", "type", "Native", "URL", zipkinURL)
			}
		}
	}
	// Determine which OpenTracing tracer to use. We'll pass the tracer to all the
	// components that use it, as a dependency.
	var tracer stdopentracing.Tracer
	{
		if zipkinBridge && zipkinTracer != nil {
			logger.Log("tracer", "Zipkin", "type", "OpenTracing", "URL", zipkinURL)
			tracer = zipkinot.Wrap(zipkinTracer)
			zipkinTracer = nil // do not instrument with both native tracer and opentracing bridge
		} else if lightstepToken != "" {
			logger.Log("tracer", "LightStep") // probably don't want to print out the token :)
			tracer = lightstep.NewTracer(lightstep.Options{
				AccessToken: lightstepToken,
			})
			defer lightstep.Flush(ctx, tracer)
		} else if appdashAddr != "" {
			logger.Log("tracer", "Appdash", "addr", appdashAddr)
			tracer = appdashot.NewTracer(appdash.NewRemoteCollector(appdashAddr))
		} else {
			tracer = stdopentracing.GlobalTracer() // no-op
		}
	}

	// Create the (sparse) metrics we'll use in the service. They, too, are
	// dependencies that we pass to components that use them.
	var (
		duration metrics.Histogram
	)
	{
		// Endpoint-level metrics.
		duration = prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "example",
			Subsystem: "addsvc",
			Name:      "request_duration_seconds",
			Help:      "Request duration in seconds.",
		}, []string{"method", "success"})
	}
	var (
		svc         = service.New(ctx)
		endpoints   = endpoint.New(svc, logger, duration, tracer, zipkinTracer)
		httpHandler = transport.NewHTTPHandler(endpoints, tracer, zipkinTracer, logger)
	)
	var g group.Group
	{
		// The debug listener mounts the http.DefaultServeMux, and serves up
		// stuff like the Prometheus metrics route, the Go debug and profiling
		// routes, and so on.
		debugListener, err := net.Listen("tcp", debugAddr)
		if err != nil {
			level.Error(logger).Log("transport", "debug/HTTP", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			level.Info(logger).Log("transport", "debug/HTTP", "addr", debugAddr)
			return http.Serve(debugListener, http.DefaultServeMux)
		}, func(error) {
			debugListener.Close()
		})
	}
	{
		// The HTTP listener mounts the Go kit HTTP handler we created.
		httpListener, err := net.Listen("tcp", httpAddr)
		if err != nil {
			level.Error(logger).Log("transport", "HTTP", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			level.Info(logger).Log("transport", "HTTP", "addr", httpAddr)
			return http.Serve(httpListener, httpHandler)
		}, func(error) {
			httpListener.Close()
		})
	}
	return g.Run()
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initLogger, initConfig)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "./idas.yaml", "config file (default is ./idas.yaml)")

	// log level and format
	flag.AddFlags(rootCmd.PersistentFlags(), &logConfig)

	rootCmd.Flags().StringVar(&debugAddr, "debug.addr", ":8080", "Debug and metrics listen address")
	rootCmd.Flags().StringVar(&httpAddr, "http-addr", ":8081", "HTTP listen address")
	rootCmd.Flags().StringVar(&zipkinURL, "zipkin-url", "", "Enable Zipkin tracing via HTTP reporter URL e.g. http://localhost:9411/api/v2/spans")
	rootCmd.Flags().BoolVar(&zipkinBridge, "zipkin-ot-bridge", false, "Use Zipkin OpenTracing bridge instead of native implementation")
	rootCmd.Flags().StringVar(&lightstepToken, "lightstep-token", "", "Enable LightStep tracing via a LightStep access token")
	rootCmd.Flags().StringVar(&appdashAddr, "appdash-addr", "", "Enable Appdash tracing via an Appdash server host:port")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile == "" {
		cfgFile = "./idas.yaml"
	}
	if err := config.ReloadConfigFromFile(logs.GetRootLogger(), cfgFile); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to load config: %s\n", err)
		os.Exit(1)
	}
}

// initLogger
func initLogger() {
	logger := logs.New(&logConfig)
	logs.SetRootLogger(logger)
}
