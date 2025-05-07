package charts

import (
	"io"
	"log"
	"os"
	"slices"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/rs/zerolog"
)

type HeatmapUtility struct {
	DataSet           *HeatmapDataset
	ChartInstructions *ChartCallParams
	ChartUtils        *ChartUtils
}

func (x *HeatmapUtility) buildHeatmapdata() error {
	aggregatedData, err := aggregateHeatmapData(x.ChartInstructions.Dataset,
		x.ChartInstructions.XAxisField,
		x.ChartInstructions.YAxisField)
	if err != nil {
		return err
	}
	x.ChartInstructions.ProcessedData = []map[string]any{{
		"heatmap": aggregatedData,
	},
	}
	var (
		xAxisVals []any
		hMap      [][]int
		yAxisVals []any
	)
	res := &HeatmapDataset{
		aggData: aggregatedData,
		XValues: make([]float64, 0),
		YValues: make([]float64, 0),
		ZValues: make([]float64, 0),
	}

	xCategoryMap := make(map[any]int)
	yCategoryMap := make(map[any]int)

	var xIndex, yIndex int

	for yKey, yMap := range aggregatedData {
		yAxisVals = append(yAxisVals, yKey)
		if _, exists := yCategoryMap[yKey]; !exists {
			yCategoryMap[yKey] = yIndex
			yIndex++
		}
		yNumericValue := float64(yCategoryMap[yKey])
		res.XValues = append(res.XValues, yNumericValue)
		for xKey, count := range yMap {
			if !slices.Contains(xAxisVals, xKey) {
				xAxisVals = append(xAxisVals, xKey)
			}
			fg := []int{int(yNumericValue), slices.Index(xAxisVals, xKey), int(count)}
			hMap = append(hMap, fg)
			if _, exists := xCategoryMap[xKey]; !exists {
				xCategoryMap[xKey] = xIndex
				xIndex++
			}
			xNumericValue := float64(xCategoryMap[xKey])

			res.XValues = append(res.XValues, xNumericValue)
			res.ZValues = append(res.ZValues, count)
		}
	}
	log.Println("XAxisVals: ", xAxisVals)
	log.Println("YAxisVals: ", yAxisVals)
	log.Println("HMap: ", hMap)
	res.xVals = xAxisVals
	res.yVals = yAxisVals
	res.hMap = hMap
	x.DataSet = res
	return nil
}

func (x *HeatmapUtility) genHeatMapData() []opts.HeatMapData {
	items := make([]opts.HeatMapData, 0)
	for i := 0; i < len(x.DataSet.hMap); i++ {
		if x.DataSet.hMap[i][2] == 0 {
			items = append(items, opts.HeatMapData{Value: [3]any{x.DataSet.hMap[i][1], x.DataSet.hMap[i][0], "-"}})
		} else {
			items = append(items, opts.HeatMapData{Value: [3]any{x.DataSet.hMap[i][1], x.DataSet.hMap[i][0], x.DataSet.hMap[i][2]}})
		}
	}
	return items
}

func (x *HeatmapUtility) heatMapBase() *charts.HeatMap {
	hm := charts.NewHeatMap()
	chartOpts := getChartGlobalOpts(x.ChartInstructions, x.ChartUtils)
	chartOpts = append(chartOpts, charts.WithXAxisOpts(opts.XAxis{
		Name:      x.ChartInstructions.YAxisField,
		Type:      "category",
		Data:      x.DataSet.xVals,
		SplitArea: &opts.SplitArea{Show: opts.Bool(true)},
	}),
		charts.WithYAxisOpts(opts.YAxis{
			Name:      x.ChartInstructions.XAxisField,
			Type:      "category",
			Data:      x.DataSet.yVals,
			SplitArea: &opts.SplitArea{Show: opts.Bool(true)},
		}),
		charts.WithVisualMapOpts(opts.VisualMap{
			Calculable: opts.Bool(true),
			Min:        0,
			Max:        10,
			InRange: &opts.VisualMapInRange{
				Color: []string{"#50a3ba", "#eac736", "#d94e5d"},
			},
		}),
	)
	hm.SetGlobalOptions(
		chartOpts...,
	)
	hm.SetXAxis(x.DataSet.yVals).AddSeries("heatmap", x.genHeatMapData())
	return hm
}

func (x *HeatmapUtility) Build(logger zerolog.Logger) error {
	err := x.buildHeatmapdata()
	if err != nil {
		logger.Error().Msg("Error building heatmap data")
		return err
	}
	if x.ChartUtils.Page == nil {
		page := components.NewPage()
		page.AddCharts(
			x.heatMapBase(),
		)

		x.ChartUtils.FileType = "html"
		x.ChartUtils.SetFullPath()

		f, err := os.Create(x.ChartUtils.FullPath)
		if err != nil {
			logger.Error().Msg("Error creating file")
			return err
		}
		page.Render(io.MultiWriter(f))
		bbytes, err := os.ReadFile(x.ChartUtils.FullPath)
		if err != nil {
			logger.Error().Msg("Error reading html chart file")
		}
		logger.Info().Msg("Heatmap chart created successfully")
		x.ChartUtils.ChartBytes = bbytes
	} else {
		x.ChartUtils.Page.AddCharts(
			x.heatMapBase(),
		)
	}
	return nil
}
