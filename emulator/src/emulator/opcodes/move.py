from more_itertools import flatten

from emulator.adressing import Immediate, ZeroPage, Absolute, ZeroPageY, AbsoluteY, IndirectY, ZeroPageX, AbsoluteX, IndirectX
from emulator.opcodes.base import OpCode


class LDA(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xA1, IndirectX, 6),
                      (0xA5, ZeroPage, 3,),
                      (0xA9, Immediate, 2,),
                      (0xAD, Absolute, 4,),
                      (0xB1, IndirectY, 5),
                      (0xB5, ZeroPageX, 4,),
                      (0xB9, AbsoluteY, 4,),
                      (0xBD, AbsoluteX, 4,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def cycle_lda(): 
            if self.addressing_mode:
                address = self.addressing_mode.fetch_address(cpu, memory)
                value = self.addressing_mode.read_from(cpu, memory, address)

                cpu.a = value
                cpu.zero = cpu.a == 0 
                cpu.negative = (cpu.a & 0b10000000) > 0

        cpu.exec_in_cycle(cycle_lda)


class STA(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x81, IndirectX, 6,),
                      (0x85, ZeroPage, 3,),
                      (0x8D, Absolute, 4,),
                      (0x91, IndirectY, 6,),
                      (0x95, ZeroPageX, 4,),
                      (0x99, AbsoluteY, 5,),
                      (0x9D, AbsoluteX, 5,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def cycle_sta():
            if self.addressing_mode:
                address = self.addressing_mode.fetch_address(cpu, memory)
                self.addressing_mode.write_to(cpu, memory, address, cpu.a)
        cpu.exec_in_cycle(cycle_sta)


class LDX(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xA2, Immediate, 2,),
                      (0xA6, ZeroPage, 3,),
                      (0xAE, Absolute, 4,),
                      (0xB6, ZeroPageY, 4,),
                      (0xBE, AbsoluteY, 4,)]
        return map(cls.create_dict_entry, variations)
    
    def exec(self, cpu, memory):
        def cycle_ldx(): 
            if self.addressing_mode:
                address = self.addressing_mode.fetch_address(cpu, memory)
                value = self.addressing_mode.read_from(cpu, memory, address)

                cpu.x = value
                cpu.zero = cpu.x == 0 
                cpu.negative = (cpu.x & 0b10000000) > 0

        cpu.exec_in_cycle(cycle_ldx)

class STX(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x86, ZeroPage, 3,),
                      (0x8E, Absolute, 4,),
                      (0x96, ZeroPageY, 4,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def cycle_stx():
            if self.addressing_mode:
                address = self.addressing_mode.fetch_address(cpu, memory)
                self.addressing_mode.write_to(cpu, memory, address, cpu.x)
        cpu.exec_in_cycle(cycle_stx)

class LDY(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xA0, Immediate, 2,),
                      (0xA4, ZeroPage, 3,),
                      (0xAC, Absolute, 4,),
                      (0xB4, ZeroPageX, 4,),
                      (0xBC, AbsoluteX, 4,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def cycle_ldy(): 
            if self.addressing_mode:
                address = self.addressing_mode.fetch_address(cpu, memory)
                value = self.addressing_mode.read_from(cpu, memory, address)

                cpu.y = value
                cpu.zero = cpu.y == 0 
                cpu.negative = (cpu.y & 0b10000000) > 0

        cpu.exec_in_cycle(cycle_ldy)


class STY(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x84, ZeroPage, 3,),
                      (0x8C, Absolute, 4,),
                      (0x94, ZeroPageX, 4,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def cycle_sty():
            if self.addressing_mode:
                address = self.addressing_mode.fetch_address(cpu, memory)
                self.addressing_mode.write_to(cpu, memory, address, cpu.y)
        cpu.exec_in_cycle(cycle_sty)

class TAX(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xAA, None, 2,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def cycle_tax():
            cpu.x = cpu.a
            cpu.zero = cpu.x == 0
            cpu.negative = (cpu.x & 0b10000000) > 0
        cpu.exec_in_cycle(cycle_tax)


class TXA(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x8A, None, 2,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def cycle_txa():
            cpu.a = cpu.x
            cpu.zero = cpu.a == 0
            cpu.negative = (cpu.a & 0b10000000) > 0
        cpu.exec_in_cycle(cycle_txa)


class TAY(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xA8, None, 2,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def cycle_tay():
            cpu.y = cpu.a
            cpu.zero = cpu.y == 0
            cpu.negative = (cpu.y & 0b10000000) > 0
        cpu.exec_in_cycle(cycle_tay)


class TYA(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x98, None, 2,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def cycle_tya():
            cpu.a = cpu.y
            cpu.zero = cpu.a == 0
            cpu.negative = (cpu.a & 0b10000000) > 0
        cpu.exec_in_cycle(cycle_tya)


class TSX(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xBA, None, 2,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        cpu.x = cpu.sp
        cpu.zero = cpu.x == 0
        cpu.negative = (cpu.x & 0b10000000) > 0
        cpu.inc_cycle()


class TXS(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x9A, None, 2,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        cpu.sp = cpu.x
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
