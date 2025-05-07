package charts

import (
	"fmt"
	"image/color"
	"log"
	"sort"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/rs/zerolog"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func GetChartFlexPage() *components.Page {
	page := components.NewPage()
	page.SetLayout(components.PageFlexLayout)
	page.CSSAssets.Add(`
	.go-echarts-page-flex {
		flex-wrap: wrap;
		gap: 20px;
	}
	`)
	return page
}

func getChartGlobalOpts(params *ChartCallParams, utils *ChartUtils) []charts.GlobalOpts {
	return []charts.GlobalOpts{
		charts.WithTitleOpts(opts.Title{
			Title: params.ChartTitle,
			Top:   "0", // places the title at the top
		}),
		charts.WithLegendOpts(opts.Legend{
			Bottom: "0",      // positions the legend at the bottom
			Left:   "center", // centers the legend horizontally
		}),
		charts.WithInitializationOpts(opts.Initialization{
			Theme:  "black",
			Height: utils.Height,
			Width:  utils.Width,
		}),
	}
}

func BuildChartFromTool(tool *ChartCallParams, utils *ChartUtils, logger zerolog.Logger) error {
	var err error
	switch tool.ChartType {
	case BarChart:
		utils.Logger.Info().Msg("Building Bar Chart")
		err = buildBarChart(tool, utils, logger)
	case LineChart:
		utils.Logger.Info().Msg("Building Line Chart")
		err = buildLineChart(tool, utils, logger)
	case Heatmap:
		utils.Logger.Info().Msg("Building Heatmap")
		err = buildHeatMap(tool, utils, logger)
	case ScatterChart:
		err = buildScatteredChart(tool, utils, logger)
	case PieChart:
		utils.Logger.Info().Msg("Building Pie Chart")
		err = buildPieChart(tool, utils, logger)
	default:
		return fmt.Errorf("unsupported chart type: %s", tool.ChartType)
	}
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error building chart")
		return err
	}
	return nil
}

func buildHeatMap(tool *ChartCallParams, utils *ChartUtils, logger zerolog.Logger) error {
	chartBuilder := &HeatmapUtility{
		ChartUtils:        utils,
		ChartInstructions: tool,
	}
	return chartBuilder.Build(logger)
}

func buildScatteredChart(tool *ChartCallParams, utils *ChartUtils, logger zerolog.Logger) error {
	chartBuilder := &ScatteredChartUtility{
		ChartUtils:        utils,
		ChartInstructions: tool,
	}
	return chartBuilder.Build(logger)
}

func buildPieChart(tool *ChartCallParams, utils *ChartUtils, logger zerolog.Logger) error {
	chartBuilder := &PieChartUtility{
		ChartUtils:        utils,
		ChartInstructions: tool,
	}
	return chartBuilder.Build(logger)
}

func buildBarChart(tool *ChartCallParams, utils *ChartUtils, logger zerolog.Logger) error {
	chartBuilder := &BarChartUtility{
		ChartUtils:        utils,
		ChartInstructions: tool,
	}
	return chartBuilder.Build(logger)
}
func buildLineChart(tool *ChartCallParams, utils *ChartUtils, logger zerolog.Logger) error {
	chartBuilder := &LineChartUtility{
		ChartUtils:        utils,
		ChartInstructions: tool,
	}
	return chartBuilder.Build(logger)
}

