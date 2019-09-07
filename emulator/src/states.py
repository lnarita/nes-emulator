from constants import MemoryPositions


class BaseCPUState:
    def __init__(self, pc=MemoryPositions.INITIAL.start, sp=MemoryPositions.STACK.start, a=0, x=0, y=0, p=0, addr=None, data=None):
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
        return "| pc = 0x{:04x} | a = 0x{:02x} | x = 0x{:02x} | y = 0x{:02x} | sp = 0x{:04x} | p[NV-BDIZC] = {:08b} |{}".format(
            self.pc, self.a, self.x, self.y, self.sp, self.p, self.__load_store_str())

    def __load_store_str(self):
        return " MEM[0x%04x] = 0x%02x |" % (self.addr, self.data) if (self.addr and self.data) else ""
