package keepass

import (
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"io"
	"log"
	"reflect"

	"github.com/SHyx0rmZ/kenkyuu/crypto/salsa20"
)

func ReadFile(f io.Reader, key Key) (File, error) {
	headerHash := sha256.New()

	opts, err := ReadOptions(io.TeeReader(f, headerHash))
	if err != nil {
		return File{}, err
	}

	mk, err := masterKey(key, opts)
	if err != nil {
		return File{}, err
	}

	aesBlock, err := aes.NewCipher(mk)
	if err != nil {
		log.Fatalln(err)
	}

	headerChecksum := headerHash.Sum(nil)

	bs, err := io.ReadAll(io.LimitReader(f, 32))
	if err != nil {
		return File{}, err
	}

	bm := cipher.NewCBCDecrypter(aesBlock, opts.EncryptionIV[:])
	bm.CryptBlocks(bs, bs)

	if !reflect.DeepEqual(bs[:len(opts.StreamStartBytes)], opts.StreamStartBytes[:]) {
		return File{}, errors.New("invalid start bytes")
	}

	bmr := NewBlockModeReader(bm, f)

	var reader io.Reader = &PrefixReader{
		Reader:   bmr,
		Prefixer: Block{},
	}

	if opts.CompressionAlgorithm == CompressionAlgorithmGZIP {
		gr, err := gzip.NewReader(reader)
		if err != nil {
			return File{}, err
		}
		defer gr.Close()
		reader = gr
	}

	var kpf File
	err = xml.NewDecoder(reader).Decode(&kpf)
	if err != nil {
		return File{}, err
	}

	if string(kpf.Meta.HeaderHash) != base64.StdEncoding.EncodeToString(headerChecksum) {
		return File{}, errors.New("invalid header hash")
	}

	kpf.Options = opts
	kpf.Header = Header{
		Signatures: [2]uint32{signature1, signature2},
		Version:    fileVersion,
	}

	return kpf, nil
}

func Protection(protectedStreamKey []byte) cipher.Stream {
	k := sha256.Sum256(protectedStreamKey)
	return salsa20.New(&[8]byte{0xe8, 0x30, 0x09, 0x4b, 0x97, 0x20, 0x5d, 0x2a}, &k)
}
