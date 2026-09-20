package inspect

import "testing"

func FuzzParseSubnet(f *testing.F) {
	for _, s := range []string{"192.168.1.0/24", "10.0.0.0/8", "0.0.0.0/0", "1.2.3.4/32", "1.2.3.4/31",
		"::ffff:192.168.1.0/120", "2001:db8::/32", "", "/", "1.2.3.4", "1.2.3.4/33", "1.2.3.4/-1"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		info, err := ParseSubnet(s)
		if err != nil {
			return
		}
		// A successful parse must produce a self-consistent result.
		if info.Hosts == 0 && info.First == "" {
			t.Fatalf("empty result for %q: %+v", s, info)
		}
	})
}

func FuzzParseURL(f *testing.F) {
	for _, s := range []string{"https://example.com", "http://a:1/b?c=d#e", "://", "https://", "",
		"https://[::1]:8080/x", "https://user:pw@host/p"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		info, err := ParseURL(s)
		if err != nil {
			return
		}
		_ = URLInfoTable(info)
		_, _ = URLInfoJSON(info)
	})
}

func FuzzCertDecode(f *testing.F) {
	f.Add("-----BEGIN CERTIFICATE-----\nAAAA\n-----END CERTIFICATE-----\n")
	f.Add("")
	f.Add("not a pem")
	f.Fuzz(func(t *testing.T, s string) {
		_, _ = CertDecode([]byte(s))
	})
}
