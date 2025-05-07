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

type LineChartUtility struct {
	ChartInstructions *ChartCallParams
	ChartUtils        *ChartUtils
}

func (x *LineChartUtility) Build(logger zerolog.Logger) error {
	aggregatedData, err := GroupAndAggregate(x.ChartInstructions, x.ChartUtils)
	if err != nil || len(aggregatedData) == 0 {
		return fmt.Errorf("error aggregating data for Line chart: %v", err)
	}
	x.ChartInstructions.ProcessedData = aggregatedData
	line := charts.NewLine()
	var xAxis []string
	var seriesData []opts.LineData

	for _, item := range aggregatedData {
		if group, ok := item["group"].(string); ok {
			xAxis = append(xAxis, group)
		}
		if val, ok := item["value"].(float64); ok {
			seriesData = append(seriesData, opts.LineData{Value: val})
		}
	}

	line.SetGlobalOptions(getChartGlobalOpts(x.ChartInstructions, x.ChartUtils)...)
	line.SetXAxis(xAxis).AddSeries(x.ChartInstructions.XAxisField, seriesData)
	if x.ChartUtils.Page == nil {
		page := components.NewPage()
		page.AddCharts(
			line,
		)

		x.ChartUtils.FileType = "html"
		x.ChartUtils.SetFullPath()

		f, err := os.Create(x.ChartUtils.FullPath)
		if err != nil {
			logger.Error().Msg("Error creating file for Line chart")
			return err
		}
		page.Render(io.MultiWriter(f))
		bbytes, err := os.ReadFile(x.ChartUtils.FullPath)
		if err != nil {
			logger.Error().Msg("Error reading html Line chart file")
		}
		logger.Info().Msg("Line chart created successfully")
		x.ChartUtils.ChartBytes = bbytes
	} else {
		x.ChartUtils.Page.AddCharts(
			line,
		)
	}
	return nil
}
