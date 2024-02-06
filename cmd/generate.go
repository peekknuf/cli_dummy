// cmd/generate.go

package cmd

import (
	sc "data_generation/src_csv"
	"fmt"
	"sync"

	"github.com/spf13/cobra"
)

var (
	numRows        int
	outputFilename string
	selectedCols   []string
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate e-commerce data",
	Long:  `Generate e-commerce data based on selected columns.`,
	Run: func(cmd *cobra.Command, args []string) {
		GenerateData(numRows, outputFilename, selectedCols)
	},
}

func GenerateData(numRows int, outputFilename string, selectedCols []string) {
	ch := make(chan sc.Row, 1000)

	var wg sync.WaitGroup

	wg.Add(1)
	go sc.GenerateData(numRows, selectedCols, &wg, ch)

	wg.Add(1)
	go sc.WriteToCSV(outputFilename, ch, &wg, selectedCols)

	wg.Wait()

	fmt.Printf("Generated %d rows of e-commerce data and saved to %s\n", numRows, outputFilename)
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().IntVarP(&numRows, "rows", "r", 3000000, "Number of rows to generate")
	generateCmd.Flags().StringVarP(&outputFilename, "output", "o", "ecommerce_data.csv", "Output filename")
	generateCmd.Flags().StringSliceVarP(&selectedCols, "columns", "c", []string{"ID", "Timestamp", "ProductName", "Company", "Price", "Quantity", "Discount", "TotalPrice", "CustomerID", "FirstName", "LastName", "Email", "Address", "City", "State", "Zip", "Country"}, "Selected columns to generate")
}
