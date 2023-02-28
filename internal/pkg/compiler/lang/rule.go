package lang

import (
	"github.com/ermos/polyrule/internal/pkg/types"
	"strings"
)

type Rule func(b *strings.Builder, name string, vType types.Type, value interface{}, indent int) error
