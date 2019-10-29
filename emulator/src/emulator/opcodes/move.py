from more_itertools import flatten

from emulator.adressing import Immediate, ZeroPage, Absolute, ZeroPageY, AbsoluteY, IndirectY, ZeroPageX, AbsoluteX, IndirectX
from emulator.constants import NEGATIVE_BIT, LOW_BITS_MASK
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
                value = self.addressing_mode.read_from(cpu, memory, address,ld=1)
                if self.addressing_mode != Immediate:
                    self.addressing_mode.data = "= %02X" % memory.fetch(address)
                    cpu.addr = memory.get_effective_address(address)
                    cpu.data = value

                cpu.a = value
                cpu.zero = cpu.a == 0
                cpu.negative = (cpu.a & 0b10000000) > 0

        cycle_lda()


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
        def _stall():
            pass

        def cycle_sta():
            if self.addressing_mode:
                address = self.addressing_mode.fetch_address(cpu, memory)
                cpu.addr = memory.get_effective_address(address)
                cpu.data = cpu.a
                self.addressing_mode.data = "= %02X" % memory.fetch(address)
                self.addressing_mode.write_to(cpu, memory, address, cpu.a)

        if self.addressing_mode in [IndirectY, AbsoluteY, AbsoluteX]:
            # FIXME: this is ugly, but it works
            cycle_start = cpu.cycle
            cycle_sta()
            cycle_end = cpu.cycle
            if (cycle_end - cycle_start) < (self.cycles - 1):
                cpu.exec_in_cycle(_stall)
        else:
            cycle_sta()


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
                value = self.addressing_mode.read_from(cpu, memory, address,ld=1)
                if self.addressing_mode != Immediate:
                    self.addressing_mode.data = "= %02X" % memory.fetch(address)
                    cpu.addr = memory.get_effective_address(address)
                    cpu.data = value

                cpu.x = value
                cpu.zero = cpu.x == 0
                cpu.negative = (cpu.x & 0b10000000) > 0

        cycle_ldx()


class STX(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x86, ZeroPage, 3,),
                      (0x8E, Absolute, 4,),
                      (0x96, ZeroPageY, 4,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def cycle_stx():
            """
            #  address R/W description
           --- ------- --- ------------------------------------------
            1    PC     R  fetch opcode, increment PC
            2    PC     R  fetch low byte of address, increment PC
            3    PC     R  fetch high byte of address, increment PC
            4  address  W  write register to effective address
            """
            if self.addressing_mode:
                address = self.addressing_mode.fetch_address(cpu, memory)
                cpu.addr = memory.get_effective_address(address)
                cpu.data = cpu.x
                self.addressing_mode.data = "= %02X" % memory.fetch(address)
                self.addressing_mode.write_to(cpu, memory, address, cpu.x)

        cycle_stx()


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
                value = self.addressing_mode.read_from(cpu, memory, address,ld=1)
                if self.addressing_mode != Immediate:
                    self.addressing_mode.data = "= %02X" % memory.fetch(address)
                    cpu.addr = memory.get_effective_address(address)
                    cpu.data = value

                cpu.y = value
                cpu.zero = cpu.y == 0
                cpu.negative = (cpu.y & 0b10000000) > 0

        cycle_ldy()


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
                cpu.addr = memory.get_effective_address(address)
                cpu.data = cpu.y
                self.addressing_mode.data = "= %02X" % memory.fetch(address)
                self.addressing_mode.write_to(cpu, memory, address, cpu.y)

        cycle_sty()


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
        cpu.x = (cpu.sp & LOW_BITS_MASK)
        cpu.zero = cpu.x == 0
        cpu.negative = (cpu.x & 0b10000000) > 0
        cpu.inc_cycle()


class TXS(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x9A, None, 2,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        cpu.sp = 0x0100 | (cpu.x & LOW_BITS_MASK)
        cpu.inc_cycle()


class PLA(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x68, None, 4,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        cpu.sp += 1
        cpu.sp = cpu.sp & 0xff ^ 0x0100
        from_stack = memory.fetch(cpu.sp)
        cpu.a = from_stack
        cpu.zero = (cpu.a == 0)
        cpu.negative = (cpu.a & NEGATIVE_BIT) > 0
        cpu.inc_cycle()
        cpu.inc_cycle()
        cpu.inc_cycle()


class PHA(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x48, None, 3,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        memory.store(cpu.sp, cpu.a)
        cpu.sp -= 1
        cpu.sp = cpu.sp & 0xff ^ 0x0100
        cpu.inc_cycle()
        cpu.inc_cycle()


class PLP(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x28, None, 4,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        cpu.sp += 1
        cpu.sp = cpu.sp & 0xff ^ 0x0100
        from_stack = memory.fetch(cpu.sp)
        status = (from_stack & 0b11101111) | 0b00100000
        cpu.flags = status
        cpu.inc_cycle()
        cpu.inc_cycle()
        cpu.inc_cycle()


class PHP(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x08, None, 3,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        """
        In the byte pushed, bit 5 is always set to 1, and bit 4 is 1 if from an instruction (PHP or BRK) or 0 if from an interrupt line being pulled low (/IRQ or /NMI).
        """
        status = cpu.flags | 0b00110000
        memory.store(cpu.sp, status)
        cpu.sp -= 1
        cpu.sp = cpu.sp & 0xff ^ 0x0100
        cpu.inc_cycle()
        cpu.inc_cycle()


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
