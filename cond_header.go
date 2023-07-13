package plugin_cond_redirect

import (
	"net/http"
	"regexp"
)

type HeaderRedirectCondition struct {
	Name     string `mapstructure:"name"`
	Optional bool   `mapstructure:"optional"`
	Pattern  string `mapstructure:"pattern"`
}

type headerRedirectCondition struct {
	HeaderRedirectCondition
	pattern *regexp.Regexp
}

func (c *HeaderRedirectCondition) build() (redirectCondition, error) {
	pattern := regexp.MustCompile(c.Pattern)
	return headerRedirectCondition{
		HeaderRedirectCondition: *c,
		pattern:                 pattern,
	}, nil
}

func (c headerRedirectCondition) check(request *http.Request) bool {
	headerValue := request.Header.Get(c.Name)
	if headerValue == "" && !c.Optional {
		return false
	}

	if !c.pattern.MatchString(headerValue) {
		return false
	}

	return true
}
