from more_itertools import flatten

from constants import AddressingMode
from opcodes.base import OpCode


class BPL(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x10, AddressingMode.RELATIVE, 2,)]
        return map(lambda x: tuple((x[0], cls(*x))), variations)


class BMI(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x30, AddressingMode.RELATIVE, 2,)]
        return map(lambda x: tuple((x[0], cls(*x))), variations)


class BVC(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x50, AddressingMode.RELATIVE, 2,)]
        return map(lambda x: tuple((x[0], cls(*x))), variations)


class BVS(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x70, AddressingMode.RELATIVE, 2,)]
        return map(lambda x: tuple((x[0], cls(*x))), variations)


class BCC(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x90, AddressingMode.RELATIVE, 2,)]
        return map(lambda x: tuple((x[0], cls(*x))), variations)


class BCS(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xB0, AddressingMode.RELATIVE, 2,)]
        return map(lambda x: tuple((x[0], cls(*x))), variations)


class BNE(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xD0, AddressingMode.RELATIVE, 2,)]
        return map(lambda x: tuple((x[0], cls(*x))), variations)


class BEQ(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xF0, AddressingMode.RELATIVE, 2,)]
        return map(lambda x: tuple((x[0], cls(*x))), variations)


class BRK(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x00, None, 7,)]
        return map(lambda x: tuple((x[0], cls(*x))), variations)


class RTI(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x40, None, 6,)]
        return map(lambda x: tuple((x[0], cls(*x))), variations)


class JSR(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x20, AddressingMode.ABSOLUTE, 6,)]
        return map(lambda x: tuple((x[0], cls(*x))), variations)


class RTS(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x60, None, 6,)]
        return map(lambda x: tuple((x[0], cls(*x))), variations)


class JMP(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x4C, AddressingMode.ABSOLUTE, 3,),
                      (0x6C, AddressingMode.INDIRECT, 5,)]
        return map(lambda x: tuple((x[0], cls(*x))), variations)


class JumpOpCodes:
    opcodes = [
        BPL,
        BMI,
        BVC,
        BVS,
        BCC,
        BCS,
        BNE,
        BEQ,
        BRK,
        RTI,
        JSR,
        RTS,
        JMP
    ]

    @staticmethod
    def all_commands():
        return flatten(
            map(lambda x: list(x.create_variations()), JumpOpCodes.opcodes)
        )
