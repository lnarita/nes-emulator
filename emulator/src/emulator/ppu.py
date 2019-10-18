import time

import numpy as np

from emulator.constants import CYCLE_PERIOD
from emulator.memory import MemoryPositions

class CPU:
    def __init__(self, state=CPUState(), log_compatible_mode=False):
        self._state = state
        self._state.log_compatible_mode = log_compatible_mode

    # $2000
    @property
    def ppuctrl(self):
        return self._state.ppuctrl

    # $2001
    @property
    def ppumask(self):
        return self._state.ppumask

    # $2002 
    @property
    def ppustatus(self):
        return self._state.ppustatus

    # $2003
    @property
    def oamaddr(self):
        return self._state.a

    # $2004
    @property
    def oamdata(self):
        return self._state.y

    # $2005
    @property
    def ppuscroll(self):
        return self._state.p.negative

    # $2006
    @property
    def ppuaddr(self):
        return self._state.p.overflow

    # $2007
    @property
    def ppudata(self):
        return self._state.p.break_command

    # $2008
    @property
    def oamdma(self):
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
        start = time.monotonic()
        result = block(*args)
        done = time.monotonic()
        elapsed = done - start
        if elapsed < CYCLE_PERIOD:
            time.sleep(CYCLE_PERIOD - elapsed)
        self.inc_cycle()
        return result

    # FIXME: think of a better name
    def clear_state_mem(self):
        self._state.addr = None
        self._state.data = None

    def __str__(self):
        return self._state.__str__()
