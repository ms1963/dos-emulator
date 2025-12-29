package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type CPU struct {
	AX, BX, CX, DX uint16
	SI, DI         uint16
	SP, BP         uint16
	CS, DS, ES, SS uint16
	IP             uint16
	Flags          Flags
}

type Flags struct {
	CF, PF, AF, ZF, SF, TF, IF, DF, OF bool
}

func (f *Flags) ToUint16() uint16 {
	var result uint16 = 0x0002
	if f.CF {
		result = result | 0x0001
	}
	if f.PF {
		result = result | 0x0004
	}
	if f.AF {
		result = result | 0x0010
	}
	if f.ZF {
		result = result | 0x0040
	}
	if f.SF {
		result = result | 0x0080
	}
	if f.TF {
		result = result | 0x0100
	}
	if f.IF {
		result = result | 0x0200
	}
	if f.DF {
		result = result | 0x0400
	}
	if f.OF {
		result = result | 0x0800
	}
	return result
}

func (f *Flags) FromUint16(value uint16) {
	f.CF = (value & 0x0001) != 0
	f.PF = (value & 0x0004) != 0
	f.AF = (value & 0x0010) != 0
	f.ZF = (value & 0x0040) != 0
	f.SF = (value & 0x0080) != 0
	f.TF = (value & 0x0100) != 0
	f.IF = (value & 0x0200) != 0
	f.DF = (value & 0x0400) != 0
	f.OF = (value & 0x0800) != 0
}

func (c *CPU) GetAL() byte {
	return byte(c.AX & 0xFF)
}

func (c *CPU) SetAL(value byte) {
	c.AX = (c.AX & 0xFF00) | uint16(value)
}

func (c *CPU) GetAH() byte {
	return byte((c.AX >> 8) & 0xFF)
}

func (c *CPU) SetAH(value byte) {
	c.AX = (c.AX & 0x00FF) | (uint16(value) << 8)
}

func (c *CPU) GetBL() byte {
	return byte(c.BX & 0xFF)
}

func (c *CPU) SetBL(value byte) {
	c.BX = (c.BX & 0xFF00) | uint16(value)
}

func (c *CPU) GetBH() byte {
	return byte((c.BX >> 8) & 0xFF)
}

func (c *CPU) SetBH(value byte) {
	c.BX = (c.BX & 0x00FF) | (uint16(value) << 8)
}

func (c *CPU) GetCL() byte {
	return byte(c.CX & 0xFF)
}

func (c *CPU) SetCL(value byte) {
	c.CX = (c.CX & 0xFF00) | uint16(value)
}

func (c *CPU) GetCH() byte {
	return byte((c.CX >> 8) & 0xFF)
}

func (c *CPU) SetCH(value byte) {
	c.CX = (c.CX & 0x00FF) | (uint16(value) << 8)
}

func (c *CPU) GetDL() byte {
	return byte(c.DX & 0xFF)
}

func (c *CPU) SetDL(value byte) {
	c.DX = (c.DX & 0xFF00) | uint16(value)
}

func (c *CPU) GetDH() byte {
	return byte((c.DX >> 8) & 0xFF)
}

func (c *CPU) SetDH(value byte) {
	c.DX = (c.DX & 0x00FF) | (uint16(value) << 8)
}

func (c *CPU) UpdateZeroFlag(result uint16) {
	c.Flags.ZF = result == 0
}

func (c *CPU) UpdateZeroFlag8(result byte) {
	c.Flags.ZF = result == 0
}

func (c *CPU) UpdateSignFlag(result uint16) {
	c.Flags.SF = (result & 0x8000) != 0
}

func (c *CPU) UpdateSignFlag8(result byte) {
	c.Flags.SF = (result & 0x80) != 0
}

func (c *CPU) UpdateParityFlag(result uint16) {
	count := 0
	value := byte(result & 0xFF)
	for i := 0; i < 8; i++ {
		if (value & (1 << uint(i))) != 0 {
			count++
		}
	}
	c.Flags.PF = (count % 2) == 0
}

func (c *CPU) UpdateArithmeticFlags16(result uint16) {
	c.UpdateZeroFlag(result)
	c.UpdateSignFlag(result)
	c.UpdateParityFlag(result)
}

func (c *CPU) UpdateArithmeticFlags8(result byte) {
	c.UpdateZeroFlag8(result)
	c.UpdateSignFlag8(result)
	c.UpdateParityFlag(uint16(result))
}

type Memory struct {
	data [0x100000]byte
}

func (m *Memory) ReadByte(addr uint32) byte {
	if addr >= uint32(len(m.data)) {
		return 0
	}
	return m.data[addr]
}

func (m *Memory) WriteByte(addr uint32, value byte) {
	if addr < uint32(len(m.data)) {
		m.data[addr] = value
	}
}

func (m *Memory) ReadWord(addr uint32) uint16 {
	low := uint16(m.ReadByte(addr))
	high := uint16(m.ReadByte(addr + 1))
	return (high << 8) | low
}

func (m *Memory) WriteWord(addr uint32, value uint16) {
	m.WriteByte(addr, byte(value % 0xFF))
	m.WriteByte(addr+1, byte((value>>8) &0xFF))
}

type VideoMemory struct {
	buffer       [80 * 25 * 2]byte
	cursorX      int
	cursorY      int
	currentColor byte
	videoMode    byte
}

type EXEHeader struct {
	Signature       uint16
	BytesInLastPage uint16
	PagesInFile     uint16
	Relocations     uint16
	HeaderSize      uint16
	MinAlloc        uint16
	MaxAlloc        uint16
	InitialSS       uint16
	InitialSP       uint16
	Checksum        uint16
	InitialIP       uint16
	InitialCS       uint16
	RelocTableOff   uint16
	OverlayNumber   uint16
}

func ReadEXEHeader(data []byte) (*EXEHeader, error) {
	if len(data) < 28 {
		return nil, fmt.Errorf("file too small for EXE header")
	}

	header := &EXEHeader{
		Signature:       binary.LittleEndian.Uint16(data[0:2]),
		BytesInLastPage: binary.LittleEndian.Uint16(data[2:4]),
		PagesInFile:     binary.LittleEndian.Uint16(data[4:6]),
		Relocations:     binary.LittleEndian.Uint16(data[6:8]),
		HeaderSize:      binary.LittleEndian.Uint16(data[8:10]),
		MinAlloc:        binary.LittleEndian.Uint16(data[10:12]),
		MaxAlloc:        binary.LittleEndian.Uint16(data[12:14]),
		InitialSS:       binary.LittleEndian.Uint16(data[14:16]),
		InitialSP:       binary.LittleEndian.Uint16(data[16:18]),
		Checksum:        binary.LittleEndian.Uint16(data[18:20]),
		InitialIP:       binary.LittleEndian.Uint16(data[20:22]),
		InitialCS:       binary.LittleEndian.Uint16(data[22:24]),
		RelocTableOff:   binary.LittleEndian.Uint16(data[24:26]),
		OverlayNumber:   binary.LittleEndian.Uint16(data[26:28]),
	}

	if header.Signature != 0x5A4D && header.Signature != 0x4D5A {
		return nil, fmt.Errorf("invalid EXE signature")
	}

	return header, nil
}

type FileSystem struct {
	currentDir   string
	currentDrive byte
	drives       map[byte]string
}

type FileHandle struct {
	file     *os.File
	handle   uint16
	position int64
}

type DTA struct {
	reserved    [21]byte
	attribute   byte
	time        uint16
	date        uint16
	size        uint32
	name        [13]byte
	searchPath  string
	searchIndex int
	searchFiles []os.DirEntry
}

type Instruction struct {
	Opcode    byte
	ModRM     byte
	HasModRM  bool
	Length    int
	Operand1  uint16
	Operand2  uint16
	Immediate uint16
	Name      string
}

type InstructionDecoder struct {
	memory *Memory
}

func NewInstructionDecoder(memory *Memory) *InstructionDecoder {
	return &InstructionDecoder{memory: memory}
}

func (d *InstructionDecoder) calculateModRMLength(modrm byte) int {
	mod := (modrm >> 6) & 0x03
	rm := modrm & 0x07

	switch mod {
	case 0:
		if rm == 6 {
			return 2
		}
		return 0
	case 1:
		return 1
	case 2:
		return 2
	case 3:
		return 0
	}
	return 0
}

