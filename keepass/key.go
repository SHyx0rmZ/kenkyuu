package keepass

import (
	"crypto/aes"
	"crypto/sha256"
	"io"
	"os"
)

type Key [32]byte

func (k Key) Transform(seed []byte, rounds int) ([]byte, error) {
	cl, err := aes.NewCipher(seed)
	if err != nil {
		return nil, err
	}
	cr, err := aes.NewCipher(seed)
	if err != nil {
		return nil, err
	}
	key := make([]byte, len(k))
	copy(key, k[:])
	for rounds > 0 {
		cl.Encrypt(key[0:16], key[0:16])
		cr.Encrypt(key[16:32], key[16:32])
		rounds--
	}
	return key, nil
}

func CompositeKey(keys ...[32]byte) Key {
	h := sha256.New()
	for _, k := range keys {
		h.Write(k[:])
	}
	var r [32]byte
	copy(r[:], h.Sum(nil))
	return r
}

func FileKey(path string) ([32]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return [32]byte{}, err
	}
	defer f.Close()

	h := sha256.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return [32]byte{}, err
	}

	var r [32]byte
	copy(r[:], h.Sum(nil))
	return r, nil
}

func PasswordKey(password []byte) [32]byte {
	return sha256.Sum256(password)
}

func masterKey(key Key, opts Options) ([]byte, error) {
	transformed, err := key.Transform(opts.TransformSeed[:], int(opts.TransformRounds))
	if err != nil {
		return nil, err
	}

	transformedHash := sha256.New()
	transformedHash.Write(transformed)

	masterHash := sha256.New()
	masterHash.Write(opts.MasterSeed[:])
	masterHash.Write(transformedHash.Sum(nil))

	return masterHash.Sum(nil), nil
}
