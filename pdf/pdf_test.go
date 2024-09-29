package pdf

import (
	"os"
	"testing"
	"time"

	"github.com/klippa-app/go-pdfium/webassembly"
)

func TestPdfToImages(t *testing.T) {
    // pdfium init
	var err error
	poolConfig := webassembly.Config{
		MinIdle:  1,
		MaxIdle:  1,
		MaxTotal: 1,
	}
	pool, err = webassembly.Init(poolConfig)
	if err != nil {
		t.Error(err)
	}

	instance, err = pool.GetInstance(time.Second * 30)
	if err != nil {
		t.Error(err)
	}
	defer instance.Close()

    // opts is autofilled by cobra
    opts := FlattenPDFOptions{
        InputFiles: []string{"../assets/main.pdf"},
        OutputDir: "../assets",
        RequiresFlatSuffix: true,
        ImageDensity: 600,
        ImageQuality: 99,
    }
    tempDir := os.TempDir()

    _, err = PdfToImages("../assets/main.pdf", tempDir, opts)
    if err != nil {
        t.Error(err)
    }

    err = os.Remove(tempDir)
    if err != nil {
        t.Error(err)
    }
}

