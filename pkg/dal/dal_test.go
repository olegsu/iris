package dal

import (
	"testing"
)

func TestNewDalFromFilePathWhenPathNoSet(t *testing.T) {
	d := NewDalFromFilePath("")
	if len(d.Destinations) > 0 || len(d.Filters) > 0 || len(d.Integrations) > 0 {
		t.Errorf("Should be empty")
	}
}
