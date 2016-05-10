package tweethour

import (
	"net/url"
	"strings"
	"fmt"
	"encoding/base64"
)

func urlEncode(s string) string {
	v := url.Values{}
	v.Set("p", s)
	
	f := func(c rune) bool {
		return string(c) == "="
	}
	return strings.FieldsFunc(v.Encode(), f)[1]
}

func getEncodedKey(key string,secret string) string {

	combinedKey  := fmt.Sprintf("%s:%s",urlEncode(key),urlEncode(secret))
	return  base64.StdEncoding.EncodeToString([]byte(combinedKey))
}
