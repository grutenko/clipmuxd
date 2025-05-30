package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"sync"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/curve25519"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GRPC server for handshake service

const (
	StateInit = iota
	StatePending
	StateTrust
	StateLast
)

type State int

var ErrNotFound = errors.New("session not found")

type HandshakeSession struct {
	State               State
	Id                  string
	deviceInitiatorName string
	deviceInitiatorUUID string
	pubKeyInitiator     []byte
	pubKey              []byte
	privKey             []byte
	sharedKey           []byte
	Expiry              time.Time
}

// HandshakeGrpcServer implements the HandshakeServiceServer interface.
type HandshakeGrpcServer struct {
	Ctx        *AppCtx
	sessionsMu sync.RWMutex
	Sessions   map[string]*HandshakeSession
	Control    *ControlGrpcServer
	UnimplementedHandshakeServer
}

func NewHandshakeGrpcServer(ctx *AppCtx, control *ControlGrpcServer) HandshakeGrpcServer {
	return HandshakeGrpcServer{
		Ctx:      ctx,
		Sessions: make(map[string]*HandshakeSession),
		Control:  control,
	}
}

func (h *HandshakeGrpcServer) findSession(id string) (*HandshakeSession, error) {
	session, ok := h.Sessions[id]
	if !ok {
		return nil, ErrNotFound
	}
	return session, nil
}

func (h *HandshakeGrpcServer) removeSession(id string) error {
	delete(h.Sessions, id)
	return nil
}

func (h *HandshakeGrpcServer) createSession(
	id string,
	deviceInitiatorUUID string,
	deviceInitiatorName string,
	pubKeyInitiator []byte,
	pubKey []byte,
	privKey []byte,
	sharedKey []byte,
) *HandshakeSession {
	session := &HandshakeSession{
		State:               StateInit,
		Id:                  id,
		deviceInitiatorUUID: deviceInitiatorUUID,
		deviceInitiatorName: deviceInitiatorName,
		pubKeyInitiator:     pubKeyInitiator,
		pubKey:              pubKey,
		privKey:             privKey,
		sharedKey:           sharedKey,
		Expiry:              time.Now().Add(5 * time.Minute),
	}
	h.Sessions[id] = session
	return session
}

func (h *HandshakeGrpcServer) sessionGc() {
	for id, session := range h.Sessions {
		if time.Now().After(session.Expiry) {
			delete(h.Sessions, id)
		}
	}
}

func (h *HandshakeGrpcServer) Init(ctx context.Context, req *HandshakeInitRequest) (*HandshakeInitResponse, error) {
	pair, err := h.Ctx.Storage.GetPair(req.DeviceId)
	if pair != nil {
		return nil, status.Error(codes.AlreadyExists, "device already paired")
	}
	var privKey, pubKey [32]byte
	rand.Read(privKey[:])
	curve25519.ScalarBaseMult(&pubKey, &privKey)
	sharedKey, err := curve25519.X25519(privKey[:], req.PublicKey[:])
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate shared key")
	}
	id := make([]byte, 16)
	_, err = rand.Read(id)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate session id")
	}
	idBase64 := base64.StdEncoding.EncodeToString(id)

	h.sessionsMu.Lock()
	defer h.sessionsMu.Unlock()

	session := h.createSession(idBase64, req.DeviceId, req.DeviceName, req.PublicKey, pubKey[:], privKey[:], sharedKey[:])
	return &HandshakeInitResponse{
		Id:         idBase64,
		PublicKey:  pubKey[:],
		DeviceId:   h.Ctx.Config.DeviceId,
		DeviceName: h.Ctx.Config.DeviceName,
		Expiry:     timestamppb.New(session.Expiry),
	}, nil
}

func (h *HandshakeGrpcServer) Code(ctx context.Context, req *HandshakeCodeRequest) (*HandshakeCodeResponse, error) {
	h.sessionsMu.Lock()
	h.sessionGc()
	session, err := h.findSession(req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "session not found")
	}
	if session.State != StateInit {
		return nil, status.Error(codes.FailedPrecondition, "session not in init state")
	}
	session.State = StatePending

	hash := sha256.Sum256(session.sharedKey)
	code := binary.BigEndian.Uint16(hash[0:2])

	err = h.Control.SendDeviceHandshakeEvent(session.Id, session.deviceInitiatorName, code)
	if err != nil {
		h.removeSession(session.Id)
		h.sessionsMu.Unlock()
		return nil, status.Error(codes.Internal, "failed to send handshake event")
	}
	h.sessionsMu.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	accepted, err := h.Control.WaitDeviceHandshakeReason(ctx, session.Id)
	h.sessionsMu.Lock()
	defer h.sessionsMu.Unlock()
	if err != nil {
		h.removeSession(session.Id)
		return nil, status.Error(codes.Internal, "failed to accept handshake event")
	}

	if !accepted {
		h.removeSession(session.Id)
		return nil, status.Error(codes.PermissionDenied, "device handshake rejected")
	}
	session.State = StateTrust

	return &HandshakeCodeResponse{}, nil
}

func encrypt(sharedKey []byte, plaintext []byte) (nonce, ciphertext []byte, err error) {

}

func decrypt(sharedKey []byte, nonce []byte, ciphertext []byte) ([]byte, error) {

}

func (h *HandshakeGrpcServer) Cancel(ctx context.Context, req *HandshakeCancelRequest) (*HandshakeCancelResponse, error) {
	h.sessionsMu.Lock()
	defer h.sessionsMu.Unlock()

	h.removeSession(req.Id)
	return nil, nil
}

func (h *HandshakeGrpcServer) ExchangeTrust(ctx context.Context, req *HandshakeExchangeTrustRequest) (*HandshakeExchangeTrustResponse, error) {
	h.sessionsMu.Lock()
	defer h.sessionsMu.Unlock()

	h.sessionGc()
	session, err := h.findSession(req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "session not found")
	}
	if session.State != StateTrust {
		return nil, status.Error(codes.FailedPrecondition, "session not in trust state")
	}
	return nil, nil
}
