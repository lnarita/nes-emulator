import numpy as np
from more_itertools import flatten

from emulator.adressing import Relative, Absolute, Indirect, AddressMode
from emulator.constants import HIGH_BITS_MASK, LOW_BITS_MASK
from emulator.opcodes.base import OpCode

"""
Branching cycles:

 #   address  R/W description
--- --------- --- ---------------------------------------------
1     PC      R  fetch opcode, increment PC
2     PC      R  fetch operand, increment PC
3     PC      R  Fetch opcode of next instruction,
                 If branch is taken, add operand to PCL.
                 Otherwise increment PC.
4+    PC*     R  Fetch opcode of next instruction.
                 Fix PCH. If it did not change, increment PC.
5!    PC      R  Fetch opcode of next instruction,
                 increment PC.

Notes: The opcode fetch of the next instruction is included to
      this diagram for illustration purposes. When determining
      real execution times, remember to subtract the last
      cycle.

      * The high byte of Program Counter (PCH) may be invalid
        at this time, i.e. it may be smaller or bigger by $100.

      + If branch is taken, this cycle will be executed.

      ! If branch occurs to different page, this cycle will be
        executed.
"""


def _calculate_new_pc(cpu, operand):
    new_pc = int(np.uint16(np.int16(cpu.pc) + np.int16(np.int8(operand))))
    overflow = (cpu.pc & HIGH_BITS_MASK) != (new_pc & HIGH_BITS_MASK)
    return new_pc, overflow


