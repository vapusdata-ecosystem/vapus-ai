package charts

import (
	"path/filepath"
	"strings"

	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/rs/zerolog"
)

const (
	PieChart      = "pie"
	LineChart     = "line"
	BarChart      = "bar"
	DoughnutChart = "doughnut"
	Heatmap       = "heatmap"
	ScatterChart  = "scatter"
)

const (
	ChartDataAvgMethod   = "average"
	ChartDataCountMethod = "count"
	ChartDataSumMethod   = "sum"
)

type ChartTool struct {
	Charts []*ChartCallParams `json:"charts"`
}

type ChartCallParams struct {
	ChartType       string           `json:"chartType"`
	GroupByFields   []string         `json:"groupFields"`
	XAxisField      string           `json:"xAxisField"`
	YAxisField      string           `json:"yAxisField"`
	ValueField      string           `json:"valueField"`
	AggregateMethod string           `json:"aggregateMethod"`
	Dataset         []map[string]any `json:"dataset"`
	ChartTitle      string           `json:"chartTitle"`
	ProcessedData   []map[string]any `json:"processedData"`
	Description     string           `json:"description"`
}

type ChartUtils struct {
	Filename   string
	FileType   string
	Path       string
	Height     string
	Width      string
	Logger     zerolog.Logger
	FullPath   string
	ChartBytes []byte
	Page       *components.Page
}

func (utils *ChartUtils) SetFullPath() {
	if utils.FileType == "" {
		utils.Filename = utils.Filename + ".png"
	} else {
		utils.Filename = utils.Filename + "." + strings.ToLower(utils.FileType)
	}
	utils.FullPath = filepath.Join(utils.Path, utils.Filename)
}

type Slices struct {
	Label string
	Value float64
}

type HeatmapDataset struct {
	XValues []float64
	YValues []float64
	ZValues []float64
	aggData map[interface{}]map[interface{}]float64
	xVals   []interface{}
	yVals   []interface{}
	hMap    [][]int
}
