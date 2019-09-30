from more_itertools import flatten

from emulator.adressing import Immediate, ZeroPage, Absolute, ZeroPageY, AbsoluteY, IndirectY, ZeroPageX, AbsoluteX, IndirectX, Accumulator
from emulator.constants import NEGATIVE_BIT, LOW_BITS_MASK
from emulator.cpu import StatusRegisterFlags
from emulator.opcodes.base import OpCode
import numpy as np


def unofficial_opcode_str(self, name=None):
    if name == None:
        name = type(self).__name__

    def __str_addr():
        if self.addressing_mode is not None:
            if self.addressing_mode.low is not None and self.addressing_mode.high is not None:
                return "{:02X} {:02X}".format(self.addressing_mode.low, (self.addressing_mode.high >> 8))
            elif self.addressing_mode.low is not None:
                return "{:02X}".format(self.addressing_mode.low)
        return ""

    def __str_addr_2():
        if self.addressing_mode is not None:
            if self.addressing_mode.addr is not None and self.addressing_mode.data is not None:
                return "*{} {} {}".format(name, self.addressing_mode.addr, self.addressing_mode.data)
            elif self.addressing_mode.addr is not None:
                return "*{} {}".format(name, self.addressing_mode.addr)
        return "*{}".format(name)
    with_space = "{:02X} {:<6}{:<30} ".format(
        self.id, __str_addr(), __str_addr_2())
    without_space = "{:02X} {:<6}{:<30}".format(
        self.id, __str_addr(), __str_addr_2())

    if len(with_space) > 40:
        return without_space
    else:
        return with_space


class IGN(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x0C, Absolute, 4),
                      (0x04, ZeroPage, 3,),
                      (0x44, ZeroPage, 3,),
                      (0x64, ZeroPage, 3,),
                      (0x14, ZeroPageX, 4),
                      (0x34, ZeroPageX, 4),
                      (0x54, ZeroPageX, 4),
                      (0x74, ZeroPageX, 4),
                      (0xD4, ZeroPageX, 4),
                      (0xF4, ZeroPageX, 4),
                      (0x1C, AbsoluteX, 4),
                      (0x3C, AbsoluteX, 4),
                      (0x5C, AbsoluteX, 4),
                      (0x7C, AbsoluteX, 4),
                      (0xDC, AbsoluteX, 4),
                      (0xFC, AbsoluteX, 4)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def _cycle():
            if self.addressing_mode:
                address = self.addressing_mode.fetch_address(cpu, memory)
                value = self.addressing_mode.read_from(cpu, memory, address)
                if self.addressing_mode != Immediate:
                    self.addressing_mode.data = "= %02X" % memory.fetch(
                        address)
                    cpu.addr = address
                    cpu.data = value
        _cycle()

    def __str__(self):
        return unofficial_opcode_str(self, "NOP")


class SKB(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x80, Immediate, 2),
                      (0x82, Immediate, 2,),
                      (0x89, Immediate, 2,),
                      (0xC2, Immediate, 2,),
                      (0xE2, Immediate, 2)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        _address = self.addressing_mode.fetch_address(cpu, memory)

    def __str__(self):
        return unofficial_opcode_str(self, "NOP")


class NOP(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x1A, None, 2),
                      (0x3A, None, 2),
                      (0x5A, None, 2),
                      (0x7A, None, 2),
                      (0xDA, None, 2),
                      (0xFA, None, 2)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def _stall():
            pass

        cpu.exec_in_cycle(_stall)

    def __str__(self):
        return unofficial_opcode_str(self)


