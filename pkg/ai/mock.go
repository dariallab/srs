package ai

import (
	"context"
	"fmt"
	"time"
)

type Mock struct {
	CorrectionFn func(string) (string, error)
	ResponseFn   func(string) (string, error)
}

func NewMock() *Mock {
	corrections := map[string]string{
		"halo, wie gehts?":                       "Hallo, wie geht's?",
		"Lass uns über learning of german reden": "Lass uns über das Lernen von Deutsch reden.",
	}
	responses := map[string]string{
		"halo, wie gehts?":                       "Hallo, mir geht es gut, danke! Wie kann ich Ihnen helfen?",
		"Lass uns über learning of german reden": "Natürlich! Was möchtest du über das Erlernen der deutschen Sprache wissen? Hast du spezifische Fragen oder brauchst du Hilfe bei bestimmten Themen? Ich stehe dir gerne zur Verfügung.",
	}

	return &Mock{
		CorrectionFn: func(s string) (string, error) {
			time.Sleep(1 * time.Second)
			if c, ok := corrections[s]; ok {
				return c, nil
			}
			return "", fmt.Errorf("can't correct")
		},
		ResponseFn: func(s string) (string, error) {
			time.Sleep(1 * time.Second)
			if c, ok := responses[s]; ok {
				return c, nil
			}
			return "", fmt.Errorf("can't respond")
		},
	}
}

func (m *Mock) Correct(ctx context.Context, in string) (string, error) {
	if m.CorrectionFn == nil {
		return "", fmt.Errorf("CorrectionFn is not implemented")
	}
	return m.CorrectionFn(in)
}

func (m *Mock) Response(ctx context.Context, in string) (string, error) {
	if m.ResponseFn == nil {
		return "", fmt.Errorf("ResponseFn is not implemented")
	}
	return m.ResponseFn(in)
}