func (d *InstructionDecoder) Decode(addr uint32) *Instruction {
	inst := &Instruction{
		Opcode: d.memory.ReadByte(addr),
		Length: 1,
	}

	switch inst.Opcode {
	case 0x90:
		inst.Name = "NOP"
	case 0xB0, 0xB1, 0xB2, 0xB3, 0xB4, 0xB5, 0xB6, 0xB7:
		inst.Operand1 = uint16(d.memory.ReadByte(addr + 1))
		inst.Length = 2
		inst.Name = fmt.Sprintf("MOV r8, 0x%02X", inst.Operand1)
	case 0xB8, 0xB9, 0xBA, 0xBB, 0xBC, 0xBD, 0xBE, 0xBF:
		inst.Operand1 = d.memory.ReadWord(addr + 1)
		inst.Length = 3
		inst.Name = fmt.Sprintf("MOV r16, 0x%04X", inst.Operand1)
	case 0x88, 0x89, 0x8A, 0x8B, 0x8C, 0x8E:
		inst.ModRM = d.memory.ReadByte(addr + 1)
		inst.HasModRM = true
		inst.Length = 2 + d.calculateModRMLength(inst.ModRM)
		inst.Name = "MOV"
	case 0xA0, 0xA1, 0xA2, 0xA3:
		inst.Operand1 = d.memory.ReadWord(addr + 1)
		inst.Length = 3
		inst.Name = "MOV"
	case 0xC6, 0xC7:
		inst.ModRM = d.memory.ReadByte(addr + 1)
		inst.HasModRM = true
		if inst.Opcode == 0xC6 {
			inst.Length = 3 + d.calculateModRMLength(inst.ModRM)
		} else {
			inst.Length = 4 + d.calculateModRMLength(inst.ModRM)
		}
		inst.Name = "MOV"
	case 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57:
		inst.Name = "PUSH"
	case 0x58, 0x59, 0x5A, 0x5B, 0x5C, 0x5D, 0x5E, 0x5F:
		inst.Name = "POP"
	case 0x06, 0x0E, 0x16, 0x1E:
		inst.Name = "PUSH SEG"
	case 0x07, 0x17, 0x1F:
		inst.Name = "POP SEG"
	case 0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47:
		inst.Name = "INC"
	case 0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F:
		inst.Name = "DEC"
	case 0x00, 0x01, 0x02, 0x03, 0x04, 0x05:
		if inst.Opcode <= 0x03 {
			inst.ModRM = d.memory.ReadByte(addr + 1)
			inst.HasModRM = true
			inst.Length = 2 + d.calculateModRMLength(inst.ModRM)
		} else if inst.Opcode == 0x04 {
			inst.Immediate = uint16(d.memory.ReadByte(addr + 1))
			inst.Length = 2
		} else {
			inst.Immediate = d.memory.ReadWord(addr + 1)
			inst.Length = 3
		}
		inst.Name = "ADD"
	case 0x10, 0x11, 0x12, 0x13, 0x14, 0x15:
		if inst.Opcode <= 0x13 {
			inst.ModRM = d.memory.ReadByte(addr + 1)
			inst.HasModRM = true
			inst.Length = 2 + d.calculateModRMLength(inst.ModRM)
		} else if inst.Opcode == 0x14 {
			inst.Immediate = uint16(d.memory.ReadByte(addr + 1))
			inst.Length = 2
		} else {
			inst.Immediate = d.memory.ReadWord(addr + 1)
			inst.Length = 3
		}
		inst.Name = "ADC"
	case 0x28, 0x29, 0x2A, 0x2B, 0x2C, 0x2D:
		if inst.Opcode <= 0x2B {
			inst.ModRM = d.memory.ReadByte(addr + 1)
			inst.HasModRM = true
			inst.Length = 2 + d.calculateModRMLength(inst.ModRM)
		} else if inst.Opcode == 0x2C {
			inst.Immediate = uint16(d.memory.ReadByte(addr + 1))
			inst.Length = 2
		} else {
			inst.Immediate = d.memory.ReadWord(addr + 1)
			inst.Length = 3
		}
		inst.Name = "SUB"
	case 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D:
		if inst.Opcode <= 0x1B {
			inst.ModRM = d.memory.ReadByte(addr + 1)
			inst.HasModRM = true
			inst.Length = 2 + d.calculateModRMLength(inst.ModRM)
		} else if inst.Opcode == 0x1C {
			inst.Immediate = uint16(d.memory.ReadByte(addr + 1))
			inst.Length = 2
		} else {
			inst.Immediate = d.memory.ReadWord(addr + 1)
			inst.Length = 3
		}
		inst.Name = "SBB"
	case 0x38, 0x39, 0x3A, 0x3B, 0x3C, 0x3D:
		if inst.Opcode <= 0x3B {
			inst.ModRM = d.memory.ReadByte(addr + 1)
			inst.HasModRM = true
			inst.Length = 2 + d.calculateModRMLength(inst.ModRM)
		} else if inst.Opcode == 0x3C {
			inst.Immediate = uint16(d.memory.ReadByte(addr + 1))
			inst.Length = 2
		} else {
			inst.Immediate = d.memory.ReadWord(addr + 1)
			inst.Length = 3
		}
		inst.Name = "CMP"
	case 0x20, 0x21, 0x22, 0x23, 0x24, 0x25:
		if inst.Opcode <= 0x23 {
			inst.ModRM = d.memory.ReadByte(addr + 1)
			inst.HasModRM = true
			inst.Length = 2 + d.calculateModRMLength(inst.ModRM)
		} else if inst.Opcode == 0x24 {
			inst.Immediate = uint16(d.memory.ReadByte(addr + 1))
			inst.Length = 2
		} else {
			inst.Immediate = d.memory.ReadWord(addr + 1)
			inst.Length = 3
		}
		inst.Name = "AND"
	case 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D:
		if inst.Opcode <= 0x0B {
			inst.ModRM = d.memory.ReadByte(addr + 1)
			inst.HasModRM = true
			inst.Length = 2 + d.calculateModRMLength(inst.ModRM)
		} else if inst.Opcode == 0x0C {
			inst.Immediate = uint16(d.memory.ReadByte(addr + 1))
			inst.Length = 2
		} else {
			inst.Immediate = d.memory.ReadWord(addr + 1)
			inst.Length = 3
		}
		inst.Name = "OR"
	case 0x30, 0x31, 0x32, 0x33, 0x34, 0x35:
		if inst.Opcode <= 0x33 {
			inst.ModRM = d.memory.ReadByte(addr + 1)
			inst.HasModRM = true
			inst.Length = 2 + d.calculateModRMLength(inst.ModRM)
		} else if inst.Opcode == 0x34 {
			inst.Immediate = uint16(d.memory.ReadByte(addr + 1))
			inst.Length = 2
		} else {
			inst.Immediate = d.memory.ReadWord(addr + 1)
			inst.Length = 3
		}
		inst.Name = "XOR"
	case 0x84, 0x85, 0xA8, 0xA9:
		if inst.Opcode <= 0x85 {
			inst.ModRM = d.memory.ReadByte(addr + 1)
			inst.HasModRM = true
			inst.Length = 2 + d.calculateModRMLength(inst.ModRM)
		} else if inst.Opcode == 0xA8 {
			inst.Immediate = uint16(d.memory.ReadByte(addr + 1))
			inst.Length = 2
		} else {
			inst.Immediate = d.memory.ReadWord(addr + 1)
			inst.Length = 3
		}
		inst.Name = "TEST"
	case 0x86, 0x87:
		inst.ModRM = d.memory.ReadByte(addr + 1)
		inst.HasModRM = true
		inst.Length = 2 + d.calculateModRMLength(inst.ModRM)
		inst.Name = "XCHG"
	case 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97:
		inst.Name = "XCHG"
	case 0x8D:
		inst.ModRM = d.memory.ReadByte(addr + 1)
		inst.HasModRM = true
		inst.Length = 2 + d.calculateModRMLength(inst.ModRM)
		inst.Name = "LEA"
	case 0xA4:
		inst.Name = "MOVSB"
	case 0xA5:
		inst.Name = "MOVSW"
	case 0xA6:
		inst.Name = "CMPSB"
	case 0xA7:
		inst.Name = "CMPSW"
	case 0xAA:
		inst.Name = "STOSB"
	case 0xAB:
		inst.Name = "STOSW"
	case 0xAC:
		inst.Name = "LODSB"
	case 0xAD:
		inst.Name = "LODSW"
	case 0xAE:
		inst.Name = "SCASB"
	case 0xAF:
		inst.Name = "SCASW"
	case 0xE8:
		inst.Operand1 = d.memory.ReadWord(addr + 1)
		inst.Length = 3
		inst.Name = "CALL"
	case 0xFF:
		inst.ModRM = d.memory.ReadByte(addr + 1)
		inst.HasModRM = true
		inst.Length = 2 + d.calculateModRMLength(inst.ModRM)
		reg := (inst.ModRM >> 3) & 0x07
		if reg == 2 {
			inst.Name = "CALL"
		} else if reg == 4 {
			inst.Name = "JMP"
		} else if reg == 6 {
			inst.Name = "PUSH"
		}
	case 0x9A:
		inst.Operand1 = d.memory.ReadWord(addr + 1)
		inst.Operand2 = d.memory.ReadWord(addr + 3)
		inst.Length = 5
		inst.Name = "CALL FAR"
	case 0xE9:
		inst.Operand1 = d.memory.ReadWord(addr + 1)
		inst.Length = 3
		inst.Name = "JMP"
	case 0xEB:
		inst.Operand1 = uint16(d.memory.ReadByte(addr + 1))
		inst.Length = 2
		inst.Name = "JMP SHORT"
	case 0xEA:
		inst.Operand1 = d.memory.ReadWord(addr + 1)
		inst.Operand2 = d.memory.ReadWord(addr + 3)
		inst.Length = 5
		inst.Name = "JMP FAR"
	case 0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77:
		inst.Operand1 = uint16(d.memory.ReadByte(addr + 1))
		inst.Length = 2
		jmpNames := []string{"JO", "JNO", "JB", "JNB", "JZ", "JNZ", "JBE", "JA"}
		inst.Name = jmpNames[inst.Opcode-0x70]
	case 0x78, 0x79, 0x7A, 0x7B, 0x7C, 0x7D, 0x7E, 0x7F:
		inst.Operand1 = uint16(d.memory.ReadByte(addr + 1))
		inst.Length = 2
		jmpNames := []string{"JS", "JNS", "JP", "JNP", "JL", "JGE", "JLE", "JG"}
		inst.Name = jmpNames[inst.Opcode-0x78]
	case 0xE0:
		inst.Operand1 = uint16(d.memory.ReadByte(addr + 1))
		inst.Length = 2
		inst.Name = "LOOPNE"
	case 0xE1:
		inst.Operand1 = uint16(d.memory.ReadByte(addr + 1))
		inst.Length = 2
		inst.Name = "LOOPE"
	case 0xE2:
		inst.Operand1 = uint16(d.memory.ReadByte(addr + 1))
		inst.Length = 2
		inst.Name = "LOOP"
	case 0xE3:
		inst.Operand1 = uint16(d.memory.ReadByte(addr + 1))
		inst.Length = 2
		inst.Name = "JCXZ"
	case 0xC2:
		inst.Operand1 = d.memory.ReadWord(addr + 1)
		inst.Length = 3
		inst.Name = "RET"
	case 0xC3:
		inst.Name = "RET"
	case 0xCA:
		inst.Operand1 = d.memory.ReadWord(addr + 1)
		inst.Length = 3
		inst.Name = "RETF"
	case 0xCB:
		inst.Name = "RETF"
	case 0xCD:
		inst.Operand1 = uint16(d.memory.ReadByte(addr + 1))
		inst.Length = 2
		inst.Name = fmt.Sprintf("INT 0x%02X", inst.Operand1)
	case 0xCC:
		inst.Name = "INT 3"
	case 0xCE:
		inst.Name = "INTO"
	case 0xCF:
		inst.Name = "IRET"
	case 0x9C:
		inst.Name = "PUSHF"
	case 0x9D:
		inst.Name = "POPF"
	case 0x98:
		inst.Name = "CBW"
	case 0x99:
		inst.Name = "CWD"
	case 0x9E:
		inst.Name = "SAHF"
	case 0x9F:
		inst.Name = "LAHF"
	case 0xF4:
		inst.Name = "HLT"
	case 0xF5:
		inst.Name = "CMC"
	case 0xF8:
		inst.Name = "CLC"
	case 0xF9:
		inst.Name = "STC"
	case 0xFA:
		inst.Name = "CLI"
	case 0xFB:
		inst.Name = "STI"
	case 0xFC:
		inst.Name = "CLD"
	case 0xFD:
		inst.Name = "STD"
	case 0xD0, 0xD1, 0xD2, 0xD3:
		inst.ModRM = d.memory.ReadByte(addr + 1)
		inst.HasModRM = true
		inst.Length = 2 + d.calculateModRMLength(inst.ModRM)
		reg := (inst.ModRM >> 3) & 0x07
		shiftNames := []string{"ROL", "ROR", "RCL", "RCR", "SHL", "SHR", "SAL", "SAR"}
		inst.Name = shiftNames[reg]
	case 0xF6, 0xF7:
		inst.ModRM = d.memory.ReadByte(addr + 1)
		inst.HasModRM = true
		inst.Length = 2 + d.calculateModRMLength(inst.ModRM)
		reg := (inst.ModRM >> 3) & 0x07
		if reg == 0 || reg == 1 {
			inst.Name = "TEST"
			if inst.Opcode == 0xF6 {
				inst.Length++
			} else {
				inst.Length += 2
			}
		} else if reg == 2 {
			inst.Name = "NOT"
		} else if reg == 3 {
			inst.Name = "NEG"
		} else if reg == 4 {
			inst.Name = "MUL"
		} else if reg == 5 {
			inst.Name = "IMUL"
		} else if reg == 6 {
			inst.Name = "DIV"
		} else if reg == 7 {
			inst.Name = "IDIV"
		}
	case 0xFE:
		inst.ModRM = d.memory.ReadByte(addr + 1)
		inst.HasModRM = true
		inst.Length = 2 + d.calculateModRMLength(inst.ModRM)
		reg := (inst.ModRM >> 3) & 0x07
		if reg == 0 {
			inst.Name = "INC"
		} else if reg == 1 {
			inst.Name = "DEC"
		}
	case 0x80, 0x81, 0x82, 0x83:
		inst.ModRM = d.memory.ReadByte(addr + 1)
		inst.HasModRM = true
		if inst.Opcode == 0x80 || inst.Opcode == 0x82 {
			inst.Length = 3 + d.calculateModRMLength(inst.ModRM)
		} else if inst.Opcode == 0x83 {
			inst.Length = 3 + d.calculateModRMLength(inst.ModRM)
		} else {
			inst.Length = 4 + d.calculateModRMLength(inst.ModRM)
		}
		reg := (inst.ModRM >> 3) & 0x07
		opNames := []string{"ADD", "OR", "ADC", "SBB", "AND", "SUB", "XOR", "CMP"}
		inst.Name = opNames[reg]
	case 0xD4:
		inst.Operand1 = uint16(d.memory.ReadByte(addr + 1))
		inst.Length = 2
		inst.Name = "AAM"
	case 0xD5:
		inst.Operand1 = uint16(d.memory.ReadByte(addr + 1))
		inst.Length = 2
		inst.Name = "AAD"
	case 0xD7:
		inst.Name = "XLAT"
	case 0xF2, 0xF3:
		inst.Name = "REP"
	default:
		inst.Name = fmt.Sprintf("UNKNOWN (0x%02X)", inst.Opcode)
	}

	return inst
}

