package plugin_cond_redirect

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"regexp"
)

type redirectRule struct {
	withHost    bool
	source      *regexp.Regexp
	destination string
	condition   redirectCondition
}

type redirectCondition interface {
	check(request *http.Request) bool
}

func (raw *RawRedirectCondition) refine() (RedirectCondition, error) {
	if raw.T == "header" {
		var result HeaderRedirectCondition
		err := mapstructure.Decode(raw.Data, &result)
		return &result, err
	} else if raw.T == "cookie" {
		var result CookieRedirectCondition
		err := mapstructure.Decode(raw.Data, &result)
		return &result, err
	} else if raw.T == "and" {
		var result AndRedirectCondition
		err := mapstructure.Decode(raw.Data, &result)
		return &result, err
	} else if raw.T == "or" {
		var result OrRedirectCondition
		err := mapstructure.Decode(raw.Data, &result)
		return &result, err
	} else if raw.T == "not" {
		var result NotRedirectCondition
		err := mapstructure.Decode(raw.Data, &result)
		return &result, err
	}

	return nil, fmt.Errorf("unknown condition type: %s", raw.T)
}
