from more_itertools import flatten

from emulator.adressing import ZeroPage, Absolute
from emulator.opcodes.base import OpCode


class BIT(OpCode):
    """
    A & M, N = M7, V = M6
    This instruction is used to test if one or more bits are set in a target memory location.
    The mask pattern in A is ANDed with the value in memory to set or clear the zero flag, but the result is not kept.
    Bits 7 and 6 of the value from memory are copied into the N and V flags.
    """

    @classmethod
    def create_variations(cls):
        variations = [(0x24, ZeroPage, 3,),
                      (0x2C, Absolute, 4,)]
        return map(cls.create_dict_entry, variations)


class CLC(OpCode):
    """
    C = 0
    Set the carry flag to zero.
    """

    @classmethod
    def create_variations(cls):
        variations = [(0x18, None, 2,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def _set_flag():
            cpu.carry = False

        cpu.clear_state_mem()
        cpu.exec_in_cycle(_set_flag)


class SEC(OpCode):
    """
    C = 1
    Set the carry flag to one.
    """

    @classmethod
    def create_variations(cls):
        variations = [(0x38, None, 2,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def _set_flag():
            cpu.carry = True

        cpu.clear_state_mem()
        cpu.exec_in_cycle(_set_flag)


class CLD(OpCode):
    """
    D = 0
    Sets the decimal mode flag to zero.
    """

    @classmethod
    def create_variations(cls):
        variations = [(0xD8, None, 2,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def _set_flag():
            cpu.decimal = False

        cpu.clear_state_mem()
        cpu.exec_in_cycle(_set_flag)


class SED(OpCode):
    """
    D = 1
    Set the decimal mode flag to one.
    """

    @classmethod
    def create_variations(cls):
        variations = [(0xF8, None, 2,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def _set_flag():
            cpu.decimal = True

        cpu.clear_state_mem()
        cpu.exec_in_cycle(_set_flag)


class CLI(OpCode):
    """
    I = 0
    Clears the interrupt disable flag allowing normal interrupt requests to be serviced.
    """

    @classmethod
    def create_variations(cls):
        variations = [(0x58, None, 2,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def _set_flag():
            cpu.interrupts_disabled = False

        cpu.clear_state_mem()
        cpu.exec_in_cycle(_set_flag)


class SEI(OpCode):
    """
    I = 1
    Set the interrupt disable flag to one.
    """

    @classmethod
    def create_variations(cls):
        variations = [(0x78, None, 2,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def _set_flag():
            cpu.interrupts_disabled = True

        cpu.clear_state_mem()
        cpu.exec_in_cycle(_set_flag)


class CLV(OpCode):
    """
    V = 0
    Clears the overflow flag.
    """

    @classmethod
    def create_variations(cls):
        variations = [(0xB8, None, 2,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def _set_flag():
            cpu.overflow = False

        cpu.clear_state_mem()
        cpu.exec_in_cycle(_set_flag)


class NOP(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xEA, None, 2)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        cpu.clear_state_mem()
        pass


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