type DOSEmulator struct {
	cpu              *CPU
	memory           *Memory
	video            *VideoMemory
	fs               *FileSystem
	decoder          *InstructionDecoder
	dta              *DTA
	running          bool
	debugMode        bool
	stepMode         bool
	traceMode        bool
	breakpoints      map[uint32]bool
	fileHandles      map[uint16]*FileHandle
	nextHandle       uint16
	stack            []uint16
	instructionCount uint64
	startTime        time.Time
	interruptVectors [256]uint32
	environment      map[string]string
	psp              uint16
	repeatPrefix     byte
	programType      string
}

func NewDOSEmulator() *DOSEmulator {
	currentDir, _ := os.Getwd()
	memory := &Memory{}

	emulator := &DOSEmulator{
		cpu:     &CPU{},
		memory:  memory,
		video:   &VideoMemory{currentColor: 0x07, videoMode: 0x03},
		decoder: NewInstructionDecoder(memory),
		dta:     &DTA{},
		fs: &FileSystem{
			currentDir:   currentDir,
			currentDrive: 0,
			drives: map[byte]string{
				0: currentDir,
			},
		},
		running:      true,
		debugMode:    false,
		stepMode:     false,
		traceMode:    false,
		breakpoints:  make(map[uint32]bool),
		fileHandles:  make(map[uint16]*FileHandle),
		nextHandle:   5,
		stack:        make([]uint16, 0),
		startTime:    time.Now(),
		environment:  make(map[string]string),
		psp:          0x1000,
		repeatPrefix: 0,
	}

	emulator.environment["PATH"] = "A:\\"
	emulator.environment["COMSPEC"] = "A:\\COMMAND.COM"

	for i := 0; i < 256; i++ {
		emulator.interruptVectors[i] = 0
	}

	emulator.fileHandles[0] = &FileHandle{file: os.Stdin, handle: 0}
	emulator.fileHandles[1] = &FileHandle{file: os.Stdout, handle: 1}
	emulator.fileHandles[2] = &FileHandle{file: os.Stderr, handle: 2}

	return emulator
}

func CalculateAddress(segment, offset uint16) uint32 {
	return (uint32(segment) << 4) + uint32(offset)
}

func (e *DOSEmulator) Push(value uint16) {
	e.cpu.SP -= 2
	addr := CalculateAddress(e.cpu.SS, e.cpu.SP)
	e.memory.WriteWord(addr, value)
	if len(e.stack) < 1000 {
		e.stack = append(e.stack, value)
	}
}

func (e *DOSEmulator) Pop() uint16 {
	addr := CalculateAddress(e.cpu.SS, e.cpu.SP)
	value := e.memory.ReadWord(addr)
	e.cpu.SP += 2
	if len(e.stack) > 0 {
		e.stack = e.stack[:len(e.stack)-1]
	}
	return value
}

func (e *DOSEmulator) SetupPSP(segment uint16) {
	pspAddr := CalculateAddress(segment, 0)
	e.memory.WriteByte(pspAddr+0, 0xCD)
	e.memory.WriteByte(pspAddr+1, 0x20)
	e.memory.WriteWord(pspAddr+2, 0xA000)
	e.memory.WriteByte(pspAddr+4, 0)
	e.memory.WriteByte(pspAddr+5, 0x9A)
	e.memory.WriteWord(pspAddr+6, 0x0000)
	e.memory.WriteWord(pspAddr+8, 0x0000)
	e.memory.WriteWord(pspAddr+0x0A, 0x0000)
	e.memory.WriteWord(pspAddr+0x0C, segment)
	e.memory.WriteWord(pspAddr+0x16, segment)

	for i := uint32(0); i < 20; i++ {
		e.memory.WriteByte(pspAddr+0x18+i, 0xFF)
	}

	e.memory.WriteWord(pspAddr+0x2C, segment+0x10)
	e.memory.WriteByte(pspAddr+0x50, 0xCD)
	e.memory.WriteByte(pspAddr+0x51, 0x21)
	e.memory.WriteByte(pspAddr+0x52, 0xCB)

	for i := uint32(0); i < 16; i++ {
		e.memory.WriteByte(pspAddr+0x5C+i, 0)
	}

	for i := uint32(0); i < 16; i++ {
		e.memory.WriteByte(pspAddr+0x6C+i, 0)
	}

	e.memory.WriteByte(pspAddr+0x80, 0)
	e.memory.WriteByte(pspAddr+0x81, 0x0D)
}

func (e *DOSEmulator) LoadCOMFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	if len(data) > 65280 {
		return fmt.Errorf("COM file too large")
	}

	e.SetupPSP(e.psp)

	comStart := CalculateAddress(e.psp, 0x100)
	for i, b := range data {
		e.memory.WriteByte(comStart+uint32(i), b)
	}

	e.cpu.CS = e.psp
	e.cpu.DS = e.psp
	e.cpu.ES = e.psp
	e.cpu.SS = e.psp
	e.cpu.IP = 0x100
	e.cpu.SP = 0xFFFE
	e.cpu.AX = 0
	e.cpu.BX = 0
	e.cpu.CX = 0
	e.cpu.DX = 0
	e.cpu.SI = 0
	e.cpu.DI = 0
	e.cpu.BP = 0
	e.cpu.Flags = Flags{IF: true}

	e.programType = "COM"

	if e.debugMode {
		fmt.Printf("Loaded COM file: %s (%d bytes)\n", filename, len(data))
		fmt.Printf("Entry point: %04X:%04X\n", e.cpu.CS, e.cpu.IP)
	}

	return nil
}

func (e *DOSEmulator) LoadEXEFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	if len(data) < 28 {
		return fmt.Errorf("file too small to be an EXE")
	}

	header, err := ReadEXEHeader(data)
	if err != nil {
		return err
	}

	imageSize := int(header.PagesInFile) * 512
	if header.BytesInLastPage != 0 {
		imageSize = imageSize - 512 + int(header.BytesInLastPage)
	}

	headerSize := int(header.HeaderSize) * 16
	loadSegment := e.psp

	e.SetupPSP(loadSegment)

	programSegment := loadSegment + 0x10
	programData := data[headerSize:imageSize]
	loadAddr := CalculateAddress(programSegment, 0)

	for i, b := range programData {
		e.memory.WriteByte(loadAddr+uint32(i), b)
	}

	if header.Relocations > 0 && header.RelocTableOff > 0 {
		relocTableAddr := int(header.RelocTableOff)
		for i := 0; i < int(header.Relocations); i++ {
			if relocTableAddr+4 > len(data) {
				break
			}

			offset := binary.LittleEndian.Uint16(data[relocTableAddr : relocTableAddr+2])
			segment := binary.LittleEndian.Uint16(data[relocTableAddr+2 : relocTableAddr+4])

			relocAddr := CalculateAddress(programSegment+segment, offset)
			currentValue := e.memory.ReadWord(relocAddr)
			newValue := currentValue + programSegment
			e.memory.WriteWord(relocAddr, newValue)

			relocTableAddr += 4
		}
	}

	e.cpu.CS = programSegment + header.InitialCS
	e.cpu.IP = header.InitialIP
	e.cpu.SS = programSegment + header.InitialSS
	e.cpu.SP = header.InitialSP
	e.cpu.DS = loadSegment
	e.cpu.ES = loadSegment
	e.cpu.AX = 0
	e.cpu.BX = 0
	e.cpu.CX = 0
	e.cpu.DX = 0
	e.cpu.SI = 0
	e.cpu.DI = 0
	e.cpu.BP = 0
	e.cpu.Flags = Flags{IF: true}

	e.programType = "EXE"

	if e.debugMode {
		fmt.Printf("Loaded EXE file: %s\n", filename)
		fmt.Printf("Image size: %d bytes\n", len(programData))
		fmt.Printf("Relocations: %d\n", header.Relocations)
		fmt.Printf("Entry point: %04X:%04X\n", e.cpu.CS, e.cpu.IP)
		fmt.Printf("Initial stack: %04X:%04X\n", e.cpu.SS, e.cpu.SP)
	}

	return nil
}

func (e *DOSEmulator) LoadFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	if len(data) >= 2 {
		signature := binary.LittleEndian.Uint16(data[0:2])
		if signature == 0x5A4D || signature == 0x4D5A {
			return e.LoadEXEFile(filename)
		}
	}

	return e.LoadCOMFile(filename)
}

func (e *DOSEmulator) HandleInterrupt(intNum byte) {
	switch intNum {
	case 0x10:
		e.handleInt10()
	case 0x11:
		e.cpu.AX = 0x0021
	case 0x12:
		e.cpu.AX = 640
	case 0x13:
		e.handleInt13()
	case 0x16:
		e.handleInt16()
	case 0x1A:
		e.handleInt1A()
	case 0x20:
		e.running = false
	case 0x21:
		e.handleInt21()
	case 0x33:
		e.handleInt33()
	default:
		if e.debugMode {
			fmt.Printf("Unhandled interrupt: 0x%02X (AH=0x%02X)\n", intNum, e.cpu.GetAH())
		}
	}
}

