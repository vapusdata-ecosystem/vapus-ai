package router

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"slices"
	"strings"
	"text/template"

	"github.com/bytedance/sonic"
	"github.com/vapusdata-ecosystem/vapusdata/core/pkgs/pbtools"
	dmutils "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/utils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	"gopkg.in/yaml.v3"
)

func limitWords(text string, limit int) string {
	words := strings.Fields(text)
	if len(words) > limit {
		return strings.Join(words[:limit], " ") + "..."
	}
	return text
}

func limitletters(text string, limit int) string {
	if len(text) == 0 {
		return "--"
	}
	letters := strings.Split(text, "")
	if len(letters) > limit {
		return strings.Join(letters[:limit], "") + "..."
	}
	return text
}

func limitlettersWD(text string, limit int) string {
	if len(text) == 0 {
		return "--"
	}
	letters := strings.Split(text, "")
	if len(letters) > limit {
		return strings.Join(letters[:limit], "")
	}
	return text
}

func EpochConverter(epoch int64) string {
	if epoch == 0 {
		return "--"
	}
	return dmutils.GetFormattedTime(epoch, "2006-01-02")
}

func EpochConverterFull(epoch int64) string {
	if epoch == 0 {
		return "--"
	}
	return dmutils.GetFormattedTime(epoch, "2006-01-02 15:04")
}

func EpochConverterFullSeconds(epoch int64) string {
	if epoch == 0 {
		return "--"
	}
	return dmutils.GetFormattedTime(epoch, "2006-01-02 15:04:05")
}

func EpochConverterTextDate(epoch int64) string {
	if epoch == 0 {
		return "--"
	}
	return dmutils.GetFormattedTime(epoch, "1 January 2006")
}

func InSlice(value string, list ...string) bool {
	return slices.Contains(list, value)
}

func SliceContains(value string, list []string) bool {
	return slices.Contains(list, value)
}

func toJSON(v interface{}) string {
	a, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(a)
}

func protoToJSON(v any) string {
	l, ok := v.(protoreflect.ProtoMessage)
	if !ok {
		a, err := sonic.Marshal(v)
		if err != nil {
			return ""
		}
		return string(a)
	}
	a, err := pbtools.ProtoJsonMarshal(l)
	if err != nil {
		return ""
	}
	return string(a)
}

func parseJSON(input string) string {
	var result interface{}
	err := json.Unmarshal([]byte(input), &result)
	if err != nil {
		return "Invalid JSON"
	}

	pretty, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "Error formatting JSON"
	}

	return string(pretty)
}

func stringCheck(s string) string {
	if s == "" {
		return "--"
	}
	return s
}

func escapeHTML(input string) string {
	return template.HTMLEscapeString(input)
}

func escapeJSON(input string) string {
	return template.JSEscapeString(input)
}

func addRand(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func randBool() bool {
	return rand.Intn(2) == 0
}

func marshalToYaml(v any) string {
	a, err := yaml.Marshal(v)
	if err != nil {
		return ""
	}
	return string(a)
}

func strContains(s string, substrL string) bool {
	substrs := strings.Split(substrL, ",")
	valid := false
	for _, substr := range substrs {
		if strings.Contains(s, substr) {
			valid = true
			break
		}
	}
	return valid
}

func sliceLen[T any](slice []T, expectedLen int, condition string) bool {
	if condition == "==" {
		return len(slice) == expectedLen
	} else if condition == ">" {
		return len(slice) > expectedLen
	} else if condition == "<" {
		return len(slice) < expectedLen
	}
	return false
}

func getSlicelen[T any](slice []T) int64 {
	return int64(len(slice))
}

func slugToTitle(s string) string {
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.ReplaceAll(s, "_", " ")
	return cases.Title(language.English).String(s)
}

func enumoTitle(tr string, s any) string {
	d := s.(protoreflect.Enum)
	value := d.Descriptor().Values().ByNumber(d.Number())
	cc := string(value.Name())
	return strTitle(tr, cc)
}

func strTitle(tr, cc string) string {
	cc = strings.ReplaceAll(cc, "_", " ")
	cc = cases.Upper(language.English).String(cc)
	if tr != "" {
		cc = strings.ReplaceAll(cc, tr, "")
	}
	return cases.Title(language.English).String(cc)
}

func strUpper(str string) string {
	str = strings.ReplaceAll(str, "-", " ")
	str = strings.ReplaceAll(str, "_", " ")
	return strings.ToUpper(str)
}

func joinSlice[T ~string](s []T, separator string) string {
	if len(s) == 0 {
		return ""
	}
	strSlice := make([]string, len(s))
	for i, v := range s {
		strSlice[i] = string(v)
	}
	return strings.Join(strSlice, separator)
}
func joinElements[T comparable](elements []T, separator string) string {
	var b strings.Builder
	for i, elem := range elements {
		if i > 0 {
			b.WriteString(separator)
		}
		b.WriteString(fmt.Sprint(elem))
	}
	return b.String()
}

func intCheck(i any) string {
	if i == 0 {
		return "NA"
	}
	return fmt.Sprint(i)
}

func floatTruncate[D any](f D) string {
	v := any(f)
	f32, ok := v.(float32)
	if ok {
		return fmt.Sprintf("%.2f", f32)
	}
	f64, ok := v.(float64)
	if ok {
		return fmt.Sprintf("%.2f", f64)
	}
	return ""
}

func typeOf(v interface{}) string {
	return reflect.TypeOf(v).String()
}
