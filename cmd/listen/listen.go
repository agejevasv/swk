package listen

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/agejevasv/swk/internal/ioutil"
	listenLib "github.com/agejevasv/swk/internal/listen"
)

// Cmd is the top-level listen command.
var Cmd = &cobra.Command{
	Use:   "listen",
	Args:  cobra.NoArgs,
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

	port, err := ioutil.IntInRange(cmd, "port", 0, 65535)
	if err != nil {
		return err
	}

	// net/http panics on a status outside 100-999.
	// 1xx is an interim response: the client would wait for a final one.
	status, err := ioutil.IntInRange(cmd, "status", 200, 599)
	if err != nil {
		return err
	}

	opts := listenLib.Options{
		Status: status,
		Body:   ioutil.MustGetString(cmd, "body"),
		NoBody: ioutil.MustGetBool(cmd, "no-body"),
		Writer: cmd.ErrOrStderr(),
	}

	handler := listenLib.Handler(opts)

	// JoinHostPort brackets IPv6 literals; "%s:%d" would produce "::1:8080".
	ln, err := net.Listen("tcp", net.JoinHostPort(host, strconv.Itoa(port)))
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.ErrOrStderr(), "Listening on http://%s\n", displayAddr(ln.Addr().String()))

	server := &http.Server{
		Handler: handler,
		// Without a header timeout a single idle connection can hold the
		// server open indefinitely.
		ReadHeaderTimeout: 10 * time.Second,
	}

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

// displayAddr turns a wildcard bind address into one that can be pasted into
// a browser: "[::]:8080" is where the server listens, not where to reach it.
func displayAddr(addr string) string {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return addr
	}
	if host == "" || host == "::" || host == "0.0.0.0" {
		return net.JoinHostPort("localhost", port)
	}
	return net.JoinHostPort(host, port)
}
