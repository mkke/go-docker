package attach

import (
	"encoding/binary"
	"io"

	"docker.io/go-docker/api/types"
	"github.com/mkke/go-mlog"
)

type Handler struct {
	hijackedResponse types.HijackedResponse
	stdOut           io.Writer
	stdErr           io.Writer
	stdIn            io.Reader
	closeChs         []chan struct{}
	log              mlog.Logger
}

func NewHandler(hijackedResponse types.HijackedResponse) *Handler {
	return &Handler{
		hijackedResponse: hijackedResponse,
		log:              mlog.NewNopLogger()}
}

func (ah *Handler) WithStdout(stdout io.Writer) *Handler {
	ah.stdOut = stdout
	return ah
}

func (ah *Handler) WithStderr(stderr io.Writer) *Handler {
	ah.stdErr = stderr
	return ah
}

func (ah *Handler) WithStdin(stdin io.Reader) *Handler {
	ah.stdIn = stdin
	return ah
}

func (ah *Handler) WithLogger(log mlog.Logger) *Handler {
	ah.log = log
	return ah
}

func (ah *Handler) AddCloseListener(closeCh chan struct{}) *Handler {
	ah.closeChs = append(ah.closeChs, closeCh)
	return ah
}

func (ah *Handler) Close() {
	ah.hijackedResponse.Close()
}

type StreamType uint8

const (
	StreamStdin  = StreamType(0)
	StreamStdout = StreamType(1)
	StreamStderr = StreamType(2)
)

func (ah *Handler) Start() {
	go func() {
		if err := ah.WriteLoop(); err != nil {
			ah.log.Printf("write loop: %v", err)
		}
	}()
	go func() {
		if err := ah.ReadLoop(); err != nil {
			ah.log.Printf("read loop: %v", err)
		}
	}()
}

func (ah *Handler) WriteLoop() error {
	if ah.stdIn == nil {
		return nil
	}

	_, err := io.Copy(ah.hijackedResponse.Conn, ah.stdIn)
	_ = ah.hijackedResponse.CloseWrite()

	return err
}

func (ah *Handler) ReadLoop() error {
	defer func() {
		for _, closeCh := range ah.closeChs {
			ah.log.Printf("closing channel %v", closeCh)
			close(closeCh)
		}
	}()

	for {
		streamType, err := ah.hijackedResponse.Reader.ReadByte()
		if err != nil {
			return err
		}
		ah.log.Printf("received frame with type %d", streamType)

		if _, err := ah.hijackedResponse.Reader.Discard(3); err != nil {
			return err
		}

		var size uint32
		if err := binary.Read(ah.hijackedResponse.Reader, binary.BigEndian, &size); err != nil {
			return err
		}

		switch StreamType(streamType) {
		case StreamStdout:
			if ah.stdOut != nil {
				if _, err := io.CopyN(ah.stdOut, ah.hijackedResponse.Reader, int64(size)); err != nil {
					return err
				}
				if syncer, ok := ah.stdOut.(Syncer); ok {
					_ = syncer.Sync()
				}
				continue
			}
		case StreamStderr:
			if ah.stdErr != nil {
				if _, err := io.CopyN(ah.stdErr, ah.hijackedResponse.Reader, int64(size)); err != nil {
					return err
				}
				if syncer, ok := ah.stdErr.(Syncer); ok {
					_ = syncer.Sync()
				}
				continue
			}
		}

		// drain frame
		if _, err := ah.hijackedResponse.Reader.Discard(int(size)); err != nil {
			return err
		}
	}
}
