package transformer

import (
	"fmt"
	"strconv"
)

func StringToInt(s string) (int, error) {
	value, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %s", err)
	}
	return value, nil
}

func IntToString(i int) string {
	return strconv.Itoa(i)

}
