package php

import (
	"fmt"
	"github.com/ermos/polyrule/internal/pkg/types"
)

func mapType(vType types.Type) string {
	if vType == types.String {
		return "string"
	} else if vType == types.Int {
		return "int"
	} else if vType == types.Bool {
		return "bool"
	} else if vType == types.Float {
		return "float"
	} else {
		panic(fmt.Errorf("unsupported type %s with php", vType))
	}
}
