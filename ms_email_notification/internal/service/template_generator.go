package service

import (
	"bytes"
	"path/filepath"
	"text/template"
)

func TemplateGenerator(input AuthNotificationInput) (*bytes.Buffer, error) {
	path, err := filepath.Abs("./internal/template/base.html")
	if err != nil {
		return nil, err
	}

	t := template.Must(template.New("base.html").ParseFiles(path))

	var buff bytes.Buffer
	if err := t.Execute(&buff, input); err != nil {
		return nil, err
	}

	return &buff, nil
}
