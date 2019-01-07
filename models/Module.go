package models

import "bytes"

type Module struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Provider  string `json:"provider"`
	Version   string `json:"version"`

	Data      func() (*bytes.Buffer, error) `json:"-"`
}
