package plugin_cond_redirect

import "net/http"

type OrRedirectCondition struct {
	Children []RawRedirectCondition `mapstructure:"children"`
}

type orRedirectCondition struct {
	children []redirectCondition
}

func (c *OrRedirectCondition) build() (redirectCondition, error) {
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
	return orRedirectCondition{
		children: children,
	}, nil
}

func (c orRedirectCondition) check(request *http.Request) bool {
	for _, c := range c.children {
		if c.check(request) {
			return true
		}
	}
	return false
}
