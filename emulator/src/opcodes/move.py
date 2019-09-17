from more_itertools import flatten

from constants import AddressingMode
from opcodes.base import OpCode


class LDA(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xA1, AddressingMode.INDIRECT_X, 6),
                      (0xA5, AddressingMode.ZERO_PAGE, 3,),
                      (0xA9, AddressingMode.IMMEDIATE, 2,),
                      (0xAD, AddressingMode.ABSOLUTE, 4,),
                      (0xB1, AddressingMode.INDIRECT_Y, 5),
                      (0xB5, AddressingMode.ZERO_PAGE_X, 4,),
                      (0xB9, AddressingMode.ABSOLUTE_Y, 4,),
                      (0xBD, AddressingMode.ABSOLUTE_X, 4,)]
        return map(cls.create_dict_entry, variations)


class STA(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x81, AddressingMode.INDIRECT_X, 6,),
                      (0x85, AddressingMode.ZERO_PAGE, 3,),
                      (0x8D, AddressingMode.ABSOLUTE, 4,),
                      (0x91, AddressingMode.INDIRECT_Y, 6,),
                      (0x95, AddressingMode.ZERO_PAGE_X, 4,),
                      (0x99, AddressingMode.ABSOLUTE_Y, 5,),
                      (0x9D, AddressingMode.ABSOLUTE_X, 5,)]
        return map(cls.create_dict_entry, variations)


class LDX(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xA2, AddressingMode.IMMEDIATE, 2,),
                      (0xA6, AddressingMode.ZERO_PAGE, 3,),
                      (0xAE, AddressingMode.ABSOLUTE, 4,),
                      (0xB6, AddressingMode.ZERO_PAGE_Y, 4,),
                      (0xBE, AddressingMode.ABSOLUTE_Y, 4,)]
        return map(cls.create_dict_entry, variations)


class STX(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x86, AddressingMode.ZERO_PAGE, 3,),
                      (0x8E, AddressingMode.ABSOLUTE, 4,),
                      (0x96, AddressingMode.ZERO_PAGE_Y, 4,)]
        return map(cls.create_dict_entry, variations)


class LDY(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xA0, AddressingMode.IMMEDIATE, 2,),
                      (0xA4, AddressingMode.ZERO_PAGE, 3,),
                      (0xAC, AddressingMode.ABSOLUTE, 4,),
                      (0xB4, AddressingMode.ZERO_PAGE_X, 4,),
                      (0xBC, AddressingMode.ABSOLUTE_X, 4,)]
        return map(cls.create_dict_entry, variations)


class STY(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x84, AddressingMode.ZERO_PAGE, 3,),
                      (0x8C, AddressingMode.ABSOLUTE, 4,),
                      (0x94, AddressingMode.ZERO_PAGE_X, 4,)]
        return map(cls.create_dict_entry, variations)


class TAX(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xAA, None, 2,)]
        return map(cls.create_dict_entry, variations)
    def exec(self, cpu, memory):
        cpu.x = cpu.a
        cpu.zero = cpu.x == 0
        cpu.negative = cpu.x & 0b10000000
        cpu.inc_cycle()
        

class TXA(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x8A, None, 2,)]
        return map(cls.create_dict_entry, variations)
    def exec(self, cpu, memory):
        cpu.a = cpu.x
        cpu.zero = cpu.a == 0
        cpu.negative = cpu.a & 0b10000000
        cpu.inc_cycle()
        

class TAY(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xA8, None, 2,)]
        return map(cls.create_dict_entry, variations)
    def exec(self, cpu, memory):
        cpu.y = cpu.a
        cpu.zero = cpu.y == 0
        cpu.negative = cpu.y & 0b10000000
        cpu.inc_cycle()
        

class TYA(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x98, None, 2,)]
        return map(cls.create_dict_entry, variations)
    def exec(self, cpu, memory):
        cpu.a = cpu.y
        cpu.zero = cpu.a == 0
        cpu.negative = cpu.a & 0b10000000
        cpu.inc_cycle()
        

class TSX(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xBA, None, 2,)]
        return map(cls.create_dict_entry, variations)
    def exec(self, cpu, memory):
        cpu.x = cpu.sp
        cpu.zero = cpu.x == 0
        cpu.negative = cpu.x & 0b10000000
        cpu.inc_cycle()

class TXS(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x9A, None, 2,)]
        return map(cls.create_dict_entry, variations)
    def exec(self, cpu, memory):
        cpu.sp = cpu.x
        cpu.zero = cpu.sp == 0
        cpu.negative = cpu.sp & 0b10000000
        cpu.inc_cycle()

class PLA(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x68, None, 4,)]
        return map(cls.create_dict_entry, variations)


class PHA(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x48, None, 3,)]
        return map(cls.create_dict_entry, variations)


class PLP(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x28, None, 4,)]
        return map(cls.create_dict_entry, variations)


class PHP(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x08, None, 3,)]
        return map(cls.create_dict_entry, variations)


class MoveOpCodes:
    opcodes = [
        LDA,
        STA,
        LDX,
        STX,
        LDY,
        STY,
        TAX,
        TXA,
        TAY,
        TYA,
        TSX,
        TXS,
        PLA,
        PHA,
        PLP,
        PHP
    ]

    @staticmethod
    def all_commands():
        return flatten(
            map(lambda x: list(x.create_variations()), MoveOpCodes.opcodes)
        )
