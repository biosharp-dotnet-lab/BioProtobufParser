package generators

import (
	"bytes"
	"encoding/base64"
	"strings"

	bioproto "github.com/ag-computational-bio/BioProtobufSchemas/go"
)

//GenerateGBfromproto Genbank protobuf to genbank file
//Generates a genbank file from a given protobuf genbank
func GenerateGBfromproto(record *bioproto.Genbank) string {
	var stringbuffer bytes.Buffer

	stringbuffer.WriteString(generateHeaderString(record))
	stringbuffer.WriteString("FEATURES             Location/Qualifiers\n")
	stringbuffer.WriteString(generateQualifierString(record))
	if record.FEATURES != nil {

	}
	if record.CONTIG != "" {
		stringbuffer.WriteString("CONTIG      " + record.CONTIG + "\n")
	}
	stringbuffer.WriteString("//\n")
	return stringbuffer.String()
}

func generateHeaderString(record *bioproto.Genbank) (HeadString string) {
	var buffer bytes.Buffer
	buffer.WriteString("LOCUS       " + record.LOCUS + "\n")
	buffer.WriteString(formatStringWithNewlineChars("DEFINITION  "+record.DEFINITION, "            ", true))
	if len(record.ACCESSION) > 1 {
		buffer.WriteString(formatStringWithNewlineChars("ACCESSION   "+strings.Join(record.ACCESSION[:], " "), "            ", true))
	} else {
		buffer.WriteString("ACCESSION   " + record.ACCESSION[0] + "\n")
	}
	buffer.WriteString("VERSION     " + record.VERSION + "\n")
	if len(record.DBLINK) > 0 {
		for i, line := range record.DBLINK {
			if i == 0 {
				buffer.WriteString("DBLINK      " + line + "\n")
			} else {
				buffer.WriteString("            " + line + "\n")
			}
		}
	}
	buffer.WriteString("KEYWORDS    " + record.KEYWORDS + "\n")
	buffer.WriteString(formatStringWithNewlineChars("SOURCE      "+record.SOURCE, "            ", true))
	for i, line := range record.ORGANISM {
		if i == 0 {
			buffer.WriteString("  ORGANISM  " + line + "\n")
		} else {
			buffer.WriteString("            " + line + "\n")
		}
	}

	for _, ref := range record.REFERENCES {
		if ref.ORIGIN != "" {
			buffer.WriteString("REFERENCE   " + ref.ORIGIN + "\n")
			buffer.WriteString(formatStringWithNewlineChars("  AUTHORS   "+ref.AUTHORS, "            ", true))
			if ref.CONSRTM != "" {
				buffer.WriteString(formatStringWithNewlineChars("  CONSRTM   "+ref.CONSRTM, "            ", true))
			}
			buffer.WriteString(formatStringWithNewlineChars("  TITLE     "+ref.TITLE, "            ", true))
			buffer.WriteString(formatStringWithNewlineChars("  JOURNAL   "+ref.JOURNAL, "            ", true))
			if ref.PUBMED != "" {
				buffer.WriteString("   PUBMED   " + ref.PUBMED + "\n")
			}
		}
	}
	b64Decode, _ := base64.RawStdEncoding.DecodeString(record.COMMENT)
	buffer.WriteString("COMMENT     " + addSpacesSpecialHeader(string(b64Decode)) + "\n")

	return buffer.String()
}

func generateQualifierString(record *bioproto.Genbank) (returnstring string) {
	var buffer bytes.Buffer
	spacestring := "                "
	for _, feature := range record.FEATURES {
		if feature.IsCompliment {
			if feature.IsJoined {
				buffer.WriteString("     " + feature.TYPE + spacestring[len(feature.TYPE):] + "complement(join(" + feature.START + "," + feature.STOP + ")\n")
			} else {
				buffer.WriteString("     " + feature.TYPE + spacestring[len(feature.TYPE):] + "complement(" + feature.START + ".." + feature.STOP + ")\n")
			}

		} else {
			if feature.IsJoined {
				buffer.WriteString("     " + feature.TYPE + spacestring[len(feature.TYPE):] + "join(" + feature.START + "," + feature.STOP + ")\n")
			} else {
				buffer.WriteString("     " + feature.TYPE + spacestring[len(feature.TYPE):] + feature.START + ".." + feature.STOP + "\n")
			}
		}
		for _, qualifier := range feature.QUALIFIERS {
			if qualifier.Key == "/pseudo" {
				buffer.WriteString("                     /pseudo\n")
			} else if qualifier.Key != "" {
				buffer.WriteString(formatStringWithNewlineChars(qualifier.Key+"="+qualifier.Value, "                     ", false))
			}
		}

	}
	return buffer.String()
}

func generateSequenceString(record *bioproto.Genbank) {

}

func addSpacesSpecialHeader(inputString string) (Output string) {
	var returnbuffer bytes.Buffer
	for _, char := range inputString {
		returnbuffer.WriteRune(char)
		if char == '\n' {
			returnbuffer.WriteString("            ")
		}
	}
	return returnbuffer.String()
}

func formatStringWithNewlineChars(Splittedstring string, newlineinsertion string, hasKeyword bool) (result string) {
	var buffer bytes.Buffer
	// DEBUG: fmt.Println(Splittedstring)
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
			if currentlength >= 79-len(newlineinsertion) {
				buffer.WriteString(newlineinsertion + Splittedstring[lastsplitindex:lastspaceindex] + "\n")
				lastsplitindex = lastspaceindex + 1
				currentlength = i - lastspaceindex - 1
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
	if len(buffer.String()) > 0 {
		return keyword + buffer.String()[len(keyword):]
	}

	return keyword + "\n"
}
