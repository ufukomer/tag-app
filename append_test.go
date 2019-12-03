package main

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	suffolk = `As by your high imperial majesty
I had in charge at my depart for France,
As procurator to your excellence,
To marry Princess Margaret for your grace,
So, in the famous ancient city, Tours,
In presence of the Kings of France and Sicil,
The Dukes of Orleans, Calaber, Bretagne and Alencon,
Seven earls, twelve barons and twenty reverend bishops,
I have perform'd my task and was espoused:
And humbly now upon my bended knee,
In sight of England and her lordly peers,
Deliver up my title in the queen
To your most gracious hands, that are the substance
Of that great shadow I did represent;
The happiest gift that ever marquess gave,
The fairest queen that ever king received.`

	king = `// +build king !majesty
Suffolk, arise. Welcome, Queen Margaret:
I can express no kinder sign of love
Than this kind kiss. O Lord, that lends me life,
Lend me a heart replete with thankfulness!
For thou hast given me in this beauteous face
A world of earthly blessings to my soul,
If sympathy of love unite our thoughts.`

	queen = ``

	all = `[Kneeling] Long live Queen Margaret, England's
happiness!`
)

func TestAppendTags(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "tag-app-")
	if err != nil {
		return
	}
	defer func(dir string) {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatal(err)
		}
	}(tempDir)

	// Paths of files that are going to be used for testing.
	suffolkPath := filepath.Join(tempDir, "suffolk_test.go")
	kingPath := filepath.Join(tempDir, "king_test.go")
	queenPath := filepath.Join(tempDir, "queen_test.go")
	allPath := filepath.Join(tempDir, "all.go")

	// Initialize files with unrelated content.
	createAndCopy(t, suffolkPath, suffolk)
	createAndCopy(t, kingPath, king)
	createAndCopy(t, queenPath, queen)
	createAndCopy(t, allPath, all)

	// Add build tags to files.
	err = appendTags(tempDir, "all suffolk !king !queen", "_test.go")
	if err != nil {
		t.Fatal(err)
	}

	lines := readContent(t, suffolkPath)
	base := filepath.Base(suffolkPath)
	if length := len(lines); length != 17 {
		t.Errorf("Lenght of file (%v) is incorrect got: %v want: %v", base, length, 17)
	}
	if firstLine := lines[0]; firstLine != "// +build all suffolk !king !queen" {
		t.Errorf("First line of file (%v) is incorrect got: %v want: %v", base, firstLine, "// +build all suffolk !king !queen")
	}
	if secondLine := lines[1]; secondLine != "As by your high imperial majesty" {
		t.Errorf("Second line of file (%v) is incorrect got: %v want: %v", base, secondLine, "As by your high imperial majesty")
	}
	if lastLine := lines[len(lines)-1]; lastLine != "The fairest queen that ever king received." {
		t.Errorf("Last line of file (%v) is incorrect got: %v want: %v", base, lastLine, "The fairest queen that ever king received.")
	}

	lines = readContent(t, kingPath)
	base = filepath.Base(kingPath)
	if length := len(lines); length != 8 {
		t.Errorf("Lenght of file (%v) is incorrect got: %v want: %v", base, length, 8)
	}
	if firstLine := lines[0]; firstLine != "// +build all suffolk !king !queen" {
		t.Errorf("First line of file (%v) is incorrect got: %v want: %v", base, firstLine, "// +build all suffolk !king !queen")
	}
	if secondLine := lines[1]; secondLine != "Suffolk, arise. Welcome, Queen Margaret:" {
		t.Errorf("Second line of file (%v) is incorrect got: %v want: %v", base, secondLine, "Suffolk, arise. Welcome, Queen Margaret:")
	}
	if lastLine := lines[len(lines)-1]; lastLine != "If sympathy of love unite our thoughts." {
		t.Errorf("Last line of file (%v) is incorrect got: %v want: %v", base, lastLine, "If sympathy of love unite our thoughts.")
	}

	lines = readContent(t, queenPath)
	base = filepath.Base(queenPath)
	if length := len(lines); length != 2 {
		t.Errorf("Lenght of file (%v) is incorrect got: %v want: %v", base, length, 1)
	}
	if firstLine := lines[0]; firstLine != "// +build all suffolk !king !queen" {
		t.Errorf("First line of file (%v) is incorrect got: %v want: %v", base, firstLine, "// +build all suffolk !king !queen")
	}
	if lastLine := lines[len(lines)-1]; lastLine != "" {
		t.Errorf("Last line of file (%v) is incorrect got: %v want: %v", base, lastLine, "")
	}

	lines = readContent(t, allPath)
	base = filepath.Base(allPath)
	if length := len(lines); length != 2 {
		t.Errorf("Lenght of file (%v) is incorrect got: %v want: %v", base, length, 1)
	}
	if firstLine := lines[0]; firstLine != "[Kneeling] Long live Queen Margaret, England's" {
		t.Errorf("First line of file (%v) is incorrect got: %v want: %v", base, firstLine, "[Kneeling] Long live Queen Margaret, England's")
	}
	if lastLine := lines[len(lines)-1]; lastLine != "happiness!" {
		t.Errorf("Last line of file (%v) is incorrect got: %v want: %v", base, lastLine, "happiness!")
	}
}

func createAndCopy(t *testing.T, path, content string) {
	t.Helper()

	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.Copy(file, strings.NewReader(content))
	if err != nil {
		t.Fatal(err)
	}
	if err = file.Close(); err != nil {
		t.Fatal(err)
	}
}

func readContent(t *testing.T, path string) []string {
	t.Helper()

	file, err := os.OpenFile(path, os.O_RDONLY, 0600)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(string(b), "\n")
	if !(len(lines) > 0) {
		t.Fatalf("File (%v) content is empty", filepath.Base(path))
	}

	return lines
}