func (e *DOSEmulator) handleInt10() {
	ah := e.cpu.GetAH()

	switch ah {
	case 0x00:
		mode := e.cpu.GetAL()
		e.video.videoMode = mode
	case 0x02:
		page := e.cpu.GetBH()
		row := e.cpu.GetDH()
		col := e.cpu.GetDL()
		if page == 0 {
			e.video.cursorY = int(row)
			e.video.cursorX = int(col)
		}
	case 0x03:
		e.cpu.SetDH(byte(e.video.cursorY))
		e.cpu.SetDL(byte(e.video.cursorX))
		e.cpu.SetCH(0)
		e.cpu.SetCL(7)
	case 0x06:
		lines := e.cpu.GetAL()
		attr := e.cpu.GetBH()
		if lines == 0 {
			e.clearScreen(attr)
		}
	case 0x09:
		char := e.cpu.GetAL()
		count := e.cpu.CX
		for i := uint16(0); i < count; i++ {
			fmt.Printf("%c", char)
		}
	case 0x0E:
		char := e.cpu.GetAL()
		e.teletypeOutput(char)
	case 0x0F:
		e.cpu.SetAL(e.video.videoMode)
		e.cpu.SetAH(80)
		e.cpu.SetBH(0)
	default:
		if e.debugMode {
			fmt.Printf("Unhandled INT 10h function: AH=0x%02X\n", ah)
		}
	}
}

func (e *DOSEmulator) teletypeOutput(char byte) {
	switch char {
	case '\r':
		e.video.cursorX = 0
	case '\n':
		e.video.cursorY++
		if e.video.cursorY >= 25 {
			e.scrollScreen()
			e.video.cursorY = 24
		}
	case '\b':
		if e.video.cursorX > 0 {
			e.video.cursorX--
		}
	case '\t':
		e.video.cursorX = (e.video.cursorX + 8) & ^7
		if e.video.cursorX >= 80 {
			e.video.cursorX = 0
			e.video.cursorY++
		}
	case 7:
		fmt.Print("\a")
	default:
		fmt.Printf("%c", char)
		e.video.cursorX++
		if e.video.cursorX >= 80 {
			e.video.cursorX = 0
			e.video.cursorY++
			if e.video.cursorY >= 25 {
				e.scrollScreen()
				e.video.cursorY = 24
			}
		}
	}
}

func (e *DOSEmulator) clearScreen(attr byte) {
	for i := 0; i < len(e.video.buffer); i += 2 {
		e.video.buffer[i] = ' '
		e.video.buffer[i+1] = attr
	}
	e.video.cursorX = 0
	e.video.cursorY = 0
}

func (e *DOSEmulator) scrollScreen() {
	copy(e.video.buffer[0:], e.video.buffer[160:])
	for i := 80 * 24 * 2; i < len(e.video.buffer); i += 2 {
		e.video.buffer[i] = ' '
		e.video.buffer[i+1] = e.video.currentColor
	}
}

func (e *DOSEmulator) handleInt13() {
	ah := e.cpu.GetAH()

	switch ah {
	case 0x00:
		e.cpu.SetAH(0)
		e.cpu.Flags.CF = false
	case 0x02:
		e.cpu.SetAH(0)
		e.cpu.SetAL(e.cpu.GetAL())
		e.cpu.Flags.CF = false
	case 0x08:
		e.cpu.SetAH(0)
		e.cpu.SetCH(79)
		e.cpu.SetCL(18)
		e.cpu.SetDH(1)
		e.cpu.SetDL(2)
		e.cpu.Flags.CF = false
	default:
		if e.debugMode {
			fmt.Printf("Unhandled INT 13h function: AH=0x%02X\n", ah)
		}
		e.cpu.Flags.CF = true
	}
}

func (e *DOSEmulator) handleInt16() {
	ah := e.cpu.GetAH()

	switch ah {
	case 0x00, 0x10:
		reader := bufio.NewReader(os.Stdin)
		char, _ := reader.ReadByte()
		e.cpu.SetAL(char)
		e.cpu.SetAH(0)
	case 0x01, 0x11:
		e.cpu.Flags.ZF = true
	case 0x02, 0x12:
		e.cpu.SetAL(0)
	default:
		if e.debugMode {
			fmt.Printf("Unhandled INT 16h function: AH=0x%02X\n", ah)
		}
	}
}

func (e *DOSEmulator) handleInt1A() {
	ah := e.cpu.GetAH()

	switch ah {
	case 0x00:
		now := time.Now()
		midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		ticks := uint32(now.Sub(midnight).Seconds() * 18.2)
		e.cpu.CX = uint16((ticks >> 16) & 0xFFFF)
		e.cpu.DX = uint16(ticks & 0xFFFF)
		e.cpu.SetAL(0)
	case 0x02:
		now := time.Now()
		e.cpu.SetCH(byte(now.Hour()))
		e.cpu.SetCL(byte(now.Minute()))
		e.cpu.SetDH(byte(now.Second()))
		e.cpu.Flags.CF = false
	case 0x04:
		now := time.Now()
		e.cpu.SetCH(byte(now.Year() / 100))
		e.cpu.SetCL(byte(now.Year() % 100))
		e.cpu.SetDH(byte(now.Month()))
		e.cpu.SetDL(byte(now.Day()))
		e.cpu.Flags.CF = false
	default:
		if e.debugMode {
			fmt.Printf("Unhandled INT 1Ah function: AH=0x%02X\n", ah)
		}
	}
}

func (e *DOSEmulator) handleInt21() {
	ah := e.cpu.GetAH()

	switch ah {
	case 0x01:
		reader := bufio.NewReader(os.Stdin)
		char, _ := reader.ReadByte()
		fmt.Printf("%c", char)
		e.cpu.SetAL(char)
	case 0x02:
		fmt.Printf("%c", e.cpu.GetDL())
	case 0x06:
		dl := e.cpu.GetDL()
		if dl == 0xFF {
			reader := bufio.NewReader(os.Stdin)
			char, err := reader.ReadByte()
			if err == nil {
				e.cpu.SetAL(char)
				e.cpu.Flags.ZF = false
			} else {
				e.cpu.Flags.ZF = true
			}
		} else {
			fmt.Printf("%c", dl)
		}
	case 0x07, 0x08:
		reader := bufio.NewReader(os.Stdin)
		char, _ := reader.ReadByte()
		e.cpu.SetAL(char)
	case 0x09:
		addr := CalculateAddress(e.cpu.DS, e.cpu.DX)
		for {
			ch := e.memory.ReadByte(addr)
			if ch == '$' {
				break
			}
			fmt.Printf("%c", ch)
			addr++
		}
	case 0x0A:
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimRight(input, "\r\n")
		addr := CalculateAddress(e.cpu.DS, e.cpu.DX)
		maxLen := e.memory.ReadByte(addr)
		if len(input) > int(maxLen) {
			input = input[:maxLen]
		}
		e.memory.WriteByte(addr+1, byte(len(input)))
		for i, ch := range input {
			e.memory.WriteByte(addr+2+uint32(i), byte(ch))
		}
	case 0x0E:
		e.fs.currentDrive = e.cpu.GetDL()
		e.cpu.SetAL(26)
	case 0x19:
		e.cpu.SetAL(e.fs.currentDrive)
	case 0x25:
		intNum := e.cpu.GetAL()
		offset := e.cpu.DX
		segment := e.cpu.DS
		e.interruptVectors[intNum] = CalculateAddress(segment, offset)
	case 0x2A:
		now := time.Now()
		e.cpu.CX = uint16(now.Year())
		e.cpu.SetDH(byte(now.Month()))
		e.cpu.SetDL(byte(now.Day()))
		e.cpu.SetAL(byte(now.Weekday()))
	case 0x2C:
		now := time.Now()
		e.cpu.SetCH(byte(now.Hour()))
		e.cpu.SetCL(byte(now.Minute()))
		e.cpu.SetDH(byte(now.Second()))
		e.cpu.SetDL(byte(now.Nanosecond() / 10000000))
	case 0x30:
		e.cpu.SetAL(5)
		e.cpu.SetAH(0)
		e.cpu.BX = 0
		e.cpu.CX = 0
	case 0x35:
		intNum := e.cpu.GetAL()
		addr := e.interruptVectors[intNum]
		e.cpu.BX = uint16(addr & 0xFFFF)
		e.cpu.ES = uint16(addr >> 16)
	case 0x39:
		addr := CalculateAddress(e.cpu.DS, e.cpu.DX)
		dirname := e.readNullTerminatedString(addr)
		err := os.Mkdir(dirname, 0755)
		if err != nil {
			e.cpu.Flags.CF = true
			e.cpu.AX = 3
		} else {
			e.cpu.Flags.CF = false
		}
	case 0x3A:
		addr := CalculateAddress(e.cpu.DS, e.cpu.DX)
		dirname := e.readNullTerminatedString(addr)
		err := os.Remove(dirname)
		if err != nil {
			e.cpu.Flags.CF = true
			e.cpu.AX = 3
		} else {
			e.cpu.Flags.CF = false
		}
	case 0x3B:
		addr := CalculateAddress(e.cpu.DS, e.cpu.DX)
		dirname := e.readNullTerminatedString(addr)
		err := os.Chdir(dirname)
		if err != nil {
			e.cpu.Flags.CF = true
			e.cpu.AX = 3
		} else {
			e.fs.currentDir, _ = os.Getwd()
			e.cpu.Flags.CF = false
		}
	case 0x3C:
		e.handleCreateFile()
	case 0x3D:
		e.handleOpenFile()
	case 0x3E:
		e.handleCloseFile()
	case 0x3F:
		e.handleReadFile()
	case 0x40:
		e.handleWriteFile()
	case 0x41:
		e.handleDeleteFile()
	case 0x42:
		e.handleSeekFile()
	case 0x43:
		e.handleFileAttributes()
	case 0x47:
		e.handleGetCurrentDir()
	case 0x4C:
		e.running = false
		exitCode := e.cpu.GetAL()
		if e.debugMode {
			fmt.Printf("\nProgram exited with code: %d\n", exitCode)
		}
	case 0x4E:
		e.handleFindFirst()
	case 0x4F:
		e.handleFindNext()
	case 0x51, 0x62:
		e.cpu.BX = e.psp
	case 0x56:
		e.handleRenameFile()
	default:
		if e.debugMode {
			fmt.Printf("Unhandled INT 21h function: AH=0x%02X\n", ah)
		}
	}
}

func (e *DOSEmulator) handleInt33() {
	ax := e.cpu.AX

	switch ax {
	case 0x00:
		e.cpu.AX = 0xFFFF
		e.cpu.BX = 2
	case 0x03:
		e.cpu.BX = 0
		e.cpu.CX = 0
		e.cpu.DX = 0
	default:
		if e.debugMode {
			fmt.Printf("Unhandled INT 33h function: AX=0x%04X\n", ax)
		}
	}
}

func (e *DOSEmulator) handleCreateFile() {
	addr := CalculateAddress(e.cpu.DS, e.cpu.DX)
	filename := e.readNullTerminatedString(addr)

	file, err := os.Create(filename)
	if err != nil {
		e.cpu.Flags.CF = true
		e.cpu.AX = 3
		return
	}

	handle := e.nextHandle
	e.nextHandle++
	e.fileHandles[handle] = &FileHandle{file: file, handle: handle}

	e.cpu.AX = handle
	e.cpu.Flags.CF = false
}

