/*
 Copyright © 2022 MicroOps-cn.

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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	stdlog "log"
	"net"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	//revive:disable:blank-imports
	_ "net/http/pprof"

	"github.com/MicroOps-cn/fuck/log"
	"github.com/MicroOps-cn/fuck/log/flag"
	"github.com/MicroOps-cn/fuck/signals"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	kitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/go-logr/stdr"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/oklog/oklog/pkg/group"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"gopkg.in/yaml.v3"

	"github.com/MicroOps-cn/fuck/clients/tracing"
	"github.com/MicroOps-cn/idas/config"
	"github.com/MicroOps-cn/idas/pkg/endpoint"
	"github.com/MicroOps-cn/idas/pkg/global"
	"github.com/MicroOps-cn/idas/pkg/service"
	"github.com/MicroOps-cn/idas/pkg/transport"
	"github.com/MicroOps-cn/idas/pkg/utils/httputil"
)

var (
	cfgFile         string
	configDisplay   bool
	debugAddr       string
	httpExternalURL httputil.URL
	webPrefix       string
	httpAddr        string
	proxyHTTPAddr   string
	openapiPath     string
	swaggerPath     string
	radiusAddr      string
	swaggerFilePath string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gateway",
	Short: "The idas gateway server.",
	Long:  `The idas gateway server.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := log.GetContextLogger(cmd.Context())
		ch := signals.SetupSignalHandler(logger)
		ctx, cancelFunc := context.WithCancel(cmd.Context())
		go func() {
			<-ch.Channel()
			cancelFunc()
		}()
		ctx = context.WithValue(ctx, "command", cmd.Use)
		return Run(ctx, logger, ch)
	},
}

func Run(ctx context.Context, logger kitlog.Logger, stopCh *signals.Handler) (err error) {
	var tracer *sdktrace.TracerProvider
	{
		otel.SetLogger(stdr.New(stdlog.New(kitlog.NewStdlibAdapter(level.Info(logger)), "[restful]", stdlog.LstdFlags|stdlog.Lshortfile)))
		tracer, err = tracing.NewTraceProvider(ctx, &config.Get().Trace)
		if err != nil {
			return err
		}
		otel.SetTracerProvider(tracer)
		if http.DefaultClient.Transport == nil {
			http.DefaultClient.Transport = otelhttp.NewTransport(http.DefaultTransport)
		}
		tracing.SetTraceOptions(&config.Get().Trace)
	}

	// Create the (sparse) metrics we'll use in the service. They, too, are
	// dependencies that we pass to components that use them.
	var (
		duration metrics.Histogram
	)
	{
		// Endpoint-level metrics.
		duration = prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Name: "endpoint_invoke_duration_seconds",
			Help: "Tracks the latencies for Invoke endpoints.",
		}, []string{"method", "success"})
	}
	httpLoginURL := httpExternalURL
	httpLoginURL.Path = path.Join(httpLoginURL.Path, webPrefix, "account/login")
	ctx = context.WithValue(ctx, global.HTTPLoginURLKey, httpLoginURL.String())
	ctx = context.WithValue(ctx, global.HTTPExternalURLKey, httpExternalURL.String())
	ctx = context.WithValue(ctx, global.HTTPWebPrefixKey, webPrefix)
	level.Info(logger).Log("msg", "Start service", "externalUrl", httpExternalURL, "webPrefix", webPrefix, "loginUrl", httpLoginURL)
	var (
		svc          = service.New(ctx)
		endpoints    = endpoint.New(ctx, svc, duration)
		httpHandler  = transport.NewHTTPHandler(ctx, logger, endpoints, openapiPath)
		proxyHandler = transport.NewProxyHandler(ctx, logger, endpoints)
		httpServer   = http.NewServeMux()
	)
	if err = svc.LoadSystemConfig(ctx); err != nil {
		panic(fmt.Errorf("failed to load system config: %s", err))
	}
	if len(swaggerPath) > 0 && len(openapiPath) > 0 && len(swaggerFilePath) > 0 {
		stat, err := os.Stat(swaggerFilePath)
		if err != nil {
			level.Error(logger).Log("err", err, "msg", "Failed to get swagger UI directory status, so disable that.")
		} else if stat.IsDir() {
			httpServer.Handle(swaggerPath, http.StripPrefix(swaggerPath, http.FileServer(http.Dir(swaggerFilePath))))
			level.Info(logger).Log("msg", fmt.Sprintf("enable Swagger UI on `%s`", swaggerPath))
		} else {
			level.Error(logger).Log("msg", " swagger UI local path is not directory, so disable that.")
		}
	}

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
		http.DefaultServeMux.Handle("/metrics", promhttp.Handler())
		g.Add(func() error {
			level.Info(logger).Log("msg", "Listening port", "transport", "debug/HTTP", "addr", debugAddr)
			return http.Serve(debugListener, http.DefaultServeMux)
		}, func(error) {
			debugListener.Close()
			level.Debug(logger).Log("msg", "Listen closed", "transport", "debug/HTTP", "addr", debugAddr)
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
			level.Info(logger).Log("msg", "Listening port", "transport", "HTTP", "addr", httpAddr)
			httpServer.Handle("/", httpHandler)
			serv := http.Server{Handler: httpServer, BaseContext: func(listener net.Listener) context.Context {
				return ctx
			}}
			return serv.Serve(httpListener)
		}, func(error) {
			httpListener.Close()
			level.Debug(logger).Log("msg", "Listen closed", "transport", "HTTP", "addr", httpAddr)
		})
	}
	{
		// The HTTP listener mounts the Go kit HTTP handler we created.
		proxyHTTPListener, err := net.Listen("tcp", proxyHTTPAddr)
		if err != nil {
			level.Error(logger).Log("transport", "Proxy/HTTP", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			level.Info(logger).Log("msg", "Listening port", "transport", "Proxy/HTTP", "addr", proxyHTTPAddr)
			serv := http.Server{Handler: proxyHandler, BaseContext: func(listener net.Listener) context.Context {
				return ctx
			}}
			return serv.Serve(proxyHTTPListener)
		}, func(error) {
			proxyHTTPListener.Close()
			level.Debug(logger).Log("msg", "Listen closed", "transport", "Proxy/HTTP", "addr", proxyHTTPAddr)
		})
	}
	{
		if len(radiusAddr) > 0 {
			radiusService := transport.NewRadiusService(ctx, endpoints)
			radiusService.Addr = radiusAddr
			g.Add(func() error {
				level.Info(logger).Log("msg", "Listening port", "transport", "Radius", "addr", radiusService.Addr)
				return radiusService.ListenAndServe()
			}, func(error) {
				_ = radiusService.Shutdown(ctx)
				level.Debug(logger).Log("msg", "Listen closed", "transport", "Radius", "addr", radiusService.Addr)
			})
		}
	}
	{
		stopCh.Add(1)
		g.Add(func() error {
			stopCh.WaitRequest()
			return nil
		}, func(error) {
			if tracer != nil {
				timeoutCtx, closeCh := context.WithTimeout(context.Background(), time.Second*3)
				defer closeCh()
				if err = tracer.ForceFlush(timeoutCtx); err != nil {
					level.Debug(logger).Log("msg", "failed to force flush trace", "err", err)
					return
				}
				if err = tracer.Shutdown(timeoutCtx); err != nil {
					level.Debug(logger).Log("msg", "failed to force close trace", "err", err)
					return
				}
				stopCh.Done()
			}
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
	cobra.OnInitialize(initParameter, initConfig)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "./idas.yaml", "config file")
	rootCmd.PersistentFlags().BoolVar(&configDisplay, "config.display", false, "display config")

	// log level and format
	flag.AddFlags(rootCmd.PersistentFlags(), nil)

	rootCmd.Flags().StringVar(&radiusAddr, "radius.listen-address", "", "Radius listen address")
	rootCmd.Flags().StringVar(&debugAddr, "debug.listen-address", ":8080", "Debug and metrics listen address")
	rootCmd.Flags().StringVar(&proxyHTTPAddr, "proxy.listen-address", ":8082", "HTTP proxy listen address")
	rootCmd.Flags().StringVar(&httpAddr, "http.listen-address", ":8081", "HTTP listen address")
	rootCmd.Flags().StringVar(&openapiPath, "http.openapi-path", "", "path of openapi")
	rootCmd.Flags().StringVar(&swaggerPath, "http.swagger-path", "/apidocs/", "path of swagger ui. If the value is empty, the swagger UI is disabled.")
	rootCmd.Flags().Var(&httpExternalURL, "http.external-url", "The URL under which IDAS is externally reachable (for example, if IDAS is served via a reverse proxy). Used for generating relative and absolute links back to IDAS itself. If the URL has a path portion, it will be used to prefix all HTTP endpoints served by IDAS. If omitted, relevant URL components will be derived automatically.")
	rootCmd.Flags().StringVar(&webPrefix, "http.web-prefix", "/admin/", "The path prefix of the static page. The default is the path of http.external-url.")
	rootCmd.Flags().StringVar(&swaggerFilePath, "swagger.file-path", "", "path of swagger ui local file. If the value is empty, the swagger UI is disabled.")
}

func initParameter() {
	logger := log.NewTraceLogger()

	if httpExternalURL.Scheme == "" {
		httpExternalURL.Scheme = "http"
	}
	if httpExternalURL.Path == "" {
		httpExternalURL.Path = "/"
	}
	if httpExternalURL.Path[len(httpExternalURL.Path)-1:] != "/" {
		httpExternalURL.Path = httpExternalURL.Path + "/"
	}

	if httpExternalURL.Host == "" {
		port := "80"
		if h, p, err := net.SplitHostPort(httpAddr); err == nil {
			port = p
			ip := net.ParseIP(h)
			if ip.IsLoopback() || ip.IsGlobalUnicast() {
				httpExternalURL.Host = httpAddr
			}
		}
		if httpExternalURL.Host == "" {
			httpExternalURL.Host = net.JoinHostPort("localhost", port)
			interfaces, err := net.Interfaces()
			if err != nil {
				level.Error(logger).Log("msg", "failed to get interface, please specify a valid http.external-url.")
			} else {
			loop:
				for _, iface := range interfaces {
					addrs, err := iface.Addrs()
					if err == nil {
						for _, addr := range addrs {
							ip, _, _ := net.ParseCIDR(addr.String())
							if ip.IsGlobalUnicast() {
								httpExternalURL.Host = net.JoinHostPort(ip.String(), port)
								break loop

							}
						}
					}
				}
			}
		}
	}
	if !strings.HasPrefix(webPrefix, "/") {
		webPrefix = "/" + webPrefix
	}
	if !strings.HasSuffix(webPrefix, "/") {
		webPrefix = webPrefix + "/"
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile == "" {
		cfgFile = "./idas.yaml"
	}
	logger := log.NewTraceLogger()
	if err := config.ReloadConfigFromFile(logger, cfgFile); err != nil {
		level.Error(logger).Log("msg", "failed to load config", "err", err)
		os.Exit(1)
	}
	if configDisplay {
		var buf bytes.Buffer
		err := (&jsonpb.Marshaler{OrigName: true}).Marshal(&buf, config.Get())
		if err != nil {
			level.Error(logger).Log("msg", "failed to marshaller config", "err", err)
			os.Exit(1)
		}
		var tmpObj map[string]interface{}
		err = json.NewDecoder(&buf).Decode(&tmpObj)
		if err != nil {
			level.Error(logger).Log("msg", "failed to marshaller config", "err", err)
			os.Exit(1)
		}
		encoder := yaml.NewEncoder(os.Stdout)
		encoder.SetIndent(2)
		if err = encoder.Encode(tmpObj); err != nil {
			level.Error(logger).Log("msg", "failed to encode config", "err", err)
			os.Exit(1)
		}
		os.Exit(0)
	}
}
