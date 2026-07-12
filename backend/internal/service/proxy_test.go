package service

import (
	"net/url"
	"testing"
)

func TestProxyURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		proxy Proxy
		want  string
	}{
		{
			name: "without auth",
			proxy: Proxy{
				Protocol: "http",
				Host:     "proxy.example.com",
				Port:     8080,
			},
			want: "http://proxy.example.com:8080",
		},
		{
			name: "with auth",
			proxy: Proxy{
				Protocol: "socks5",
				Host:     "socks.example.com",
				Port:     1080,
				Username: "user",
				Password: "pass",
			},
			want: "socks5://user:pass@socks.example.com:1080",
		},
		{
			name: "username only keeps no auth for compatibility",
			proxy: Proxy{
				Protocol: "http",
				Host:     "proxy.example.com",
				Port:     8080,
				Username: "user-only",
			},
			want: "http://proxy.example.com:8080",
		},
		{
			name: "with special characters in credentials",
			proxy: Proxy{
				Protocol: "http",
				Host:     "proxy.example.com",
				Port:     3128,
				Username: "first last@corp",
				Password: "p@ ss:#word",
			},
			want: "http://first%20last%40corp:p%40%20ss%3A%23word@proxy.example.com:3128",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := tc.proxy.URL(); got != tc.want {
				t.Fatalf("Proxy.URL() mismatch: got=%q want=%q", got, tc.want)
			}
		})
	}
}

func TestProxyURL_SpecialCharactersRoundTrip(t *testing.T) {
	t.Parallel()

	proxy := Proxy{
		Protocol: "http",
		Host:     "proxy.example.com",
		Port:     3128,
		Username: "first last@corp",
		Password: "p@ ss:#word",
	}

	parsed, err := url.Parse(proxy.URL())
	if err != nil {
		t.Fatalf("parse proxy URL failed: %v", err)
	}
	if got := parsed.User.Username(); got != proxy.Username {
		t.Fatalf("username mismatch after parse: got=%q want=%q", got, proxy.Username)
	}
	pass, ok := parsed.User.Password()
	if !ok {
		t.Fatal("password missing after parse")
	}
	if pass != proxy.Password {
		t.Fatalf("password mismatch after parse: got=%q want=%q", pass, proxy.Password)
	}
}

func TestProxyURLForAccount_Resin(t *testing.T) {
	t.Parallel()

	proxy := Proxy{
		Protocol: "socks5h",
		Host:     "127.0.0.1",
		Port:     2260,
		Username: "Default",
		Password: "my-token",
		IsResin:  true,
	}

	if got := proxy.URL(); got != "socks5h://Default:my-token@127.0.0.1:2260" {
		t.Fatalf("template URL mismatch: got=%q", got)
	}
	if got := proxy.URLForAccount(13); got != "socks5h://Default.13:my-token@127.0.0.1:2260" {
		t.Fatalf("account 13 URL mismatch: got=%q", got)
	}
	if got := proxy.URLForAccount(47); got != "socks5h://Default.47:my-token@127.0.0.1:2260" {
		t.Fatalf("account 47 URL mismatch: got=%q", got)
	}
	// invalid account id falls back to static template
	if got := proxy.URLForAccount(0); got != proxy.URL() {
		t.Fatalf("account 0 should fall back to template: got=%q", got)
	}
}

func TestProxyURLForAccount_NonResinUnchanged(t *testing.T) {
	t.Parallel()

	proxy := Proxy{
		Protocol: "http",
		Host:     "proxy.example.com",
		Port:     8080,
		Username: "user",
		Password: "pass",
		IsResin:  false,
	}
	if got := proxy.URLForAccount(13); got != proxy.URL() {
		t.Fatalf("non-resin should ignore account id: got=%q want=%q", got, proxy.URL())
	}
}

func TestProxyURLForAccount_ResinDefaultPlatform(t *testing.T) {
	t.Parallel()

	proxy := Proxy{
		Protocol: "socks5h",
		Host:     "127.0.0.1",
		Port:     2260,
		Password: "tok",
		IsResin:  true,
	}
	if got := proxy.URLForAccount(9); got != "socks5h://Default.9:tok@127.0.0.1:2260" {
		t.Fatalf("empty platform should default: got=%q", got)
	}
}

func TestAccountProxyURL(t *testing.T) {
	t.Parallel()

	acc := &Account{
		ID: 13,
		Proxy: &Proxy{
			Protocol: "socks5h",
			Host:     "127.0.0.1",
			Port:     2260,
			Username: "Default",
			Password: "my-token",
			IsResin:  true,
		},
	}
	if got := acc.ProxyURL(); got != "socks5h://Default.13:my-token@127.0.0.1:2260" {
		t.Fatalf("Account.ProxyURL mismatch: got=%q", got)
	}
	if got := (*Account)(nil).ProxyURL(); got != "" {
		t.Fatalf("nil account should return empty: got=%q", got)
	}
}
