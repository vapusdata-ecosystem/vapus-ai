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

type PieChartUtility struct {
	ChartInstructions *ChartCallParams
	ChartUtils        *ChartUtils
}

func (x *PieChartUtility) Build(logger zerolog.Logger) error {
	aggregatedData, err := GroupAndAggregate(x.ChartInstructions, x.ChartUtils)
	if err != nil || len(aggregatedData) == 0 {
		return fmt.Errorf("error aggregating data for pie chart: %v", err)
	}
	x.ChartInstructions.ProcessedData = aggregatedData
	pie := charts.NewPie()
	var pieData []opts.PieData

	for _, item := range aggregatedData {
		if group, ok := item["group"].(string); ok {
			if val, ok := item["value"].(float64); ok {
				pieData = append(pieData, opts.PieData{Name: group, Value: val})
			}
		}
	}

	pie.SetGlobalOptions(
		getChartGlobalOpts(x.ChartInstructions, x.ChartUtils)...,
	)
	pie.AddSeries(x.ChartInstructions.YAxisField, pieData).
		SetSeriesOptions(charts.WithLabelOpts(opts.Label{
			Show:      opts.Bool(true),
			Formatter: "{b}: {d}%",
		}))
	if x.ChartUtils.Page == nil {
		page := components.NewPage()
		page.AddCharts(
			pie,
		)
		page.SetLayout(components.PageCenterLayout)
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
			logger.Error().Msg("Error reading html pie chart file")
		}
		logger.Info().Msg("Pie chart created successfully")
		x.ChartUtils.ChartBytes = bbytes
	} else {
		x.ChartUtils.Page.AddCharts(
			pie,
		)
	}
	return nil
}
