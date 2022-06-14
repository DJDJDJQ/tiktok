package utils

import "regexp"

func FilterSQLInject(str string) bool {
	//p := `(?:')|(?:--)|(/\\*(?:.|[\\n\\r])*?\\*/)|(\b(select|update|and|or|delete|insert|trancate|char|chr|into|substr|ascii|declare|exec|count|master|into|drop|execute)\b)`
	regEx := "[`~!@#$%^&*()+=|{}':;',\\[\\].<>/?~！@#￥%……&*（）——+|{}【】‘；：”“’。，、？]"
	re, err := regexp.Compile(regEx)
	if err != nil {
		return false
	}
	return re.MatchString(str)
}
