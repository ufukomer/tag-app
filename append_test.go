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

	gloucester = `// +build king !majesty

Pardon me, gracious lord;
Some sudden qualm hath struck me at the heart
And dimm'd mine eyes, that I can read no further.`
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
	gloucesterPath := filepath.Join(tempDir, "gloucester_test.go")

	// Initialize files with unrelated content.
	createAndCopy(t, suffolkPath, suffolk)
	createAndCopy(t, kingPath, king)
	createAndCopy(t, queenPath, queen)
	createAndCopy(t, allPath, all)
	createAndCopy(t, gloucesterPath, gloucester)

	// Add build tags to files.
	err = appendTags(tempDir, "all suffolk !king !queen", "_test.go")
	if err != nil {
		t.Fatal(err)
	}

	lines := readContent(t, suffolkPath)
	base := filepath.Base(suffolkPath)
	if length := len(lines); length != 18 {
		t.Errorf("Lenght of file (%v) is incorrect got: %v want: %v", base, length, 18)
	}
	if firstLine := lines[0]; firstLine != "// +build all suffolk !king !queen" {
		t.Errorf("1st line of file (%v) is incorrect got: %v want: %v", base, firstLine, "// +build all suffolk !king !queen")
	}
	if secondLine := lines[1]; secondLine != "" {
		t.Errorf("2nd line of file (%v) is incorrect got: %v want: %v", base, secondLine, "")
	}
	if thirdLine := lines[2]; thirdLine != "As by your high imperial majesty" {
		t.Errorf("3rd line of file (%v) is incorrect got: %v want: %v", base, thirdLine, "As by your high imperial majesty")
	}
	if fourthLine := lines[3]; fourthLine != "I had in charge at my depart for France," {
		t.Errorf("4th line of file (%v) is incorrect got: %v want: %v", base, fourthLine, "I had in charge at my depart for France,")
	}
	if fifthLine := lines[4]; fifthLine != "As procurator to your excellence," {
		t.Errorf("5th line of file (%v) is incorrect got: %v want: %v", base, fifthLine, "As procurator to your excellence,")
	}
	if sixthLine := lines[5]; sixthLine != "To marry Princess Margaret for your grace," {
		t.Errorf("6th line of file (%v) is incorrect got: %v want: %v", base, sixthLine, "To marry Princess Margaret for your grace,")
	}
	if seventhLine := lines[6]; seventhLine != "So, in the famous ancient city, Tours," {
		t.Errorf("7th line of file (%v) is incorrect got: %v want: %v", base, seventhLine, "To marry Princess Margaret for your grace,")
	}
	if eighthLine := lines[7]; eighthLine != "In presence of the Kings of France and Sicil," {
		t.Errorf("8th line of file (%v) is incorrect got: %v want: %v", base, eighthLine, "So, in the famous ancient city, Tours,")
	}
	if ninthLine := lines[8]; ninthLine != "The Dukes of Orleans, Calaber, Bretagne and Alencon," {
		t.Errorf("9th line of file (%v) is incorrect got: %v want: %v", base, ninthLine, "The Dukes of Orleans, Calaber, Bretagne and Alencon,")
	}
	if tenthLine := lines[9]; tenthLine != "Seven earls, twelve barons and twenty reverend bishops," {
		t.Errorf("10th line of file (%v) is incorrect got: %v want: %v", base, tenthLine, "Seven earls, twelve barons and twenty reverend bishops,")
	}
	if eleventhLine := lines[10]; eleventhLine != "I have perform'd my task and was espoused:" {
		t.Errorf("11th line of file (%v) is incorrect got: %v want: %v", base, eleventhLine, "I have perform'd my task and was espoused:")
	}
	if twelfthLine := lines[11]; twelfthLine != "And humbly now upon my bended knee," {
		t.Errorf("12th line of file (%v) is incorrect got: %v want: %v", base, twelfthLine, "And humbly now upon my bended knee,")
	}
	if thirteenthLine := lines[12]; thirteenthLine != "In sight of England and her lordly peers," {
		t.Errorf("13th line of file (%v) is incorrect got: %v want: %v", base, thirteenthLine, "In sight of England and her lordly peers,")
	}
	if fourteenthLine := lines[13]; fourteenthLine != "Deliver up my title in the queen" {
		t.Errorf("14th line of file (%v) is incorrect got: %v want: %v", base, fourteenthLine, "Deliver up my title in the queen")
	}
	if fifteenthLine := lines[14]; fifteenthLine != "To your most gracious hands, that are the substance" {
		t.Errorf("15th line of file (%v) is incorrect got: %v want: %v", base, fifteenthLine, "To your most gracious hands, that are the substance")
	}
	if sixteenthLine := lines[15]; sixteenthLine != "Of that great shadow I did represent;" {
		t.Errorf("16th line of file (%v) is incorrect got: %v want: %v", base, sixteenthLine, "Of that great shadow I did represent;")
	}
	if seventeenthLine := lines[16]; seventeenthLine != "The happiest gift that ever marquess gave," {
		t.Errorf("17th line of file (%v) is incorrect got: %v want: %v", base, seventeenthLine, "The happiest gift that ever marquess gave,")
	}
	if lastLine := lines[len(lines)-1]; lastLine != "The fairest queen that ever king received." {
		t.Errorf("Last line of file (%v) is incorrect got: %v want: %v", base, lastLine, "The fairest queen that ever king received.")
	}
	lines = readContent(t, kingPath)
	base = filepath.Base(kingPath)
	if length := len(lines); length != 9 {
		t.Errorf("Lenght of file (%v) is incorrect got: %v want: %v", base, length, 9)
	}
	if firstLine := lines[0]; firstLine != "// +build all suffolk !king !queen" {
		t.Errorf("First line of file (%v) is incorrect got: %v want: %v", base, firstLine, "// +build all suffolk !king !queen")
	}
	if secondLine := lines[1]; secondLine != "" {
		t.Errorf("Second line of file (%v) is incorrect got: %v want: %v", base, secondLine, "")
	}
	if thirdLine := lines[2]; thirdLine != "Suffolk, arise. Welcome, Queen Margaret:" {
		t.Errorf("First line of file (%v) is incorrect got: %v want: %v", base, thirdLine, "Suffolk, arise. Welcome, Queen Margaret:")
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

	lines = readContent(t, gloucesterPath)
	base = filepath.Base(gloucesterPath)
	if length := len(lines); length != 5 {
		t.Errorf("Lenght of file (%v) is incorrect got: %v want: %v", base, length, 5)
	}
	if firstLine := lines[0]; firstLine != "// +build all suffolk !king !queen" {
		t.Errorf("1st line of file (%v) is incorrect got: %v want: %v", base, firstLine, "// +build all suffolk !king !queen")
	}
	if secondLine := lines[1]; secondLine != "" {
		t.Errorf("2nd line of file (%v) is incorrect got: %v want: %v", base, secondLine, "")
	}
	if thirdLine := lines[2]; thirdLine != "Pardon me, gracious lord;" {
		t.Errorf("3rd line of file (%v) is incorrect got: %v want: %v", base, thirdLine, "Pardon me, gracious lord;")
	}
	if lastLine := lines[len(lines)-1]; lastLine != "And dimm'd mine eyes, that I can read no further." {
		t.Errorf("Last line of file (%v) is incorrect got: %v want: %v", base, lastLine, "And dimm'd mine eyes, that I can read no further.")
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
