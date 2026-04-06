package debug

import (
	"encoding/json"
	"fmt"
)

func Dump(v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Printf("%+v\n", v)
		return
	}
	fmt.Println(string(b))
}
