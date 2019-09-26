from more_itertools import flatten

from emulator.adressing import Relative, Absolute, Indirect
from emulator.opcodes.base import OpCode


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
            self.adressing_mode.write_to(cpu, memory, cpu.sp, returnAdress)
            cpu.sp -= 1
            
            destinationAddress = self.addressing_mode.read_from(cpu, memory, address)
            cpu.pc = destinationAddress

class RTS(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x60, None, 6,)]
        return map(cls.create_dict_entry, variations)
    def exec(self, cpu, memory):
        if self.addressing_mode:
            cpu.sp += 1
            returnAdress = self.adressing_mode.read_from(cpu, memory, cpu.sp)
            cpu.pc = returnAdress

class JMP(OpCode):
    @classmethod
    def create_variations(cls):
        variations = [(0x4C, Absolute, 3,),
                      (0x6C, Indirect, 5,)]
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
