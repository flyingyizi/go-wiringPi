package rpi

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

//ref http://xillybus.com/tutorials/device-tree-zynq-1
//ref  https://github.com/stianeikeland/go-rpio/blob/master/rpio.go

const (
	PeripheralBaseUnknown int64 = 0
	PeripheralBase2835    int64 = 0x20000000
	PeripheralBase2836    int64 = 0x3f000000
	PeripheralBase2837    int64 = 0x3f000000
)

// ModelT :  Raspberry Pi Revision :: Model
type ModelT int

const (
	ModelA       ModelT = 0  //  RPIModelA "Model A",	//  0
	ModelB       ModelT = 1  //   "Model B",	//  1
	ModelAPlus   ModelT = 2  //   "Model A+",	//  2
	ModelBPlus   ModelT = 3  //   "Model B+",	//  3
	Model2B      ModelT = 4  //   "Pi 2",	//  4
	ModelAlpha   ModelT = 5  //   "Alpha",	//  5
	ModelCM      ModelT = 6  //   "CM",		//  6
	ModelUnknown ModelT = 7  //   "Unknown07",	// 07
	Model3B      ModelT = 8  //   "Pi 3",	// 08
	ModelZero    ModelT = 9  //   "Pi Zero",	// 09
	ModelCM3     ModelT = 10 //   "CM3",	// 10
	ModelZeroW   ModelT = 12 //   "Pi Zero-W",	// 12
)

// MemoryT :  Raspberry Pi Revision :: memeory type
type MemoryT int

// the value is from PRI post PI2 revision
const (
	RpiUnknownMB MemoryT = -1
	Rpi256MB     MemoryT = 0
	Rpi512MB     MemoryT = 1
	Rpi1024MB    MemoryT = 2
)

type ProcessorT int

const (
	BroadcomUnknown ProcessorT = -1
	Broadcom2835    ProcessorT = 0
	Broadcom2836    ProcessorT = 1
	Broadcom2837    ProcessorT = 2
)

type MakerT int

const (
	MakerUnknown   MakerT = -1
	MakerSony      MakerT = 0
	MakerEgoman    MakerT = 1
	MakerEmbest    MakerT = 2
	MakerSonyJapan MakerT = 3
	MakerEmbest1   MakerT = 4
)

type RpiInfoT struct {
	model        ModelT
	mem          MemoryT
	processor    ProcessorT
	manufacturer MakerT
	pcbRev       PcbRevT
	overVolted   bool //
	revision     uint64
}

var ModelName = map[ModelT]string{
	ModelA:       "Model A",   //  0
	ModelB:       "Model B",   //  1
	ModelAPlus:   "Model A+",  //  2
	ModelBPlus:   "Model B+",  //  3
	Model2B:      "Pi 2",      //
	ModelAlpha:   "Alpha",     //  5
	ModelCM:      "CM",        //  6
	ModelUnknown: "Unknown",   //
	Model3B:      "Pi 3",      // 08
	ModelZero:    "Pi Zero",   // 09
	ModelCM3:     "CM3",       // 10
	ModelZeroW:   "Pi Zero-W", // 12
}

var MakerName = map[MakerT]string{
	MakerUnknown:   "maker unknown",
	MakerSony:      "maker SONY",
	MakerEgoman:    "maker EGOMAN",
	MakerEmbest:    "maker EMBEST",
	MakerSonyJapan: "maker SONY",
	MakerEmbest1:   "maker EMBEST",
}

type PcbRevT int

const (
	PcbRevUnknown PcbRevT = -1
	PcbRev1       PcbRevT = 0
	PcbRev1_1     PcbRevT = 1
	PcbRev1_2     PcbRevT = 2
	PcbRev2       PcbRevT = 3
)

func getRevision() (revision string, err error) {
	cpuinfo, err := ioutil.ReadFile("/proc/cpuinfo")
	if err != nil {
		return
	}

	lines := strings.Split(string(cpuinfo), "\n")
	for _, l := range lines {
		fields := strings.Split(l, ":")
		if len(fields) == 2 {
			k := strings.TrimSpace(fields[0])
			v := strings.TrimSpace(fields[1])
			if k == "Revision" {
				revision = v
				break
			}
		}
	}
	return
}

