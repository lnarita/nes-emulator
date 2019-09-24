from more_itertools import flatten

from emulator.adressing import IndirectX, ZeroPage, Immediate, Absolute, IndirectY, ZeroPageX, AbsoluteY, AbsoluteX, Accumulator
from emulator.constants import NEGATIVE_BIT
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
        if self.addressing_mode:
            # FIXME: count cycles
            address = self.addressing_mode.fetch_address(cpu, memory)
            cpu.addr = address
            cpu.data = self.addressing_mode.read_from(cpu, memory, address)
            cpu.a = cpu.a | cpu.data
            cpu.zero = (cpu.a == 0)
            cpu.negative = (cpu.a & NEGATIVE_BIT) > 0


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


class CPX(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xE0, Immediate, 2,),
                      (0xE4, ZeroPage, 3,),
                      (0xEC, Absolute, 4,)]
        return map(cls.create_dict_entry, variations)


class CPY(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xC0, Immediate, 2,),
                      (0xC4, ZeroPage, 3,),
                      (0xCC, Absolute, 4,)]
        return map(cls.create_dict_entry, variations)


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

            self.addressing_mode.write_to(cpu,memory, address, value)


class DEX(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xCA, None, 2,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def _dec_x():
            cpu.x -= 1
            cpu.zero = (cpu.x == 0)
            cpu.negative = cpu.x < 0

        cpu.exec_in_cycle(_dec_x)


class DEY(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x88, None, 2,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def _dec_y():
            # TODO: hmmm, won't this f*ck up? python integers are not bound and 6502 works only with 8 bit integers :/
            cpu.y -= 1
            cpu.zero = (cpu.y == 0)
            cpu.negative = cpu.y < 0

        cpu.exec_in_cycle(_dec_y)


class INC(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xE6, ZeroPage, 5,),
                      (0xEE, Absolute, 6,),
                      (0xF6, ZeroPageX, 6,),
                      (0xFE, AbsoluteX, 7,)]
        return map(cls.create_dict_entry, variations)

    def exec (self, cpu, memory):
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

            self.addressing_mode.write_to(cpu,memory, address, value)



class INX(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xE8, None, 2,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def _dec_y():
            # TODO: hmmm, won't this f*ck up? python integers are not bound and 6502 works only with 8 bit integers :/
            cpu.x += 1
            cpu.zero = (cpu.y == 0)
            cpu.negative = cpu.y < 0

        cpu.exec_in_cycle(_dec_y)


class INY(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xC8, None, 2,)]
        return map(cls.create_dict_entry, variations)


class ASL(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x06, ZeroPage, 5,),
                      (0x0A, Accumulator, 2,),
                      (0x0E, Absolute, 6,),
                      (0x16, ZeroPageX, 6,),
                      (0x1E, AbsoluteX, 7,)]
        return map(cls.create_dict_entry, variations)

    def exec (self, cpu, memory):
        if self.addressing_mode:
            address = self.addressing_mode.fetch_address(cpu, memory)
            value = self.addressing_mode.read_from(cpu, memory, address)
            value = (value << 1)
            
            cpu.carry = (value & 0b100000000) > 0
            cpu.negative = (value & 0b10000000) > 0
            value = (value & 0b11111111) # truncates
            cpu.zero = (value == 0)

            self.addressing_mode.write_to(cpu,memory, address, value)


class ROL(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x26, ZeroPage, 5,),
                      (0x2A, Accumulator, 2,),
                      (0x2E, Absolute, 6,),
                      (0x36, ZeroPageX, 6,),
                      (0x3E, AbsoluteX, 7,)]
        return map(cls.create_dict_entry, variations)


class LSR(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x46, ZeroPage, 5,),
                      (0x4A, Accumulator, 2,),
                      (0x4E, Absolute, 6,),
                      (0x56, ZeroPageX, 6,),
                      (0x5E, AbsoluteX, 7,)]
        return map(cls.create_dict_entry, variations)


class ROR(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x66, ZeroPage, 5,),
                      (0x6A, Accumulator, 2,),
                      (0x6E, Absolute, 6,),
                      (0x76, ZeroPageX, 6,),
                      (0x7E, AbsoluteX, 7,)]
        return map(cls.create_dict_entry, variations)


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
