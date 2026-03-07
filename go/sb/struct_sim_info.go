package sb

import (
	"bytes"
	"fmt"
	"slices"
)

type SimInfo struct {
	Id uint32 `bson:"id" json:"id"`
	Title string `bson:"title" json:"title"`
	Content string `bson:"content" json:"content"`
	A bool `bson:"a" json:"a"`
	B bool `bson:"b" json:"b"`
	C bool `bson:"c" json:"c"`
	D bool `bson:"d" json:"d"`
	Zip []byte `bson:"zip" json:"zip"`
}

func sizeSimInfoBody(s *SimInfo, bits []byte) (int, error) {
	if s == nil { return 0, fmt.Errorf("sizeSimInfo: nil value") }
	bodySize := 0
	if s.Id != 0 {
		bodySize += 4
		if bits != nil { SetBit(bits, uint8(0), true) }
	}
	if s.Title != "" {
		fieldSize, err := sizeText(s.Title)
		if err != nil { return 0, fmt.Errorf("sizeSimInfo Title: %w", err) }
		bodySize += fieldSize
		if bits != nil { SetBit(bits, uint8(1), true) }
	}
	if s.Content != "" {
		fieldSize, err := sizeText(s.Content)
		if err != nil { return 0, fmt.Errorf("sizeSimInfo Content: %w", err) }
		bodySize += fieldSize
		if bits != nil { SetBit(bits, uint8(2), true) }
	}
	if bits != nil { SetBit(bits, uint8(3), s.A) }
	if bits != nil { SetBit(bits, uint8(4), s.B) }
	if bits != nil { SetBit(bits, uint8(5), s.C) }
	if bits != nil { SetBit(bits, uint8(6), s.D) }
	if s.Zip != nil {
		fieldSize, err := sizeBin(s.Zip)
		if err != nil { return 0, fmt.Errorf("sizeSimInfo Zip: %w", err) }
		bodySize += fieldSize
		if bits != nil { SetBit(bits, uint8(7), true) }
	}
	return bodySize, nil
}

func sizeSimInfo(s *SimInfo) (int, error) {
	bodySize, err := sizeSimInfoBody(s, nil)
	if err != nil { return 0, err }
	return 1 + bodySize, nil
}

func sizeSimInfoList(v SimInfoList) (int, error) {
	if len(v) > 65535 { return 0, fmt.Errorf("list length exceeds uint16 max") }
	totalSize := 2
	for i, item := range v {
		itemSize, err := sizeSimInfo(item)
		if err != nil { return 0, fmt.Errorf("sizeSimInfoList[%d]: %w", i, err) }
		totalSize += itemSize
	}
	return totalSize, nil
}

func GetSimInfo(buf *bytes.Buffer, s *SimInfo) error {
	if s == nil { return nil }
	const bitSize = 1
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
	if err := ValidateSimInfo(s); err != nil { return fmt.Errorf("ValidateSimInfo: %w", err) }
	return nil
}

func SetSimInfo(buf *bytes.Buffer, s *SimInfo) error {
	if s == nil { return fmt.Errorf("SetSimInfo: nil value") }
	if err := ValidateSimInfo(s); err != nil { return fmt.Errorf("ValidateSimInfo: %w", err) }
	const bitSize = 1
	var bits [1]byte
	bodySize, err := sizeSimInfoBody(s, bits[:])
	if err != nil { return err }
	buf.Grow(bitSize + bodySize)
	if _, err := buf.Write(bits[:]); err != nil { return fmt.Errorf("SetSimInfo write bitmask: %w", err) }
	if s.Id != 0 {
		if err := SetU32(buf, s.Id); err != nil { return fmt.Errorf("SetSimInfo Id: %w", err) }
	}
	if s.Title != "" {
		if err := SetText(buf, s.Title); err != nil { return fmt.Errorf("SetSimInfo Title: %w", err) }
	}
	if s.Content != "" {
		if err := SetText(buf, s.Content); err != nil { return fmt.Errorf("SetSimInfo Content: %w", err) }
	}
	if s.Zip != nil {
		if err := SetBin(buf, s.Zip); err != nil { return fmt.Errorf("SetSimInfo Zip: %w", err) }
	}
	return nil
}

func ValidateSimInfo(s *SimInfo) error {
	if s == nil { return nil }
	return nil
}

func EqSimInfo(a, b *SimInfo) bool {
	if a == b { return true }
	if a == nil || b == nil { return false }
	if !EqU32(a.Id, b.Id) { return false }
	if !EqText(a.Title, b.Title) { return false }
	if !EqText(a.Content, b.Content) { return false }
	if !EqBool(a.A, b.A) { return false }
	if !EqBool(a.B, b.B) { return false }
	if !EqBool(a.C, b.C) { return false }
	if !EqBool(a.D, b.D) { return false }
	if !EqBin(a.Zip, b.Zip) { return false }
	return true
}

// Standalone functions
func ReadSimInfo(buf *bytes.Buffer) (*SimInfo, error) {
	s := new(SimInfo)
	return s, GetSimInfo(buf, s)
}

type SimInfoList []*SimInfo
func GetSimInfoList(buf *bytes.Buffer) (SimInfoList, error) {
	count, err := GetU16(buf)
	if err != nil { return nil, err }
	list := make([]*SimInfo, count)
	for i := range list {
		item, err := ReadSimInfo(buf)
		if err != nil { return nil, err }
		list[i] = item
	}
	return SimInfoList(list), nil
}
func SetSimInfoList(buf *bytes.Buffer, v SimInfoList) error {
	totalSize, err := sizeSimInfoList(v)
	if err != nil { return err }
	buf.Grow(totalSize)
	if err := SetU16(buf, uint16(len(v))); err != nil { return err }
	for _, item := range v {
		if err := SetSimInfo(buf, item); err != nil { return err }
	}
	return nil
}
func ValidateSimInfoList(v SimInfoList) error {
	for i, item := range v {
		if item == nil { return fmt.Errorf("SimInfoList[%d]: nil item", i) }
		if err := ValidateSimInfo(item); err != nil { return fmt.Errorf("SimInfoList[%d]: %w", i, err) }
	}
	return nil
}
func EqSimInfoList(a, b SimInfoList) bool { return slices.EqualFunc(a, b, EqSimInfo) }
