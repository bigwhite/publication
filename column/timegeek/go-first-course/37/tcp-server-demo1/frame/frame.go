package frame

import (
	"encoding/binary"
	"errors"
	"io"
)

/*
Frame定义

frameHeader + framePayload(packet)

frameHeader
	4 bytes: length 整型，帧总长度(含头及payload)

framePayload
	Packet
*/

type FramePayload []byte

type StreamFrameCodec interface {
	Encode(io.Writer, FramePayload) error   // data -> frame，并写入io.Writer
	Decode(io.Reader) (FramePayload, error) // 从io.Reader中提取frame payload，并返回给上层
}

var ErrShortWrite = errors.New("short write")
var ErrShortRead = errors.New("short read")

type myFrameCodec struct{}

func NewMyFrameCodec() StreamFrameCodec {
	return &myFrameCodec{}
}

func (p *myFrameCodec) Encode(w io.Writer, framePayload FramePayload) error {
	var f = framePayload
	var totalLen int32 = int32(len(framePayload)) + 4 // 4对应Frame中totalLength的4个字节。总长度 = payload字节长度+totalLength的字节长度

	// 先在frame中写上totalLength
	err := binary.Write(w, binary.BigEndian, &totalLen)
	if err != nil {
		return err
	}

	// 再在frame中写上payload
	n, err := w.Write([]byte(f)) // write the frame payload to outbound stream
	if err != nil {
		return err
	}

	if n != len(framePayload) {
		return ErrShortWrite
	}
	return nil
}

func (p *myFrameCodec) Decode(r io.Reader) (FramePayload, error) {
	var totalLen int32                                 // int32为4个字节，对应了frame中totalLength的4个字节
	err := binary.Read(r, binary.BigEndian, &totalLen) // 先从frame中读出前面的4个字节，该4个字节为frame中的totalLength
	if err != nil {
		return nil, err
	}

	buf := make([]byte, totalLen-4) // 再从frame中读出totalLength以后的剩下的字节，这些剩下的字节为payload
	n, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}

	if n != int(totalLen-4) {
		return nil, ErrShortRead
	}

	return FramePayload(buf), nil
}