func (e *DOSEmulator) handleOpenFile() {
	addr := CalculateAddress(e.cpu.DS, e.cpu.DX)
	filename := e.readNullTerminatedString(addr)
	mode := e.cpu.GetAL()

	var file *os.File
	var err error

	switch mode & 0x03 {
	case 0:
		file, err = os.Open(filename)
	case 1:
		file, err = os.OpenFile(filename, os.O_WRONLY, 0644)
	case 2:
		file, err = os.OpenFile(filename, os.O_RDWR, 0644)
	}

	if err != nil {
		e.cpu.Flags.CF = true
		e.cpu.AX = 2
		return
	}

	handle := e.nextHandle
	e.nextHandle++
	e.fileHandles[handle] = &FileHandle{file: file, handle: handle}

	e.cpu.AX = handle
	e.cpu.Flags.CF = false
}

func (e *DOSEmulator) handleCloseFile() {
	handle := e.cpu.BX

	if handle <= 2 {
		e.cpu.Flags.CF = false
		return
	}

	if fh, ok := e.fileHandles[handle]; ok {
		fh.file.Close()
		delete(e.fileHandles, handle)
		e.cpu.Flags.CF = false
	} else {
		e.cpu.Flags.CF = true
		e.cpu.AX = 6
	}
}

func (e *DOSEmulator) handleReadFile() {
	handle := e.cpu.BX
	count := e.cpu.CX
	addr := CalculateAddress(e.cpu.DS, e.cpu.DX)

	if fh, ok := e.fileHandles[handle]; ok {
		buffer := make([]byte, count)
		n, _ := fh.file.Read(buffer)

		for i := 0; i < n; i++ {
			e.memory.WriteByte(addr+uint32(i), buffer[i])
		}

		e.cpu.AX = uint16(n)
		e.cpu.Flags.CF = false
	} else {
		e.cpu.Flags.CF = true
		e.cpu.AX = 6
	}
}

func (e *DOSEmulator) handleWriteFile() {
	handle := e.cpu.BX
	count := e.cpu.CX
	addr := CalculateAddress(e.cpu.DS, e.cpu.DX)

	if handle == 1 || handle == 2 {
		for i := uint16(0); i < count; i++ {
			ch := e.memory.ReadByte(addr + uint32(i))
			fmt.Printf("%c", ch)
		}
		e.cpu.AX = count
		e.cpu.Flags.CF = false
		return
	}

	if fh, ok := e.fileHandles[handle]; ok {
		buffer := make([]byte, count)
		for i := uint16(0); i < count; i++ {
			buffer[i] = e.memory.ReadByte(addr + uint32(i))
		}

		n, _ := fh.file.Write(buffer)
		e.cpu.AX = uint16(n)
		e.cpu.Flags.CF = false
	} else {
		e.cpu.Flags.CF = true
		e.cpu.AX = 6
	}
}

func (e *DOSEmulator) handleDeleteFile() {
	addr := CalculateAddress(e.cpu.DS, e.cpu.DX)
	filename := e.readNullTerminatedString(addr)

	err := os.Remove(filename)
	if err != nil {
		e.cpu.Flags.CF = true
		e.cpu.AX = 2
	} else {
		e.cpu.Flags.CF = false
	}
}

func (e *DOSEmulator) handleSeekFile() {
	handle := e.cpu.BX
	method := e.cpu.GetAL()
	offset := int64(uint32(e.cpu.CX)<<16 | uint32(e.cpu.DX))

	if fh, ok := e.fileHandles[handle]; ok {
		var whence int
		switch method {
		case 0:
			whence = io.SeekStart
		case 1:
			whence = io.SeekCurrent
		case 2:
			whence = io.SeekEnd
		}

		newPos, err := fh.file.Seek(offset, whence)
		if err != nil {
			e.cpu.Flags.CF = true
			e.cpu.AX = 1
		} else {
			e.cpu.DX = uint16((newPos >> 16) & 0xFFFF)
			e.cpu.AX = uint16(newPos & 0xFFFF)
			e.cpu.Flags.CF = false
		}
	} else {
		e.cpu.Flags.CF = true
		e.cpu.AX = 6
	}
}

func (e *DOSEmulator) handleFileAttributes() {
	al := e.cpu.GetAL()
	addr := CalculateAddress(e.cpu.DS, e.cpu.DX)
	filename := e.readNullTerminatedString(addr)

	if al == 0 {
		info, err := os.Stat(filename)
		if err != nil {
			e.cpu.Flags.CF = true
			e.cpu.AX = 2
		} else {
			attr := uint16(0)
			if info.IsDir() {
				attr = attr | 0x10
			}
			if info.Mode().Perm() & 0200 == 0 {
				attr = attr | 0x01
			}
			e.cpu.CX = attr
			e.cpu.Flags.CF = false
		}
	} else {
		e.cpu.Flags.CF = false
	}
}

func (e *DOSEmulator) handleGetCurrentDir() {
	drive := e.cpu.GetDL()
	addr := CalculateAddress(e.cpu.DS, e.cpu.SI)

	currentDir := e.fs.currentDir
	if drive != 0 {
		if path, ok := e.fs.drives[drive-1]; ok {
			currentDir = path
		}
	}

	currentDir = strings.TrimPrefix(currentDir, e.fs.drives[e.fs.currentDrive])
	currentDir = strings.TrimPrefix(currentDir, string(filepath.Separator))

	for i, ch := range currentDir {
		e.memory.WriteByte(addr+uint32(i), byte(ch))
	}
	e.memory.WriteByte(addr+uint32(len(currentDir)), 0)

	e.cpu.Flags.CF = false
}

func (e *DOSEmulator) handleFindFirst() {
	addr := CalculateAddress(e.cpu.DS, e.cpu.DX)
	pattern := e.readNullTerminatedString(addr)

	dir := filepath.Dir(pattern)
	if dir == "." {
		dir, _ = os.Getwd()
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		e.cpu.Flags.CF = true
		e.cpu.AX = 2
		return
	}

	e.dta.searchPath = pattern
	e.dta.searchIndex = 0
	e.dta.searchFiles = files

	e.handleFindNext()
}

func (e *DOSEmulator) handleFindNext() {
	if e.dta.searchIndex >= len(e.dta.searchFiles) {
		e.cpu.Flags.CF = true
		e.cpu.AX = 18
		return
	}

	file := e.dta.searchFiles[e.dta.searchIndex]
	e.dta.searchIndex++

	info, _ := file.Info()

	e.dta.attribute = 0
	if file.IsDir() {
		e.dta.attribute = e.dta.attribute | 0x10
	}

	modTime := info.ModTime()
	e.dta.time = uint16((modTime.Hour() << 11) | (modTime.Minute() << 5) | (modTime.Second() / 2))
	e.dta.date = uint16(((modTime.Year() - 1980) << 9) | (int(modTime.Month()) << 5) | modTime.Day())
	e.dta.size = uint32(info.Size())

	name := file.Name()
	if len(name) > 12 {
		name = name[:12]
	}
	copy(e.dta.name[:], name)

	e.cpu.Flags.CF = false
}

func (e *DOSEmulator) handleRenameFile() {
	addr1 := CalculateAddress(e.cpu.DS, e.cpu.DX)
	addr2 := CalculateAddress(e.cpu.ES, e.cpu.DI)
	oldName := e.readNullTerminatedString(addr1)
	newName := e.readNullTerminatedString(addr2)

	err := os.Rename(oldName, newName)
	if err != nil {
		e.cpu.Flags.CF = true
		e.cpu.AX = 2
	} else {
		e.cpu.Flags.CF = false
	}
}

func (e *DOSEmulator) readNullTerminatedString(addr uint32) string {
	var result []byte
	for {
		ch := e.memory.ReadByte(addr)
		if ch == 0 {
			break
		}
		result = append(result, ch)
		addr++
		if len(result) > 256 {
			break
		}
	}
	return string(result)
}

