package pkg

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/jedib0t/go-pretty/v6/list"
	table "github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
)

var defaultTableStyle = table.StyleRounded

func NewTableWritter() table.Writer {
	tw := table.NewWriter()
	tw.SetStyle(defaultTableStyle)
	tw.SetAutoIndex(true)
	tw.SuppressTrailingSpaces()
	tw.SetOutputMirror(os.Stdout)
	tw.SetColumnConfigs([]table.ColumnConfig{
		{
			VAlign: text.VAlignMiddle,
		},
	})
	return tw
}

func NewListWritter(items []interface{}, style list.Style) list.Writer {
	lw := list.NewWriter()
	lw.SetStyle(style)
	lw.AppendItems(items)
	return lw
}

func GetSpinner(charSet int) *spinner.Spinner {
	if charSet == 0 {
		charSet = 43
	}
	s := spinner.New(spinner.CharSets[35], 120*time.Millisecond)
	s.Color("white")
	s.Prefix = "Executing, please wait  "
	return s
}

func LogTitles(mess string, logger zerolog.Logger) {
	xx := text.FormatTitle.Apply(mess)
	xx = text.Underline.Sprintf(xx)
	logger.Info().Msg("\n")
	logger.Info().Msg(xx)
}
func LogTitlesf(mess string, logger zerolog.Logger, args ...interface{}) {
	xx := text.FormatTitle.Apply(mess)
	xx = text.Underline.Sprintf(xx)
	logger.Info().Msg("\n")
	for _, arg := range args {
		xx = fmt.Sprintf(xx+" : %v", arg)
	}
	logger.Info().Msg(xx)
}

func LogDescription(mess string, logger zerolog.Logger, args ...interface{}) {
	xx := text.FormatDefault.Apply(mess)
	logger.Info().Msg("\n")
	logger.Info().Msg(fmt.Sprintf(xx, args...))
	logger.Info().Msg("\n")
}

func UnixTransformer(tt int64) string {
	t := time.Unix(tt, 0)
	// tf := text.NewUnixTimeTransformer("dd-MM-yyyy", time.Local)
	return t.Format("2006-01-02 15:04:05")
}

func ParseAndBuildYamlTable(val []byte) error {
	obj := make(map[string]interface{})
	err := yaml.Unmarshal(val, &obj)
	if err != nil {
		return err
	}
	tableData := [][]map[string]any{}
	parseNestedMap("", obj, &tableData, 0)
	tw := NewTableWritter()
	tw.SetAutoIndex(false)
	tw.AppendHeader(table.Row{"Key", "Value"})
	for _, row := range tableData {
		for _, cols := range row {
			for k, v := range cols {
				if v == "" {
					continue
				}
				tw.AppendRow(table.Row{k, v})
			}
		}
	}

	tw.Render()
	return nil
}

func parseNestedMap(prefix string, data interface{}, tableData *[][]map[string]any, level int) {
	indent := strings.Repeat("  ", level) // Indentation based on level of nesting
	if data == nil {
		return
	}
	switch reflect.TypeOf(data).Kind() {
	case reflect.Map:
		// Iterate through map elements
		for key, value := range data.(map[string]interface{}) {
			fullKey := indent + key
			if prefix != "" {
				fullKey = prefix + "." + key
			}
			*tableData = append(*tableData, []map[string]any{{fullKey: ""}})
			parseNestedMap(fullKey, value, tableData, level+1)
		}
	case reflect.Slice:
		// Iterate through slice elements
		for i, item := range data.([]interface{}) {
			fullKey := fmt.Sprintf("%s[%d]", prefix, i)
			*tableData = append(*tableData, []map[string]any{{fullKey: ""}})
			parseNestedMap(fullKey, item, tableData, level+1)
		}
	default:
		// Add leaf node (actual value)
		*tableData = append(*tableData, []map[string]any{{indent + prefix: fmt.Sprintf("%v", data)}})
	}
}
