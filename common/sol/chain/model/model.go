package model

import (
	bin "github.com/gagliardetto/binary"
)

type DistriEvent interface {
	UnmarshalWithDecoder(decoder *bin.Decoder) error
}

func DecodeDistriEvent(bytes []byte, event DistriEvent) error {
	decoder := bin.NewBorshDecoder(bytes)
	if err := event.UnmarshalWithDecoder(decoder); err != nil {
		return err
	}
	return nil
}
