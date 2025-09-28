package processor

import (
	"encoding/json"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func CompareJson(json1, json2 []byte) (bool, string) {
	var obj1, obj2 any
	if err := json.Unmarshal(json1, &obj1); err != nil { //TODO jsoniter
		return false, err.Error()
	}
	if err := json.Unmarshal(json2, &obj2); err != nil { //TODO jsoniter
		return false, err.Error()
	}

	opts := []cmp.Option{
		cmpopts.SortSlices(func(a, b any) bool {
			return fmt.Sprintf("%v", a) < fmt.Sprintf("%v", b)
		}),
		cmpopts.EquateEmpty(),
		cmpopts.EquateNaNs(),
	}

	if diff := cmp.Diff(obj1, obj2, opts...); diff != "" {
		return false, diff
	}
	return true, ""
}