func buildLineChartGonum(tool *ChartCallParams, utils *ChartUtils) (*plot.Plot, error) {
	// Typically, "line" ignores aggregator & uses (x,y) pairs directly
	rsDataset := []map[string]any{}
	var err error
	if len(tool.GroupByFields) > 0 && tool.AggregateMethod != "" && tool.AggregateMethod != "none" {
		rsDataset, err = groupAndAggregate(tool.Dataset, tool.GroupByFields, tool.YAxisField, tool.AggregateMethod, utils)
		if err != nil {
			return nil, err
		}
		// We'll treat "label" as category, "value" as numeric
		// return buildPlainBar(grouped, "label", "value", tool.ChartTitle)
	} else {
		rsDataset = tool.Dataset
	}
	xKey := tool.XAxisField
	if xKey == "" {
		xKey = "x"
	}
	yKey := tool.YAxisField
	if yKey == "" {
		yKey = "y"
	}

	pts := make(plotter.XYs, 0, len(rsDataset))
	for i, row := range rsDataset {
		rawX, ok := row[xKey]
		if !ok {
			return nil, fmt.Errorf("missing '%s' in row %d", xKey, i)
		}
		xVal, okF := toFloat(rawX)
		if !okF {
			return nil, fmt.Errorf("x value '%s' in row %d not numeric", xKey, i)
		}

		rawY, ok := row[yKey]
		if !ok {
			return nil, fmt.Errorf("missing '%s' in row %d", yKey, i)
		}
		yVal, okF2 := toFloat(rawY)
		if !okF2 {
			return nil, fmt.Errorf("y value '%s' in row %d not numeric", yKey, i)
		}

		pts = append(pts, plotter.XY{X: xVal, Y: yVal})
	}

	// Sort by X to get a proper line
	sort.Slice(pts, func(i, j int) bool { return pts[i].X < pts[j].X })

	p := plot.New()
	p.Title.Text = "Line Chart"
	p.X.Label.Text = xKey
	p.Y.Label.Text = yKey

	line, err := plotter.NewLine(pts)
	if err != nil {
		return nil, err
	}
	line.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255} // red line

	p.Add(line)
	return p, nil
}

func buildBarChartGonum(tool *ChartCallParams, utils *ChartUtils) (*plot.Plot, error) {
	// If grouping is requested:
	if len(tool.GroupByFields) > 0 && tool.AggregateMethod != "" && tool.AggregateMethod != "none" {
		grouped, err := groupAndAggregate(tool.Dataset, tool.GroupByFields, tool.YAxisField, tool.AggregateMethod, utils)
		if err != nil {
			utils.Logger.Error().Err(err).Msg("Error grouping and aggregating data")
			return nil, err
		}
		// We'll treat "label" as category, "value" as numeric
		return buildPlainBarGonum(grouped, "group", "value", tool.ChartTitle, utils)
	}
	xKey := tool.XAxisField
	if xKey == "" {
		xKey = "x"
	}
	yKey := tool.YAxisField
	if yKey == "" {
		yKey = "y"
	}
	return buildPlainBarGonum(tool.Dataset, xKey, yKey, tool.ChartTitle, utils)
}

func buildPlainBarGonum(
	data []map[string]any,
	xKey, yKey, chartTitle string, utils *ChartUtils) (*plot.Plot, error) {

	var categories []string
	vals := make(plotter.Values, 0, len(data))

	for i, row := range data {
		rawX, ok := row[xKey]
		if !ok {
			return nil, fmt.Errorf("missing '%s' in row %d", xKey, i)
		}
		category := fmt.Sprintf("%v", rawX)

		rawY, ok := row[yKey]
		if !ok {
			return nil, fmt.Errorf("missing '%s' in row %d", yKey, i)
		}
		yVal, okF := toFloat(rawY)
		if !okF {
			return nil, fmt.Errorf("value for '%s' in row %d not numeric", yKey, i)
		}

		categories = append(categories, category)
		vals = append(vals, yVal)
	}
	log.Println(utils)
	p := plot.New()
	p.Title.Text = chartTitle
	p.Y.Label.TextStyle.Rotation = 45
	p.X.Label.TextStyle.Color = color.RGBA{R: 100, G: 140, B: 240, A: 255}
	p.Y.Label.TextStyle.Color = color.RGBA{R: 100, G: 140, B: 240, A: 255}
	// Use nominal X for string categories
	p.NominalX(categories...)

	bars, err := plotter.NewBarChart(vals, vg.Points(14))
	if err != nil {
		return nil, err
	}
	bars.Color = color.RGBA{R: 0, G: 0, B: 0, A: 255}

	p.Add(bars)
	return p, nil
}
