package conf

import (
	"encoding/json"
	"fmt"
)

type themeElement struct {
	Size  *int    `json:"size,omitempty"`
	Color *string `json:"color,omitempty"`
}

func (e *themeElement) UnmarshalJSON(data []byte) error {
	// Define an alias to avoid recursion
	type Alias themeElement
	// Temporary struct to hold the unmarshaled data
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(e),
	}
	// Unmarshal into the temporary struct
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if e.Size != nil && *e.Size < 0 {
		return fmt.Errorf("size must be non-negative")
	}
	return nil
}

type Theme struct {
	Background   *string       `json:"background,omitempty"`
	Header       *themeElement `json:"header,omitempty"`
	SubHeader    *themeElement `json:"subHeader,omitempty"`
	SubSubHeader *themeElement `json:"subSubHeader,omitempty"`
	Link         *themeElement `json:"link,omitempty"`
	Default      *themeElement `json:"default,omitempty"`
}