//-------------------------------------------------------------------------
// refer : https://github.com/AndrewFromMelbourne/raspberry_pi_revision  and wiringPI archive
//-------------------------------------------------------------------------
//
// The file /proc/cpuinfo contains a line such as:-
//
// Revision    : 0003
//
// that holds the revision number of the Raspberry Pi.
// Known revisions (prior to the Raspberry Pi 2) are:
//
//     +----------+---------+---------+--------+-------------+
//     | Revision |  Model  | PCB Rev | Memory | Manufacture |
//     +----------+---------+---------+--------+-------------+
//     |   0000   |         |         |        |             |
//     |   0001   |         |         |        |             |
//     |   0002   |    B    |    1    | 256 MB |  EGOMAN     |
//     |   0003   |    B    |    1-1  | 256 MB |   EGOMAN    |
//     |   0004   |    B    |    1-2  | 256 MB |   Sony      |
//     |   0005   |    B    |    1-2  | 256 MB |   EGOMAN    |
//     |   0006   |    B    |    1-2  | 256 MB |   EGOMAN    |
//     |   0007   |    A    |    1-2  | 256 MB |   Egoman    |
//     |   0008   |    A    |    1-2  | 256 MB |   Sony      |
//     |   0009   |    A    |    1-2  | 256 MB |   Egoman    |
//     |   000d   |    B    |    1-2  | 512 MB |   Egoman    |
//     |   000e   |    B    |    1-2  | 512 MB |   Sony      |
//     |   000f   |    B    |    1-2  | 512 MB |   Egoman    |
//     |   0010   |    B+   |    1-2  | 512 MB |   Sony      |
//     |   0011   | compute |    1-1  | 512 MB |   Sony      |
//     |   0012   |    A+   |    1-1  | 256 MB |   Sony      |
//     |   0013   |    B+   |    1-2  | 512 MB |   Embest    |
//     |   0014   | compute |    1-1  | 512 MB |   Embest    |
//     |   0015   |    A+   |    1-1  | 256 MB |   Embest    |
//     |   0016   |    A+   |    1-1  | 256 MB |   Embest    |
//     |   0017   | compute |    1-1  | 512 MB |   sony      |
//     |   0018   | A+      |    1-1  | 256 MB |   sony      |
//     |   0019   |    B+   |    1-2  | 512 MB |   egoman    |
//     |   001a   | compute |    1-1  | 512 MB |   egoman    |
//     |   001b   | A+      |    1-1  | 256 MB |   egoman    |
//     +----------+---------+---------+--------+-------------+
//
// If the Raspberry Pi has been over-volted (voiding the warranty) the
// revision number will have 100 at the front. e.g. 1000002.
func getPreRPI2FromRevision(revision string) (info RpiInfoT, err error) {
	vision, err := strconv.ParseUint(revision, 16, 32) //hex number without 0x lea
	if err != nil {
		return
	}

	encoded := (vision & (1 << 23)) >> 23
	if encoded != 0 {
		err = errors.New("the encoding flag must be 0 in pre PI2")
		return
	}

	info.revision = vision

	warantybit := (vision & (1 << 24)) >> 24
	if warantybit == 1 {
		info.overVolted = true
	}

	if len(revision) > 4 {
		revision = revision[len(revision)-4:] //keep the last 4 byte
	}

	info.processor = Broadcom2835

	//     +----------+---------+---------+--------+-------------+
	//     | Revision |  Model  | PCB Rev | Memory | Manufacture |
	//     +----------+---------+---------+--------+-------------+
	//     |   0000   |         |         |        |             |
	//     |   0001   |         |         |        |             |
	switch revision {
	case "0002":
		//     |   0002   |    B    |    1    | 256 MB |  EGOMAN     |
		info.manufacturer = MakerEgoman
		info.mem = Rpi256MB
		info.model = ModelB
		info.pcbRev = PcbRev1
	case "0003":
		//     |   0003   |    B    |    1-1  | 256 MB |   EGOMAN    |
		info.manufacturer = MakerEgoman
		info.mem = Rpi256MB
		info.model = ModelB
		info.pcbRev = PcbRev1_1
	case "0004":
		//     |   0004   |    B    |    1-2  | 256 MB |   Sony      |
		info.manufacturer = MakerSony
		info.mem = Rpi256MB
		info.model = ModelB
		info.pcbRev = PcbRev1_2
	case "0005", "0006":
		//     |   0005   |    B    |    1-2  | 256 MB |   EGOMAN    |
		//     |   0006   |    B    |    1-2  | 256 MB |   EGOMAN    |
		info.manufacturer = MakerEgoman
		info.mem = Rpi256MB
		info.model = ModelB
		info.pcbRev = PcbRev1_2
	case "0007":
		//     |   0007   |    A    |    1-2  | 256 MB |   Egoman    |
		info.manufacturer = MakerEgoman
		info.mem = Rpi256MB
		info.model = ModelA
		info.pcbRev = PcbRev1_2
	case "0008":
		//     |   0008   |    A    |    1-2  | 256 MB |   Sony      |
		info.manufacturer = MakerSony
		info.mem = Rpi256MB
		info.model = ModelA
		info.pcbRev = PcbRev1_2
	case "0009":
		//     |   0009   |    A    |    1-2  | 256 MB |   Egoman     |
		info.manufacturer = MakerEgoman
		info.mem = Rpi256MB
		info.model = ModelA
		info.pcbRev = PcbRev1_2
	case "000d":
		//     |   000d   |    B    |    1-2  | 512 MB |   Egoman    |
		info.manufacturer = MakerEgoman
		info.mem = Rpi512MB
		info.model = ModelB
		info.pcbRev = PcbRev1_2
	case "000e":
		//     |   000e   |    B    |    1-2  | 512 MB |   Sony      |
		info.manufacturer = MakerSony
		info.mem = Rpi512MB
		info.model = ModelB
		info.pcbRev = PcbRev1_2
	case "000f":
		//     |   000f   |    B    |    1-2  | 512 MB |   Egoman    |
		info.manufacturer = MakerEgoman
		info.mem = Rpi512MB
		info.model = ModelB
		info.pcbRev = PcbRev1_2
	case "0010":
		//     |   0010   |    B+   |    1-2  | 512 MB |   Sony      |
		info.manufacturer = MakerSony
		info.mem = Rpi512MB
		info.model = ModelBPlus
		info.pcbRev = PcbRev1_2
	case "0011":
		//     |   0011   | compute |    1-1  | 512 MB |   Sony      |
		info.manufacturer = MakerSony
		info.mem = Rpi512MB
		info.model = ModelCM
		info.pcbRev = PcbRev1_1
	case "0012":
		//     |   0012   |    A+   |    1-1  | 256 MB |   Sony      |
		info.manufacturer = MakerSony
		info.mem = Rpi256MB
		info.model = ModelAPlus
		info.pcbRev = PcbRev1_1
	case "0013":
		//     |   0013   |    B+   |    1-2  | 512 MB |   Embest    |
		info.manufacturer = MakerEmbest
		info.mem = Rpi512MB
		info.model = ModelBPlus
		info.pcbRev = PcbRev1_2
	case "0014":
		//     |   0014   | compute |    1-1  | 512 MB |   Embest    |
		info.manufacturer = MakerEmbest
		info.mem = Rpi512MB
		info.model = ModelCM
		info.pcbRev = PcbRev1_1
	case "0015":
		//     |   0015   |    A+   |    1-1  | 256 MB |   Embest    |
		info.manufacturer = MakerEmbest
		info.mem = Rpi256MB
		info.model = ModelAPlus
		info.pcbRev = PcbRev1_1
	case "0016":
		//     |   0016   |    A+   |    1-1  | 256 MB |   Embest    |
		info.manufacturer = MakerEmbest
		info.mem = Rpi256MB
		info.model = ModelAPlus
		info.pcbRev = PcbRev1_1
	case "0017":
		//     |   0017   | compute |    1-1  | 512 MB |   sony      |
		info.manufacturer = MakerSony
		info.mem = Rpi512MB
		info.model = ModelCM
		info.pcbRev = PcbRev1_1
	case "0018":
		//     |   0018   | A+      |    1-1  | 256 MB |   sony      |
		info.manufacturer = MakerSony
		info.mem = Rpi256MB
		info.model = ModelAPlus
		info.pcbRev = PcbRev1_1
	case "0019":
		//     |   0019   |    B+   |    1-2  | 512 MB |   egoman    |
		info.manufacturer = MakerEgoman
		info.mem = Rpi512MB
		info.model = ModelBPlus
		info.pcbRev = PcbRev1_2
	case "001a":
		//     |   001a   | compute |    1-1  | 512 MB |   egoman    |
		info.manufacturer = MakerEgoman
		info.mem = Rpi512MB
		info.model = ModelCM
		info.pcbRev = PcbRev1_1
	case "001b":
		//     |   001b   | A+     |    1-1  | 256 MB |   egoman    |
		info.manufacturer = MakerEgoman
		info.mem = Rpi256MB
		info.model = ModelAPlus
		info.pcbRev = PcbRev1_1
	default:
		info.manufacturer = MakerUnknown
		info.mem = RpiUnknownMB
		info.model = ModelUnknown
		info.pcbRev = PcbRevUnknown
	}
	return
}

