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


class BPL(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x10, Relative, 2,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def _stall():
            pass

        def _calculate_new_pc(operand):
            pcl = ((cpu.pc & LOW_BITS_MASK) + operand)
            overflow = (pcl & HIGH_BITS_MASK) > 0
            pcl &= LOW_BITS_MASK
            new_pc = (cpu.pc & HIGH_BITS_MASK) | pcl
            return new_pc, overflow

        def _add_cycle():
            operand = self.addressing_mode.fetch_address(cpu, memory)
            next_instruction = memory.fetch(cpu.pc)
            should_take_branch = not cpu.negative
            overflow = False
            if should_take_branch:
                new_pc, overflow = cpu.exec_in_cycle(_calculate_new_pc, operand)
            else:
                new_pc = cpu.pc + 1
            if ((cpu.pc + operand) & 0xF0000) > 0:
                self.addressing_mode.addr = "$%04X" % (new_pc)
            else:
                self.addressing_mode.addr = "$%04X" % (cpu.pc + operand)
            new_next_instruction = memory.fetch(new_pc)
            if should_take_branch:
                cpu.inc_pc_by(operand)
                if ((cpu.pc + operand) & 0xF0000) > 0:
                    cpu.pc = new_pc
                if overflow:
                    cpu.exec_in_cycle(_stall)

        cpu.exec_in_cycle(_add_cycle)


class BMI(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x30, Relative, 2,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def _stall():
            pass

        def _calculate_new_pc(operand):
            pcl = ((cpu.pc & LOW_BITS_MASK) + operand)
            overflow = (pcl & HIGH_BITS_MASK) > 0
            pcl &= LOW_BITS_MASK
            new_pc = (cpu.pc & HIGH_BITS_MASK) | pcl
            return new_pc, overflow

        def _add_cycle():
            operand = self.addressing_mode.fetch_address(cpu, memory)
            next_instruction = memory.fetch(cpu.pc)
            should_take_branch = cpu.negative
            overflow = False
            if should_take_branch:
                new_pc, overflow = cpu.exec_in_cycle(_calculate_new_pc, operand)
            else:
                new_pc = cpu.pc + 1
            if ((cpu.pc + operand) & 0xF0000) > 0:
                self.addressing_mode.addr = "$%04X" % (new_pc)
            else:
                self.addressing_mode.addr = "$%04X" % (cpu.pc + operand)
            new_next_instruction = memory.fetch(new_pc)
            if should_take_branch:
                cpu.inc_pc_by(operand)
                if ((cpu.pc + operand) & 0xF0000) > 0:
                    cpu.pc = new_pc
                if overflow:
                    cpu.exec_in_cycle(_stall)

        cpu.exec_in_cycle(_add_cycle)


class BVC(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x50, Relative, 2,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def _stall():
            pass

        def _calculate_new_pc(operand):
            pcl = ((cpu.pc & LOW_BITS_MASK) + operand)
            overflow = (pcl & HIGH_BITS_MASK) > 0
            pcl &= LOW_BITS_MASK
            new_pc = (cpu.pc & HIGH_BITS_MASK) | pcl
            return new_pc, overflow

        def _add_cycle():
            operand = self.addressing_mode.fetch_address(cpu, memory)
            next_instruction = memory.fetch(cpu.pc)
            should_take_branch = not cpu.overflow
            overflow = False
            if should_take_branch:
                new_pc, overflow = cpu.exec_in_cycle(_calculate_new_pc, operand)
            else:
                new_pc = cpu.pc + 1
            if ((cpu.pc + operand) & 0xF0000) > 0:
                self.addressing_mode.addr = "$%04X" % (new_pc)
            else:
                self.addressing_mode.addr = "$%04X" % (cpu.pc + operand)
            new_next_instruction = memory.fetch(new_pc)
            if should_take_branch:
                cpu.inc_pc_by(operand)
                if ((cpu.pc + operand) & 0xF0000) > 0:
                    cpu.pc = new_pc
                if overflow:
                    cpu.exec_in_cycle(_stall)

        cpu.exec_in_cycle(_add_cycle)


class BVS(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x70, Relative, 2,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def _stall():
            pass

        def _calculate_new_pc(operand):
            pcl = ((cpu.pc & LOW_BITS_MASK) + operand)
            overflow = (pcl & HIGH_BITS_MASK) > 0
            pcl &= LOW_BITS_MASK
            new_pc = (cpu.pc & HIGH_BITS_MASK) | pcl
            return new_pc, overflow

        def _add_cycle():
            operand = self.addressing_mode.fetch_address(cpu, memory)
            next_instruction = memory.fetch(cpu.pc)
            should_take_branch = cpu.overflow
            overflow = False
            if should_take_branch:
                new_pc, overflow = cpu.exec_in_cycle(_calculate_new_pc, operand)
            else:
                new_pc = cpu.pc + 1
            if ((cpu.pc + operand) & 0xF0000) > 0:
                self.addressing_mode.addr = "$%04X" % (new_pc)
            else:
                self.addressing_mode.addr = "$%04X" % (cpu.pc + operand)
            new_next_instruction = memory.fetch(new_pc)
            if should_take_branch:
                cpu.inc_pc_by(operand)
                if ((cpu.pc + operand) & 0xF0000) > 0:
                    cpu.pc = new_pc
                if overflow:
                    cpu.exec_in_cycle(_stall)

        cpu.exec_in_cycle(_add_cycle)


class BCC(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x90, Relative, 2,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def _stall():
            pass

        def _calculate_new_pc(operand):
            pcl = ((cpu.pc & LOW_BITS_MASK) + operand)
            overflow = (pcl & HIGH_BITS_MASK) > 0
            pcl &= LOW_BITS_MASK
            new_pc = (cpu.pc & HIGH_BITS_MASK) | pcl
            return new_pc, overflow

        def _add_cycle():
            operand = self.addressing_mode.fetch_address(cpu, memory)
            next_instruction = memory.fetch(cpu.pc)
            should_take_branch = not cpu.carry
            overflow = False
            if should_take_branch:
                new_pc, overflow = cpu.exec_in_cycle(_calculate_new_pc, operand)
            else:
                new_pc = cpu.pc + 1
            if ((cpu.pc + operand) & 0xF0000) > 0:
                self.addressing_mode.addr = "$%04X" % (new_pc)
            else:
                self.addressing_mode.addr = "$%04X" % (cpu.pc + operand)
            new_next_instruction = memory.fetch(new_pc)
            if should_take_branch:
                cpu.inc_pc_by(operand)
                if ((cpu.pc + operand) & 0xF0000) > 0:
                    cpu.pc = new_pc
                if overflow:
                    cpu.exec_in_cycle(_stall)

        cpu.exec_in_cycle(_add_cycle)


class BCS(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xB0, Relative, 2,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def _stall():
            pass

        def _calculate_new_pc(operand):
            pcl = ((cpu.pc & LOW_BITS_MASK) + operand)
            overflow = (pcl & HIGH_BITS_MASK) > 0
            pcl &= LOW_BITS_MASK
            new_pc = (cpu.pc & HIGH_BITS_MASK) | pcl
            return new_pc, overflow

        def _add_cycle():
            operand = self.addressing_mode.fetch_address(cpu, memory)
            next_instruction = memory.fetch(cpu.pc)
            should_take_branch = cpu.carry
            overflow = False
            if should_take_branch:
                new_pc, overflow = cpu.exec_in_cycle(_calculate_new_pc, operand)
            else:
                new_pc = cpu.pc + 1
            if ((cpu.pc + operand) & 0xF0000) > 0:
                self.addressing_mode.addr = "$%04X" % (new_pc)
            else:
                self.addressing_mode.addr = "$%04X" % (cpu.pc + operand)
            new_next_instruction = memory.fetch(new_pc)
            if should_take_branch:
                cpu.inc_pc_by(operand)
                if ((cpu.pc + operand) & 0xF0000) > 0:
                    cpu.pc = new_pc
                if overflow:
                    cpu.exec_in_cycle(_stall)

        cpu.exec_in_cycle(_add_cycle)


class BNE(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xD0, Relative, 2,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def _stall():
            pass

        def _calculate_new_pc(operand):
            pcl = ((cpu.pc & LOW_BITS_MASK) + operand)
            overflow = (pcl & HIGH_BITS_MASK) > 0
            pcl &= LOW_BITS_MASK
            new_pc = (cpu.pc & HIGH_BITS_MASK) | pcl
            return new_pc, overflow

        def _add_cycle():
            operand = self.addressing_mode.fetch_address(cpu, memory)
            next_instruction = memory.fetch(cpu.pc)
            should_take_branch = not cpu.zero
            overflow = False
            if should_take_branch:
                new_pc, overflow = cpu.exec_in_cycle(_calculate_new_pc, operand)
            else:
                new_pc = cpu.pc + 1
            if ((cpu.pc + operand) & 0xF0000) > 0:
                self.addressing_mode.addr = "$%04X" % (new_pc)
            else:
                self.addressing_mode.addr = "$%04X" % (cpu.pc + operand)
            new_next_instruction = memory.fetch(new_pc)
            if should_take_branch:
                cpu.inc_pc_by(operand)
                if ((cpu.pc + operand) & 0xF0000) > 0:
                    cpu.pc = new_pc
                if overflow:
                    cpu.exec_in_cycle(_stall)

        cpu.exec_in_cycle(_add_cycle)


class BEQ(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xF0, Relative, 2,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def _stall():
            pass

        def _calculate_new_pc(operand):
            pcl = ((cpu.pc & LOW_BITS_MASK) + operand)
            overflow = (pcl & HIGH_BITS_MASK) > 0
            pcl &= LOW_BITS_MASK
            new_pc = (cpu.pc & HIGH_BITS_MASK) | pcl
            return new_pc, overflow

        def _add_cycle():
            operand = self.addressing_mode.fetch_address(cpu, memory)
            next_instruction = memory.fetch(cpu.pc)
            should_take_branch = cpu.zero
            overflow = False
            if should_take_branch:
                new_pc, overflow = cpu.exec_in_cycle(_calculate_new_pc, operand)
            else:
                new_pc = cpu.pc + 1
            if ((cpu.pc + operand) & 0xF0000) > 0:
                self.addressing_mode.addr = "$%04X" % (new_pc)
            else:
                self.addressing_mode.addr = "$%04X" % (cpu.pc + operand)
            new_next_instruction = memory.fetch(new_pc)
            if should_take_branch:
                cpu.inc_pc_by(operand)
                if ((cpu.pc + operand) & 0xF0000) > 0:
                    cpu.pc = new_pc
                if overflow:
                    cpu.exec_in_cycle(_stall)

        cpu.exec_in_cycle(_add_cycle)


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
        cpu.exec_in_cycle(memory.stack_push, cpu, (cpu.pc & HIGH_BITS_MASK) >> 8)
        cpu.exec_in_cycle(memory.stack_push, cpu, (cpu.pc & LOW_BITS_MASK))

        status = cpu.flags
        cpu.exec_in_cycle(memory.stack_push, cpu, (status | 0b00010000))

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
        def _stall():
            pass
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
        from_stack = cpu.exec_in_cycle(memory.stack_pop, cpu)
        cpu.exec_in_cycle(_stall)
        status = (from_stack & 0b11101111) | 0b00100000
        cpu.flags = status
        cpu.exec_in_cycle(_stall)
        pcl = cpu.exec_in_cycle(memory.stack_pop,cpu)
        pch = cpu.exec_in_cycle(memory.stack_pop,cpu)
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

        def _cycle():
            addr = self.addressing_mode.fetch_address(cpu, memory)
            cpu.exec_in_cycle(memory.stack_push, cpu, ((cpu.pc - 1) & HIGH_BITS_MASK) >> 8)
            cpu.exec_in_cycle(memory.stack_push, cpu, ((cpu.pc - 1) & LOW_BITS_MASK))
            cpu.pc = addr

        cpu.exec_in_cycle(_cycle)


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

        def _stall():
            pass

        addr_low = cpu.exec_in_cycle(memory.stack_pop, cpu)
        addr_high = cpu.exec_in_cycle(memory.stack_pop, cpu)
        addr_high = addr_high << 8
        addr = AddressMode.get_16_bits_addr_from_high_low(addr_high, addr_low)
        cpu.pc = addr + 1
        cpu.exec_in_cycle(_stall)
        cpu.exec_in_cycle(_stall)
        cpu.exec_in_cycle(_stall)


class JMP(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x4C, Absolute, 3,),
                      (0x6C, Indirect, 5,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def cycle_jmp():
            if self.addressing_mode:
                jmp_address = self.addressing_mode.fetch_address(cpu, memory)
                cpu.pc = jmp_address

        cycle_jmp()


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
