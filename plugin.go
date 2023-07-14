package plugin_cond_redirect

import (
	"context"
	"net/http"
	"regexp"
)

type RedirectRule struct {
	WithHost           bool                 `mapstructure:"withHost,omitempty"`
	SourcePattern      string               `mapstructure:"sourcePattern"`
	DestinationPattern string               `mapstructure:"destinationPattern"`
	Condition          RawRedirectCondition `mapstructure:"condition"`
}

type RawRedirectCondition struct {
	T    string                 `mapstructure:"type"`
	Data map[string]interface{} `mapstructure:",remain"`
}

type RedirectCondition interface {
	build() (redirectCondition, error)
}

type Config struct {
	StatusCode int            `mapstructure:"statusCode,omitempty"`
	Rules      []RedirectRule `mapstructure:"rules,omitempty"`
}

type ConditionalRedirect struct {
	next       http.Handler
	config     *Config
	name       string
	statusCode int
	rules      []redirectRule
}

func CreateConfig() *Config {
	return &Config{
		StatusCode: 0,
		Rules:      make([]RedirectRule, 0),
	}
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	rules := make([]redirectRule, 0)
	for _, r := range config.Rules {
		refined, err := r.Condition.refine()
		if err != nil {
			return nil, err
		}
		condition, err := refined.build()
		if err != nil {
			return nil, err
		}
		rules = append(rules, redirectRule{
			withHost:    r.WithHost,
			source:      regexp.MustCompile(r.SourcePattern),
			destination: r.DestinationPattern,
			condition:   condition,
		})
	}
	statusCode := config.StatusCode
	if statusCode == 0 {
		statusCode = 302
	}
	return &ConditionalRedirect{
		next:       next,
		config:     config,
		name:       name,
		statusCode: statusCode,
		rules:      rules,
	}, nil
}

func (c *ConditionalRedirect) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	url := req.URL.String()
	uri := req.URL.RequestURI()
	for _, r := range c.rules {
		src := uri
		if r.withHost {
			src = url
		}
		if r.source.MatchString(src) && r.condition.check(req) {
			rw.Header().Set("Location", r.source.ReplaceAllString(src, r.destination))
			rw.WriteHeader(c.statusCode)
			return
		}
	}

	c.next.ServeHTTP(rw, req)
}
