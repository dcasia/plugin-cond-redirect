package plugin_cond_redirect

import (
	"net/http"
	"regexp"
)

type CookieRedirectCondition struct {
	Name     string `mapstructure:"name"`
	Optional bool   `mapstructure:"optional"`
	Path     string `mapstructure:"path"`
	Pattern  string `mapstructure:"pattern"`
}

type cookieRedirectCondition struct {
	CookieRedirectCondition
	pattern *regexp.Regexp
}

func (c CookieRedirectCondition) build() (redirectCondition, error) {
	pattern := regexp.MustCompile(c.Pattern)
	return cookieRedirectCondition{
		CookieRedirectCondition: c,
		pattern:                 pattern,
	}, nil
}

func (c cookieRedirectCondition) check(request *http.Request) bool {
	cookie, err := request.Cookie(c.Name)
	if err != nil && !c.Optional {
		return false
	}

	if c.Path != "" && cookie.Path != c.Path {
		return false
	}

	if !c.pattern.MatchString(cookie.Value) {
		return false
	}

	return true
}
