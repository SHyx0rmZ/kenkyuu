package keepass

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"reflect"
	"time"

	"github.com/SHyx0rmZ/kenkyuu/uuid"
)

const signature1 = 0x9aa2d903
const signature2 = 0xb54bfb67
const fileVersion = 0x00030001
const fileVersionMin = 0x00020000
const fileVersionCriticalMask = 0xffff0000

type headerField uint8

const (
	EndOfHeader headerField = iota
	Comment
	CipherID
	CompressionFlags
	MasterSeed
	TransformSeed
	TransformRounds
	EncryptionIV
	ProtectedStreamKey
	StreamStartBytes
	InnerRandomStreamID
)

var cipherAES = uuid.UUID{0x31, 0xc1, 0xf2, 0xe6, 0xbf, 0x71, 0x43, 0x50, 0xbe, 0x58, 0x05, 0x21, 0x6a, 0xfc, 0x5a, 0xff}

type CompressionAlgorithm uint32

const (
	CompressionAlgorithmNone CompressionAlgorithm = iota
	CompressionAlgorithmGZIP
)

type EncryptionAlgorithm uint32

const (
	EncryptionNone EncryptionAlgorithm = iota
	EncryptionARC4
	EncryptionSalsa20
)

type Header struct {
	Signatures [2]uint32
	Version    uint32
}

type Options struct {
	Cipher uuid.UUID
	CompressionAlgorithm
	MasterSeed          [32]byte
	TransformSeed       [32]byte
	TransformRounds     int64
	EncryptionIV        [16]byte
	ProtectedStreamKey  [32]byte
	StreamStartBytes    [32]byte
	InnerRandomStreamID EncryptionAlgorithm
}

func ReadOptions(f io.Reader) (Options, error) {
	var header Header

	err := binary.Read(f, binary.LittleEndian, &header)
	if err != nil {
		return Options{}, err
	}

	if header.Signatures[0] != signature1 || header.Signatures[1] != signature2 {
		return Options{}, errors.New("invalid signature")
	}

	version := header.Version & fileVersionCriticalMask
	if version < fileVersionMin || version > (fileVersion&fileVersionCriticalMask) {
		return Options{}, errors.New("unsupported file version")
	}

	var opts Options

	var end bool
	for !end {
		var field struct {
			ID     headerField
			Length uint16
		}

		err = binary.Read(f, binary.LittleEndian, &field)
		if err != nil {
			return Options{}, err
		}

		switch field.ID {
		case EndOfHeader:
			end = true
			_, err = io.Copy(io.Discard, io.LimitReader(f, int64(field.Length)))
			if err != nil {
				return Options{}, err
			}
		case CipherID:
			bs, err := io.ReadAll(io.LimitReader(f, int64(field.Length)))
			if err != nil {
				return Options{}, err
			}
			if len(bs) != len(opts.Cipher) {
				return Options{}, errors.New("too few bytes")
			}
			for i := range opts.Cipher {
				opts.Cipher[i] = bs[i]
			}
			if !reflect.DeepEqual(opts.Cipher, cipherAES) {
				return Options{}, errors.New("unknown cipher")
			}
		case CompressionFlags:
			err = binary.Read(f, binary.LittleEndian, &opts.CompressionAlgorithm)
			if err != nil {
				return Options{}, err
			}
			if opts.CompressionAlgorithm > 1 {
				return Options{}, errors.New("unknown compression algorithm")
			}
		case MasterSeed:
			err = binary.Read(f, binary.LittleEndian, &opts.MasterSeed)
			if err != nil {
				return Options{}, err
			}
		case TransformSeed:
			err = binary.Read(f, binary.LittleEndian, &opts.TransformSeed)
			if err != nil {
				return Options{}, err
			}
		case TransformRounds:
			err = binary.Read(f, binary.LittleEndian, &opts.TransformRounds)
			if err != nil {
				return Options{}, err
			}
		case EncryptionIV:
			err = binary.Read(f, binary.LittleEndian, &opts.EncryptionIV)
			if err != nil {
				return Options{}, err
			}
		case ProtectedStreamKey:
			err = binary.Read(f, binary.LittleEndian, &opts.ProtectedStreamKey)
			if err != nil {
				return Options{}, err
			}
		case StreamStartBytes:
			err = binary.Read(f, binary.LittleEndian, &opts.StreamStartBytes)
			if err != nil {
				return Options{}, err
			}
		case InnerRandomStreamID:
			err = binary.Read(f, binary.LittleEndian, &opts.InnerRandomStreamID)
			if err != nil {
				return Options{}, err
			}
		default:
			_, err = io.Copy(io.Discard, io.LimitReader(f, int64(field.Length)))
			if err != nil {
				return Options{}, err
			}
			fmt.Println(field)
		}
	}

	return opts, nil
}

