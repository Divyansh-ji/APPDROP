package handlers

var AllowedWidgetTypes = map[string]bool{
	"banner":       true,
	"product_grid": true,
	"text":         true,
	"image":        true,
	"spacer":       true,
}

func IsAllowedWidgetType(t string) bool {
	return AllowedWidgetTypes[t]
}
