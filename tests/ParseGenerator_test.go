package tests

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/ag-computational-bio/BioProtobufParser/gbparse"
	"github.com/ag-computational-bio/BioProtobufParser/generators"
	bioproto "github.com/ag-computational-bio/BioProtobufSchemas/go"
)

func TestGBFFParserAndGenerator(t *testing.T) {

	// Add Waitgroup
	var wg sync.WaitGroup
	// Open testfile
	gbff, err := os.Open("../testfiles/archaea.1.genomic.gbff")
	// Read file as bytebuffer for comparison
	filecontent, err := ioutil.ReadFile("../testfiles/archaea.1.genomic.gbff")

	if err != nil {
		log.Fatal(err)
	}

	output := make(chan *bioproto.Genbank, 1000000)

	log.Println("Parsing gbff file...")
	defer gbff.Close()

	parser := gbparse.GBParser{}

	wg.Add(1)
	go func() {
		parser.ReadAndParseFile(gbff, output)
		wg.Done()
	}()
	wg.Wait()

	log.Println("Parsing complete, reading protobuf...")
	result := ""
	// Close Output channel before reading
	close(output)
	for record := range output {
		result += generators.GenerateGBfromproto(record)
	}

	// DEBUG: Writefile
	//err = ioutil.WriteFile("demoresult.txt", []byte(result), 777)
	//if err != nil{
	//	fmt.Println(err.Error())
	//}

	// compare resultstring from protobuf object against raw string from file
	origlines := strings.Split(string(filecontent), "\n")
	genlines := strings.Split(result, "\n")

	for index, lin := range genlines {
		if lin != origlines[index] {
			t.Errorf("Parsed and generated file not equal in line: %v", index)
			t.Errorf(lin)
			t.Errorf(origlines[index])
			break
		}
	}
}

func TestFASTAParserAndGenerator(t *testing.T) {

	// Add Waitgroup
	var wg sync.WaitGroup

	// Initialize Parser

	// Open testfile
	fasta, err := os.Open("../testfiles/Test50Entries.fasta")
	// Read file as bytebuffer for comparison
	filecontent, err := ioutil.ReadFile("../testfiles/Test50Entries.fasta")

	if err != nil {
		log.Fatal(err)
	}

	parser := gbparse.FASTAParser{}

	log.Println("Parsing fasta file...")
	defer fasta.Close()

	output := make(chan *bioproto.Fasta, 1000000)

	wg.Add(1)
	go func() {
		parser.ReadAndParseFile(fasta, output)
		wg.Done()
	}()
	wg.Wait()

	log.Println("Parsing complete, reading protobuf...")
	result := ""

	// Close Output channel before reading
	close(output)
	for record := range output {
		result += generators.GenerateFastafromproto(record)
	}

	// compare resultstring from protobuf object against raw string from file
	if result != string(filecontent) {
		t.Errorf("Parsed and generated file not equal!")
	}

}
