package grep

type Grep struct {
	cfg     Config
	matcher Matcher
}

func NewGrep(cfg Config) (*Grep, error) {
	if cfg.Context > 0 {
		if cfg.After == 0 {
			cfg.After = cfg.Context
		}
		if cfg.Before == 0 {
			cfg.Before = cfg.Context
		}
	}

	var m Matcher
	var err error

	if cfg.Fixed {
		m = NewFixedMatcher(cfg)
	} else {
		m, err = NewRegexMatcher(cfg)
		if err != nil {
			return nil, err
		}
	}

	return &Grep{cfg: cfg, matcher: m}, nil
}