class BPL(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x10, Relative, 2,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        operand = self.addressing_mode.fetch_address(cpu, memory)
        next_instruction = memory.fetch(cpu.pc)
        should_take_branch = not cpu.negative
        overflow = False
        if should_take_branch:
            new_pc, overflow = _calculate_new_pc(cpu, operand)
            cpu.exec_in_cycle()
        else:
            new_pc = cpu.pc + 1
        self.addressing_mode.addr = "$%04X" % (_calculate_new_pc(cpu, operand)[0])
        new_next_instruction = memory.fetch(new_pc)
        if should_take_branch:
            cpu.pc = new_pc
            if overflow:
                cpu.exec_in_cycle()

        cpu.exec_in_cycle()


class BMI(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x30, Relative, 2,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        operand = self.addressing_mode.fetch_address(cpu, memory)
        next_instruction = memory.fetch(cpu.pc)
        should_take_branch = cpu.negative
        overflow = False
        if should_take_branch:
            new_pc, overflow = _calculate_new_pc(cpu, operand)
            cpu.exec_in_cycle()
        else:
            new_pc = cpu.pc + 1
        self.addressing_mode.addr = "$%04X" % (_calculate_new_pc(cpu, operand)[0])
        new_next_instruction = memory.fetch(new_pc)
        if should_take_branch:
            cpu.pc = new_pc
            if overflow:
                cpu.exec_in_cycle()

        cpu.exec_in_cycle()


class BVC(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x50, Relative, 2,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        operand = self.addressing_mode.fetch_address(cpu, memory)
        next_instruction = memory.fetch(cpu.pc)
        should_take_branch = not cpu.overflow
        overflow = False
        if should_take_branch:
            new_pc, overflow = _calculate_new_pc(cpu, operand)
            cpu.exec_in_cycle()
        else:
            new_pc = cpu.pc + 1
        self.addressing_mode.addr = "$%04X" % (_calculate_new_pc(cpu, operand)[0])
        new_next_instruction = memory.fetch(new_pc)
        if should_take_branch:
            cpu.pc = new_pc
            if overflow:
                cpu.exec_in_cycle()

        cpu.exec_in_cycle()


class BVS(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x70, Relative, 2,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        operand = self.addressing_mode.fetch_address(cpu, memory)
        next_instruction = memory.fetch(cpu.pc)
        should_take_branch = cpu.overflow
        overflow = False
        if should_take_branch:
            new_pc, overflow = _calculate_new_pc(cpu, operand)
            cpu.exec_in_cycle()
        else:
            new_pc = cpu.pc + 1
        self.addressing_mode.addr = "$%04X" % (_calculate_new_pc(cpu, operand)[0])
        new_next_instruction = memory.fetch(new_pc)
        if should_take_branch:
            cpu.pc = new_pc
            if overflow:
                cpu.exec_in_cycle()

        cpu.exec_in_cycle()


class BCC(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x90, Relative, 2,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        operand = self.addressing_mode.fetch_address(cpu, memory)
        next_instruction = memory.fetch(cpu.pc)
        should_take_branch = not cpu.carry
        overflow = False
        if should_take_branch:
            new_pc, overflow = _calculate_new_pc(cpu, operand)
            cpu.exec_in_cycle()
        else:
            new_pc = cpu.pc + 1
        self.addressing_mode.addr = "$%04X" % (_calculate_new_pc(cpu, operand)[0])
        new_next_instruction = memory.fetch(new_pc)
        if should_take_branch:
            cpu.pc = new_pc
            if overflow:
                cpu.exec_in_cycle()

        cpu.exec_in_cycle()


class BCS(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xB0, Relative, 2,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        operand = self.addressing_mode.fetch_address(cpu, memory)
        next_instruction = memory.fetch(cpu.pc)
        should_take_branch = cpu.carry
        overflow = False
        if should_take_branch:
            new_pc, overflow = _calculate_new_pc(cpu, operand)
            cpu.exec_in_cycle()
        else:
            new_pc = cpu.pc + 1
        self.addressing_mode.addr = "$%04X" % (_calculate_new_pc(cpu, operand)[0])
        new_next_instruction = memory.fetch(new_pc)
        if should_take_branch:
            cpu.pc = new_pc
            if overflow:
                cpu.exec_in_cycle()

        cpu.exec_in_cycle()


class BNE(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xD0, Relative, 2,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        operand = self.addressing_mode.fetch_address(cpu, memory)
        next_instruction = memory.fetch(cpu.pc)
        should_take_branch = not cpu.zero
        overflow = False
        if should_take_branch:
            new_pc, overflow = _calculate_new_pc(cpu, operand)
            cpu.exec_in_cycle()
        else:
            new_pc = cpu.pc + 1
        self.addressing_mode.addr = "$%04X" % (_calculate_new_pc(cpu, operand)[0])
        new_next_instruction = memory.fetch(new_pc)
        if should_take_branch:
            cpu.pc = new_pc
            if overflow:
                cpu.exec_in_cycle()

        cpu.exec_in_cycle()


class BEQ(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xF0, Relative, 2,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        operand = self.addressing_mode.fetch_address(cpu, memory)
        next_instruction = memory.fetch(cpu.pc)
        should_take_branch = cpu.zero
        overflow = False
        if should_take_branch:
            new_pc, overflow = _calculate_new_pc(cpu, operand)
            cpu.exec_in_cycle()
        else:
            new_pc = cpu.pc + 1
        self.addressing_mode.addr = "$%04X" % (_calculate_new_pc(cpu, operand)[0])
        new_next_instruction = memory.fetch(new_pc)
        if should_take_branch:
            cpu.pc = new_pc
            if overflow:
                cpu.exec_in_cycle()

        cpu.exec_in_cycle()


class BRK(OpCode):
    """
    The BRK instruction forces the generation of an interrupt request.
    The program counter and processor status are pushed on the stack then the IRQ interrupt vector at $FFFE/F is loaded into the PC and the break flag in the status set to one.
    """

    @classmethod
    def create_variations(cls):
        variations = [(0x00, None, 7,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        memory.stack_push(cpu, (cpu.pc & HIGH_BITS_MASK) >> 8)
        cpu.exec_in_cycle()
        memory.stack_push(cpu, (cpu.pc & LOW_BITS_MASK))
        cpu.exec_in_cycle()

        status = cpu.flags
        memory.stack_push(cpu, (status | 0b00010000))
        cpu.exec_in_cycle()

        pcl = Absolute.read_from(cpu, memory, 0xFFFE)
        pch = Absolute.read_from(cpu, memory, 0xFFFF)
        new_pc = AddressMode.get_16_bits_addr_from_high_low((pch << 8), pcl)

        cpu.pc = new_pc
        cpu.break_command = True


class RTI(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x40, None, 6,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        """
        #  address R/W description
       --- ------- --- -----------------------------------------------
        1    PC     R  fetch opcode, increment PC
        2    PC     R  read next instruction byte (and throw it away)
        3  $0100,S  R  increment S
        4  $0100,S  R  pull P from stack, increment S
        5  $0100,S  R  pull PCL from stack, increment S
        6  $0100,S  R  pull PCH from stack
        """
        from_stack = memory.stack_pop(cpu)
        cpu.exec_in_cycle()
        cpu.exec_in_cycle()
        status = (from_stack & 0b11101111) | 0b00100000
        cpu.flags = status
        cpu.exec_in_cycle()
        
        pcl = memory.stack_pop(cpu)
        cpu.exec_in_cycle()
        pch =memory.stack_pop(cpu)
        cpu.exec_in_cycle()
        cpu.pc = AddressMode.get_16_bits_addr_from_high_low((pch << 8), pcl)


class JSR(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x20, Absolute, 6,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        """
          #  address R/W description
       --- ------- --- -------------------------------------------------
        1    PC     R  fetch opcode, increment PC
        2    PC     R  fetch low address byte, increment PC
        3  $0100,S  R  internal operation (predecrement S?)
        4  $0100,S  W  push PCH on stack, decrement S
        5  $0100,S  W  push PCL on stack, decrement S
        6    PC     R  copy low address byte to PCL, fetch high address
                       byte to PCH
        """
        addr = self.addressing_mode.fetch_address(cpu, memory)
        memory.stack_push(cpu, ((cpu.pc - 1) & HIGH_BITS_MASK) >> 8)
        cpu.exec_in_cycle()
        memory.stack_push(cpu, ((cpu.pc - 1) & LOW_BITS_MASK))
        cpu.exec_in_cycle()
        cpu.pc = addr

        cpu.exec_in_cycle()


class RTS(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x60, None, 6,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        """
        #  address R/W description
       --- ------- --- -----------------------------------------------
        1    PC     R  fetch opcode, increment PC
        2    PC     R  read next instruction byte (and throw it away)
        3  $0100,S  R  increment S
        4  $0100,S  R  pull PCL from stack, increment S
        5  $0100,S  R  pull PCH from stack
        6    PC     R  increment PC
        """
        addr_low = memory.stack_pop(cpu)
        cpu.exec_in_cycle()
        addr_high = memory.stack_pop(cpu)
        cpu.exec_in_cycle()
        addr_high = addr_high << 8
        addr = AddressMode.get_16_bits_addr_from_high_low(addr_high, addr_low)
        cpu.pc = addr + 1
        cpu.exec_in_cycle()
        cpu.exec_in_cycle()
        cpu.exec_in_cycle()


class JMP(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x4C, Absolute, 3,),
                      (0x6C, Indirect, 5,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        if self.addressing_mode:
            jmp_address = self.addressing_mode.fetch_address(cpu, memory)
            cpu.pc = jmp_address

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
