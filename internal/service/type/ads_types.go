package ntype

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

type Float64Array []float64
type StringArray []string

func (a Float64Array) Value() (driver.Value, error) {
	return fmt.Sprintf("{%s}", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(a)), ","), "[]")), nil
}

func (a *Float64Array) Scan(value interface{}) error {
	if value == nil {
		*a = Float64Array{}
		return nil
	}

	strVal, ok := value.(string)
	if !ok {
		return fmt.Errorf("unsupported data type: %T", value)
	}

	strVal = strings.Trim(strVal, "{}")
	strArr := strings.Split(strVal, ",")

	var result []float64
	for _, v := range strArr {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}
		result = append(result, f)
	}

	*a = result
	return nil
}

func (a StringArray) Value() (driver.Value, error) {
	return fmt.Sprintf("{%s}", strings.Join(a, ",")), nil
}

func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = StringArray{}
		return nil
	}

	strVal, ok := value.(string)
	if !ok {
		return fmt.Errorf("unsupported data type: %T", value)
	}

	strVal = strings.Trim(strVal, "{}")
	*a = strings.Split(strVal, ",")

	return nil
}
