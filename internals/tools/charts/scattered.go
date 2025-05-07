package charts

import (
	"fmt"
	"io"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/rs/zerolog"
)

type ScatteredChartUtility struct {
	ChartInstructions *ChartCallParams
	ChartUtils        *ChartUtils
}

func (x *ScatteredChartUtility) Build(logger zerolog.Logger) error {
	scatterData, err := AggregateScatterData(x.ChartInstructions)
	if err != nil || len(scatterData) == 0 {
		return fmt.Errorf("error aggregating data for Scattered chart: %v", err)
	}

	scatter := charts.NewScatter()

	scatter.SetGlobalOptions(getChartGlobalOpts(x.ChartInstructions, x.ChartUtils)...)
	scatter.AddSeries(x.ChartInstructions.XAxisField, scatterData)
	if x.ChartUtils.Page == nil {
		page := components.NewPage()
		page.AddCharts(
			scatter,
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
			logger.Error().Msg("Error reading html Scattered chart file")
		}
		logger.Info().Msg("Scattered chart created successfully")
		x.ChartUtils.ChartBytes = bbytes
	} else {
		x.ChartUtils.Page.AddCharts(
			scatter,
		)
	}
	return nil
}
