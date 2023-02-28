package errors

import (
	"fmt"
	"github.com/ermos/polyrule/internal/pkg/types"
)

func UnsupportedTypeForRule(vType types.Type, name string) error {
	return fmt.Errorf("unsupported type \"%s\" for rule %s", vType, name)
}
