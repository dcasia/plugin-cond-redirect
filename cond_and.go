package plugin_cond_redirect

import (
	"net/http"
)

type AndRedirectCondition struct {
	Children []RawRedirectCondition `mapstructure:"children"`
}

type andRedirectCondition struct {
	children []redirectCondition
}

func (c *AndRedirectCondition) build() (redirectCondition, error) {
	children := make([]redirectCondition, 0)
	for _, child := range c.Children {
		refined, err := child.refine()
		if err != nil {
			return nil, err
		}
		built, err := refined.build()
		if err != nil {
			return nil, err
		}
		children = append(children, built)
	}
	return andRedirectCondition{
		children: children,
	}, nil
}

func (c andRedirectCondition) check(request *http.Request) bool {
	for _, c := range c.children {
		if !c.check(request) {
			return false
		}
	}
	return true
}
