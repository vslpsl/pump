package pipe

import (
	"errors"
	"fmt"
	"io"
	"log"
)

type Limiter interface {
	Lease(size int64)
}

type Pipe struct {
	source     io.Reader
	target     io.Writer
	limiter    Limiter
	bufferSize int64

	stat *Stat
}

func New(source io.Reader, target io.Writer, limiter Limiter, bufferSize int64) *Pipe {
	return &Pipe{
		source:     source,
		target:     target,
		limiter:    limiter,
		bufferSize: bufferSize,
		stat:       &Stat{},
	}
}

func (p *Pipe) Pump() error {
	log.Printf("buffer size: %d bytes\n", p.bufferSize)
	p.stat.Start()
	defer func() {
		p.stat.Stop()
		fmt.Println(p.stat)
	}()

	if p.limiter == nil {
		return p.copy()
	}

	return p.LimitedCopy()
}

func (p *Pipe) copy() error {
	buffer := make([]byte, p.bufferSize)
	bytesWritten, err := io.CopyBuffer(p.target, p.source, buffer)
	p.stat.BytesPiped = bytesWritten
	if err != nil {
		return err
	}

	return nil
}

func (p *Pipe) LimitedCopy() error {
	buffer := make([]byte, p.bufferSize)
	var written int64
	var err error

	for {
		p.limiter.Lease(p.bufferSize)

		bytesRead, readErr := p.source.Read(buffer)
		if bytesRead > 0 {
			bytesWritten, writeErr := p.target.Write(buffer[0:bytesRead])
			if bytesWritten < 0 || bytesRead < bytesWritten {
				bytesWritten = 0
				if writeErr == nil {
					writeErr = errors.New("invalid write result")
				}
			}
			written += int64(bytesWritten)
			if writeErr != nil {
				err = writeErr
				break
			}
			if bytesRead != bytesWritten {
				err = io.ErrShortWrite
				break
			}
		}
		if readErr != nil {
			if readErr != io.EOF {
				err = readErr
			}
			break
		}
	}

	p.stat.BytesPiped = written

	return err
}
