from more_itertools import flatten

from constants import AddressingMode
from opcodes.base import OpCode


class BPL(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x10, AddressingMode.RELATIVE, 2,)]
        return map(cls.create_dict_entry, variations)


class BMI(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x30, AddressingMode.RELATIVE, 2,)]
        return map(cls.create_dict_entry, variations)


class BVC(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x50, AddressingMode.RELATIVE, 2,)]
        return map(cls.create_dict_entry, variations)


class BVS(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x70, AddressingMode.RELATIVE, 2,)]
        return map(cls.create_dict_entry, variations)


class BCC(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x90, AddressingMode.RELATIVE, 2,)]
        return map(cls.create_dict_entry, variations)


class BCS(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xB0, AddressingMode.RELATIVE, 2,)]
        return map(cls.create_dict_entry, variations)


class BNE(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xD0, AddressingMode.RELATIVE, 2,)]
        return map(cls.create_dict_entry, variations)


class BEQ(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xF0, AddressingMode.RELATIVE, 2,)]
        return map(cls.create_dict_entry, variations)


class BRK(OpCode):
    """
    The BRK instruction forces the generation of an interrupt request.
    The program counter and processor status are pushed on the stack then the IRQ interrupt vector at $FFFE/F is loaded into the PC and the break flag in the status set to one.
    """

    @classmethod
    def create_variations(cls):
        variations = [(0x00, None, 7,)]
        return map(cls.create_dict_entry, variations)

    @classmethod
    def exec(cls, cpu, memory):
        cpu.break_command = True
        # TODO: pushes


class RTI(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x40, None, 6,)]
        return map(cls.create_dict_entry, variations)


class JSR(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x20, AddressingMode.ABSOLUTE, 6,)]
        return map(cls.create_dict_entry, variations)


class RTS(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x60, None, 6,)]
        return map(cls.create_dict_entry, variations)


class JMP(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x4C, AddressingMode.ABSOLUTE, 3,),
                      (0x6C, AddressingMode.INDIRECT, 5,)]
        return map(cls.create_dict_entry, variations)


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
