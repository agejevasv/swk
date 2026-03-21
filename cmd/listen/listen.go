package listen

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/agejevasv/swk/internal/ioutil"
	listenLib "github.com/agejevasv/swk/internal/listen"
)

// Cmd is the top-level listen command.
var Cmd = &cobra.Command{
	Use:   "listen",
	Short: "Log incoming HTTP requests",
	Long:  "Start an HTTP server that logs all incoming requests. Useful for testing webhooks and callbacks.",
	RunE:  runListen,
}

func init() {
	Cmd.Flags().IntP("port", "p", 8080, "listen port (0 for random)")
	Cmd.Flags().StringP("host", "H", "0.0.0.0", "bind address")
	Cmd.Flags().IntP("status", "s", 200, "response status code")
	Cmd.Flags().StringP("body", "b", "", "response body")
	Cmd.Flags().Bool("no-body", false, "don't log request bodies")
}

func runListen(cmd *cobra.Command, args []string) error {
	host := ioutil.MustGetString(cmd, "host")
	port := ioutil.MustGetInt(cmd, "port")

	opts := listenLib.Options{
		Status: ioutil.MustGetInt(cmd, "status"),
		Body:   ioutil.MustGetString(cmd, "body"),
		NoBody: ioutil.MustGetBool(cmd, "no-body"),
		Writer: cmd.ErrOrStderr(),
	}

	handler := listenLib.Handler(opts)

	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.ErrOrStderr(), "Listening on http://%s\n", ln.Addr())

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
