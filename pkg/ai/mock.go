package ai

import (
	"context"
	"fmt"
	"time"
)

type Mock struct {
	CorrectionFn func(string) (string, error)
}

func NewMock() *Mock {
	source := map[string]string{
		"halo, wie gehts?":                       "Hallo, wie geht's?",
		"Lass uns über learning of german reden": "Lass uns über das Lernen von Deutsch reden.",
	}

	return &Mock{
		CorrectionFn: func(s string) (string, error) {
			time.Sleep(1 * time.Second)
			if c, ok := source[s]; ok {
				return c, nil
			}
			return "", fmt.Errorf("can't correct")
		},
	}
}

func (m *Mock) Correct(ctx context.Context, in string) (string, error) {
	if m.CorrectionFn == nil {
		return "", fmt.Errorf("CorrectionFn is not implemented")
	}
	return m.CorrectionFn(in)
}
