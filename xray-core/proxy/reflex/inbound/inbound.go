package inbound

import (
    "context"
    "errors"
    "io"

    "github.com/xtls/xray-core/common"
    "github.com/xtls/xray-core/common/net"
    "github.com/xtls/xray-core/common/protocol"
    "github.com/xtls/xray-core/features/routing"
    "github.com/xtls/xray-core/proxy/reflex/encoding"
    "github.com/xtls/xray-core/transport/internet/stat"
)

func init() {
    common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
        return New(ctx, config.(*Config))
    }))
}

type Handler struct {
    // ToDo
}

func New(ctx context.Context, config *Config) (*Handler, error) {
    return &Handler{}, nil
}

func (h *Handler) Network() []net.Network {
    return []net.Network{net.Network_TCP}
}

func (h *Handler) Process(ctx context.Context, network net.Network, conn stat.Connection, dispatcher routing.Dispatcher) error {
    // will be implemented in step 4
    return nil
}

// handleSession: processes the encrypted session
func (h *Handler) handleSession(ctx context.Context, reader io.Reader, conn stat.Connection, dispatcher routing.Dispatcher, sess *encoding.Session, user *protocol.MemoryUser) error {
    for {
        frame, err := sess.ReadFrame(reader)
        if err != nil {
            if err == io.EOF {
                return nil
            }
            return err
        }

        switch frame.Type {
        case encoding.FrameTypeData:
            // TODO: handle data frames 
            continue

        case encoding.FrameTypePadding:
            // ignored for now
            continue

        case encoding.FrameTypeTiming:
            // ignored for now
            continue

        case encoding.FrameTypeClose:
            return nil

        default:
            return errors.New("unknown frame type")
        }
    }
}