class LAX(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [
            (0xA3, IndirectX, 6),
            (0xA7, ZeroPage, 3),
            (0xAF, Absolute, 4),
            (0xB3, IndirectY, 5),
            (0xB7, ZeroPageY, 4),
            (0xBF, AbsoluteY, 4)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def _lda():
            if self.addressing_mode:
                address = self.addressing_mode.fetch_address(cpu, memory)
                value = self.addressing_mode.read_from(cpu, memory, address)
                if self.addressing_mode != Immediate:
                    self.addressing_mode.data = "= %02X" % memory.fetch(
                        address)
                    cpu.addr = address
                    cpu.data = value

                cpu.a = value
                cpu.zero = cpu.a == 0
                cpu.negative = (cpu.a & 0b10000000) > 0

        def _tax():
            cpu.x = cpu.a
            cpu.zero = cpu.x == 0
            cpu.negative = (cpu.x & 0b10000000) > 0

        _lda()
        _tax()

    def __str__(self):
        return unofficial_opcode_str(self)

# SAX (d,X) ($83 dd; 6 cycles)
# SAX d ($87 dd; 3 cycles)
# SAX a ($8F aa aa; 4 cycles)
# SAX d,Y ($97 dd; 4 cycles)
class SAX(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x83, IndirectX, 6),
                      (0x87, ZeroPage, 3),
                      (0x8F, Absolute, 4),
                      (0x97, ZeroPageY, 4)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def _stall():
            pass

        def cycle_sta():
            if self.addressing_mode:
                address = self.addressing_mode.fetch_address(cpu, memory)
                cpu.addr = address
                cpu.data = cpu.a & cpu.x
                self.addressing_mode.data = "= %02X" % memory.fetch(address)
                self.addressing_mode.write_to(cpu, memory, address, cpu.a & cpu.x)

        cycle_sta()

    def __str__(self):
        return unofficial_opcode_str(self)

class SBC(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xEB, Immediate, 2,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def _cycle():
            def wrap_sub(a, b):
                return (a - b) % 0x100

            address = self.addressing_mode.fetch_address(cpu, memory)
            subtrahend = self.addressing_mode.read_from(cpu, memory, address)
            minuend = cpu.a
            if self.addressing_mode != Immediate:
                cpu.addr = address
                cpu.data = subtrahend
                self.addressing_mode.data = "= %02X" % memory.fetch(address)

            new_a = (wrap_sub(wrap_sub(minuend, subtrahend), (0 if cpu.carry else 1)))
            cpu.carry = (new_a & 0b01111111) == new_a
            cpu.overflow = ((cpu.a ^ subtrahend) & NEGATIVE_BIT > 0) and (((cpu.a ^ new_a) & LOW_BITS_MASK) & NEGATIVE_BIT > 0)
            cpu.a = new_a
            cpu.a &= 0xFF
            cpu.zero = cpu.a == 0
            cpu.negative = cpu.a >> 7 == 1

        _cycle()
        
    def __str__(self):
        return unofficial_opcode_str(self)


# DCP (d,X) ($C3 dd; 8 cycles)
# DCP d ($C7 dd; 5 cycles)
# DCP a ($CF aa aa; 6 cycles)
# DCP (d),Y ($D3 dd; 8 cycles)
# DCP d,X ($D7 dd; 6 cycles)
# DCP a,Y ($DB aa aa; 7 cycles)
# DCP a,X ($DF aa aa; 7 cycles)
class DCP(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [
            (0xC3, IndirectX, 8),
            (0xC7, ZeroPage, 5),
            (0xCF, Absolute, 6),
            (0xD3, IndirectY, 8),
            (0xD7, ZeroPageX, 6),
            (0xDB, AbsoluteY, 7),
            (0xDF, AbsoluteX, 7),
        ]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def _stall():
            pass
        def _exec():
            def _dec_m(value):
                if (value == 0):
                    value = 0b11111111
                else:
                    value -= 1
                return value

            def _cycle():
                address = self.addressing_mode.fetch_address(cpu, memory)
                value = self.addressing_mode.read_from(cpu, memory, address)
                value = cpu.exec_in_cycle(_dec_m, value)
                cpu.negative = (value & 0b10000000) > 0
                cpu.zero = (value == 0)

                self.addressing_mode.data = "= %02X" % memory.fetch(address)
                self.addressing_mode.write_to(cpu, memory, address, value)
                cpu.addr = address
                cpu.data = value

                minuend = cpu.a
                if self.addressing_mode != Immediate:
                    cpu.addr = address
                    cpu.data = value
                # Two's complement
                subatrend = value
                subatrend = abs(~subatrend ^ 0xFF) & 0xFF
                tmp = minuend + subatrend
                tmp &= 0xFF
                cpu.carry = (minuend >= tmp)
                cpu.zero = tmp == 0
                cpu.negative = tmp >> 7 == 1

            if self.addressing_mode == AbsoluteX:
                # FIXME: this is ugly, but it works
                cycle_start = cpu.cycle
                _cycle()
                cycle_end = cpu.cycle
                if (cycle_end - cycle_start) < (self.cycles - 1):
                    cpu.exec_in_cycle(_stall)
            else:
                _cycle()

        _exec()
        
    def __str__(self):
        return unofficial_opcode_str(self)

class ISC(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [
            (0xE3, IndirectX, 8),
            (0xE7, ZeroPage, 5),
            (0xEF, Absolute, 6),
            (0xF3, IndirectY, 8),
            (0xF7, ZeroPageX, 6),
            (0xFB, AbsoluteY, 7),
            (0xFF, AbsoluteX, 7),
        ]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def _stall():
            pass
        def _cycle():
            def _inc_m(value):
                if (value == 0b11111111):
                    value = 0
                else:
                    value += 1
                return value

            address = self.addressing_mode.fetch_address(cpu, memory)
            value = self.addressing_mode.read_from(cpu, memory, address)
            value = cpu.exec_in_cycle(_inc_m, value)
            cpu.negative = (value & 0b10000000) > 0
            cpu.zero = (value == 0)

            self.addressing_mode.data = "= %02X" % memory.fetch(address)
            self.addressing_mode.write_to(cpu, memory, address, value)
            cpu.addr = address
            cpu.data = value

            subtrahend = value
            minuend = cpu.a
            if self.addressing_mode != Immediate:
                cpu.addr = address
                cpu.data = subtrahend

            n = np.int16(minuend) - np.int16(subtrahend) - np.int16(0 if cpu.carry else 1)
            a = np.uint8(n)
            cpu.zero = a == 0
            cpu.negative = (a & NEGATIVE_BIT) > 0
            cpu.carry = n >= 0
            cpu.overflow = ((minuend ^ subtrahend) & 0x80 > 0) and ((minuend ^ a) & 0x80 > 0)
            cpu.a = int(a)

        if self.addressing_mode == AbsoluteX:
            # FIXME: this is ugly, but it works
            cycle_start = cpu.cycle
            _cycle()
            cycle_end = cpu.cycle
            if (cycle_end - cycle_start) < (self.cycles - 1):
                cpu.exec_in_cycle(_stall)
        else:
            _cycle()
        
    def __str__(self):
        return unofficial_opcode_str(self, "ISB")


class SLO(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [
            (0x03, IndirectX, 8),
            (0x07, ZeroPage, 5),
            (0x0F, Absolute, 6),
            (0x13, IndirectY, 8),
            (0x17, ZeroPageX, 6),
            (0x1B, AbsoluteY, 7),
            (0x1F, AbsoluteX, 7),
        ]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):

        def _stall():
            pass

        def _exec_asl(value):
            new_value = (value << 1)

            cpu.carry = (new_value & 0b100000000) > 0
            cpu.negative = (new_value & 0b10000000) > 0
            new_value = (new_value & 0b11111111)  # truncates
            cpu.zero = (new_value == 0)
            return new_value

        def _cycle():
            address = self.addressing_mode.fetch_address(cpu, memory)
            value = self.addressing_mode.read_from(cpu, memory, address)
            self.addressing_mode.write_to(cpu, memory, address, value)
            new_value = _exec_asl(value)
            if self.addressing_mode != Accumulator:
                self.addressing_mode.data = "= %02X" % memory.fetch(address)
            self.addressing_mode.write_to(cpu, memory, address, new_value)

            if self.addressing_mode != Accumulator:
                cpu.addr = address
                cpu.data = new_value

            value = new_value
            cpu.a |= value
            cpu.zero = (cpu.a == 0)
            cpu.negative = (cpu.a & NEGATIVE_BIT) > 0
            if self.addressing_mode != Immediate:
                cpu.data = value
                cpu.addr = address

        if self.addressing_mode == AbsoluteX:
            # FIXME: this is ugly, but it works
            cycle_start = cpu.cycle
            _cycle()
            cycle_end = cpu.cycle
            if (cycle_end - cycle_start) < (self.cycles - 1):
                cpu.exec_in_cycle(_stall)
        else:
            _cycle()
        
    def __str__(self):
        return unofficial_opcode_str(self)

class RLA(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [
            (0x23, IndirectX, 8),
            (0x27, ZeroPage, 5),
            (0x2F, Absolute, 6),
            (0x33, IndirectY, 8),
            (0x37, ZeroPageX, 6),
            (0x3B, AbsoluteY, 7),
            (0x3F, AbsoluteX, 7),
        ]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def _stall():
            pass

        def _exec_rol(value):
            new_value = (value << 1) & LOW_BITS_MASK
            if cpu.carry:
                new_value |= 0b00000001
            else:
                new_value &= 0b111111110
            cpu.carry = (value & 0b10000000) > 0
            cpu.zero = (new_value == 0)
            cpu.negative = (new_value & NEGATIVE_BIT) > 0
            return new_value

        def _cycle():
            address = self.addressing_mode.fetch_address(cpu, memory)
            value = self.addressing_mode.read_from(cpu, memory, address)
            self.addressing_mode.write_to(cpu, memory, address, value)
            new_value = _exec_rol(value)
            if self.addressing_mode != Accumulator:
                self.addressing_mode.data = "= %02X" % memory.fetch(address)
            self.addressing_mode.write_to(cpu, memory, address, new_value)
            if self.addressing_mode != Accumulator:
                cpu.data = new_value
                cpu.addr = address

            value = new_value
            cpu.a &= value
            cpu.zero = (cpu.a == 0)
            cpu.negative = (cpu.a & NEGATIVE_BIT) > 0
            if self.addressing_mode != Immediate:
                cpu.data = value
                cpu.addr = address

        if self.addressing_mode == AbsoluteX:
            # FIXME: this is ugly, but it works
            cycle_start = cpu.cycle
            _cycle()
            cycle_end = cpu.cycle
            if (cycle_end - cycle_start) < (self.cycles - 1):
                cpu.exec_in_cycle(_stall)
        else:
            _cycle()
        
    def __str__(self):
        return unofficial_opcode_str(self)

class SRE(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [
            (0x43, IndirectX, 8),
            (0x47, ZeroPage, 5),
            (0x4F, Absolute, 6),
            (0x53, IndirectY, 8),
            (0x57, ZeroPageX, 6),
            (0x5B, AbsoluteY, 7),
            (0x5F, AbsoluteX, 7),
        ]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):

        def _stall():
            pass

        def _exec_lsr(value):
            new_value = (value >> 1) & LOW_BITS_MASK
            cpu.carry = (value & 0b00000001) > 0
            cpu.zero = (new_value == 0)
            cpu.negative = (new_value & NEGATIVE_BIT) > 0
            return new_value

        def _cycle():
            address = self.addressing_mode.fetch_address(cpu, memory)
            value = self.addressing_mode.read_from(cpu, memory, address)
            self.addressing_mode.write_to(cpu, memory, address, value)
            new_value = _exec_lsr(value)
            if self.addressing_mode != Accumulator:
                self.addressing_mode.data = "= %02X" % memory.fetch(address)
            self.addressing_mode.write_to(cpu, memory, address, new_value)
            if self.addressing_mode != Accumulator:
                cpu.data = new_value
                cpu.addr = address

            value = new_value
            cpu.a ^= value
            cpu.zero = (cpu.a == 0)
            cpu.negative = (cpu.a & NEGATIVE_BIT) > 0
            if self.addressing_mode != Immediate:
                cpu.data = value
                cpu.addr = address

        if self.addressing_mode == AbsoluteX:
            # FIXME: this is ugly, but it works
            cycle_start = cpu.cycle
            _cycle()
            cycle_end = cpu.cycle
            if (cycle_end - cycle_start) < (self.cycles - 1):
                cpu.exec_in_cycle(_stall)
        else:
            _cycle()
        
    def __str__(self):
        return unofficial_opcode_str(self)



class RRA(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [
            (0x63, IndirectX, 8),
            (0x67, ZeroPage, 5),
            (0x6F, Absolute, 6),
            (0x73, IndirectY, 8),
            (0x77, ZeroPageX, 6),
            (0x7B, AbsoluteY, 7),
            (0x7F, AbsoluteX, 7),
        ]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def _stall():
            pass

        def _exec_ror(value):
            new_value = (value >> 1) & LOW_BITS_MASK
            if cpu.carry:
                new_value |= 0b10000000
            else:
                new_value &= 0b011111111
            cpu.carry = (value & 0b00000001) > 0
            cpu.zero = (new_value == 0)
            cpu.negative = (new_value & NEGATIVE_BIT) > 0
            return new_value

        def _cycle():
            address = self.addressing_mode.fetch_address(cpu, memory)
            value = self.addressing_mode.read_from(cpu, memory, address)
            self.addressing_mode.write_to(cpu, memory, address, value)
            new_value = _exec_ror(value)
            if self.addressing_mode != Accumulator:
                self.addressing_mode.data = "= %02X" % memory.fetch(address)
            self.addressing_mode.write_to(cpu, memory, address, new_value)
            if self.addressing_mode != Accumulator:
                cpu.data = new_value
                cpu.addr = address

            addend2 = new_value
            addend1 = cpu.a
            cpu.a = addend1 + addend2 + (1 if cpu.carry else 0)
            cpu.carry = (cpu.a >> 8) != 0
            cpu.a &= 0xff
            cpu.overflow = addend1 >> 7 == addend2 >> 7 and addend1 >> 7 != cpu.a >> 7
            cpu.negative = cpu.a >> 7 == 1
            cpu.zero = cpu.a == 0

            if self.addressing_mode != Immediate:
                cpu.addr = address
                cpu.data = addend2


        if self.addressing_mode == AbsoluteX:
            # FIXME: this is ugly, but it works
            cycle_start = cpu.cycle
            _cycle()
            cycle_end = cpu.cycle
            if (cycle_end - cycle_start) < (self.cycles - 1):
                cpu.exec_in_cycle(_stall)
        else:
            _cycle()
        
    def __str__(self):
        return unofficial_opcode_str(self)


class UnofficialOpcodes:
    opcodes = [
        IGN,
        NOP,
        SKB,
        LAX,
        SAX,
        SBC,
        DCP,
        ISC,
        SLO,
        RLA,
        SRE,
        RRA
    ]

    @staticmethod
    def all_commands():
        return flatten(
            map(lambda x: list(x.create_variations()), UnofficialOpcodes.opcodes)
        )
