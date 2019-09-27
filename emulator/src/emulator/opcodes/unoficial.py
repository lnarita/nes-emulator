from more_itertools import flatten

from emulator.adressing import Immediate, ZeroPage, Absolute, ZeroPageY, AbsoluteY, IndirectY, ZeroPageX, AbsoluteX, IndirectX
from emulator.constants import NEGATIVE_BIT, LOW_BITS_MASK
from emulator.cpu import StatusRegisterFlags
from emulator.opcodes.base import OpCode


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
                    return "{} {} {}".format("*NOP", self.addressing_mode.addr, self.addressing_mode.data)
                elif self.addressing_mode.addr is not None:
                    return "{} {}".format("*NOP", self.addressing_mode.addr)
            return "{}".format(type(self).__name__)

        return "{:02X} {:<6}{:<30} ".format(self.id, __str_addr(), __str_addr_2())

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
        address = self.addressing_mode.fetch_address(cpu, memory)

    def __str__(self):
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
                    return "{} {} {}".format("*NOP", self.addressing_mode.addr, self.addressing_mode.data)
                elif self.addressing_mode.addr is not None:
                    return "{} {}".format("*NOP", self.addressing_mode.addr)
            return "{}".format(type(self).__name__)

        return "{:02X} {:<6}{:<30} ".format(self.id, __str_addr(), __str_addr_2())


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
        def __str_addr():
            return ""

        def __str_addr_2():
            return "{}".format("*NOP")

        return "{:02X} {:<6}{:<30} ".format(self.id, __str_addr(), __str_addr_2())


class UnoficialOpcodes:
    opcodes = [
        IGN,
        NOP, 
        SKB
    ]

    @staticmethod
    def all_commands():
        return flatten(
            map(lambda x: list(x.create_variations()), UnoficialOpcodes.opcodes)
        )
