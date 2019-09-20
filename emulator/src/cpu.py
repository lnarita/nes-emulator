import time

from constants import CYCLE_PERIOD
from memory import MemoryPositions


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
    def __init__(self, pc=MemoryPositions.PRG_ROM_START.start, sp=MemoryPositions.STACK.start, a=0, x=0, y=0, p=StatusRegisterFlags(), addr=None, data=None,
                 cycle=1):
        super().__init__()
        self.pc = pc
        self.sp = sp
        self.x = x
        self.a = a
        self.y = y
        self.p = p
        self.addr = addr
        self.data = data
        self.cycle = cycle

    def __str__(self):
        return "| pc = 0x{:04x} | a = 0x{:02x} | x = 0x{:02x} | y = 0x{:02x} | sp = 0x{:04x} | p[NV-BDIZC] = {} |{}".format(
            self.pc, self.a, self.x, self.y, self.sp, self.p, self.__load_store_str())

    def __load_store_str(self):
        return " MEM[0x%04x] = 0x%02x |" % (self.addr, self.data) if (self.addr and self.data) else ""


class CPU:
    def __init__(self, state=CPUState()):
        self._state = state

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
        self._state.cycle += 1

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
