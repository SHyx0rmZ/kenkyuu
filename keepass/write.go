package keepass

import (
	"bufio"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/xml"
	"io"
	"log"
)

func randomBytes32() (r [32]byte) {
	var err error
	_, err = rand.Read(r[:])
	if err != nil {
		log.Fatalln(err)
	}
	return r
}

func randomBytes16() (r [16]byte) {
	x := randomBytes32()
	copy(r[:], x[:])
	return r
}

func UpdateProtectedValues(f *File, newProtectedStreamKey [32]byte) {
	pr := Protection(f.Options.ProtectedStreamKey[:])
	pw := Protection(newProtectedStreamKey[:])

	update := func(s *string) {
		pv, err := base64.StdEncoding.DecodeString(*s)
		if err != nil {
			panic(err)
		}
		pr.XORKeyStream(pv, pv)
		pw.XORKeyStream(pv, pv)
		*s = base64.StdEncoding.EncodeToString(pv)
	}
	var walk func(*Group)
	walk = func(g *Group) {
		for _, e := range g.Entry {
			for i, f := range e.String {
				if !f.Value.Protected || len(f.Value.Value) == 0 {
					continue
				}

				update(&e.String[i].Value.Value)
			}

			for i, f := range e.Binary {
				if !f.Value.Protected || len(f.Value.Value) == 0 {
					continue
				}

				update(&e.Binary[i].Value.Value)
			}

			for _, h := range e.History {
				for i, f := range h.String {
					if !f.Value.Protected || len(f.Value.Value) == 0 {
						continue
					}

					update(&h.String[i].Value.Value)
				}

				for i, f := range h.Binary {
					if !f.Value.Protected || len(f.Value.Value) == 0 {
						continue
					}

					update(&h.Binary[i].Value.Value)
				}
			}
		}

		for _, c := range g.Group {
			walk(c)
		}
	}
	walk(f.Root.Group)

	f.Options.ProtectedStreamKey = newProtectedStreamKey
}

func WriteFile(w io.Writer, f *File, key Key) error {
	hh := sha256.New()
	ff := io.MultiWriter(w, hh)

	UpdateProtectedValues(f, randomBytes32())

	f.Options.MasterSeed = randomBytes32()
	f.Options.EncryptionIV = randomBytes16()
	f.Options.StreamStartBytes = randomBytes32()

	err := WriteOptions(ff, f.Options)
	if err != nil {
		log.Fatalln(err)
	}

	mk, err := masterKey(key, f.Options)
	if err != nil {
		return err
	}

	aesBlock, err := aes.NewCipher(mk)
	if err != nil {
		return err
	}

	hv := hh.Sum(nil)
	bm := cipher.NewCBCEncrypter(aesBlock, f.Options.EncryptionIV[:])

	sb := make([]byte, len(f.Options.StreamStartBytes))
	bm.CryptBlocks(sb, f.Options.StreamStartBytes[:])
	_, err = ff.Write(sb)

	bw := NewBlockModeWriter(bm, ff)
	defer bw.Close()

	pw := &PrefixWriter{
		w: bw,
	}
	defer pw.Write(nil)

	bfw := bufio.NewWriter(pw)
	defer bfw.Flush()

	gw, _ := gzip.NewWriterLevel(bfw, gzip.DefaultCompression)
	defer gw.Close()

	_, err = gw.Write([]byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`))
	if err != nil {
		return err
	}

	f.Meta.HeaderHash = make([]byte, base64.StdEncoding.EncodedLen(len(hv)))
	base64.StdEncoding.Encode(f.Meta.HeaderHash, hv)

	enc := xml.NewEncoder(gw)
	enc.Indent("", "    ")
	err = enc.Encode(&f)
	if err != nil {
		return err
	}
	return nil
}
