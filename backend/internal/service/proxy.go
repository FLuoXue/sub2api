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

	resinAccountIDPlaceholder = "<accountid>"
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
}

func (p *Proxy) IsActive() bool {
	return p.Status == StatusActive
}

// IsExpired 报告代理是否已过期（基于 expires_at，与 status 无关）。
func (p *Proxy) IsExpired(now time.Time) bool {
	return p.ExpiresAt != nil && !p.ExpiresAt.After(now)
}

// URL returns the proxy URL using the stored credentials and template intact.
func (p *Proxy) URL() string {
	return p.buildURL(p.Username, p.Password)
}

// URLForAccount expands the optional <accountid> username placeholder for a
// specific upstream account. Proxies without the placeholder are unchanged.
func (p *Proxy) URLForAccount(accountID int64) string {
	if p == nil {
		return ""
	}
	if accountID <= 0 || !strings.Contains(strings.ToLower(p.Username), resinAccountIDPlaceholder) {
		return p.URL()
	}
	username := expandResinUsername(p.Username, accountID)
	return p.buildURL(username, p.Password)
}

func expandResinUsername(username string, accountID int64) string {
	placeholder := resinAccountIDPlaceholder
	replacement := strconv.FormatInt(accountID, 10)
	lower := strings.ToLower(username)
	expanded := make([]byte, 0, len(username))
	for {
		index := strings.Index(lower, placeholder)
		if index < 0 {
			expanded = append(expanded, username...)
			break
		}
		expanded = append(expanded, username[:index]...)
		expanded = append(expanded, replacement...)
		username = username[index+len(placeholder):]
		lower = lower[index+len(placeholder):]
	}
	return string(expanded)
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