func (e *DOSEmulator) Execute(inst *Instruction) {
	switch inst.Opcode {
	case 0x90:
		e.cpu.IP++
	case 0xB0:
		e.cpu.SetAL(byte(inst.Operand1))
		e.cpu.IP += uint16(inst.Length)
	case 0xB1:
		e.cpu.SetCL(byte(inst.Operand1))
		e.cpu.IP += uint16(inst.Length)
	case 0xB2:
		e.cpu.SetDL(byte(inst.Operand1))
		e.cpu.IP += uint16(inst.Length)
	case 0xB3:
		e.cpu.SetBL(byte(inst.Operand1))
		e.cpu.IP += uint16(inst.Length)
	case 0xB4:
		e.cpu.SetAH(byte(inst.Operand1))
		e.cpu.IP += uint16(inst.Length)
	case 0xB5:
		e.cpu.SetCH(byte(inst.Operand1))
		e.cpu.IP += uint16(inst.Length)
	case 0xB6:
		e.cpu.SetDH(byte(inst.Operand1))
		e.cpu.IP += uint16(inst.Length)
	case 0xB7:
		e.cpu.SetBH(byte(inst.Operand1))
		e.cpu.IP += uint16(inst.Length)
	case 0xB8:
		e.cpu.AX = inst.Operand1
		e.cpu.IP += uint16(inst.Length)
	case 0xB9:
		e.cpu.CX = inst.Operand1
		e.cpu.IP += uint16(inst.Length)
	case 0xBA:
		e.cpu.DX = inst.Operand1
		e.cpu.IP += uint16(inst.Length)
	case 0xBB:
		e.cpu.BX = inst.Operand1
		e.cpu.IP += uint16(inst.Length)
	case 0xBC:
		e.cpu.SP = inst.Operand1
		e.cpu.IP += uint16(inst.Length)
	case 0xBD:
		e.cpu.BP = inst.Operand1
		e.cpu.IP += uint16(inst.Length)
	case 0xBE:
		e.cpu.SI = inst.Operand1
		e.cpu.IP += uint16(inst.Length)
	case 0xBF:
		e.cpu.DI = inst.Operand1
		e.cpu.IP += uint16(inst.Length)
	case 0x50:
		e.Push(e.cpu.AX)
		e.cpu.IP++
	case 0x51:
		e.Push(e.cpu.CX)
		e.cpu.IP++
	case 0x52:
		e.Push(e.cpu.DX)
		e.cpu.IP++
	case 0x53:
		e.Push(e.cpu.BX)
		e.cpu.IP++
	case 0x54:
		e.Push(e.cpu.SP)
		e.cpu.IP++
	case 0x55:
		e.Push(e.cpu.BP)
		e.cpu.IP++
	case 0x56:
		e.Push(e.cpu.SI)
		e.cpu.IP++
	case 0x57:
		e.Push(e.cpu.DI)
		e.cpu.IP++
	case 0x06:
		e.Push(e.cpu.ES)
		e.cpu.IP++
	case 0x0E:
		e.Push(e.cpu.CS)
		e.cpu.IP++
	case 0x16:
		e.Push(e.cpu.SS)
		e.cpu.IP++
	case 0x1E:
		e.Push(e.cpu.DS)
		e.cpu.IP++
	case 0x58:
		e.cpu.AX = e.Pop()
		e.cpu.IP++
	case 0x59:
		e.cpu.CX = e.Pop()
		e.cpu.IP++
	case 0x5A:
		e.cpu.DX = e.Pop()
		e.cpu.IP++
	case 0x5B:
		e.cpu.BX = e.Pop()
		e.cpu.IP++
	case 0x5C:
		e.cpu.SP = e.Pop()
		e.cpu.IP++
	case 0x5D:
		e.cpu.BP = e.Pop()
		e.cpu.IP++
	case 0x5E:
		e.cpu.SI = e.Pop()
		e.cpu.IP++
	case 0x5F:
		e.cpu.DI = e.Pop()
		e.cpu.IP++
	case 0x07:
		e.cpu.ES = e.Pop()
		e.cpu.IP++
	case 0x17:
		e.cpu.SS = e.Pop()
		e.cpu.IP++
	case 0x1F:
		e.cpu.DS = e.Pop()
		e.cpu.IP++
	case 0x9C:
		e.Push(e.cpu.Flags.ToUint16())
		e.cpu.IP++
	case 0x9D:
		e.cpu.Flags.FromUint16(e.Pop())
		e.cpu.IP++
	case 0x40:
		e.cpu.AX++
		e.cpu.UpdateArithmeticFlags16(e.cpu.AX)
		e.cpu.IP++
	case 0x41:
		e.cpu.CX++
		e.cpu.UpdateArithmeticFlags16(e.cpu.CX)
		e.cpu.IP++
	case 0x42:
		e.cpu.DX++
		e.cpu.UpdateArithmeticFlags16(e.cpu.DX)
		e.cpu.IP++
	case 0x43:
		e.cpu.BX++
		e.cpu.UpdateArithmeticFlags16(e.cpu.BX)
		e.cpu.IP++
	case 0x44:
		e.cpu.SP++
		e.cpu.IP++
	case 0x45:
		e.cpu.BP++
		e.cpu.IP++
	case 0x46:
		e.cpu.SI++
		e.cpu.IP++
	case 0x47:
		e.cpu.DI++
		e.cpu.IP++
	case 0x48:
		e.cpu.AX--
		e.cpu.UpdateArithmeticFlags16(e.cpu.AX)
		e.cpu.IP++
	case 0x49:
		e.cpu.CX--
		e.cpu.UpdateArithmeticFlags16(e.cpu.CX)
		e.cpu.IP++
	case 0x4A:
		e.cpu.DX--
		e.cpu.UpdateArithmeticFlags16(e.cpu.DX)
		e.cpu.IP++
	case 0x4B:
		e.cpu.BX--
		e.cpu.UpdateArithmeticFlags16(e.cpu.BX)
		e.cpu.IP++
	case 0x4C:
		e.cpu.SP--
		e.cpu.IP++
	case 0x4D:
		e.cpu.BP--
		e.cpu.IP++
	case 0x4E:
		e.cpu.SI--
		e.cpu.IP++
	case 0x4F:
		e.cpu.DI--
		e.cpu.IP++
	case 0x04:
		result := uint16(e.cpu.GetAL()) + inst.Immediate
		e.cpu.Flags.CF = result > 0xFF
		e.cpu.SetAL(byte(result))
		e.cpu.UpdateArithmeticFlags8(byte(result))
		e.cpu.IP += uint16(inst.Length)
	case 0x05:
		result := uint32(e.cpu.AX) + uint32(inst.Immediate)
		e.cpu.Flags.CF = result > 0xFFFF
		e.cpu.AX = uint16(result)
		e.cpu.UpdateArithmeticFlags16(e.cpu.AX)
		e.cpu.IP += uint16(inst.Length)
	case 0x2C:
		result := int16(e.cpu.GetAL()) - int16(inst.Immediate)
		e.cpu.Flags.CF = result < 0
		e.cpu.SetAL(byte(result))
		e.cpu.UpdateArithmeticFlags8(byte(result))
		e.cpu.IP += uint16(inst.Length)
	case 0x2D:
		result := int32(e.cpu.AX) - int32(inst.Immediate)
		e.cpu.Flags.CF = result < 0
		e.cpu.AX = uint16(result)
		e.cpu.UpdateArithmeticFlags16(e.cpu.AX)
		e.cpu.IP += uint16(inst.Length)
	case 0x3C:
		result := int16(e.cpu.GetAL()) - int16(inst.Immediate)
		e.cpu.Flags.CF = result < 0
		e.cpu.UpdateZeroFlag8(byte(result))
		e.cpu.UpdateSignFlag8(byte(result))
		e.cpu.IP += uint16(inst.Length)
	case 0x3D:
		result := int32(e.cpu.AX) - int32(inst.Immediate)
		e.cpu.Flags.CF = result < 0
		e.cpu.UpdateZeroFlag(uint16(result))
		e.cpu.UpdateSignFlag(uint16(result))
		e.cpu.IP += uint16(inst.Length)
	case 0x24:
		e.cpu.SetAL(e.cpu.GetAL() & byte(inst.Immediate))
		e.cpu.Flags.CF = false
		e.cpu.Flags.OF = false
		e.cpu.UpdateArithmeticFlags8(e.cpu.GetAL())
		e.cpu.IP += uint16(inst.Length)
	case 0x25:
		e.cpu.AX = e.cpu.AX & inst.Immediate
		e.cpu.Flags.CF = false
		e.cpu.Flags.OF = false
		e.cpu.UpdateArithmeticFlags16(e.cpu.AX)
		e.cpu.IP += uint16(inst.Length)
	case 0x0C:
		e.cpu.SetAL(e.cpu.GetAL() | byte(inst.Immediate))
		e.cpu.Flags.CF = false
		e.cpu.Flags.OF = false
		e.cpu.UpdateArithmeticFlags8(e.cpu.GetAL())
		e.cpu.IP += uint16(inst.Length)
	case 0x0D:
		e.cpu.AX = e.cpu.AX | inst.Immediate
		e.cpu.Flags.CF = false
		e.cpu.Flags.OF = false
		e.cpu.UpdateArithmeticFlags16(e.cpu.AX)
		e.cpu.IP += uint16(inst.Length)
	case 0x34:
		e.cpu.SetAL(e.cpu.GetAL() ^ byte(inst.Immediate))
		e.cpu.Flags.CF = false
		e.cpu.Flags.OF = false
		e.cpu.UpdateArithmeticFlags8(e.cpu.GetAL())
		e.cpu.IP += uint16(inst.Length)
	case 0x35:
		e.cpu.AX = e.cpu.AX ^ inst.Immediate
		e.cpu.Flags.CF = false
		e.cpu.Flags.OF = false
		e.cpu.UpdateArithmeticFlags16(e.cpu.AX)
		e.cpu.IP += uint16(inst.Length)
	case 0xA4:
		srcAddr := CalculateAddress(e.cpu.DS, e.cpu.SI)
		dstAddr := CalculateAddress(e.cpu.ES, e.cpu.DI)
		e.memory.WriteByte(dstAddr, e.memory.ReadByte(srcAddr))
		if e.cpu.Flags.DF {
			e.cpu.SI--
			e.cpu.DI--
		} else {
			e.cpu.SI++
			e.cpu.DI++
		}
		e.cpu.IP++
	case 0xA5:
		srcAddr := CalculateAddress(e.cpu.DS, e.cpu.SI)
		dstAddr := CalculateAddress(e.cpu.ES, e.cpu.DI)
		e.memory.WriteWord(dstAddr, e.memory.ReadWord(srcAddr))
		if e.cpu.Flags.DF {
			e.cpu.SI -= 2
			e.cpu.DI -= 2
		} else {
			e.cpu.SI += 2
			e.cpu.DI += 2
		}
		e.cpu.IP++
	case 0xAA:
		dstAddr := CalculateAddress(e.cpu.ES, e.cpu.DI)
		e.memory.WriteByte(dstAddr, e.cpu.GetAL())
		if e.cpu.Flags.DF {
			e.cpu.DI--
		} else {
			e.cpu.DI++
		}
		e.cpu.IP++
	case 0xAB:
		dstAddr := CalculateAddress(e.cpu.ES, e.cpu.DI)
		e.memory.WriteWord(dstAddr, e.cpu.AX)
		if e.cpu.Flags.DF {
			e.cpu.DI -= 2
		} else {
			e.cpu.DI += 2
		}
		e.cpu.IP++
	case 0xAC:
		srcAddr := CalculateAddress(e.cpu.DS, e.cpu.SI)
		e.cpu.SetAL(e.memory.ReadByte(srcAddr))
		if e.cpu.Flags.DF {
			e.cpu.SI--
		} else {
			e.cpu.SI++
		}
		e.cpu.IP++
	case 0xAD:
		srcAddr := CalculateAddress(e.cpu.DS, e.cpu.SI)
		e.cpu.AX = e.memory.ReadWord(srcAddr)
		if e.cpu.Flags.DF {
			e.cpu.SI -= 2
		} else {
			e.cpu.SI += 2
		}
		e.cpu.IP++
	case 0xEB:
		offset := int8(inst.Operand1)
		e.cpu.IP = uint16(int32(e.cpu.IP) + int32(offset) + int32(inst.Length))
	case 0xE9:
		offset := int16(inst.Operand1)
		e.cpu.IP = uint16(int32(e.cpu.IP) + int32(offset) + int32(inst.Length))
	case 0xEA:
		e.cpu.IP = inst.Operand1
		e.cpu.CS = inst.Operand2
	case 0x74, 0x75:
		if (inst.Opcode == 0x74 && e.cpu.Flags.ZF) || (inst.Opcode == 0x75 && !e.cpu.Flags.ZF) {
			offset := int8(inst.Operand1)
			e.cpu.IP = uint16(int32(e.cpu.IP) + int32(offset) + int32(inst.Length))
		} else {
			e.cpu.IP += uint16(inst.Length)
		}
	case 0x72, 0x73:
		if (inst.Opcode == 0x72 && e.cpu.Flags.CF) || (inst.Opcode == 0x73 && !e.cpu.Flags.CF) {
			offset := int8(inst.Operand1)
			e.cpu.IP = uint16(int32(e.cpu.IP) + int32(offset) + int32(inst.Length))
		} else {
			e.cpu.IP += uint16(inst.Length)
		}
	case 0x76, 0x77:
		if (inst.Opcode == 0x76 && (e.cpu.Flags.CF || e.cpu.Flags.ZF)) || (inst.Opcode == 0x77 && !(e.cpu.Flags.CF || e.cpu.Flags.ZF)) {
			offset := int8(inst.Operand1)
			e.cpu.IP = uint16(int32(e.cpu.IP) + int32(offset) + int32(inst.Length))
		} else {
			e.cpu.IP += uint16(inst.Length)
		}
	case 0x78, 0x79:
		if (inst.Opcode == 0x78 && e.cpu.Flags.SF) || (inst.Opcode == 0x79 && !e.cpu.Flags.SF) {
			offset := int8(inst.Operand1)
			e.cpu.IP = uint16(int32(e.cpu.IP) + int32(offset) + int32(inst.Length))
		} else {
			e.cpu.IP += uint16(inst.Length)
		}
	case 0x7C, 0x7D:
		sfNeOF := e.cpu.Flags.SF != e.cpu.Flags.OF
		if (inst.Opcode == 0x7C && sfNeOF) || (inst.Opcode == 0x7D && !sfNeOF) {
			offset := int8(inst.Operand1)
			e.cpu.IP = uint16(int32(e.cpu.IP) + int32(offset) + int32(inst.Length))
		} else {
			e.cpu.IP += uint16(inst.Length)
		}
	case 0x7E, 0x7F:
		sfNeOF := e.cpu.Flags.SF != e.cpu.Flags.OF
		if (inst.Opcode == 0x7E && (e.cpu.Flags.ZF || sfNeOF)) || (inst.Opcode == 0x7F && !(e.cpu.Flags.ZF || sfNeOF)) {
			offset := int8(inst.Operand1)
			e.cpu.IP = uint16(int32(e.cpu.IP) + int32(offset) + int32(inst.Length))
		} else {
			e.cpu.IP += uint16(inst.Length)
		}
	case 0x70, 0x71:
		if (inst.Opcode == 0x70 && e.cpu.Flags.OF) || (inst.Opcode == 0x71 && !e.cpu.Flags.OF) {
			offset := int8(inst.Operand1)
			e.cpu.IP = uint16(int32(e.cpu.IP) + int32(offset) + int32(inst.Length))
		} else {
			e.cpu.IP += uint16(inst.Length)
		}
	case 0x7A, 0x7B:
		if (inst.Opcode == 0x7A && e.cpu.Flags.PF) || (inst.Opcode == 0x7B && !e.cpu.Flags.PF) {
			offset := int8(inst.Operand1)
			e.cpu.IP = uint16(int32(e.cpu.IP) + int32(offset) + int32(inst.Length))
		} else {
			e.cpu.IP += uint16(inst.Length)
		}
	case 0xE0:
		e.cpu.CX--
		if e.cpu.CX != 0 && !e.cpu.Flags.ZF {
			offset := int8(inst.Operand1)
			e.cpu.IP = uint16(int32(e.cpu.IP) + int32(offset) + int32(inst.Length))
		} else {
			e.cpu.IP += uint16(inst.Length)
		}
	case 0xE1:
		e.cpu.CX--
		if e.cpu.CX != 0 && e.cpu.Flags.ZF {
			offset := int8(inst.Operand1)
			e.cpu.IP = uint16(int32(e.cpu.IP) + int32(offset) + int32(inst.Length))
		} else {
			e.cpu.IP += uint16(inst.Length)
		}
	case 0xE2:
		e.cpu.CX--
		if e.cpu.CX != 0 {
			offset := int8(inst.Operand1)
			e.cpu.IP = uint16(int32(e.cpu.IP) + int32(offset) + int32(inst.Length))
		} else {
			e.cpu.IP += uint16(inst.Length)
		}
	case 0xE3:
		if e.cpu.CX == 0 {
			offset := int8(inst.Operand1)
			e.cpu.IP = uint16(int32(e.cpu.IP) + int32(offset) + int32(inst.Length))
		} else {
			e.cpu.IP += uint16(inst.Length)
		}
	case 0xE8:
		e.Push(e.cpu.IP + uint16(inst.Length))
		offset := int16(inst.Operand1)
		e.cpu.IP = uint16(int32(e.cpu.IP) + int32(offset) + int32(inst.Length))
	case 0x9A:
		e.Push(e.cpu.CS)
		e.Push(e.cpu.IP + uint16(inst.Length))
		e.cpu.IP = inst.Operand1
		e.cpu.CS = inst.Operand2
	case 0xC3:
		e.cpu.IP = e.Pop()
	case 0xC2:
		e.cpu.IP = e.Pop()
		e.cpu.SP += inst.Operand1
	case 0xCB:
		e.cpu.IP = e.Pop()
		e.cpu.CS = e.Pop()
	case 0xCA:
		e.cpu.IP = e.Pop()
		e.cpu.CS = e.Pop()
		e.cpu.SP += inst.Operand1
	case 0xCD:
		e.HandleInterrupt(byte(inst.Operand1))
		e.cpu.IP += uint16(inst.Length)
	case 0xCC:
		e.HandleInterrupt(3)
		e.cpu.IP++
	case 0xCE:
		if e.cpu.Flags.OF {
			e.HandleInterrupt(4)
		}
		e.cpu.IP++
	case 0xCF:
		e.cpu.IP = e.Pop()
		e.cpu.CS = e.Pop()
		e.cpu.Flags.FromUint16(e.Pop())
	case 0xF8:
		e.cpu.Flags.CF = false
		e.cpu.IP++
	case 0xF9:
		e.cpu.Flags.CF = true
		e.cpu.IP++
	case 0xFA:
		e.cpu.Flags.IF = false
		e.cpu.IP++
	case 0xFB:
		e.cpu.Flags.IF = true
		e.cpu.IP++
	case 0xFC:
		e.cpu.Flags.DF = false
		e.cpu.IP++
	case 0xFD:
		e.cpu.Flags.DF = true
		e.cpu.IP++
	case 0xF5:
		e.cpu.Flags.CF = !e.cpu.Flags.CF
		e.cpu.IP++
	case 0x98:
		if (e.cpu.GetAL() & 0x80) != 0 {
			e.cpu.SetAH(0xFF)
		} else {
			e.cpu.SetAH(0x00)
		}
		e.cpu.IP++
	case 0x99:
		if (e.cpu.AX & 0x8000) != 0 {
			e.cpu.DX = 0xFFFF
		} else {
			e.cpu.DX = 0x0000
		}
		e.cpu.IP++
	case 0x9E:
		flags := e.cpu.GetAH()
		e.cpu.Flags.CF = (flags & 0x01) != 0
		e.cpu.Flags.PF = (flags & 0x04) != 0
		e.cpu.Flags.AF = (flags & 0x10) != 0
		e.cpu.Flags.ZF = (flags & 0x40) != 0
		e.cpu.Flags.SF = (flags & 0x80) != 0
		e.cpu.IP++
	case 0x9F:
		flags := byte(0x02)
		if e.cpu.Flags.CF {
			flags = flags | 0x01
		}
		if e.cpu.Flags.PF {
			flags = flags | 0x04
		}
		if e.cpu.Flags.AF {
			flags = flags | 0x10
		}
		if e.cpu.Flags.ZF {
			flags = flags | 0x40
		}
		if e.cpu.Flags.SF {
			flags = flags | 0x80
		}
		e.cpu.SetAH(flags)
		e.cpu.IP++
	case 0xD4:
		base := byte(inst.Operand1)
		if base == 0 {
			base = 10
		}
		al := e.cpu.GetAL()
		e.cpu.SetAH(al / base)
		e.cpu.SetAL(al % base)
		e.cpu.UpdateArithmeticFlags16(e.cpu.AX)
		e.cpu.IP += uint16(inst.Length)
	case 0xD5:
		base := byte(inst.Operand1)
		if base == 0 {
			base = 10
		}
		al := e.cpu.GetAL() + e.cpu.GetAH()*base
		e.cpu.SetAL(al)
		e.cpu.SetAH(0)
		e.cpu.UpdateArithmeticFlags8(al)
		e.cpu.IP += uint16(inst.Length)
	case 0xD7:
		addr := CalculateAddress(e.cpu.DS, e.cpu.BX+uint16(e.cpu.GetAL()))
		e.cpu.SetAL(e.memory.ReadByte(addr))
		e.cpu.IP++
	case 0xF4:
		e.running = false
		if e.debugMode {
			fmt.Println("CPU halted")
		}
	case 0xF2, 0xF3:
		e.repeatPrefix = inst.Opcode
		e.cpu.IP++
	default:
		if e.debugMode {
			fmt.Printf("Unimplemented opcode: 0x%02X at %04X:%04X\n", inst.Opcode, e.cpu.CS, e.cpu.IP)
		}
		e.cpu.IP += uint16(inst.Length)
	}

	if e.repeatPrefix != 0 && e.cpu.CX > 0 {
		nextAddr := CalculateAddress(e.cpu.CS, e.cpu.IP)
		nextOpcode := e.memory.ReadByte(nextAddr)

		if nextOpcode >= 0xA4 && nextOpcode <= 0xAF {
			e.cpu.IP -= uint16(inst.Length)
			e.cpu.CX--

			if e.cpu.CX == 0 {
				e.repeatPrefix = 0
				e.cpu.IP += uint16(inst.Length)
			}
		} else {
			e.repeatPrefix = 0
		}
	}
}

