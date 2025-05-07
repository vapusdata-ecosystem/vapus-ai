package charts

import (
	"fmt"
	"io"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/rs/zerolog"
)

type BarChartUtility struct {
	ChartInstructions *ChartCallParams
	ChartUtils        *ChartUtils
}

func (x *BarChartUtility) Build(logger zerolog.Logger) error {
	aggregatedData, err := GroupAndAggregate(x.ChartInstructions, x.ChartUtils)
	if err != nil || len(aggregatedData) == 0 {
		return fmt.Errorf("error aggregating data for bar chart: %v", err)
	}
	x.ChartInstructions.ProcessedData = aggregatedData
	bar := charts.NewBar()
	var xAxis []string
	var seriesData []opts.BarData

	for _, item := range aggregatedData {
		if group, ok := item["group"].(string); ok {
			xAxis = append(xAxis, group)
		}
		if val, ok := item["value"].(float64); ok {
			seriesData = append(seriesData, opts.BarData{Value: val})
		}
	}
	bar.SetGlobalOptions(getChartGlobalOpts(x.ChartInstructions, x.ChartUtils)...)
	bar.SetXAxis(xAxis).AddSeries(x.ChartInstructions.YAxisField, seriesData)
	if x.ChartUtils.Page == nil {
		page := components.NewPage()
		page.AddCharts(
			bar,
		)

		x.ChartUtils.FileType = "html"
		x.ChartUtils.SetFullPath()

		f, err := os.Create(x.ChartUtils.FullPath)
		if err != nil {
			logger.Error().Msg("Error creating file for bar chart")
			return err
		}
		page.Render(io.MultiWriter(f))
		bbytes, err := os.ReadFile(x.ChartUtils.FullPath)
		if err != nil {
			logger.Error().Msg("Error reading html bar chart file")
		}
		logger.Info().Msg("Barchart chart created successfully")
		x.ChartUtils.ChartBytes = bbytes
	} else {
		x.ChartUtils.Page.AddCharts(
			bar,
		)
	}
	return nil
}
