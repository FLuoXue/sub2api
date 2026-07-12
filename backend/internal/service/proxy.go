package service

import (
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	FallbackModeNone   = "none"
	FallbackModeProxy  = "proxy"
	FallbackModeDirect = "direct"

	// defaultResinPlatform is used when a Resin proxy has an empty username.
	// Resin V1 identity is "{Platform}.{AccountID}:TOKEN".
	defaultResinPlatform = "Default"
)

type Proxy struct {
	ID             int64
	Name           string
	Protocol       string
	Host           string
	Port           int
	Username       string
	Password       string
	Status         string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	ExpiresAt      *time.Time
	FallbackMode   string
	BackupProxyID  *int64
	ExpiryWarnDays int
	// IsResin marks a shared Resin forward-proxy gateway.
	// When true, request-time auth username becomes "{Username}.{accountID}"
	// (e.g. Default.13) while password stays the Resin token.
	IsResin bool
}

func (p *Proxy) IsActive() bool {
	return p.Status == StatusActive
}

// IsExpired 报告代理是否已过期（基于 expires_at，与 status 无关）。
func (p *Proxy) IsExpired(now time.Time) bool {
	return p.ExpiresAt != nil && !p.ExpiresAt.After(now)
}

// URL returns the static proxy URL using stored credentials.
// For Resin proxies this is the gateway template (Platform:TOKEN), mainly
// used by admin probe/export. Request paths must use URLForAccount.
func (p *Proxy) URL() string {
	return p.buildURL(p.Username, p.Password)
}

// URLForAccount returns the proxy URL for a specific upstream account.
// Non-Resin proxies ignore accountID and behave like URL().
// Resin proxies expand username to "{Platform}.{accountID}".
func (p *Proxy) URLForAccount(accountID int64) string {
	if p == nil {
		return ""
	}
	if !p.IsResin || accountID <= 0 {
		return p.URL()
	}
	platform := strings.TrimSpace(p.Username)
	if platform == "" {
		platform = defaultResinPlatform
	}
	username := platform + "." + strconv.FormatInt(accountID, 10)
	return p.buildURL(username, p.Password)
}

func (p *Proxy) buildURL(username, password string) string {
	if p == nil {
		return ""
	}
	u := &url.URL{
		Scheme: p.Protocol,
		Host:   net.JoinHostPort(p.Host, strconv.Itoa(p.Port)),
	}
	if username != "" && password != "" {
		u.User = url.UserPassword(username, password)
	}
	return u.String()
}

type ProxyWithAccountCount struct {
	Proxy
	AccountCount   int64
	LatencyMs      *int64
	LatencyStatus  string
	LatencyMessage string
	IPAddress      string
	Country        string
	CountryCode    string
	Region         string
	City           string
	QualityStatus  string
	QualityScore   *int
	QualityGrade   string
	QualitySummary string
	QualityChecked *int64
}

type ProxyAccountSummary struct {
	ID       int64
	Name     string
	Platform string
	Type     string
	Notes    *string
}
