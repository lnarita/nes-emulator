from enum import unique, Enum, auto

@unique
class PPUMemoryPositions(Enum):
    PATTERN_TABLES = (auto(), 0x0000, 0x1FFF)
    NAMETABLES = (auto(), 0x2000, 0x2FFF)
    NAMETABLES_MIRROR = (auto(), 0x3000, 0x3EFF)
    PALLETES = (auto(), 0x3F00, 0x3F1F)
    PALLETES_MIRROR = (auto(), 0x3F20, 0x3FFF)
    MIRROR = (auto(), 0x4000, 0xFFFF)

    def __init__(self, value, start, end):
        self._value_ = value
        self.start = start
        self.end = end

    def contains(self, addr):
        return self.start <= addr <= self.end

    def wrap(self, addr):
        return (self.start + (addr - self.start) % ((self.end + 1) - self.start))


class PPU:
    def __init__(self, ppuctrl=0x0, ppumask=0x0, ppustatus=0x0, oamaddr=0x0, oamdata=0x0, ppuscroll=0x0, ppuaddr=0x0, ppudata=0x0, oamdma=0x0, hi_lo_latch=False, ‭mirroring=False):
        self.ppuctrl = ppuctrl
        self.ppumask = ppumask
        self.ppustatus = ppustatus
        self.oamaddr = oamaddr
        self.oamdata = oamdata
        self.ppuscroll = ppuscroll
        self.ppuaddr = ppuaddr
        self.ppudata = ppudata
        self.oamdma = oamdma
        self.hi_lo_latch = hi_lo_latch
        self.oam = [0x0] * 256  # 64 sprites * 4 bytes
        # self.ram = [0x0] * 0x4000
        self.nametables = [0x0] * 0x800
        self.palletes = [0x0] * 0x20
        self.mirroring = mirroring

    def nametable_addr(addr):
        addr -= MemoryPositions.NAMETABLES.start

        # Vertical mirroring
        if self.mirroring:
            # Nametables 2 and 3
            if addr >= 0x800:
                addr -= 0x800

        # Horizontal mirroring
        else:
            # Nametables 1 and 3
            if addr % 0x800 >= 0x400:
                addr -= 0x400

        return addr


    def fetch(self, addr):
        if MemoryPositions.PATTERN_TABLES.contains(addr)
            pass

        elif MemoryPositions.NAMETABLES.contains(addr):
            return self.nametables[nametable_addr(addr)]

        elif MemoryPositions.NAMETABLES_MIRROR.contains(addr):
            return self.fetch(addr % 0x1000 + MemoryPositions.NAMETABLES.start)

        elif MemoryPositions.PALLETES.contains(addr):
            return self.palletes[addr - MemoryPositions.PALLETES.start]

        elif MemoryPositions.PALLETE_MIRROR.contains(addr):
            return self.fetch(addr % 0x20 + MemoryPositions.PALLETES.start)

        elif MemoryPositions.MIRROR.contains(addr):
            return self.fetch(addr % 0x4000)

        else:
            raise IndexError("Invalid Address 0x{:04x}".format(addr))


    def store(self, addr, value):
        if MemoryPositions.PATTERN_TABLES.contains(addr)
            pass
        elif MemoryPositions.NAMETABLES.contains(addr):
            pass
        elif MemoryPositions.NAMETABLES_MIRROR.contains(addr):
            pass
        elif MemoryPositions.PALLETES.contains(addr):
            pass
        elif MemoryPositions.PALLETE_MIRROR.contains(addr):
            pass
        elif MemoryPositions.MIRROR.contains(addr):
            pass
        else:
            raise IndexError("Invalid Address 0x{:04x}".format(addr))