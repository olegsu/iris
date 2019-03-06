package filter

import (
	"fmt"
	"testing"

	"github.com/olegsu/iris/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestNewFactory(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Get filter factory",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message := fmt.Sprintf("Failed to get factory")
			factory := NewFactory(logger.New(nil))
			assert.IsType(t, &f{}, factory, message)
		})
	}
}