/*
//-------------------------------------------------------------------------
//
// With the release of the Raspberry Pi 2, there is a new encoding of the
// Revision field in /proc/cpuinfo. The bit fields are as follows
//
//     +----+----+----+----+----+----+----+----+
//     |FEDC|BA98|7654|3210|FEDC|BA98|7654|3210|
//     +----+----+----+----+----+----+----+----+
//     |    |    |    |    |    |    |    |AAAA|
//     |    |    |    |    |    |BBBB|BBBB|    |
//     |    |    |    |    |CCCC|    |    |    |
//     |    |    |    |DDDD|    |    |    |    |
//     |    |    | EEE|    |    |    |    |    |
//     |    |    |F   |    |    |    |    |    |
//     |    |   G|    |    |    |    |    |    |
//     |    |  H |    |    |    |    |    |    |
//     +----+----+----+----+----+----+----+----+
//     |1098|7654|3210|9876|5432|1098|7654|3210|
//     +----+----+----+----+----+----+----+----+
//
// +---+-------+--------------+--------------------------------------------+
// | # | bits  |   contains   | values                                     |
// +---+-------+--------------+--------------------------------------------+
// | A | 00-03 | PCB Revision | (the pcb revision number)                  |
// | B | 04-11 | Model name   | A, B, A+, B+, B Pi2, Alpha, Compute Module |
// |   |       |              | unknown, B Pi3, Zero                       |
// | C | 12-15 | Processor    | BCM2835, BCM2836, BCM2837                  |
// | D | 16-19 | Manufacturer | Sony, Egoman, Embest, unknown, Embest      |
// | E | 20-22 | Memory size  | 256 MB, 512 MB, 1024 MB                    |
// | F | 23-23 | encoded flag | (if set, revision is a bit field)          |
// | G | 24-24 | waranty bit  | (if set, warranty void - Pre Pi2)          |
// | H | 25-25 | waranty bit  | (if set, warranty void - Post Pi2)         |
// +---+-------+--------------+--------------------------------------------+
//
// Also, due to some early issues the warranty bit has been move from bit
// 24 to bit 25 of the revision number (i.e. 0x2000000).
*/
func getPostRPI2FromRevision(revision string) (info RpiInfoT, err error) {
	vision, err := strconv.ParseUint(revision, 16, 32) //hex number without 0x lea
	if err != nil {
		return
	}

	encoded := (vision & (1 << 23)) >> 23
	if encoded != 1 {
		err = errors.New("the encoding flag must be 1 in post PI2")
		return
	}

	info.revision = vision

	warantybit := (vision & (1 << 25)) >> 25
	if warantybit == 1 {
		info.overVolted = true
	}

	mem := (vision & (7 << 20)) >> 20
	memindex := MemoryT(mem)
	switch memindex {
	case Rpi256MB, Rpi512MB, Rpi1024MB, RpiUnknownMB:
	default:
		memindex = RpiUnknownMB
	}
	info.mem = memindex

	Manufacturer := (vision & (0xf << 16)) >> 16
	Manufacturerindex := MakerT(Manufacturer)
	switch Manufacturerindex {
	case MakerSony, MakerEgoman, MakerEmbest, MakerSonyJapan, MakerEmbest1:
	default:
		Manufacturerindex = MakerUnknown

	}
	info.manufacturer = Manufacturerindex

	process := (vision & (0xf << 12)) >> 12
	processindex := ProcessorT(process)
	switch processindex {
	case BroadcomUnknown, Broadcom2835, Broadcom2836, Broadcom2837:
	default:
		processindex = BroadcomUnknown
	}
	info.processor = processindex

	model := (vision & (0xff << 4)) >> 4
	modelindex := ModelT(model)
	switch modelindex {
	case ModelA, ModelB, ModelAPlus, ModelBPlus, ModelAlpha, ModelCM, Model2B, ModelUnknown, Model3B, ModelZero, ModelCM3, ModelZeroW:
	default:
		modelindex = ModelUnknown
	}
	info.model = modelindex

	pcbrev := (vision & (0xf))
	pcbrevindex := PcbRevT(pcbrev)
	switch pcbrevindex {
	case PcbRevUnknown, PcbRev1, PcbRev1_1, PcbRev1_2, PcbRev2:
	default:
		pcbrevindex = PcbRevUnknown

	}
	info.pcbRev = pcbrevindex

	return
}

