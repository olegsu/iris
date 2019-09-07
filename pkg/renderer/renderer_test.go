package renderer_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/olegsu/iris/pkg/renderer"
	"github.com/stretchr/testify/assert"
)

func Test_renderer_Render(t *testing.T) {
	tests := []struct {
		name            string
		templateReaders map[string]io.Reader
		valueReaders    map[string][]io.Reader
		wantErr         bool
		result          *bytes.Buffer
	}{
		{
			name:    "No readers should retrun empty []byte",
			wantErr: false,
			result:  bytes.NewBufferString(""),
		},
		{
			templateReaders: map[string]io.Reader{
				"key": strings.NewReader("{{ .Values.key }}\n"),
			},
			valueReaders: map[string][]io.Reader{
				"Values": []io.Reader{strings.NewReader("key: value")},
			},
			name:    "Render",
			wantErr: false,
			result:  bytes.NewBufferString("value\n"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			render := renderer.New(&renderer.Options{
				TemplateReaders: test.templateReaders,
				ValueReaders:    test.valueReaders,
			})
			out, err := render.Render()
			if (err != nil) != test.wantErr {
				t.Errorf("renderer.Render() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			assert.Equal(tt, test.result.String(), out.String())
		})

	}
}
