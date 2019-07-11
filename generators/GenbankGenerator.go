package generators

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gbparsertest2/gbparse"
	"io/ioutil"
	"os"
	"strings"
)

func GenerateGBfromproto(record *gbparse.Genbank) (fastarecord string) {
	var stringbuffer bytes.Buffer

	stringbuffer.WriteString(generateHeaderString(record))
	stringbuffer.WriteString("FEATURES             Location/Qualifiers\n")
	stringbuffer.WriteString(generateQualifierString(record, read_json()))
	if record.FEATURES != nil {

	}
	if record.CONTIG != "" {
		stringbuffer.WriteString("CONTIG      " + record.CONTIG + "\n")
	}
	stringbuffer.WriteString("//\n")
	return stringbuffer.String()
}

func generateHeaderString(record *gbparse.Genbank) (HeadString string) {
	var buffer bytes.Buffer
	buffer.WriteString("LOCUS       " + record.LOCUS + "\n")
	buffer.WriteString(formatStringWithNewlineChars("DEFINITION  "+record.DEFINITION, "            ", true))
	buffer.WriteString("ACCESSION   " + record.ACCESSION + "\n")
	buffer.WriteString("VERSION     " + record.ACCESSION + "\n")
	buffer.WriteString("DBLINK      " + addSpacesSpecialHeader(record.DBLINK) + "\n")
	buffer.WriteString("KEYWORDS    " + record.KEYWORDS + "\n")
	buffer.WriteString("SOURCE      " + record.SOURCE + "\n")
	buffer.WriteString("  ORGANISM  " + addSpacesSpecialHeader(record.ORGANISM) + "\n")

	for _, ref := range record.REFERENCES {
		buffer.WriteString("REFERENCE   " + ref.ORIGIN + "\n")
		buffer.WriteString(formatStringWithNewlineChars("  AUTHORS   "+ref.AUTHORS, "            ", true))
		buffer.WriteString(formatStringWithNewlineChars("  TITLE     "+ref.TITLE, "            ", true))
		buffer.WriteString("  JOURNAL   " + ref.JOURNAL + "\n")
		buffer.WriteString("   PUBMED   " + ref.PUBMED + "\n")
	}
	buffer.WriteString("COMMENT     " + addSpacesSpecialHeader(record.COMMENT) + "\n")

	return buffer.String()
}
func read_json() (result map[string][]string) {
	jsonFile, err := os.Open("generators/categorys_by_occurence.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue), &result)
	return result
}

func generateQualifierString(record *gbparse.Genbank, jsonmap map[string][]string) (QualString string) {
	var buffer bytes.Buffer
	spacestring := "                "
	for _, feature := range record.FEATURES {
		if feature.IsCompliment {
			buffer.WriteString("     " + feature.TYPE + spacestring[len(feature.TYPE):] + "complement(" + feature.START + ".." + feature.STOP + ")\n")
		} else {
			buffer.WriteString("     " + feature.TYPE + spacestring[len(feature.TYPE):] + feature.START + ".." + feature.STOP + "\n")
		}
		for _, occurence := range jsonmap[feature.TYPE] {
			if val, inMap := feature.QUALIFIERS[occurence]; inMap {
				buffer.WriteString(formatStringWithNewlineChars(occurence+"="+val, "                     ", false))
			}
		}

	}
	return buffer.String()
}

func generateSequenceString(record *gbparse.Genbank) {

}

func addSpacesSpecialHeader(inputString string) (Output string) {
	var returnbuffer bytes.Buffer
	for _, char := range inputString {
		if char == '\n' {
			returnbuffer.WriteString("            ")
		}
		returnbuffer.WriteRune(char)
	}
	return returnbuffer.String()
}

func formatStringWithNewlineChars(Splittedstring string, newlineinsertion string, hasKeyword bool) (result string) {
	var buffer bytes.Buffer
	keyword := ""
	if hasKeyword {
		keyword = Splittedstring[:len(newlineinsertion)]
		Splittedstring = Splittedstring[len(newlineinsertion):]
	}
	lastsplitindex := 0
	currentlength := 0
	if strings.ContainsRune(Splittedstring, ' ') {
		currentlength := 0
		lastspaceindex := 0
		for i, char := range Splittedstring {
			if char == ' ' {
				lastspaceindex = i
			}
			if currentlength >= 80-len(newlineinsertion)-1 {
				buffer.WriteString(newlineinsertion + Splittedstring[lastsplitindex:lastspaceindex] + "\n")
				lastsplitindex = lastspaceindex + 1
				currentlength = 0
			}
			if i == len(Splittedstring)-1 {
				buffer.WriteString(newlineinsertion + Splittedstring[lastsplitindex:] + "\n")
			}
			currentlength++
		}
	} else {
		for i := range Splittedstring {
			if currentlength >= 80-len(newlineinsertion)-1 {
				buffer.WriteString(newlineinsertion + Splittedstring[lastsplitindex:i] + "\n")
				lastsplitindex = i
				currentlength = 0
			}
			if i == len(Splittedstring)-1 {
				buffer.WriteString(newlineinsertion + Splittedstring[lastsplitindex:] + "\n")
			}
			currentlength++

		}
	}
	return keyword + buffer.String()[len(keyword):]
}