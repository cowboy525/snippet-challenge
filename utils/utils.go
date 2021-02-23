package utils

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/iancoleman/strcase"
	"github.com/topoface/snippet-challenge/types"
	"golang.org/x/net/idna"
)

func GetHostnameFromSiteURL(siteURL string) string {
	u, err := url.Parse(siteURL)
	if err != nil {
		return ""
	}

	return u.Hostname()
}

func GetUrlFromRequest(r *http.Request) string {
	scheme := "https"
	if r.TLS == nil {
		scheme = "http"
	}
	return fmt.Sprintf("%s://%s%s?%s", scheme, r.Host, r.URL.Path, r.URL.Query().Encode())
}

func SplitUintsByComma(s string) []uint64 {
	idlist := strings.Split(s, ",")
	ids := []uint64{}
	for _, v := range idlist {
		if val, err := strconv.ParseUint(v, 10, 64); err == nil {
			ids = append(ids, val)
		}
	}
	return ids
}

func SplitStringByComma(s string) []string {
	list := strings.Split(s, ",")
	res := []string{}
	for i := range list {
		s := strings.TrimSpace(list[i])
		if len(s) > 0 {
			res = append(res, s)
		}
	}
	return res
}

func StructToMap(item interface{}) map[string]interface{} {
	out := structs.Map(item)
	res := map[string]interface{}{}
	for k, v := range out {
		res[strcase.ToSnake(k)] = v
	}
	return res
}

func StructToMapForGorm(item interface{}) map[string]interface{} {
	out := structs.Map(item)
	res := map[string]interface{}{}
	for k, v := range out {
		res[strcase.ToSnake(k)] = v
	}

	fields := structs.Fields(item)
	for _, field := range fields {
		curFieldName := strcase.ToSnake(field.Name())
		gorm := field.Tag("gorm")
		// this field is ignored
		if strings.TrimSpace(gorm) == "-" {
			delete(res, curFieldName)
		}
		//
		if len(gorm) > 0 {
			tags := strings.Split(gorm, ";")
			for _, value := range tags {
				v := strings.Split(value, ":")
				k := strings.TrimSpace(strings.ToLower(v[0]))
				if k == "column" && v[1] != curFieldName {
					res[v[1]] = res[curFieldName]
					delete(res, curFieldName)
				}
			}
		}
	}
	return res
}

func ConvertToString(item interface{}) string {
	b, _ := json.Marshal(item)
	return string(b)
}

func GetNowWithoutMicroseconds() types.DateTimeWithoutMicroseconds {
	now := types.DateTimeWithoutMicroseconds{Time: time.Now().UTC().Round(time.Millisecond)}
	return now
}

func GetToday() *types.Date {
	now := (types.Date)(time.Now().UTC().Round(time.Hour * 24))
	return &now
}

// GetOrderBy return orderBy
func GetOrderBy(input string, orderingFields, ordering []string, fieldMap map[string]string) []string {
	var fields []string
	var result []string

	for _, order := range SplitStringByComma(input) {
		if len(order) == 0 {
			continue
		}

		desc := order[:1] == "-"
		if desc {
			order = order[1:]
		}
		if Contains(orderingFields, order) {
			fields = append(fields, order)
			orderBy := order
			if field, ok := fieldMap[order]; ok {
				orderBy = field
			}
			if desc {
				orderBy += " desc"
			}
			result = append(result, orderBy)
		}
	}

	for _, order := range ordering {
		field := strings.Split(order, " ")[0]
		if !Contains(fields, field) {
			result = append(result, order)
		}
	}

	return result
}

// Contains check if array element exists.
func Contains(arrayType interface{}, item interface{}) bool {
	arr := reflect.ValueOf(arrayType)

	if arr.Kind() != reflect.Array && arr.Kind() != reflect.Slice {
		return false
	}

	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == item {
			return true
		}
	}
	return false
}

// RemoveProhibitedCharacters replace prohibited characters with full-width characters
func RemoveProhibitedCharacters(s string) string {
	replaced := strings.ReplaceAll(strings.ReplaceAll(s, "{", "｛"), "}", "｝")
	replaced = strings.ReplaceAll(strings.ReplaceAll(replaced, "[", "［"), "]", "］")
	return replaced
}

func GetStringFromInterface(v interface{}) string {
	if value, ok := v.(string); ok {
		return value
	}
	if ptr, ok := v.(*string); ok && ptr != nil {
		return *ptr
	}
	return ""
}

func GetUintFromInterface(v interface{}) uint64 {
	if value, ok := v.(uint64); ok {
		return value
	}
	if ptr, ok := v.(*uint64); ok && ptr != nil {
		return *ptr
	}
	return 0
}

func IsValidEmail(email string) bool {
	userRegex1 := regexp.MustCompile(`(?i)^\"([\001-\010\013\014\016-\037!#-\[\]-\177]|\\[\001-\011\013\014\016-\177])*\"`)
	userRegex2 := regexp.MustCompile(`(?i)^[-!#$%&'*+/=?^_` + "`" + `{}|~0-9A-Z]+(\.[-!#$%&'*+/=?^_` + "`" + `{}|~0-9A-Z]+)*`)
	domainRegex := regexp.MustCompile(`(?i)((?:[A-Z0-9](?:[A-Z0-9-]{0,61}[A-Z0-9])?\.)+)(?:[A-Z0-9-]{2,63}([^-]?))`)
	literalRegex := regexp.MustCompile(`(?i)\[([A-f0-9:\.]+)\]`)

	if len(email) == 0 || !strings.Contains(email, "@") {
		return false
	}

	list := strings.Split(email, "@")
	userPart := list[0]
	domainPart := list[1]

	if !userRegex1.MatchString(userPart) && !userRegex2.MatchString(userPart) {
		return false
	}

	validateDomainPart := func(domain_part string) bool {
		if domainRegex.MatchString(domainPart) {
			return true
		}
		if literalRegex.MatchString(domainPart) {
			ipAddress := literalRegex.FindStringSubmatch(domainPart)[1]
			if net.ParseIP(ipAddress) != nil {
				return true
			}
		}
		return false
	}

	if domainPart != "localhost" && !validateDomainPart(domainPart) {
		domainPart, err := idna.Lookup.ToASCII(domainPart)
		if err == nil && validateDomainPart(domainPart) {
			return true
		}
		return false
	}
	return true
}

func IsArrayInterface(data interface{}) bool {
	val := reflect.ValueOf(data)
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return val.Kind() == reflect.Array || val.Kind() == reflect.Slice
}

func JoinUint64sToString(data []uint64, delimeter string) string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(data)), delimeter), "[]")
}
