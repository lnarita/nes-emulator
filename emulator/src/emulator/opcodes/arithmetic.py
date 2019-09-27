from more_itertools import flatten

from emulator.adressing import IndirectX, ZeroPage, Immediate, Absolute, IndirectY, ZeroPageX, AbsoluteY, AbsoluteX, Accumulator
from emulator.constants import NEGATIVE_BIT, LOW_BITS_MASK
from emulator.opcodes.base import OpCode


class ORA(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x01, IndirectX, 6,),
                      (0x05, ZeroPage, 3,),
                      (0x09, Immediate, 2,),
                      (0x0D, Absolute, 4,),
                      (0x11, IndirectY, 5,),
                      (0x15, ZeroPageX, 4,),
                      (0x19, AbsoluteY, 4,),
                      (0x1D, AbsoluteX, 4,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def _cycle():
            address = self.addressing_mode.fetch_address(cpu, memory)
            value = self.addressing_mode.read_from(cpu, memory, address)
            cpu.a |= value
            cpu.zero = (cpu.a == 0)
            cpu.negative = (cpu.a & NEGATIVE_BIT) > 0
            if self.addressing_mode != Immediate:
                cpu.data = value
                cpu.addr = address
                self.addressing_mode.data = "= %02X" % memory.fetch(address)
        cpu.exec_in_cycle(_cycle)


class AND(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x21, IndirectX, 6,),
                      (0x25, ZeroPage, 3,),
                      (0x29, Immediate, 2,),
                      (0x2D, Absolute, 4,),
                      (0x31, IndirectY, 5,),
                      (0x35, ZeroPageX, 4,),
                      (0x39, AbsoluteY, 4,),
                      (0x3D, AbsoluteX, 4,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def _cycle():
            address = self.addressing_mode.fetch_address(cpu, memory)
            value = self.addressing_mode.read_from(cpu, memory, address)
            cpu.a &= value
            cpu.zero = (cpu.a == 0)
            cpu.negative = (cpu.a & NEGATIVE_BIT) > 0
            if self.addressing_mode != Immediate:
                cpu.data = value
                cpu.addr = address
                self.addressing_mode.data = "= %02X" % memory.fetch(address)
        cpu.exec_in_cycle(_cycle)


class EOR(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x41, IndirectX, 6,),
                      (0x45, ZeroPage, 3,),
                      (0x49, Immediate, 2,),
                      (0x4D, Absolute, 4,),
                      (0x51, IndirectY, 5,),
                      (0x55, ZeroPageX, 4,),
                      (0x59, AbsoluteY, 4,),
                      (0x5D, AbsoluteX, 4,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def _cycle():
            address = self.addressing_mode.fetch_address(cpu, memory)
            value = self.addressing_mode.read_from(cpu, memory, address)
            cpu.a ^= value
            cpu.zero = (cpu.a == 0)
            cpu.negative = (cpu.a & NEGATIVE_BIT) > 0
            if self.addressing_mode != Immediate:
                cpu.data = value
                cpu.addr = address
                self.addressing_mode.data = "= %02X" % memory.fetch(address)
        cpu.exec_in_cycle(_cycle)


class ADC(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x61, IndirectX, 6,),
                      (0x65, ZeroPage, 3,),
                      (0x69, Immediate, 2,),
                      (0x6D, Absolute, 4,),
                      (0x71, IndirectY, 5,),
                      (0x75, ZeroPageX, 4,),
                      (0x79, AbsoluteY, 4,),
                      (0x7D, AbsoluteX, 4,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def _cycle():
            address = self.addressing_mode.fetch_address(cpu, memory)
            addend2 = self.addressing_mode.read_from(cpu, memory, address)
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
                self.addressing_mode.data = "= %02X" % memory.fetch(address)
        cpu.exec_in_cycle(_cycle)


class SBC(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xE1, IndirectX, 6),
                      (0xE5, ZeroPage, 3,),
                      (0xE9, Immediate, 2,),
                      (0xED, Absolute, 4,),
                      (0xF1, IndirectY, 5,),
                      (0xF5, ZeroPageX, 4,),
                      (0xF9, AbsoluteY, 4,),
                      (0xFD, AbsoluteX, 4,)]
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
        cpu.exec_in_cycle(_cycle)


class CMP(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xC1, IndirectX, 6,),
                      (0xC5, ZeroPage, 3,),
                      (0xC9, Immediate, 2,),
                      (0xCD, Absolute, 4,),
                      (0xD1, IndirectY, 5,),
                      (0xD5, ZeroPageX, 4,),
                      (0xD9, AbsoluteY, 4,),
                      (0xDD, AbsoluteX, 4,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def _cycle():
            address = self.addressing_mode.fetch_address(cpu, memory)
            subtrahend = self.addressing_mode.read_from(cpu, memory, address)
            minuend = cpu.a
            if self.addressing_mode != Immediate:
                cpu.addr = address
                cpu.data = subtrahend
                self.addressing_mode.data = "= %02X" % memory.fetch(address)
            # Two's complement
            subtrahend = abs(~subtrahend ^ 0xFF) & 0xFF
            tmp = minuend + subtrahend
            tmp &= 0xFF
            cpu.carry = (minuend >= tmp)
            cpu.zero = tmp == 0
            cpu.negative = tmp >> 7 == 1

        cpu.exec_in_cycle(_cycle)


class CPX(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xE0, Immediate, 2,),
                      (0xE4, ZeroPage, 3,),
                      (0xEC, Absolute, 4,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def _cycle():
            address = self.addressing_mode.fetch_address(cpu, memory)
            subtrahend = self.addressing_mode.read_from(cpu, memory, address)
            minuend = cpu.x
            if self.addressing_mode != Immediate:
                cpu.addr = address
                cpu.data = subtrahend
                self.addressing_mode.data = "= %02X" % memory.fetch(address)
            # Two's complement
            subtrahend = abs(~subtrahend ^ 0xFF) & 0xFF
            tmp = minuend + subtrahend
            tmp &= 0xFF
            cpu.carry = (minuend >= tmp)
            cpu.zero = tmp == 0
            cpu.negative = tmp >> 7 == 1
        cpu.exec_in_cycle(_cycle)


class CPY(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xC0, Immediate, 2,),
                      (0xC4, ZeroPage, 3,),
                      (0xCC, Absolute, 4,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def _cycle():
            address = self.addressing_mode.fetch_address(cpu, memory)
            subtrahend = self.addressing_mode.read_from(cpu, memory, address)
            minuend = cpu.y
            if self.addressing_mode != Immediate:
                cpu.addr = address
                cpu.data = subtrahend
                self.addressing_mode.data = "= %02X" % memory.fetch(address)

            # Two's complement
            subtrahend = abs(~subtrahend ^ 0xFF) & 0xFF
            tmp = minuend + subtrahend
            tmp &= 0xFF
            cpu.carry = (minuend >= tmp)
            cpu.zero = tmp == 0
            cpu.negative = tmp >> 7 == 1
        cpu.exec_in_cycle(_cycle)


class DEC(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xC6, ZeroPage, 5,),
                      (0xCE, Absolute, 6,),
                      (0xD6, ZeroPageX, 6,),
                      (0xDE, AbsoluteX, 7,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        if self.addressing_mode:
            def _dec_m(value):
                if (value == 0):
                    value = 0b11111111
                else:
                    value -= 1
                return value

            address = self.addressing_mode.fetch_address(cpu, memory)
            value = self.addressing_mode.read_from(cpu, memory, address)
            value = cpu.exec_in_cycle(_dec_m, value)
            cpu.negative = (value & 0b10000000) > 0
            cpu.zero = (value == 0)

            self.addressing_mode.data = "= %02X" % memory.fetch(address)
            self.addressing_mode.write_to(cpu, memory, address, value)


class DEX(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xCA, None, 2,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def _dec_x():
            cpu.x -= 1
            cpu.x &= LOW_BITS_MASK
            cpu.zero = (cpu.x == 0)
            cpu.negative = (cpu.x & NEGATIVE_BIT) > 0

        cpu.exec_in_cycle(_dec_x)


class DEY(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x88, None, 2,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def _dec_y():
            cpu.y -= 1
            cpu.y &= LOW_BITS_MASK
            cpu.zero = (cpu.y == 0)
            cpu.negative = (cpu.y & NEGATIVE_BIT) > 0

        cpu.exec_in_cycle(_dec_y)


class INC(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xE6, ZeroPage, 5,),
                      (0xEE, Absolute, 6,),
                      (0xF6, ZeroPageX, 6,),
                      (0xFE, AbsoluteX, 7,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        if self.addressing_mode:
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


class INX(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xE8, None, 2,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def _inc_x():
            cpu.x += 1
            cpu.x &= LOW_BITS_MASK
            cpu.zero = (cpu.x == 0)
            cpu.negative = (cpu.x & NEGATIVE_BIT) > 0

        cpu.exec_in_cycle(_inc_x)


class INY(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xC8, None, 2,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def _inc_y():
            cpu.y += 1
            cpu.y &= LOW_BITS_MASK
            cpu.zero = (cpu.y == 0)
            cpu.negative = (cpu.y & NEGATIVE_BIT) > 0

        cpu.exec_in_cycle(_inc_y)


class ASL(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x06, ZeroPage, 5,),
                      (0x0A, Accumulator, 2,),
                      (0x0E, Absolute, 6,),
                      (0x16, ZeroPageX, 6,),
                      (0x1E, AbsoluteX, 7,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def _exec_asl(value):
            new_value = (value << 1)

            cpu.carry = (new_value & 0b100000000) > 0
            cpu.negative = (new_value & 0b10000000) > 0
            new_value = (new_value & 0b11111111)  # truncates
            cpu.zero = (new_value == 0)
            return new_value

        address = self.addressing_mode.fetch_address(cpu, memory)
        value = self.addressing_mode.read_from(cpu, memory, address)
        new_value = cpu.exec_in_cycle(_exec_asl, value)
        if self.addressing_mode != Accumulator:
            self.addressing_mode.data = "= %02X" % memory.fetch(address)
        self.addressing_mode.write_to(cpu, memory, address, new_value)
        if self.addressing_mode != Accumulator:
            cpu.addr = address
            cpu.data = new_value


class ROL(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x26, ZeroPage, 5,),
                      (0x2A, Accumulator, 2,),
                      (0x2E, Absolute, 6,),
                      (0x36, ZeroPageX, 6,),
                      (0x3E, AbsoluteX, 7,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
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

        address = self.addressing_mode.fetch_address(cpu, memory)
        value = self.addressing_mode.read_from(cpu, memory, address)
        self.addressing_mode.write_to(cpu, memory, address, value)
        new_value = cpu.exec_in_cycle(_exec_rol, value)
        if self.addressing_mode != Accumulator:
            self.addressing_mode.data = "= %02X" % memory.fetch(address)
        self.addressing_mode.write_to(cpu, memory, address, new_value)
        if self.addressing_mode != Accumulator:
            cpu.data = new_value
            cpu.addr = address


class LSR(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x46, ZeroPage, 5,),
                      (0x4A, Accumulator, 2,),
                      (0x4E, Absolute, 6,),
                      (0x56, ZeroPageX, 6,),
                      (0x5E, AbsoluteX, 7,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def _exec_lsr(value):
            new_value = (value >> 1) & LOW_BITS_MASK
            cpu.carry = (value & 0b00000001) > 0
            cpu.zero = (new_value == 0)
            cpu.negative = (new_value & NEGATIVE_BIT) > 0
            return new_value

        address = self.addressing_mode.fetch_address(cpu, memory)
        value = self.addressing_mode.read_from(cpu, memory, address)
        self.addressing_mode.write_to(cpu, memory, address, value)
        new_value = cpu.exec_in_cycle(_exec_lsr, value)
        if self.addressing_mode != Accumulator:
            self.addressing_mode.data = "= %02X" % memory.fetch(address)
        self.addressing_mode.write_to(cpu, memory, address, new_value)
        if self.addressing_mode != Accumulator:
            cpu.data = new_value
            cpu.addr = address


class ROR(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x66, ZeroPage, 5,),
                      (0x6A, Accumulator, 2,),
                      (0x6E, Absolute, 6,),
                      (0x76, ZeroPageX, 6,),
                      (0x7E, AbsoluteX, 7,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
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

        address = self.addressing_mode.fetch_address(cpu, memory)
        value = self.addressing_mode.read_from(cpu, memory, address)
        self.addressing_mode.write_to(cpu, memory, address, value)
        new_value = cpu.exec_in_cycle(_exec_ror, value)
        if self.addressing_mode != Accumulator:
            self.addressing_mode.data = "= %02X" % memory.fetch(address)
        self.addressing_mode.write_to(cpu, memory, address, new_value)
        if self.addressing_mode != Accumulator:
            cpu.data = new_value
            cpu.addr = address


class ArithmeticAndLogicalOpCodes:
    opcodes = [
        ORA,
        AND,
        EOR,
        ADC,
        SBC,
        CMP,
        CPX,
        CPY,
        DEC,
        DEX,
        DEY,
        INC,
        INX,
        INY,
        ASL,
        ROL,
        LSR,
        ROR
    ]

    @staticmethod
    def all_commands():
        return flatten(
            map(lambda x: list(x.create_variations()), ArithmeticAndLogicalOpCodes.opcodes)
        )
