package configs

import (
	"context"
	"math"
	"time"

	"github.com/grid-x/modbus"
)

type ModbusTCPClient struct {
	handler           *modbus.TCPClientHandler
	client            modbus.Client
	backgroundContext context.Context
}

func NewModbusTCPClient(
	backgroundContext context.Context,
	address string,
	unitID byte,
) (*ModbusTCPClient, error) {

	handler := modbus.NewTCPClientHandler(address)
	handler.Timeout = 3 * time.Second
	handler.SlaveID = unitID

	if err := handler.Connect(backgroundContext); err != nil {
		return nil, err
	}

	client := modbus.NewClient(handler)

	return &ModbusTCPClient{
		handler:           handler,
		client:            client,
		backgroundContext: backgroundContext,
	}, nil
}

func (m *ModbusTCPClient) Close() error {
	if m.handler != nil {
		return m.handler.Close()
	}
	return nil
}

func (m *ModbusTCPClient) float32ToRegistersBE(val float32) (uint16, uint16) {
	bits := math.Float32bits(val)
	high := uint16(bits >> 16)
	low := uint16(bits & 0xFFFF)
	return high, low
}

func (m *ModbusTCPClient) WriteFloat32(
	startAddr uint16,
	value float32,
) error {

	high, low := m.float32ToRegistersBE(value)

	payload := []byte{
		byte(high >> 8), byte(high),
		byte(low >> 8), byte(low),
	}

	_, err := m.client.WriteMultipleRegisters(
		m.backgroundContext,
		startAddr,
		2,
		payload,
	)

	return err
}
