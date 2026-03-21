package inspect

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	encLib "github.com/agejevasv/swk/internal/encode"
	"github.com/agejevasv/swk/internal/ioutil"
)

var jwtCmd = &cobra.Command{
	Use:   "jwt [token]",
	Short: "Inspect a JWT token",
	Long:  "Decode a JWT token and display its header, claims, and expiry status.",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadFileInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		info, err := encLib.JWTDecode(input)
		if err != nil {
			return err
		}

		if ioutil.MustGetBool(cmd, "json") {
			out, err := encLib.JWTInfoJSON(info)
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), string(out))
		} else {
			formatJWT(cmd.OutOrStdout(), info)
		}

		if ioutil.MustGetBool(cmd, "check-expiry") && isExpired(info) {
			return ioutil.CheckFailedError{}
		}

		return nil
	},
}

func init() {
	jwtCmd.Flags().Bool("check-expiry", false, "exit with code 1 if token is expired")
	jwtCmd.Flags().Bool("json", false, "output as JSON")
	Cmd.AddCommand(jwtCmd)
}

var registeredClaims = []string{"sub", "iss", "aud", "exp", "nbf", "iat", "jti"}

func formatJWT(out io.Writer, info *encLib.JWTInfo) {
	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)

	// Header
	fmt.Fprintln(w, "Header:")
	for _, key := range []string{"alg", "typ"} {
		if v, ok := info.Header[key]; ok {
			fmt.Fprintf(w, "  %s\t%s\n", key, formatClaimValue(key, v))
		}
	}
	// Any non-standard header fields
	for key, v := range info.Header {
		if key != "alg" && key != "typ" {
			fmt.Fprintf(w, "  %s\t%s\n", key, formatClaimValue(key, v))
		}
	}

	// Payload
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Payload:")

	// Registered claims first, in order
	for _, key := range registeredClaims {
		if v, ok := info.Payload[key]; ok {
			fmt.Fprintf(w, "  %s\t%s\n", key, formatClaimValue(key, v))
		}
	}

	// Custom claims, sorted
	registered := make(map[string]bool, len(registeredClaims))
	for _, k := range registeredClaims {
		registered[k] = true
	}

	var custom []string
	for key := range info.Payload {
		if !registered[key] {
			custom = append(custom, key)
		}
	}
	sort.Strings(custom)

	for _, key := range custom {
		fmt.Fprintf(w, "  %s\t%s\n", key, formatClaimValue(key, info.Payload[key]))
	}

	// Signature
	fmt.Fprintln(w)
	fmt.Fprintf(w, "Signature:\t%s\n", info.Signature)

	w.Flush()
}

func formatClaimValue(key string, v any) string {
	if isTimestampClaim(key) {
		return formatTimestamp(key, v)
	}

	switch val := v.(type) {
	case string:
		return val
	case float64:
		if val == float64(int64(val)) {
			return fmt.Sprintf("%d", int64(val))
		}
		return fmt.Sprintf("%g", val)
	case bool:
		return fmt.Sprintf("%t", val)
	case []any:
		if key == "aud" {
			parts := make([]string, 0, len(val))
			for _, item := range val {
				parts = append(parts, fmt.Sprintf("%v", item))
			}
			return strings.Join(parts, ", ")
		}
		b, _ := json.Marshal(val)
		return string(b)
	case map[string]any:
		b, _ := json.Marshal(val)
		return string(b)
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%v", val)
	}
}

func isTimestampClaim(key string) bool {
	return key == "exp" || key == "iat" || key == "nbf"
}

func formatTimestamp(key string, v any) string {
	f, ok := v.(float64)
	if !ok {
		return fmt.Sprintf("%v", v)
	}

	t := time.Unix(int64(f), 0).UTC()
	s := t.Format("2006-01-02 15:04:05 UTC")

	if key == "exp" {
		if time.Now().After(t) {
			s += " (expired)"
		} else {
			s += " (valid)"
		}
	}

	return s
}

func isExpired(info *encLib.JWTInfo) bool {
	if info.ExpiredAt == nil {
		return false
	}
	return time.Now().After(*info.ExpiredAt)
}