func (e *DOSEmulator) Run() {
	if !e.debugMode {
		fmt.Printf("Running %s program...\n", e.programType)
	}
	e.running = true
	maxInstructions := uint64(100000000)

	for e.running && e.instructionCount < maxInstructions {
		addr := CalculateAddress(e.cpu.CS, e.cpu.IP)
		inst := e.decoder.Decode(addr)

		if e.debugMode || e.traceMode {
			fmt.Printf("%04X:%04X  %-30s  AX=%04X BX=%04X CX=%04X DX=%04X\n",
				e.cpu.CS, e.cpu.IP, inst.Name,
				e.cpu.AX, e.cpu.BX, e.cpu.CX, e.cpu.DX)
		}

		if e.stepMode {
			fmt.Print("Press Enter (c=continue, q=quit)> ")
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			if input == "c" {
				e.stepMode = false
			} else if input == "q" {
				e.running = false
				return
			}
		}

		e.Execute(inst)
		e.instructionCount++

		if e.instructionCount%100000 == 0 && !e.debugMode {
			fmt.Print(".")
		}
	}

	if e.instructionCount >= maxInstructions {
		fmt.Println("\nMaximum instruction count reached")
	}

	if !e.debugMode {
		fmt.Println()
	}
}

func (e *DOSEmulator) SimpleShell() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("MS-DOS Emulator v5.0 - Full COM & EXE Support")
	fmt.Println("Full 8086 CPU + BIOS + DOS + 150+ Instructions")
	fmt.Println()
	fmt.Println("Type 'HELP' for available commands")
	fmt.Println()

	for {
		driveLetter := string(rune('A' + e.fs.currentDrive))
		fmt.Printf("%s:\\> ", driveLetter)

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		command := strings.ToUpper(parts[0])

		switch command {
		case "HELP", "?":
			e.showHelp()
		case "CLS":
			fmt.Print("\033[H\033[2J")
		case "VER":
			fmt.Println("MS-DOS Emulator Version 5.0 - Complete COM & EXE Support")
			fmt.Printf("Instructions executed: %d\n", e.instructionCount)
		case "DIR":
			e.listDirectory()
		case "CD":
			e.changeDirectory(parts)
		case "MD", "MKDIR":
			e.makeDirectory(parts)
		case "RD", "RMDIR":
			e.removeDirectory(parts)
		case "DEL", "ERASE":
			e.deleteFile(parts)
		case "TYPE":
			e.typeFile(parts)
		case "COPY":
			e.copyFile(parts)
		case "REN", "RENAME":
			e.renameFile(parts)
		case "ECHO":
			if len(parts) > 1 {
				fmt.Println(strings.Join(parts[1:], " "))
			}
		case "DATE":
			e.showDate()
		case "TIME":
			e.showTime()
		case "MEM":
			e.showMemoryInfo()
		case "REGS":
			e.showRegisters()
		case "DEBUG":
			e.debugMode = !e.debugMode
			fmt.Printf("Debug mode: %v\n", e.debugMode)
		case "STEP":
			e.stepMode = !e.stepMode
			fmt.Printf("Step mode: %v\n", e.stepMode)
		case "TRACE":
			e.traceMode = !e.traceMode
			fmt.Printf("Trace mode: %v\n", e.traceMode)
		case "DUMP":
			e.dumpMemory(parts)
		case "STACK":
			e.showStack()
		case "STATS":
			e.showStatistics()
		case "DISASM":
			e.disassemble(parts)
		case "RUN", "EXEC":
			if len(parts) < 2 {
				fmt.Println("Usage: RUN <filename>")
				continue
			}
			if err := e.LoadFile(parts[1]); err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				e.Run()
			}
		case "EXIT", "QUIT":
			fmt.Println("Exiting emulator...")
			return
		default:
			ext := strings.ToUpper(filepath.Ext(command))
			if ext == ".COM" || ext == ".EXE" {
				if err := e.LoadFile(command); err != nil {
					fmt.Printf("Bad command or file name: %s\n", command)
				} else {
					e.Run()
				}
			} else {
				fmt.Printf("Bad command or file name: %s\n", command)
			}
		}
	}
}