func GetBoardInfo() (info RpiInfoT, periphereBase int64, err error) {
	revision, err := getRevision()

	// suggest it is postPI2 firstly, then try prePI2
	info, err = getPostRPI2FromRevision(revision)
	if err != nil {
		info, err = getPreRPI2FromRevision(revision)
	}
	if err != nil {
		return
	}

	switch info.processor {
	case Broadcom2835:
		periphereBase = PeripheralBase2835
	case Broadcom2836:
		periphereBase = PeripheralBase2836
	case Broadcom2837:
		periphereBase = PeripheralBase2837
	default:
		err = errors.New("unknown processor")
	}
	return
}

// Read /proc/device-tree/soc/ranges and determine the base address.
func get_dt_ranges(filename string) (base int64, err error) {
	ranges, err := os.Open(filename)
	defer ranges.Close()
	if err != nil {
		return
	}

	b := make([]byte, 4)
	n, err := ranges.ReadAt(b, 4)
	if n != 4 || err != nil {
		return
	}
	buf := bytes.NewReader(b)
	err = binary.Read(buf, binary.BigEndian, &base)
	if err != nil {
		return
	}
	return
}

func getPeripheralBase() (base int64, err error) {
	base, err = get_dt_ranges("/proc/device-tree/soc/ranges")
	return
}
