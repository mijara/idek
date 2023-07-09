package idek

import (
	"encoding/json"
	"reflect"
	"strings"
)

type docParams struct {
	Name   string      `json:"name"`
	Type   string      `json:"type"`
	Fields []docParams `json:"fields,omitempty"`
}

type docView struct {
	Method  string
	Path    string
	Headers []docParams `json:"headers"`
	Input   []docParams `json:"input"`
	Output  []docParams `json:"output"`
}

var docViewHandlers []docView

func registerViewHandlerDoc[H, I, O any](method, path string) {
	docViewHandlers = append(docViewHandlers, docView{
		Method:  method,
		Path:    path,
		Headers: transformToParams(elementInstance[H]()),
		Input:   transformToParams(elementInstance[I]()),
		Output:  transformToParams(elementInstance[O]()),
	})
}

func elementInstance[H any]() any {
	if new(H) == nil {
		return nil
	}
	return *new(H)
}

func transformToParams(value any) []docParams {
	result := make([]docParams, 0)
	if value == nil {
		return result
	}

	typeOf := reflect.TypeOf(value)

	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)
		name := field.Name

		headerParts := strings.Split(field.Tag.Get("header"), ",")
		if len(headerParts) > 0 && headerParts[0] != "" {
			name = headerParts[0]
		}

		tagParts := strings.Split(field.Tag.Get("json"), ",")
		if len(tagParts) > 0 && tagParts[0] != "" {
			name = tagParts[0]
		}

		typ := field.Type.String()

		var subFields []docParams
		if field.Type.Kind() == reflect.Struct {
			typ = "struct"
			subFields = transformToParams(reflect.New(field.Type).Elem().Interface())
		}

		result = append(result, docParams{
			Name:   name,
			Type:   typ,
			Fields: subFields,
		})
	}

	return result
}

type docJSON struct {
	Views []docView
}

func generateDocsJSON() []byte {
	jsonData := docJSON{Views: docViewHandlers}
	data, _ := json.Marshal(jsonData)
	return data
}