func (e *DOSEmulator) showHelp() {
	fmt.Println("\nAVAILABLE COMMANDS:")
	fmt.Println("File: DIR, CD, MD, RD, DEL, TYPE, COPY, REN")
	fmt.Println("System: CLS, VER, DATE, TIME, MEM, ECHO")
	fmt.Println("Emulator: RUN, DEBUG, STEP, TRACE, REGS, DUMP, STACK, STATS, DISASM, EXIT")
	fmt.Println("Supports: .COM and .EXE files\n")
}

func (e *DOSEmulator) listDirectory() {
	files, err := os.ReadDir(".")
	if err != nil {
		fmt.Println("Error reading directory")
		return
	}

	fmt.Println("\n Volume in drive A is EMULATOR")
	fmt.Println(" Directory of A:\\\n")

	fileCount := 0
	dirCount := 0
	totalSize := int64(0)

	for _, file := range files {
		info, _ := file.Info()

		if file.IsDir() {
			fmt.Printf("%-12s <DIR>         %s\n",
				file.Name(),
				info.ModTime().Format("01-02-06  03:04p"))
			dirCount++
		} else {
			fmt.Printf("%-12s %10d %s\n",
				file.Name(),
				info.Size(),
				info.ModTime().Format("01-02-06  03:04p"))
			fileCount++
			totalSize += info.Size()
		}
	}

	fmt.Printf("\n    %d File(s) %d bytes\n", fileCount, totalSize)
	fmt.Printf("    %d Dir(s)\n\n", dirCount)
}

func (e *DOSEmulator) changeDirectory(parts []string) {
	if len(parts) < 2 {
		fmt.Println(e.fs.currentDir)
		return
	}

	err := os.Chdir(parts[1])
	if err != nil {
		fmt.Println("Invalid directory")
	} else {
		e.fs.currentDir, _ = os.Getwd()
	}
}

func (e *DOSEmulator) makeDirectory(parts []string) {
	if len(parts) < 2 {
		fmt.Println("Usage: MD <directory>")
		return
	}
	err := os.Mkdir(parts[1], 0755)
	if err != nil {
		fmt.Println("Unable to create directory")
	}
}

func (e *DOSEmulator) removeDirectory(parts []string) {
	if len(parts) < 2 {
		fmt.Println("Usage: RD <directory>")
		return
	}
	err := os.Remove(parts[1])
	if err != nil {
		fmt.Println("Unable to remove directory")
	}
}

func (e *DOSEmulator) deleteFile(parts []string) {
	if len(parts) < 2 {
		fmt.Println("Usage: DEL <filename>")
		return
	}
	err := os.Remove(parts[1])
	if err != nil {
		fmt.Println("File not found")
	}
}

func (e *DOSEmulator) typeFile(parts []string) {
	if len(parts) < 2 {
		fmt.Println("Usage: TYPE <filename>")
		return
	}
	content, err := os.ReadFile(parts[1])
	if err != nil {
		fmt.Println("File not found")
		return
	}
	fmt.Print(string(content))
}

func (e *DOSEmulator) copyFile(parts []string) {
	if len(parts) < 3 {
		fmt.Println("Usage: COPY <source> <destination>")
		return
	}
	source, err := os.ReadFile(parts[1])
	if err != nil {
		fmt.Println("File not found")
		return
	}
	err = os.WriteFile(parts[2], source, 0644)
	if err != nil {
		fmt.Println("Unable to copy file")
		return
	}
	fmt.Println("        1 file(s) copied")
}

func (e *DOSEmulator) renameFile(parts []string) {
	if len(parts) < 3 {
		fmt.Println("Usage: REN <oldname> <newname>")
		return
	}
	err := os.Rename(parts[1], parts[2])
	if err != nil {
		fmt.Println("Unable to rename file")
	}
}

func (e *DOSEmulator) showDate() {
	fmt.Printf("Current date: %s\n", time.Now().Format("Mon 01/02/2006"))
}

func (e *DOSEmulator) showTime() {
	fmt.Printf("Current time: %s\n", time.Now().Format("15:04:05"))
}

func (e *DOSEmulator) showMemoryInfo() {
	fmt.Println("\nMemory Type        Total       Used       Free")
	fmt.Println("Conventional       640K        128K       512K")
	fmt.Println("Extended          1024K          0K      1024K")
	fmt.Println()
}

func (e *DOSEmulator) showRegisters() {
	fmt.Println("\nCPU REGISTERS:")
	fmt.Printf("AX=%04X  BX=%04X  CX=%04X  DX=%04X\n",
		e.cpu.AX, e.cpu.BX, e.cpu.CX, e.cpu.DX)
	fmt.Printf("SI=%04X  DI=%04X  BP=%04X  SP=%04X\n",
		e.cpu.SI, e.cpu.DI, e.cpu.BP, e.cpu.SP)
	fmt.Printf("CS=%04X  DS=%04X  ES=%04X  SS=%04X\n",
		e.cpu.CS, e.cpu.DS, e.cpu.ES, e.cpu.SS)
	fmt.Printf("IP=%04X  FLAGS=%04X\n", e.cpu.IP, e.cpu.Flags.ToUint16())

	flags := ""
	if e.cpu.Flags.CF {
		flags += "CF "
	}
	if e.cpu.Flags.PF {
		flags += "PF "
	}
	if e.cpu.Flags.AF {
		flags += "AF "
	}
	if e.cpu.Flags.ZF {
		flags += "ZF "
	}
	if e.cpu.Flags.SF {
		flags += "SF "
	}
	if e.cpu.Flags.TF {
		flags += "TF "
	}
	if e.cpu.Flags.IF {
		flags += "IF "
	}
	if e.cpu.Flags.DF {
		flags += "DF "
	}
	if e.cpu.Flags.OF {
		flags += "OF "
	}
	fmt.Printf("Flags: %s\n\n", flags)
}

func (e *DOSEmulator) dumpMemory(parts []string) {
	startAddr := uint32(0)
	length := uint32(256)

	if len(parts) > 1 {
		addr, _ := strconv.ParseUint(parts[1], 16, 32)
		startAddr = uint32(addr)
	}

	if len(parts) > 2 {
		l, _ := strconv.ParseUint(parts[2], 10, 32)
		length = uint32(l)
	}

	fmt.Printf("\nMemory dump from %08X:\n", startAddr)
	for i := uint32(0); i < length; i += 16 {
		addr := startAddr + i
		fmt.Printf("%08X: ", addr)

		for j := uint32(0); j < 16; j++ {
			fmt.Printf("%02X ", e.memory.ReadByte(addr+j))
		}

		fmt.Print(" | ")
		for j := uint32(0); j < 16; j++ {
			ch := e.memory.ReadByte(addr + j)
			if ch >= 32 && ch <= 126 {
				fmt.Printf("%c", ch)
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func (e *DOSEmulator) showStack() {
	fmt.Println("\nStack (top 10 entries):")
	count := 0
	for i := len(e.stack) - 1; i >= 0 && count < 10; i-- {
		fmt.Printf("  [%02d] %04X\n", count, e.stack[i])
		count++
	}
	if len(e.stack) == 0 {
		fmt.Println("  (empty)")
	}
	fmt.Println()
}

func (e *DOSEmulator) showStatistics() {
	elapsed := time.Since(e.startTime)

	fmt.Println("\nEMULATOR STATISTICS:")
	fmt.Printf("Program type:     %s\n", e.programType)
	fmt.Printf("Instructions:     %d\n", e.instructionCount)
	fmt.Printf("Running time:     %s\n", elapsed.Round(time.Millisecond))

	if elapsed.Seconds() > 0 {
		ips := float64(e.instructionCount) / elapsed.Seconds()
		fmt.Printf("IPS:              %.0f\n", ips)
	}

	fmt.Printf("Stack depth:      %d\n", len(e.stack))
	fmt.Printf("File handles:     %d\n\n", len(e.fileHandles))
}

func (e *DOSEmulator) disassemble(parts []string) {
	startAddr := CalculateAddress(e.cpu.CS, e.cpu.IP)
	count := 20

	if len(parts) > 1 {
		addr, _ := strconv.ParseUint(parts[1], 16, 32)
		startAddr = uint32(addr)
	}

	if len(parts) > 2 {
		c, _ := strconv.Atoi(parts[2])
		count = c
	}

	fmt.Printf("\nDisassembly from %08X:\n", startAddr)
	addr := startAddr
	for i := 0; i < count; i++ {
		inst := e.decoder.Decode(addr)
		fmt.Printf("%08X: %s\n", addr, inst.Name)
		addr += uint32(inst.Length)
	}
	fmt.Println()
}

func main() {
	if len(os.Args) > 1 {
		emulator := NewDOSEmulator()

		switch os.Args[1] {
		case "-h", "--help":
			fmt.Println("MS-DOS Emulator v5.0 - Complete COM & EXE Support")
			fmt.Println("\nUsage:")
			fmt.Println("  dos              Start interactive shell")
			fmt.Println("  dos <file>       Run COM or EXE file directly")
			fmt.Println("  dos -d <file>    Run in debug mode")
			fmt.Println("\nSupported file formats:")
			fmt.Println("  .COM files       - DOS COM executables")
			fmt.Println("  .EXE files       - DOS EXE executables with relocations")
			fmt.Println("\nFeatures:")
			fmt.Println("  - Full 8086 CPU emulation")
			fmt.Println("  - BIOS interrupts (INT 10h, 16h, 1Ah)")
			fmt.Println("  - DOS interrupts (INT 20h, 21h)")
			fmt.Println("  - File system operations")
			fmt.Println("  - Interactive debugger")
			return

		case "-d", "--debug":
			if len(os.Args) > 2 {
				emulator.debugMode = true
				if err := emulator.LoadFile(os.Args[2]); err != nil {
					fmt.Printf("Error: %v\n", err)
					return
				}
				emulator.Run()
				return
			}

		default:
			if err := emulator.LoadFile(os.Args[1]); err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			emulator.Run()
			return
		}
	}

	emulator := NewDOSEmulator()
	emulator.SimpleShell()
}
