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
            # B = break flag (1 when interupt was caused by a BRK)
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


class BaseCPUState:
    def __init__(self, pc=0xC000, sp=MemoryPositions.STACK.start, a=0, x=0, y=0, p=StatusRegisterFlags(), addr=None, data=None):
        super().__init__()
        self.pc = pc
        self.sp = sp
        self.x = x
        self.a = a
        self.y = y
        self.p = p
        self.addr = addr
        self.data = data

    def __str__(self):
        return "| pc = 0x{:04x} | a = 0x{:02x} | x = 0x{:02x} | y = 0x{:02x} | sp = 0x{:04x} | p[NV-BDIZC] = {} |{}".format(
            self.pc, self.a, self.x, self.y, self.sp, self.p, self.__load_store_str())

    def __load_store_str(self):
        return " MEM[0x%04x] = 0x%02x |" % (self.addr, self.data) if (self.addr and self.data) else ""
