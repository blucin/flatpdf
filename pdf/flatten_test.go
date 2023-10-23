package pdf

import (
	"testing"
)

func TestFlattenReturnsPath(t *testing.T) {
	path := "~/testFile.pdf"
	got, err := Flatten(path)

	if err != nil {
		t.Error(err)
		t.Errorf("Flatten('~/testFile.pdf') = %v; want path to flattened pdf", got)
	}
}
