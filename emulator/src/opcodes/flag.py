from more_itertools import flatten

from constants import AddressingMode
from opcodes.base import OpCode


class BIT(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x24, AddressingMode.ZERO_PAGE, 3,),
                      (0x2C, AddressingMode.ABSOLUTE, 4,)]
        return map(lambda x: tuple((x[0], cls(*x))), variations)


class CLC(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x18, None, 2,)]
        return map(lambda x: tuple((x[0], cls(*x))), variations)


class SEC(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x38, None, 2,)]
        return map(lambda x: tuple((x[0], cls(*x))), variations)


class CLD(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xD8, None, 2,)]
        return map(lambda x: tuple((x[0], cls(*x))), variations)


class SED(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xF8, None, 2,)]
        return map(lambda x: tuple((x[0], cls(*x))), variations)


class CLI(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x58, None, 2,)]
        return map(lambda x: tuple((x[0], cls(*x))), variations)


class SEI(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x78, None, 2,)]
        return map(lambda x: tuple((x[0], cls(*x))), variations)


class CLV(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xB8, None, 2,)]
        return map(lambda x: tuple((x[0], cls(*x))), variations)


class NOP(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xEA, None, 2)]
        return map(lambda x: tuple((x[0], cls(*x))), variations)


class FlagOpCodes:
    opcodes = [
        BIT,
        CLC,
        SEC,
        CLD,
        SED,
        CLI,
        SEI,
        CLV,
        NOP
    ]

    @staticmethod
    def all_commands():
        return flatten(
            map(lambda x: list(x.create_variations()), FlagOpCodes.opcodes)
        )
