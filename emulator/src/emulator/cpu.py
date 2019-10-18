import time

import numpy as np

from emulator.constants import CYCLE_PERIOD, CYCLE_PERIOD_SIZE
from emulator.memory import MemoryPositions


class StatusRegisterFlags:
    def __init__(self, n=False, v=False, b=False, d=False, i=True, z=False, c=False, int_value=None):
        if int_value:
            self.negative = (int_value & 0b10000000) != 0
            self.overflow = (int_value & 0b01000000) != 0
            self.break_command = (int_value & 0b00010000) != 0
            self.decimal = (int_value & 0b00001000) != 0
            self.interrupts_disabled = (int_value & 0b00000100) != 0
            self.zero = (int_value & 0b00000010) != 0
            self.carry = (int_value & 0b00000001) != 0
        else:
            # N = negative flag (1 when result is negative)
            self.negative = n
            # V = overflow flag (1 on signed overflow)
            self.overflow = v
            # B = break flag (1 when interrupt was caused by a BRK)
            self.break_command = b
            # D = decimal flag (1 when CPU in BCD mode)
            self.decimal = d
            # I = IRQ flag (when 1, no interrupts will occur (exceptions are IRQs forced by BRK and NMIs))
            self.interrupts_disabled = i
            # Z = zero flag (1 when all bits of a result are 0)
            self.zero = z
            # C = carry flag (1 on unsigned overflow)
            self.carry = c

    def __str__(self):
        value = ['0'] * 8
        if self.negative:
            value[0] = '1'
        if self.overflow:
            value[1] = '1'
        value[2] = '1'
        if self.break_command:
            value[3] = '1'
        if self.decimal:
            value[4] = '1'
        if self.interrupts_disabled:
            value[5] = '1'
        if self.zero:
            value[6] = '1'
        if self.carry:
            value[7] = '1'
        return "".join(value)


class CPUState:
    def __init__(self, pc=MemoryPositions.PRG_ROM_START.start, sp=MemoryPositions.STACK.end - 2, a=0, x=0, y=0, p=None, addr=None, data=None,
                 cycle=0, log_compatible_mode=False):
        super().__init__()
        self.pc = pc
        self.sp = sp
        self.x = x
        self.a = a
        self.y = y
        if not p:
            self.p = StatusRegisterFlags(int_value=0x34)
        else:
            self.p = p
        self.addr = addr
        self.data = data
        self.cycle = cycle
        self.log_compatible_mode = log_compatible_mode

    def __str__(self):
        if self.log_compatible_mode:
            return "A:%02X X:%02X Y:%02X P:%02X SP:%02X" % (
                np.uint8(self.a), np.uint8(self.x), np.uint8(self.y), int(self.p.__str__(), 2), np.uint8(self.sp & 0x00FF))
        else:
            return "| pc = 0x{:04x} | a = 0x{:02x} | x = 0x{:02x} | y = 0x{:02x} | sp = 0x{:04x} | p[NV-BDIZC] = {} |{}".format(
                np.uint16(self.pc), np.uint8(self.a), np.uint8(self.x), np.uint8(self.y), np.uint16(self.sp), self.p, self.__load_store_str())

    def __load_store_str(self):
        return " MEM[0x%04x] = 0x%02x |" % (np.uint16(self.addr), np.uint8(self.data)) if (self.addr and self.data) else ""


class CPU:
    def __init__(self, state=None, log_compatible_mode=False):
        if not state:
            self._state = CPUState()
        else:
            self._state = state
        self._state.log_compatible_mode = log_compatible_mode
        self.cycle_period_start = time.monotonic()

    @property
    def pc(self):
        return self._state.pc

    @property
    def sp(self):
        return self._state.sp

    @property
    def x(self):
        return self._state.x

    @property
    def a(self):
        return self._state.a

    @property
    def y(self):
        return self._state.y

    @property
    def negative(self):
        return self._state.p.negative

    @property
    def overflow(self):
        return self._state.p.overflow

    @property
    def break_command(self):
        return self._state.p.break_command

    @property
    def decimal(self):
        return self._state.p.decimal

    @property
    def interrupts_disabled(self):
        return self._state.p.interrupts_disabled

    @property
    def zero(self):
        return self._state.p.zero

    @property
    def carry(self):
        return self._state.p.carry

    @property
    def addr(self):
        return self._state.addr

    @property
    def data(self):
        return self._state.data

    @property
    def cycle(self):
        return self._state.cycle

    @property
    def flags(self):
        return int(self._state.p.__str__(), 2)

    @flags.setter
    def flags(self, value):
        new_status = StatusRegisterFlags(int_value=value)
        self.negative = new_status.negative
        self.overflow = new_status.overflow
        self.break_command = new_status.break_command
        self.decimal = new_status.decimal
        self.interrupts_disabled = new_status.interrupts_disabled
        self.zero = new_status.zero
        self.carry = new_status.carry

    @pc.setter
    def pc(self, value):
        self._state.pc = value

    @sp.setter
    def sp(self, value):
        self._state.sp = value

    @x.setter
    def x(self, value):
        self._state.x = value

    @a.setter
    def a(self, value):
        self._state.a = value

    @y.setter
    def y(self, value):
        self._state.y = value

    @negative.setter
    def negative(self, value):
        self._state.p.negative = value

    @overflow.setter
    def overflow(self, value):
        self._state.p.overflow = value

    @break_command.setter
    def break_command(self, value):
        self._state.p.break_command = value

    @decimal.setter
    def decimal(self, value):
        self._state.p.decimal = value

    @interrupts_disabled.setter
    def interrupts_disabled(self, value):
        self._state.p.interrupts_disabled = value

    @zero.setter
    def zero(self, value):
        self._state.p.zero = value

    @carry.setter
    def carry(self, value):
        self._state.p.carry = value

    @addr.setter
    def addr(self, value):
        self._state.addr = value

    @data.setter
    def data(self, value):
        self._state.data = value

    def inc_cycle(self):
        self.inc_cycle_by(1)

    def inc_cycle_by(self, value):
        self._state.cycle += value

    def inc_pc_by(self, value=1):
        self._state.pc += value

    def exec_in_cycle(self, block, *args):
        # FIXME: delay code
        if self.cycle % CYCLE_PERIOD_SIZE == 0 and self.cycle != 0:
            current = time.monotonic()
            elapsed = current - self.cycle_period_start
            self.cycle_period_start = current
            print("Δt = %.9f. Cycle period: %.9f" % (elapsed, CYCLE_PERIOD * CYCLE_PERIOD_SIZE))
            if elapsed < CYCLE_PERIOD * CYCLE_PERIOD_SIZE:
                time.sleep(CYCLE_PERIOD * CYCLE_PERIOD_SIZE - elapsed)
        
        self.inc_cycle()
        result = block(*args)
        return result

    # FIXME: think of a better name
    def clear_state_mem(self):
        self._state.addr = None
        self._state.data = None

    def __str__(self):
        return self._state.__str__()
