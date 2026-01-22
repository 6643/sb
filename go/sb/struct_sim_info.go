package sb

import (
	"bytes"
	"fmt"
	"math"
	"slices"
	
)

type SimInfo struct {
	Id uint32  
	Title string  
	Content string  
	A bool  
	B bool  
	C bool  
	D bool  
	Zip []byte  
}

func (s *SimInfo) Get(buf *bytes.Buffer) error {
	if buf.Len() == 0 { return nil }
	bitSize := int(math.Ceil(float64(8) / 8.0))
	if buf.Len() < bitSize { return fmt.Errorf("GetSimInfo bitmask: %d - %d", buf.Len(), bitSize) }
	bits := buf.Next(bitSize)
	if GetBit(bits, uint8(0)) {
		val, err := GetU32(buf)
		if err != nil { return fmt.Errorf("GetSimInfo Id: %w", err) }
		s.Id = val
	}
	if GetBit(bits, uint8(1)) {
		val, err := GetText(buf)
		if err != nil { return fmt.Errorf("GetSimInfo Title: %w", err) }
		s.Title = val
	}
	if GetBit(bits, uint8(2)) {
		val, err := GetText(buf)
		if err != nil { return fmt.Errorf("GetSimInfo Content: %w", err) }
		s.Content = val
	}
	s.A = GetBit(bits, uint8(3))
	s.B = GetBit(bits, uint8(4))
	s.C = GetBit(bits, uint8(5))
	s.D = GetBit(bits, uint8(6))
	if GetBit(bits, uint8(7)) {
		val, err := GetBin(buf)
		if err != nil { return fmt.Errorf("GetSimInfo Zip: %w", err) }
		s.Zip = val
	}
	return nil
}

func (s *SimInfo) Set(buf *bytes.Buffer) error {
	if s == nil { return nil }
	bits := make([]byte, uint8(math.Ceil(float64(8)/8.0)))
	body := bytes.NewBuffer(nil)
	if s.Id != 0 {
		if err := SetU32(body, s.Id); err != nil { return fmt.Errorf("SetSimInfo Id: %w", err) }
		SetBit(bits, uint8(0), true)
	}
	if s.Title != "" {
		if err := SetText(body, s.Title); err != nil { return fmt.Errorf("SetSimInfo Title: %w", err) }
		SetBit(bits, uint8(1), true)
	}
	if s.Content != "" {
		if err := SetText(body, s.Content); err != nil { return fmt.Errorf("SetSimInfo Content: %w", err) }
		SetBit(bits, uint8(2), true)
	}
	SetBit(bits, uint8(3), s.A)
	SetBit(bits, uint8(4), s.B)
	SetBit(bits, uint8(5), s.C)
	SetBit(bits, uint8(6), s.D)
	if s.Zip != nil {
		if err := SetBin(body, s.Zip); err != nil { return fmt.Errorf("SetSimInfo Zip: %w", err) }
		SetBit(bits, uint8(7), true)
	}

	if _, err := buf.Write(bits); err != nil { return fmt.Errorf("SetSimInfo write bitmask: %w", err) }
	_, err := body.WriteTo(buf); return err
}

func (s *SimInfo) Eq(other *SimInfo) bool {
	if s == other { return true }
	if s == nil || other == nil { return false }
	if !EqU32(s.Id, other.Id) { return false }
	if !EqText(s.Title, other.Title) { return false }
	if !EqText(s.Content, other.Content) { return false }
	if !EqBool(s.A, other.A) { return false }
	if !EqBool(s.B, other.B) { return false }
	if !EqBool(s.C, other.C) { return false }
	if !EqBool(s.D, other.D) { return false }
	if !EqBin(s.Zip, other.Zip) { return false }
	return true
}

// Standalone functions for compatibility
func GetSimInfo(buf *bytes.Buffer) (*SimInfo, error) {
	s := new(SimInfo); return s, s.Get(buf)
}
func SetSimInfo(buf *bytes.Buffer, s *SimInfo) error { return s.Set(buf) }
func EqSimInfo(a, b *SimInfo) bool { return a.Eq(b) }

type SimInfoList []*SimInfo
func (v SimInfoList) Set(buf *bytes.Buffer) error { return setList(buf, v, SetSimInfo) }
func (v *SimInfoList) Get(buf *bytes.Buffer) error {
	val, err := getList[*SimInfo, SimInfoList](buf, GetSimInfo)
	if err == nil { *v = val }; return err
}
func (v SimInfoList) Eq(other SimInfoList) bool { return slices.EqualFunc(v, other, EqSimInfo) }
