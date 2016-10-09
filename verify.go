package main

import (
	"time"

	"github.com/dustin/go-humanize"
	"github.com/go-errors/errors"
	"github.com/itchio/butler/comm"
	"github.com/itchio/wharf/eos"
	"github.com/itchio/wharf/pwr"
)

func verify(signaturePath string, dir string, woundsPath string) {
	must(doVerify(signaturePath, dir, woundsPath))
}

func doVerify(signaturePath string, dir string, woundsPath string) error {
	if woundsPath == "" {
		comm.Opf("Verifying %s", dir)
	} else {
		comm.Opf("Verifying %s, writing wounds to %s", dir, woundsPath)
	}
	startTime := time.Now()

	signatureReader, err := eos.Open(signaturePath)
	if err != nil {
		return errors.Wrap(err, 1)
	}
	defer signatureReader.Close()

	signature, err := pwr.ReadSignature(signatureReader)
	if err != nil {
		return errors.Wrap(err, 1)
	}

	vc := &pwr.ValidatorContext{
		Consumer:   comm.NewStateConsumer(),
		WoundsPath: woundsPath,
	}

	comm.StartProgress()

	err = vc.Validate(dir, signature)
	if err != nil {
		return errors.Wrap(err, 1)
	}

	comm.EndProgress()

	prettySize := humanize.IBytes(uint64(signature.Container.Size))
	perSecond := humanize.IBytes(uint64(float64(signature.Container.Size) / time.Since(startTime).Seconds()))
	comm.Statf("%s (%s) @ %s/s\n", prettySize, signature.Container.Stats(), perSecond)

	if vc.WoundsConsumer.HasWounds() {
		comm.Dief("%s corrupted data found", humanize.IBytes(uint64(vc.WoundsConsumer.TotalCorrupted())))
	}

	return nil
}