func WriteOptions(w io.Writer, opts Options) error {
	err := binary.Write(w, binary.LittleEndian, &Header{
		Signatures: [2]uint32{signature1, signature2},
		Version:    fileVersion,
	})
	if err != nil {
		return err
	}

	err = binary.Write(w, binary.LittleEndian, struct {
		CipherID struct {
			ID     headerField
			Length uint16
			Cipher uuid.UUID
		}
		CompressionFlags struct {
			ID     headerField
			Length uint16
			CompressionAlgorithm
		}
		MasterSeed struct {
			ID         headerField
			Length     uint16
			MasterSeed [32]byte
		}
		TransformSeed struct {
			ID            headerField
			Length        uint16
			TransformSeed [32]byte
		}
		TransformRounds struct {
			ID              headerField
			Length          uint16
			TransformRounds int64
		}
		EncryptionIV struct {
			ID           headerField
			Length       uint16
			EncryptionIV [16]byte
		}
		ProtectedStreamKey struct {
			ID                 headerField
			Length             uint16
			ProtectedStreamKey [32]byte
		}
		StreamStartBytes struct {
			ID               headerField
			Length           uint16
			StreamStartBytes [32]byte
		}
		InnerRandomStreamID struct {
			ID                  headerField
			Length              uint16
			InnerRandomStreamID EncryptionAlgorithm
		}
		EndOfHeader struct {
			ID        headerField
			Length    uint16
			Delimiter [4]byte
		}
	}{
		CipherID: struct {
			ID     headerField
			Length uint16
			Cipher uuid.UUID
		}{
			ID:     CipherID,
			Length: uint16(len(opts.Cipher)),
			Cipher: opts.Cipher,
		},
		CompressionFlags: struct {
			ID     headerField
			Length uint16
			CompressionAlgorithm
		}{
			ID:                   CompressionFlags,
			Length:               4,
			CompressionAlgorithm: opts.CompressionAlgorithm,
		},
		MasterSeed: struct {
			ID         headerField
			Length     uint16
			MasterSeed [32]byte
		}{
			ID:         MasterSeed,
			Length:     32,
			MasterSeed: opts.MasterSeed,
		},
		TransformSeed: struct {
			ID            headerField
			Length        uint16
			TransformSeed [32]byte
		}{
			ID:            TransformSeed,
			Length:        32,
			TransformSeed: opts.TransformSeed,
		},
		TransformRounds: struct {
			ID              headerField
			Length          uint16
			TransformRounds int64
		}{
			ID:              TransformRounds,
			Length:          8,
			TransformRounds: opts.TransformRounds,
		},
		EncryptionIV: struct {
			ID           headerField
			Length       uint16
			EncryptionIV [16]byte
		}{
			ID:           EncryptionIV,
			Length:       16,
			EncryptionIV: opts.EncryptionIV,
		},
		ProtectedStreamKey: struct {
			ID                 headerField
			Length             uint16
			ProtectedStreamKey [32]byte
		}{
			ID:                 ProtectedStreamKey,
			Length:             32,
			ProtectedStreamKey: opts.ProtectedStreamKey,
		},
		StreamStartBytes: struct {
			ID               headerField
			Length           uint16
			StreamStartBytes [32]byte
		}{
			ID:               StreamStartBytes,
			Length:           32,
			StreamStartBytes: opts.StreamStartBytes,
		},
		InnerRandomStreamID: struct {
			ID                  headerField
			Length              uint16
			InnerRandomStreamID EncryptionAlgorithm
		}{
			ID:                  InnerRandomStreamID,
			Length:              4,
			InnerRandomStreamID: opts.InnerRandomStreamID,
		},
		EndOfHeader: struct {
			ID        headerField
			Length    uint16
			Delimiter [4]byte
		}{
			ID:        EndOfHeader,
			Length:    4,
			Delimiter: [4]byte{'\r', '\n', '\r', '\n'},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

type File struct {
	Header  Header  `xml:"-"`
	Options Options `xml:"-"`
	Meta    Meta
	Root    struct {
		Group          *Group
		DeletedObjects string
	}
}

func (k File) Write(w io.Writer) error {
	err := WriteOptions(w, k.Options)
	if err != nil {
		return err
	}
	return nil
}

type Meta struct {
	Generator                  string
	HeaderHash                 []byte
	DatabaseName               string
	DatabaseNameChanged        time.Time
	DatabaseDescription        string
	DatabaseDescriptionChanged time.Time
	DefaultUserName            string
	DefaultUserNameChanged     time.Time
	MaintenanceHistoryDays     int
	Color                      string
	MasterKeyChanged           time.Time
	MasterKeyChangeRec         int
	MasterKeyChangedForce      int
	MemoryProtection           struct {
		ProtectTitle    Bool
		ProtectUserName Bool
		ProtectPassword Bool
		ProtectURL      Bool
		ProtectNotes    Bool
	}
	CustomIcons                []CustomIcon
	RecycleBinEnabled          Bool
	RecycleBinUUID             []byte
	RecycleBinChanged          time.Time
	EntryTemplatesGroup        []byte
	EntryTemplatesGroupChanged time.Time
	LastSelectedGroup          []byte
	LastTopVisibleGroup        []byte
	HistoryMaxItems            int
	HistoryMaxSize             int
	Binaries                   string
	CustomData                 string
}

type CustomIcon struct {
	UUID []byte
	Data []byte
}

func NewFile() File {
	return File{
		Header: Header{
			Signatures: [2]uint32{signature1, signature2},
			Version:    fileVersion,
		},
		Options: Options{
			Cipher:               cipherAES,
			CompressionAlgorithm: CompressionAlgorithmGZIP,
			MasterSeed:           randomBytes32(),
			TransformSeed:        randomBytes32(),
			TransformRounds:      100000,
			EncryptionIV:         randomBytes16(),
			ProtectedStreamKey:   randomBytes32(),
			StreamStartBytes:     randomBytes32(),
			InnerRandomStreamID:  EncryptionSalsa20,
		},
		Meta: Meta{
			Generator:                  "kenkyuu",
			DatabaseNameChanged:        time.Now(),
			DatabaseDescriptionChanged: time.Now(),
			DefaultUserNameChanged:     time.Now(),
			MasterKeyChanged:           time.Now(),
			RecycleBinChanged:          time.Now(),
			EntryTemplatesGroupChanged: time.Now(),
		},
	}
}
