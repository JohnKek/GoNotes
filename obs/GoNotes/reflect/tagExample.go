package main

import (
	"fmt"
	"reflect"
	"strings"
)

type Product struct {
	ID          int     `meta:"id,primary_key"`
	Name        string  `meta:"name,required,description:Product name"`
	Price       float64 `meta:"price,required,description:Product price"`
	Description string  `meta:"description,optional,description:Detailed product description"`
}

func main() {
	product := Product{Name: "Gadget", Price: 99.99}

	t := reflect.TypeOf(product)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		metaTag := field.Tag.Get("meta")
		metaParts := strings.Split(metaTag, ",")

		// Извлечение атрибутов
		fieldName := field.Name
		isRequired := false
		description := ""

		for _, part := range metaParts {
			if part == "required" {
				isRequired = true
			} else if strings.HasPrefix(part, "description:") {
				description = strings.TrimPrefix(part, "description:")
			}
		}

		// Пример использования
		fmt.Printf("Field: %s\n", fieldName)
		fmt.Printf("  Required: %t\n", isRequired)
		fmt.Printf("  Description: %s\n", description)
	}
}
