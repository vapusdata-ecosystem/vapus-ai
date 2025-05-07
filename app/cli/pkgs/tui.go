package pkg

// import (
// 	"fmt"
// 	"log"
// 	"os"
// 	"reflect"

// 	"github.com/rivo/tview"
// 	"gopkg.in/yaml.v3"
// )

// // Define a struct to represent the YAML data
// type Item struct {
// 	ID       string                 `yaml:"id"`
// 	Name     string                 `yaml:"name"`
// 	Category string                 `yaml:"category"`
// 	Details  map[string]interface{} `yaml:"details"`
// }

// // Function to parse a YAML file containing a list of items
// func parseYAMLFile(filePath string) ([]Item, error) {
// 	file, err := os.ReadFile(filePath)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var items []Item
// 	err = yaml.Unmarshal(file, &items)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return items, nil
// }

// // Function to display a nested map (like details) in a tree format
// func displayDetailsInTreeView(data map[string]interface{}) *tview.TreeView {
// 	root := tview.NewTreeNode("Details").SetColor(tview.ColorGreen)
// 	parseNestedMap("", data, root, 0)

// 	tree := tview.NewTreeView().
// 		SetRoot(root).
// 		SetCurrentNode(root).
// 		SetBorder(true).
// 		SetTitle("Item Details")
// 	return tree
// }

// // Recursive function to build tree nodes for nested data
// func parseNestedMap(prefix string, data interface{}, parentNode *tview.TreeNode, level int) {
// 	switch reflect.TypeOf(data).Kind() {
// 	case reflect.Map:
// 		for key, value := range data.(map[string]interface{}) {
// 			fullKey := key
// 			if prefix != "" {
// 				fullKey = prefix + "." + key
// 			}
// 			childNode := tview.NewTreeNode(fullKey).SetColor(tview.ColorYellow)
// 			parentNode.AddChild(childNode)
// 			parseNestedMap(fullKey, value, childNode, level+1)
// 		}
// 	case reflect.Slice:
// 		for i, item := range data.([]interface{}) {
// 			childNode := tview.NewTreeNode(fmt.Sprintf("[%d]", i)).SetColor(tview.ColorYellow)
// 			parentNode.AddChild(childNode)
// 			parseNestedMap(fmt.Sprintf("%s[%d]", prefix, i), item, childNode, level+1)
// 		}
// 	default:
// 		leafNode := tview.NewTreeNode(fmt.Sprintf("%s: %v", prefix, data)).SetColor(tview.ColorWhite)
// 		parentNode.AddChild(leafNode)
// 	}
// }

// func main() {
// 	app := tview.NewApplication()

// 	// Parse the YAML file with a list of items
// 	items, err := parseYAMLFile("config_list.yaml")
// 	if err != nil {
// 		log.Fatalf("Failed to parse YAML file: %v", err)
// 	}

// 	// Create a table to display the items
// 	table := tview.NewTable().
// 		SetSelectable(true, false).
// 		SetBorders(true).
// 		SetTitle("Items List").
// 		SetBorder(true)

// 	// Populate the table headers
// 	headers := []string{"ID", "Name", "Category"}
// 	for i, header := range headers {
// 		table.SetCell(0, i, tview.NewTableCell(header).
// 			SetTextColor(tview.ColorYellow).
// 			SetSelectable(false).
// 			SetAlign(tview.AlignCenter))
// 	}

// 	// Populate the table with item data
// 	for i, item := range items {
// 		table.SetCell(i+1, 0, tview.NewTableCell(item.ID).SetTextColor(tview.ColorWhite))
// 		table.SetCell(i+1, 1, tview.NewTableCell(item.Name).SetTextColor(tview.ColorWhite))
// 		table.SetCell(i+1, 2, tview.NewTableCell(item.Category).SetTextColor(tview.ColorWhite))
// 	}

// 	// Handle row selection
// 	table.SetSelectedFunc(func(row, column int) {
// 		if row == 0 {
// 			return // Ignore header row
// 		}
// 		selectedItem := items[row-1] // Get the selected item (adjust for header row)
// 		detailsView := displayDetailsInTreeView(selectedItem.Details)

// 		// Switch to details view when a row is selected
// 		app.SetRoot(detailsView, true)
// 	})

// 	// Add a back button to return to the table view from the details view
// 	backButton := tview.NewButton("Back").
// 		SetSelectedFunc(func() {
// 			app.SetRoot(table, true) // Go back to the main table view
// 		})
// 	detailsBox := tview.NewFlex().
// 		SetDirection(tview.FlexRow).
// 		AddItem(backButton, 1, 0, false)

// 	// Start the application with the table view as the root
// 	if err := app.SetRoot(table, true).Run(); err != nil {
// 		panic(err)
// 	}
// }

// // Parse YAML and recursively build tree nodes
// func parseYAMLToTreeNodes(prefix string, data interface{}, level int) *tview.TreeNode {
// 	var node *tview.TreeNode
// 	switch reflect.TypeOf(data).Kind() {
// 	case reflect.Map:
// 		// If data is a map, create a node for each key-value pair
// 		node = tview.NewTreeNode(prefix).SetColor(tview.Colors.PrimaryText)
// 		for key, value := range data.(map[string]interface{}) {
// 			child := parseYAMLToTreeNodes(key, value, level+1)
// 			node.AddChild(child)
// 		}
// 	case reflect.Slice:
// 		// If data is a slice, create a node for each item in the slice
// 		node = tview.NewTreeNode(prefix).SetColor(tview.Colors.PrimaryText)
// 		for i, item := range data.([]interface{}) {
// 			child := parseYAMLToTreeNodes(fmt.Sprintf("[%d]", i), item, level+1)
// 			node.AddChild(child)
// 		}
// 	default:
// 		// If data is a leaf value, display it as a key-value pair
// 		valueStr := fmt.Sprintf("%v", data)
// 		node = tview.NewTreeNode(fmt.Sprintf("%s: %s", prefix, valueStr)).SetColor(tview.Colors.SecondaryText)
// 	}

// 	return node
// }

// func main() {
// 	// Read the YAML file
// 	filePath := "config.yaml"
// 	file, err := os.ReadFile(filePath)
// 	if err != nil {
// 		log.Fatalf("Failed to read YAML file: %v", err)
// 	}

// 	// Parse YAML data into a generic map
// 	var yamlData map[string]interface{}
// 	err = yaml.Unmarshal(file, &yamlData)
// 	if err != nil {
// 		log.Fatalf("Failed to parse YAML file: %v", err)
// 	}

// 	// Create a new tview application
// 	app := tview.NewApplication()

// 	// Create the root node for the tree view
// 	root := tview.NewTreeNode("Root").SetColor(tview.Colors.TertiaryText)
// 	root.SetExpanded(true)

// 	// Convert parsed YAML data to tree nodes and add them to the root
// 	for key, value := range yamlData {
// 		child := parseYAMLToTreeNodes(key, value, 0)
// 		root.AddChild(child)
// 	}

// 	// Create the tree view and set its root
// 	tree := tview.NewTreeView().
// 		SetRoot(root).
// 		SetCurrentNode(root).
// 		SetBorder(true).
// 		SetTitle("YAML Data Viewer")

// 	// Run the application
// 	if err := app.SetRoot(tree, true).Run(); err != nil {
// 		panic(err)
// 	}
// }
