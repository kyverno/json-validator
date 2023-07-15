package cmd

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"

	"github.com/kyverno/json-validator/backend/pkg/model"
	"github.com/kyverno/json-validator/backend/pkg/ui"
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{}
	cmd.AddCommand(runCommand())

	return cmd
}

func runCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "run",
		RunE: func(cmd *cobra.Command, _ []string) error {
			host, _ := cmd.Flags().GetString("host")
			port, _ := cmd.Flags().GetInt("port")
			mode, _ := cmd.Flags().GetString("mode")

			gin.SetMode(mode)

			router := gin.New()
			router.Use(gin.Recovery())

			if gin.Mode() == gin.DebugMode {
				router.Use(gin.Logger())
				router.Use(cors.New(cors.Config{
					AllowOrigins:  []string{"*"},
					AllowMethods:  []string{"POST", "GET", "HEAD"},
					AllowHeaders:  []string{"Origin", "Content-Type"},
					ExposeHeaders: []string{"Content-Length"},
				}))
			}

			static(router)
			api(cmd, router)

			address := fmt.Sprintf("%v:%v", host, port)
			srv := &http.Server{
				Addr:    address,
				Handler: router,
			}

			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					panic(err)
				}
			}()

			fmt.Printf("Server started: %s\n", address)

			ctx, stop := signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM)
			defer stop()
			<-ctx.Done()

			stop()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			return srv.Shutdown(ctx)
		},
	}

	cmd.PersistentFlags().StringP("service", "s", "https://kyverno-svc.kyverno:443/validate/fail", "Kyverno validate API")
	cmd.PersistentFlags().String("host", "0.0.0.0", "Server Host")
	cmd.PersistentFlags().IntP("port", "p", 8080, "Server Port")
	cmd.PersistentFlags().String("mode", gin.ReleaseMode, "Server mode")

	flag.Parse()

	return cmd
}

func api(cmd *cobra.Command, router *gin.Engine) error {
	service, _ := cmd.Flags().GetString("service")
	client := newHTTPClient()

	router.POST("/validate", func(ctx *gin.Context) {
		var spec map[string]any

		err := ctx.BindJSON(&spec)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		templ := model.DefaultAdmissionReview
		templ.Request.Object.Spec = spec

		data, err := json.Marshal(templ)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		resp, err := client.Post(service, "application/json", bytes.NewReader(data))
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		body := model.Response{}

		err = json.NewDecoder(resp.Body).Decode(&body)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"valid":   body.Response.Allowed,
			"message": body.Response.Status.Message,
		})
	})

	return nil
}

func static(router *gin.Engine) error {
	fs, err := fs.Sub(ui.Files, "dist")
	if err != nil {
		return err
	}
	fileServer := http.FileServer(http.FS(fs))

	router.NoRoute(func(c *gin.Context) {
		fileServer.ServeHTTP(c.Writer, c.Request)
	})
	return nil
}

func newHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 60 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 10 * time.Second,
	}
}
