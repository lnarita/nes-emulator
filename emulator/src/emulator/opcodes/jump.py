from more_itertools import flatten

from emulator.adressing import Relative, Absolute, Indirect
from emulator.opcodes.base import OpCode
from emulator.constants import HIGH_BITS_MASK, LOW_BITS_MASK


class BPL(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x10, Relative, 2,)]
        return map(cls.create_dict_entry, variations)
    def exec(self, cpu, memory):
        if self.addressing_mode:
            def _add_cycle():
                def _inc_pc(value):
                    return cpu.pc + value
                if not cpu.negative:
                    address = self.addressing_mode.fetch_address(cpu, memory)
                    value = self.addressing_mode.read_from(cpu, memory, address)
                    cpu.pc = cpu.exec_in_cycle(_inc_pc, value)
            cpu.exec_in_cycle(_add_cycle)

class BMI(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x30, Relative, 2,)]
        return map(cls.create_dict_entry, variations)
    def exec(self, cpu, memory):
        if self.addressing_mode:
            def _add_cycle():
                def _inc_pc(value):
                    return cpu.pc + value
                if cpu.negative:
                    address = self.addressing_mode.fetch_address(cpu, memory)
                    value = self.addressing_mode.read_from(cpu, memory, address)
                    cpu.pc = cpu.exec_in_cycle(_inc_pc, value)
            cpu.exec_in_cycle(_add_cycle)

class BVC(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x50, Relative, 2,)]
        return map(cls.create_dict_entry, variations)
    def exec(self, cpu, memory):
        if self.addressing_mode:
            def _add_cycle():
                def _inc_pc(value):
                    return cpu.pc + value
                if not cpu.overflow:
                    address = self.addressing_mode.fetch_address(cpu, memory)
                    value = self.addressing_mode.read_from(cpu, memory, address)
                    cpu.pc = cpu.exec_in_cycle(_inc_pc, value)
            cpu.exec_in_cycle(_add_cycle)

class BVS(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x70, Relative, 2,)]
        return map(cls.create_dict_entry, variations)
    def exec(self, cpu, memory):
        if self.addressing_mode:
            def _add_cycle():
                def _inc_pc(value):
                    return cpu.pc + value
                if cpu.overflow:
                    address = self.addressing_mode.fetch_address(cpu, memory)
                    value = self.addressing_mode.read_from(cpu, memory, address)
                    cpu.pc = cpu.exec_in_cycle(_inc_pc, value)
            cpu.exec_in_cycle(_add_cycle)

class BCC(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x90, Relative, 2,)]
        return map(cls.create_dict_entry, variations)
    def exec(self, cpu, memory):
        if self.addressing_mode:
            def _add_cycle():
                def _inc_pc(value):
                    return cpu.pc + value
                if not cpu.carry:
                    address = self.addressing_mode.fetch_address(cpu, memory)
                    value = self.addressing_mode.read_from(cpu, memory, address)
                    cpu.pc = cpu.exec_in_cycle(_inc_pc, value)
            cpu.exec_in_cycle(_add_cycle)

class BCS(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xB0, Relative, 2,)]
        return map(cls.create_dict_entry, variations)
    def exec(self, cpu, memory):
        if self.addressing_mode:
            def _add_cycle():
                def _inc_pc(value):
                    return cpu.pc + value
                if cpu.carry:
                    address = self.addressing_mode.fetch_address(cpu, memory)
                    value = self.addressing_mode.read_from(cpu, memory, address)
                    cpu.pc = cpu.exec_in_cycle(_inc_pc, value)
            cpu.exec_in_cycle(_add_cycle)

class BNE(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xD0, Relative, 2,)]
        return map(cls.create_dict_entry, variations)
    def exec(self, cpu, memory):
        if self.addressing_mode:
            def _add_cycle():
                def _inc_pc(value):
                    return cpu.pc + value
                if not cpu.zero:
                    address = self.addressing_mode.fetch_address(cpu, memory)
                    value = self.addressing_mode.read_from(cpu, memory, address)
                    cpu.pc = cpu.exec_in_cycle(_inc_pc, value)
            cpu.exec_in_cycle(_add_cycle)

class BEQ(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0xF0, Relative, 2,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        if self.addressing_mode:
            def _add_cycle():
                def _inc_pc(value):
                    return cpu.pc + value
                if cpu.zero:
                    address = self.addressing_mode.fetch_address(cpu, memory)
                    value = self.addressing_mode.read_from(cpu, memory, address)

                    cpu.pc = cpu.exec_in_cycle(_inc_pc, value)
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

        status = 0b00000000
        status |= cpu.negative and 0b10000000
        status |= cpu.overflow and 0b01000000
        status |= True and 0b00100000
        status |= cpu.break_command and 0b00010000
        status |= cpu.decimal and 0b00001000
        status |= cpu.interrupts_disabled and 0b00000100
        status |= cpu.zero and 0b00000010
        status |= cpu.carry and 0b00000001

        memory.store(cpu.sp, status)
        cpu.sp -= 1

        cpu.pc = address
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
        variations = [(0x20, Absolute, 6,)]
        return map(cls.create_dict_entry, variations)
    def exec(self, cpu, memory):
        if self.addressing_mode:
            returnAdress = cpu.pc
            self.addressing_mode.write_to(cpu, memory, cpu.sp, returnAdress)
            cpu.sp -= 1
            
            destinationAddress = self.addressing_mode.fetch_address(cpu, memory)
            cpu.pc = destinationAddress

class RTS(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x60, None, 6,)]
        return map(cls.create_dict_entry, variations)
    def exec(self, cpu, memory):
        if self.addressing_mode:
            cpu.sp += 1
            returnAdress = self.addressing_mode.read_from(cpu, memory, cpu.sp)
            cpu.pc = returnAdress

class JMP(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x4C, Absolute, 3,),
                      (0x6C, Indirect, 5,)]
        return map(cls.create_dict_entry, variations)

    def exec(self, cpu, memory):
        def cycle_jmp():
            if self.addressing_mode:
                jmpAddress = self.addressing_mode.fetch_address(cpu, memory)
                cpu.pc = jmpAddress
        cpu.exec_in_cycle(cycle_jmp)

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
