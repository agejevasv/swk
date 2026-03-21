package serve

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/agejevasv/swk/internal/ioutil"
	serveLib "github.com/agejevasv/swk/internal/serve"
)

// Cmd is the top-level serve command.
var Cmd = &cobra.Command{
	Use:   "serve [dir]",
	Short: "Start a local static file server",
	Long:  "Serve a directory over HTTP (or HTTPS with --tls). Defaults to the current directory on port 8080.",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runServe,
}

func init() {
	Cmd.Flags().IntP("port", "p", 8080, "listen port (0 for random)")
	Cmd.Flags().StringP("host", "H", "0.0.0.0", "bind address")
	Cmd.Flags().Bool("cors", false, "enable permissive CORS headers")
	Cmd.Flags().Bool("no-index", false, "disable directory listing")
	Cmd.Flags().Bool("no-log", false, "disable access logging")
	Cmd.Flags().Bool("tls", false, "enable HTTPS")
	Cmd.Flags().String("cert", "cert.pem", "TLS certificate file")
	Cmd.Flags().String("key", "cert-key.pem", "TLS private key file")
}

func runServe(cmd *cobra.Command, args []string) error {
	var dir string
	if len(args) > 0 {
		abs, err := filepath.Abs(args[0])
		if err != nil {
			return fmt.Errorf("invalid path: %w", err)
		}
		dir = abs
	} else {
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("cannot determine working directory: %w", err)
		}
		dir = wd
	}

	info, err := os.Stat(dir)
	if err != nil {
		return fmt.Errorf("cannot access %s: %w", dir, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", dir)
	}

	host := ioutil.MustGetString(cmd, "host")
	port := ioutil.MustGetInt(cmd, "port")

	opts := serveLib.Options{
		Root:      dir,
		Host:      host,
		Port:      port,
		CORS:      ioutil.MustGetBool(cmd, "cors"),
		NoIndex:   ioutil.MustGetBool(cmd, "no-index"),
		NoLog:     ioutil.MustGetBool(cmd, "no-log"),
		LogWriter: cmd.ErrOrStderr(),
	}

	handler := serveLib.Handler(opts)

	useTLS := ioutil.MustGetBool(cmd, "tls")

	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return err
	}

	scheme := "http"
	if useTLS {
		scheme = "https"
		certFile := ioutil.MustGetString(cmd, "cert")
		keyFile := ioutil.MustGetString(cmd, "key")
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return fmt.Errorf("failed to load TLS certificate: %w", err)
		}
		ln = tls.NewListener(ln, &tls.Config{Certificates: []tls.Certificate{cert}})
	}

	fmt.Fprintf(cmd.ErrOrStderr(), "Serving %s on %s://%s\n", dir, scheme, ln.Addr())

	server := &http.Server{Handler: handler}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(shutdownCtx)
	}()

	err = server.Serve(ln)
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}
