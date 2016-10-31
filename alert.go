package alert

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/gordonklaus/portaudio"
	"github.com/ventu-io/slf"
)

var pwdCurr string = "KristinaEtc/github.com/alert"
var log = slf.WithContext(pwdCurr)

type Music struct {
	fDesc     *os.File
	musicFile string
	songData  *io.SectionReader
	chunk     *io.SectionReader
	id        ID
	stop      chan bool
}

var m music
var mutex = &sync.Mutex{}

// Init initialize files of Music structure
func Init(fName string) {

	m := Music{
		musicFile: fName,
	}

	var err error
	m.fDesc, err = os.Open(m.musicFile)
	if err != nil {
		log.Fatal(err.Error())
	}
}

// Close method close all opened descriptors
// of Music structure
func (m Music) Close() {
	if m.fDesc != nil {
		m.fDesc.Close()
	}
}

// Stop is manage music state:
// m.Stop(true) of m.Stop(false).
// By default it is continue playing at the beginning after every Stop() func calling.
func (m Music) Stop(stopState bool) {
	if m.fDesc != nil {
		m.fDesc.Close()
	}

	if stopState {
		go m.Run()
	} else {
		m.stop <- true
		log.Warn("stopping")
	}
}

// Run is a main process where file will be readed
// by chunk and reproduced it
func (m Music) Run() {
	mutex.Lock()
	defer mutex.Unlock()

	var err error
	m.id, m.songData, err = readChunk(m.fDesc)
	if err != nil {
		log.Fatal(err.Error())
	}

	if m.id.String() != "FORM" {
		log.Fatal("bad file format")
	}
	_, err = m.songData.Read(m.id[:])
	if err != nil {
		log.Fatal(err.Error())
	}

	if m.id.String() != "AIFF" {
		log.Fatal("bad file format")
	}

	var c commonChunk
	var audio io.Reader
	for {

		m.id, m.chunk, err = readChunk(m.songData)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err.Error())
		}

		switch m.id.String() {
		case "COMM":
			err = binary.Read(m.chunk, binary.BigEndian, &c)
			if err != nil {
				log.Fatal(err.Error())
			}
		case "SSND":
			m.chunk.Seek(8, 1) //ignore offset and block
			audio = m.chunk
		default:
			log.Warnf("ignoring unknown chunk '%s'\n", m.id)
		}
	}

	//assume 44100 sample rate, mono, 32 bit

	portaudio.Initialize()
	defer portaudio.Terminate()
	out := make([]int32, 8192)
	stream, err := portaudio.OpenDefaultStream(0, 1, 44100, len(out), &out)
	chk(err)
	defer stream.Close()

	chk(stream.Start())
	defer stream.Stop()
	for remaining := int(c.NumSamples); remaining > 0; remaining -= len(out) {
		if len(out) > remaining {
			out = out[:remaining]
		}
		err := binary.Read(audio, binary.BigEndian, out)
		if err == io.EOF {
			break
		}
		chk(err)
		chk(stream.Write())
		go func() {
			stopMusic := <-stop
			switch stopMusic {
			case false:
				fmt.Print("starting music\n")
				stream.Start()
			case true:
				fmt.Print("stopping music\n")
				stream.Stop()
				return
			default:
			}
		}()
	}
}

func readChunk(r readerAtSeeker) (id ID, data *io.SectionReader, err error) {
	_, err = r.Read(id[:])
	if err != nil {
		return
	}
	var n int32
	err = binary.Read(r, binary.BigEndian, &n)
	if err != nil {
		return
	}
	off, _ := r.Seek(0, 1)
	data = io.NewSectionReader(r, off, int64(n))
	_, err = r.Seek(int64(n), 1)
	return
}

type readerAtSeeker interface {
	io.Reader
	io.ReaderAt
	io.Seeker
}

type ID [4]byte

func (id ID) String() string {
	return string(id[:])
}

type commonChunk struct {
	NumChans      int16
	NumSamples    int32
	BitsPerSample int16
	SampleRate    [10]byte
}
