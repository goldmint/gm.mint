package transaction

import (
	"bytes"
	"io"

	sumuslib "github.com/void616/gm-sumuslib"
	"github.com/void616/gm-sumuslib/serializer"
	"golang.org/x/crypto/sha3"
)

// parser parses transaction
type parser struct {
	*serializer.Deserializer
	digestWriter *bytes.Buffer

	nonce uint64
}

// ParsedTransaction data
type ParsedTransaction struct {
	Nonce     uint64
	From      sumuslib.PublicKey
	Digest    sumuslib.Digest
	Signature sumuslib.Signature
}

func newParser(r io.Reader) (*parser, error) {
	digestWriter := bytes.NewBuffer(make([]byte, 256))
	digestWriter.Reset()
	des := serializer.NewStreamDeserializer(io.TeeReader(r, digestWriter))

	// read nonce
	txnonce := des.GetUint64()
	if err := des.Error(); err != nil {
		return nil, err
	}

	return &parser{
		Deserializer: des,
		digestWriter: digestWriter,
		nonce:        txnonce,
	}, nil
}

// Complete completes parsing and returns a parsed transaction common data
func (p *parser) Complete(from sumuslib.PublicKey) (*ParsedTransaction, error) {
	// errors?
	if err := p.Error(); err != nil {
		return nil, err
	}

	// calc tx digest
	var digest sumuslib.Digest
	{
		hasher := sha3.New256()
		_, err := hasher.Write(p.digestWriter.Bytes())
		if err != nil {
			return nil, err
		}
		b := hasher.Sum(nil)
		copy(digest[:], b)
	}

	// "signed" byte
	var signed = p.GetByte()
	if err := p.Error(); err != nil {
		return nil, err
	}

	// get signature if the tx is signed
	var signature sumuslib.Signature
	{
		if signed != 0 {
			// signature
			b := p.GetBytes(sumuslib.SignatureSize)
			if err := p.Error(); err != nil {
				return nil, err
			}
			copy(signature[:], b)
			// TODO: verify signature?
		} else {
			// digest
			_ = p.GetBytes(sumuslib.DigestSize)
			if err := p.Error(); err != nil {
				return nil, err
			}
			// TODO: compare digests?
		}
	}

	return &ParsedTransaction{
		From:      from,
		Nonce:     p.nonce,
		Digest:    digest,
		Signature: signature,
	}, nil
}
