package protlook

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

type AllType struct {
	Yn        bool
	Num       int
	Str       string
	NumArr    [3]int8
	StrArr    [3]string
	byteSlice []byte
	InfoMap   map[string]string
}

type Tx1 struct {
	AllType
	Level string
	X1    Tx1x1
	X2    *Tx1x1
}

type Tx1x1 struct {
	AllType
	Level string
}

func (self Tx1) HiTx1_AA() string {
	return "Hi, Tx1."
}

func (self Tx1) hiTx1_AB() string {
	return "Hi, Tx1."
}

func (self *Tx1) HiTx1_BA() string {
	return "Hi, Tx1."
}

func (self *Tx1) hiTx1_BB() string {
	return "Hi, Tx1x1."
}

func (self Tx1x1) HiTx1x1_AA() string {
	return "Hi, Tx1x1."
}

func (self Tx1x1) hiTx1x1_AB() string {
	return "Hi, Tx1x1."
}

func (self *Tx1x1) HiTx1x1_BA() string {
	return "Hi, Tx1x1."
}

func (self *Tx1x1) hiTx1x1_BB() string {
	return "Hi, Tx1x1."
}

// https://medium.com/@hau12a1/770209c791b4
func captureOutput(f func()) string {
	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	stdout := os.Stdout
	stderr := os.Stderr
	os.Stdout = writer
	os.Stderr = writer
	log.SetOutput(writer)

	out := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, reader)
		out <- buf.String()

		reader.Close()
		os.Stdout = stdout
		os.Stderr = stderr
		log.SetOutput(os.Stderr)
	}()
	f()
	writer.Close()
	return <-out
}

var TestSampleOutputAns string = `~~~: (struct(protlook.Tx1))
  AllType: (struct(protlook.AllType))
    Yn: (bool) false
    Num: (int(int)) 0
    Str: (string) ""
    NumArr: (array([3]int8))
      [0 0 0]
    StrArr: (array([3]string))
      [  ]
    byteSlice: (inner) (slice([]uint8)) []
    InfoMap: (map(map[string]string))
      map[]
  Level: (string) "Tx1"
  X1: (struct(protlook.Tx1x1))
    AllType: (struct(protlook.AllType))
      Yn: (bool) true
      Num: (int(int)) 3
      Str: (string) "Tx1"
      NumArr: (array([3]int8))
        [1 2 3]
      StrArr: (array([3]string))
        [C e r]
      byteSlice: (inner) (slice([]uint8)) []
      InfoMap: (map(map[string]string))
        map[InfoMap:map[string]string{} Num:Num NumArr:Array Str:Str StrArr:Array Yn:Yn byteSlice:Slice]
    Level: (string) "Tx1x1"
    HiTx1x1_AA: (func() string)
  X2: (ptr(*protlook.Tx1x1))
    HiTx1x1_AA: (func() string)
    HiTx1x1_BA: (func() string)
  HiTx1_AA: (func() string)
~~~: (ptr(*protlook.Tx1))
  HiTx1_AA: (func() string)
  HiTx1_BA: (func() string)
`

func TestSample(t *testing.T) {
	insTarget := Tx1{
		Level: "Tx1",
		X1: Tx1x1{
			Level: "Tx1x1",
			AllType: AllType{
				Yn:     true,
				Num:    3,
				Str:    "Tx1",
				NumArr: [3]int8{1, 2, 3},
				StrArr: [3]string{"C", "e", "r"},
				// byteSlice: []byte("ABc"),
				InfoMap: map[string]string{
					"Yn":        "Yn",
					"Num":       "Num",
					"Str":       "Str",
					"NumArr":    "Array",
					"StrArr":    "Array",
					"byteSlice": "Slice",
					"InfoMap":   "map[string]string{}",
				},
			},
		},
	}
	insTarget.X2 = &insTarget.X1
	outputTxt := captureOutput(func() {
		Print(insTarget, &insTarget)
	})
	if TestSampleOutputAns == outputTxt {
		fmt.Println("# Print(insTarget, &insTarget)")
		fmt.Println(outputTxt)
	} else {
		assert.Equal(t, TestSampleOutputAns, outputTxt)
	}
}
