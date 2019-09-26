from more_itertools import flatten

from emulator.adressing import Relative, Absolute, Indirect, AddressMode
from emulator.constants import HIGH_BITS_MASK, LOW_BITS_MASK
from emulator.opcodes.base import OpCode
from emulator.constants import HIGH_BITS_MASK, LOW_BITS_MASK
from emulator.cpu import StatusRegisterFlags


class BPL(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x10, Relative, 2,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def _stall():
            pass

        def _add_cycle():
            operand = self.addressing_mode.fetch_address(cpu, memory)
            next_instruction = memory.fetch(cpu.pc)
            should_take_branch = not cpu.negative
            """
            3     PC      R  Fetch opcode of next instruction,
                             If branch is taken, add operand to PCL.
                             Otherwise increment PC.
            4+    PC*     R  Fetch opcode of next instruction.
                             Fix PCH. If it did not change, increment PC.
            5!    PC      R  Fetch opcode of next instruction,
                             increment PC.
            """
            pcl = cpu.pc
            overflow = False
            new_pc = cpu.pc
            self.addressing_mode.addr = "$%04X" % (cpu.pc + operand)
            if should_take_branch:
                pcl = ((cpu.pc & LOW_BITS_MASK) + operand)
                overflow = (pcl & HIGH_BITS_MASK) > 0
                pcl &= LOW_BITS_MASK
                new_pc = (cpu.pc & HIGH_BITS_MASK) | pcl
            else:
                new_pc = cpu.pc + 1
            new_next_instruction = memory.fetch(new_pc)
            if should_take_branch:
                if next_instruction == new_next_instruction:
                    cpu.inc_pc_by(1)
                else:
                    cpu.inc_pc_by(operand)
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

        def _add_cycle():
            operand = self.addressing_mode.fetch_address(cpu, memory)
            next_instruction = memory.fetch(cpu.pc)
            should_take_branch = cpu.negative
            """
            3     PC      R  Fetch opcode of next instruction,
                             If branch is taken, add operand to PCL.
                             Otherwise increment PC.
            4+    PC*     R  Fetch opcode of next instruction.
                             Fix PCH. If it did not change, increment PC.
            5!    PC      R  Fetch opcode of next instruction,
                             increment PC.
            """
            pcl = cpu.pc
            overflow = False
            new_pc = cpu.pc
            self.addressing_mode.addr = "$%04X" % (cpu.pc + operand)
            if should_take_branch:
                pcl = ((cpu.pc & LOW_BITS_MASK) + operand)
                overflow = (pcl & HIGH_BITS_MASK) > 0
                pcl &= LOW_BITS_MASK
                new_pc = (cpu.pc & HIGH_BITS_MASK) | pcl
            else:
                new_pc = cpu.pc + 1
            new_next_instruction = memory.fetch(new_pc)
            if should_take_branch:
                if next_instruction == new_next_instruction:
                    cpu.inc_pc_by(1)
                else:
                    cpu.inc_pc_by(operand)
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

        def _add_cycle():
            operand = self.addressing_mode.fetch_address(cpu, memory)
            next_instruction = memory.fetch(cpu.pc)
            should_take_branch = not cpu.overflow
            """
            3     PC      R  Fetch opcode of next instruction,
                             If branch is taken, add operand to PCL.
                             Otherwise increment PC.
            4+    PC*     R  Fetch opcode of next instruction.
                             Fix PCH. If it did not change, increment PC.
            5!    PC      R  Fetch opcode of next instruction,
                             increment PC.
            """
            pcl = cpu.pc
            overflow = False
            new_pc = cpu.pc
            self.addressing_mode.addr = "$%04X" % (cpu.pc + operand)
            if should_take_branch:
                pcl = ((cpu.pc & LOW_BITS_MASK) + operand)
                overflow = (pcl & HIGH_BITS_MASK) > 0
                pcl &= LOW_BITS_MASK
                new_pc = (cpu.pc & HIGH_BITS_MASK) | pcl
            else:
                new_pc = cpu.pc + 1
            new_next_instruction = memory.fetch(new_pc)
            if should_take_branch:
                if next_instruction == new_next_instruction:
                    cpu.inc_pc_by(1)
                else:
                    cpu.inc_pc_by(operand)
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

        def _add_cycle():
            operand = self.addressing_mode.fetch_address(cpu, memory)
            next_instruction = memory.fetch(cpu.pc)
            should_take_branch = cpu.overflow
            """
            3     PC      R  Fetch opcode of next instruction,
                             If branch is taken, add operand to PCL.
                             Otherwise increment PC.
            4+    PC*     R  Fetch opcode of next instruction.
                             Fix PCH. If it did not change, increment PC.
            5!    PC      R  Fetch opcode of next instruction,
                             increment PC.
            """
            pcl = cpu.pc
            overflow = False
            new_pc = cpu.pc
            self.addressing_mode.addr = "$%04X" % (cpu.pc + operand)
            if should_take_branch:
                pcl = ((cpu.pc & LOW_BITS_MASK) + operand)
                overflow = (pcl & HIGH_BITS_MASK) > 0
                pcl &= LOW_BITS_MASK
                new_pc = (cpu.pc & HIGH_BITS_MASK) | pcl
            else:
                new_pc = cpu.pc + 1
            new_next_instruction = memory.fetch(new_pc)
            if should_take_branch:
                if next_instruction == new_next_instruction:
                    cpu.inc_pc_by(1)
                else:
                    cpu.inc_pc_by(operand)
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

        def _add_cycle():
            operand = self.addressing_mode.fetch_address(cpu, memory)
            next_instruction = memory.fetch(cpu.pc)
            should_take_branch = not cpu.carry
            """
            3     PC      R  Fetch opcode of next instruction,
                             If branch is taken, add operand to PCL.
                             Otherwise increment PC.
            4+    PC*     R  Fetch opcode of next instruction.
                             Fix PCH. If it did not change, increment PC.
            5!    PC      R  Fetch opcode of next instruction,
                             increment PC.
            """
            pcl = cpu.pc
            overflow = False
            new_pc = cpu.pc
            self.addressing_mode.addr = "$%04X" % (cpu.pc + operand)
            if should_take_branch:
                pcl = ((cpu.pc & LOW_BITS_MASK) + operand)
                overflow = (pcl & HIGH_BITS_MASK) > 0
                pcl &= LOW_BITS_MASK
                new_pc = (cpu.pc & HIGH_BITS_MASK) | pcl
            else:
                new_pc = cpu.pc + 1
            new_next_instruction = memory.fetch(new_pc)
            if should_take_branch:
                if next_instruction == new_next_instruction:
                    cpu.inc_pc_by(1)
                else:
                    cpu.inc_pc_by(operand)
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

        def _add_cycle():
            operand = self.addressing_mode.fetch_address(cpu, memory)
            next_instruction = memory.fetch(cpu.pc)
            should_take_branch = cpu.carry
            """
            3     PC      R  Fetch opcode of next instruction,
                             If branch is taken, add operand to PCL.
                             Otherwise increment PC.
            4+    PC*     R  Fetch opcode of next instruction.
                             Fix PCH. If it did not change, increment PC.
            5!    PC      R  Fetch opcode of next instruction,
                             increment PC.
            """
            pcl = cpu.pc
            overflow = False
            new_pc = cpu.pc
            self.addressing_mode.addr = "$%04X" % (cpu.pc + operand)
            if should_take_branch:
                pcl = ((cpu.pc & LOW_BITS_MASK) + operand)
                overflow = (pcl & HIGH_BITS_MASK) > 0
                pcl &= LOW_BITS_MASK
                new_pc = (cpu.pc & HIGH_BITS_MASK) | pcl
            else:
                new_pc = cpu.pc + 1
            new_next_instruction = memory.fetch(new_pc)
            if should_take_branch:
                if next_instruction == new_next_instruction:
                    cpu.inc_pc_by(1)
                else:
                    cpu.inc_pc_by(operand)
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

        def _add_cycle():
            operand = self.addressing_mode.fetch_address(cpu, memory)
            next_instruction = memory.fetch(cpu.pc)
            should_take_branch = not cpu.zero
            """
            3     PC      R  Fetch opcode of next instruction,
                             If branch is taken, add operand to PCL.
                             Otherwise increment PC.
            4+    PC*     R  Fetch opcode of next instruction.
                             Fix PCH. If it did not change, increment PC.
            5!    PC      R  Fetch opcode of next instruction,
                             increment PC.
            """
            pcl = cpu.pc
            overflow = False
            new_pc = cpu.pc
            self.addressing_mode.addr = "$%04X" % (cpu.pc + operand)
            if should_take_branch:
                pcl = ((cpu.pc & LOW_BITS_MASK) + operand)
                overflow = (pcl & HIGH_BITS_MASK) > 0
                pcl &= LOW_BITS_MASK
                new_pc = (cpu.pc & HIGH_BITS_MASK) | pcl
            else:
                new_pc = cpu.pc + 1
            new_next_instruction = memory.fetch(new_pc)
            if should_take_branch:
                if next_instruction == new_next_instruction:
                    cpu.inc_pc_by(1)
                else:
                    cpu.inc_pc_by(operand)
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

        def _add_cycle():
            operand = self.addressing_mode.fetch_address(cpu, memory)
            next_instruction = memory.fetch(cpu.pc)
            should_take_branch = cpu.zero
            """
            3     PC      R  Fetch opcode of next instruction,
                             If branch is taken, add operand to PCL.
                             Otherwise increment PC.
            4+    PC*     R  Fetch opcode of next instruction.
                             Fix PCH. If it did not change, increment PC.
            5!    PC      R  Fetch opcode of next instruction,
                             increment PC.
            """
            pcl = cpu.pc
            overflow = False
            new_pc = cpu.pc
            self.addressing_mode.addr = "$%04X" % (cpu.pc + operand)
            if should_take_branch:
                pcl = ((cpu.pc & LOW_BITS_MASK) + operand)
                overflow = (pcl & HIGH_BITS_MASK) > 0
                pcl &= LOW_BITS_MASK
                new_pc = (cpu.pc & HIGH_BITS_MASK) | pcl
            else:
                new_pc = cpu.pc + 1
            new_next_instruction = memory.fetch(new_pc)
            if should_take_branch:
                if next_instruction == new_next_instruction:
                    cpu.inc_pc_by(1)
                else:
                    cpu.inc_pc_by(operand)
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
        def _read_16_bits_low(memory, addr):
            addr_low = memory.fetch(addr)
            return addr_low & LOW_BITS_MASK
        def _read_16_bits_high(memory, addr):
            addr_high = memory.fetch(addr)
            return (addr_high << 8) & HIGH_BITS_MASK
        def _read_ptr_low():
            low = _read_16_bits_low(memory, 0xfffe)
            cpu.inc_pc_by(1)
            return low
        def _read_ptr_high():
            high = _read_16_bits_high(memory, 0xffff)
            cpu.inc_pc_by(1)
            return high
        def _read_addr_low(addr):
            low = _read_16_bits_low(memory, addr)
            return low
        def _get_16_bits_addr_from_high_low(high, low):
            return high | low
        def _read_addr_high_and_return_address(addr, l):
            h = _read_16_bits_high(memory, addr + 1)
            value = _get_16_bits_addr_from_high_low(l, h)
            return value

        pointer_low = cpu.exec_in_cycle(_read_ptr_low)  # 2
        pointer_high = cpu.exec_in_cycle(_read_ptr_high)  # 3
        pointer = _get_16_bits_addr_from_high_low(pointer_low, pointer_high)
        address_low = cpu.exec_in_cycle(_read_addr_low, pointer)  # 4
        address = cpu.exec_in_cycle(_read_addr_high_and_return_address, pointer, address_low)  # 5

        memory.store(cpu.sp, cpu.pc)
        cpu.sp -= 1

        status = cpu.flags

        memory.store(cpu.sp, status)
        cpu.sp -= 1

        cpu.pc = address
        cpu.break_command = True


class RTI(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x40, None, 6,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        cpu.sp += 1
        cpu.sp = cpu.sp & 0xff ^ 0x0100
        from_stack = memory.fetch(cpu.sp)
        status = (from_stack & 0b11101111) | 0b00100000
        new_status = StatusRegisterFlags(int_value=status)
        cpu.negative = new_status.negative
        cpu.overflow = new_status.overflow
        cpu.break_command = new_status.break_command
        cpu.decimal = new_status.decimal
        cpu.interrupts_disabled = new_status.interrupts_disabled
        cpu.zero = new_status.zero
        cpu.carry = new_status.carry
        cpu.inc_cycle()
        cpu.inc_cycle()
        cpu.inc_cycle()
        cpu.sp += 1
        from_stack = memory.fetch(cpu.sp)
        cpu.pc = from_stack



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
        cpu.exec_in_cycle(memory.stack_push, cpu, ((cpu.pc - 1) & HIGH_BITS_MASK) >> 8)
        cpu.exec_in_cycle(memory.stack_push, cpu, ((cpu.pc - 1) & LOW_BITS_MASK))
        cpu.pc = addr


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
        addr_low = cpu.exec_in_cycle(memory.stack_pop, cpu)
        addr_high = cpu.exec_in_cycle(memory.stack_pop, cpu)
        addr_high = addr_high << 8
        addr = AddressMode.get_16_bits_addr_from_high_low(addr_high, addr_low)
        cpu.pc = addr + 1


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
