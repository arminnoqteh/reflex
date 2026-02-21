package encoding

import (
  "encoding/binary"
  "errors"
  "io"
)

// Frame: represents a Reflex protocol frame
type Frame struct {
    Length  uint16
    Type    uint8
    Payload []byte
}

const (
    FrameTypeData =    0x01
    FrameTypePadding = 0x02
    FrameTypeTiming =  0x03
    FrameTypeClose =   0x04
)

//  ReadFrame: reads and decrypts a frame from the reader
func (s *Session) ReadFrame(r io.Reader) (*Frame, error) {
    // read frame header
    header := make([]byte, 3)
    if _, err := io.ReadFull(r, header); err != nil {
        return nil, err
    }

    frameType := header[0]
    encPayloadLen := binary.BigEndian.Uint16(header[1:3])

    // read encrypted payload
    encPayload := make([]byte, encPayloadLen)
    if _, err := io.ReadFull(r, encPayload); err != nil {
        return nil, err
    }

    // decrypt the payload
    nonce := makeNonce(s.readNonce)
    s.readNonce++

    payload, err := s.aead.Open(nil, nonce[:], encPayload, nil)
    if err != nil {
        return nil, errors.New("failed to decrypt frame: possible corruption or replay attack")
    }

    return &Frame{
        Length:  encPayloadLen,
        Type:    frameType,
        Payload: payload,
    }, nil
}


// WriteFrame: encrypts and writes a frame to the writer
func (s *Session) WriteFrame(w io.Writer, frameType uint8, payload []byte) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    // encrypt the payload
    nonce := makeNonce(s.writeNonce)
    s.writeNonce++

    encrypted := s.aead.Seal(nil, nonce[:], payload, nil)

    // write frame header
    header := make([]byte, 3)
    header[0] = frameType
    binary.BigEndian.PutUint16(header[1:3], uint16(len(encrypted)))

    if _, err := w.Write(header); err != nil {
        return err
    }

    // write encrypted payload
    if _, err := w.Write(encrypted); err != nil {
        return err
    }

    return nil
}
