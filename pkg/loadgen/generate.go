package loadgen

import (
	"fmt"
	"github.com/google/uuid"
)

func GenerateByNumberOfKeys(keys int) map[string]string {
	obj := map[string]string{}
	for i := 0; i < keys; i++ {
		obj[fmt.Sprintf("%d", i)] = uuid.New().String()
	}
	return obj
}
