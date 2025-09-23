package main

import (
	"fmt"
	"strings"
)

// Product defines a product category with its resources and data sources
type Product struct {
	Name        string   `json:"name"`
	Resources   []string `json:"resources"`
	DataSources []string `json:"data_sources"`
}

// GetIndex parses the Resources List content and returns product categories
func GetIndex(content string) ([]Product, error) {
	products := []Product{}
	lines := strings.Split(content, "\n")

	var currentProduct *Product
	currentSection := ""

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Check if this is a new product category
		if !strings.HasPrefix(line, "edgenext_") && line != "Data Source" && line != "Resource" {
			// This is a product name
			if currentProduct != nil {
				products = append(products, *currentProduct)
			}
			currentProduct = &Product{
				Name:        line,
				Resources:   []string{},
				DataSources: []string{},
			}
			currentSection = ""
			continue
		}

		// Check if this is a section header
		if line == "Data Source" {
			currentSection = "datasource"
			continue
		}
		if line == "Resource" {
			currentSection = "resource"
			continue
		}

		// This should be a resource or data source name
		if strings.HasPrefix(line, "edgenext_") {
			if currentProduct == nil {
				return nil, fmt.Errorf("found resource/data source '%s' without product category", line)
			}

			switch currentSection {
			case "datasource":
				currentProduct.DataSources = append(currentProduct.DataSources, line)
			case "resource":
				currentProduct.Resources = append(currentProduct.Resources, line)
			default:
				// If no section specified, try to determine from the name
				if strings.Contains(line, "_") {
					// For now, default to resource if unclear
					currentProduct.Resources = append(currentProduct.Resources, line)
				}
			}
		}
	}

	// Add the last product
	if currentProduct != nil {
		products = append(products, *currentProduct)
	}

	// Validate that we have some products
	if len(products) == 0 {
		return nil, fmt.Errorf("no products found in Resources List")
	}

	return products, nil
}
