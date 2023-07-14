package plugin_cond_redirect

import "net/http"

type NotRedirectCondition struct {
	Condition RawRedirectCondition `mapstructure:"condition"`
}

type notRedirectCondition struct {
	condition redirectCondition
}

func (c NotRedirectCondition) build() (redirectCondition, error) {
	refined, err := c.Condition.refine()
	if err != nil {
		return nil, err
	}
	built, err := refined.build()
	if err != nil {
		return nil, err
	}
	return notRedirectCondition{
		condition: built,
	}, nil
}

func (c notRedirectCondition) check(request *http.Request) bool {
	return !c.condition.check(request)
}
