package registertradeorder

import (
	"fmt"
	"strings"
)

type (
	InvalidInputErr struct {
		Details map[string]string
	}
)

func (e *InvalidInputErr) Error() string {
	sb := strings.Builder{}
	sb.WriteString("invalid register trade order error")
	for key, val := range e.Details {
		sb.WriteString(fmt.Sprintf(", %s: %s", key, val))
	}

	return sb.String()
}
