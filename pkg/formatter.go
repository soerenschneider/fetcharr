package pkg

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/soerenschneider/fetcharr/internal/syncer"
)

func BytesToHumanSize(bytes int64) string {
	const (
		kb = 1024
		mb = kb * 1024
		gb = mb * 1024
		tb = gb * 1024
		pb = tb * 1024
		eb = pb * 1024
	)

	switch {
	case bytes >= eb:
		return fmt.Sprintf("%.2f EB", float64(bytes)/float64(eb))
	case bytes >= pb:
		return fmt.Sprintf("%.2f PB", float64(bytes)/float64(pb))
	case bytes >= tb:
		return fmt.Sprintf("%.2f TB", float64(bytes)/float64(tb))
	case bytes >= gb:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(gb))
	case bytes >= mb:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(mb))
	case bytes >= kb:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(kb))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

func Format(templateData string, stats syncer.Stats) (string, error) {
	tmpl, err := template.New("webhook").Parse(templateData)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	if err := tmpl.Execute(&b, stats); err != nil {
		return "", err
	}

	return b.String(), nil
}
