from enum import unique, Enum, auto

from emulator.constants import KB
from emulator.ppu import PPU


@unique
class MemoryPositions(Enum):
    ZERO_PAGE = (auto(), 0x0000, 0x00FF)
    STACK = (auto(), 0x0100, 0x01FF)
    RAM = (auto(), 0x0200, 0x07FF)
    RAM_MIRROR_1 = (auto(), 0x0800, 0x0FFF)
    RAM_MIRROR_2 = (auto(), 0x1000, 0x17FF)
    RAM_MIRROR_3 = (auto(), 0x1800, 0x1FFF)
    PPU_REGISTERS = (auto(), 0x2000, 0x2007)
    PPU_REGISTERS_MIRROR = (auto(), 0x2008, 0x3FFF)
    APU_IO_REGISTERS = (auto(), 0x4000, 0x4015)
    CONTROLLERS = (auto(), 0x4016, 0x4017)
    APU_IO_EXTRAS = (auto(), 0x4018, 0x401F)
    CARTRIDGE = (auto(), 0x4020, 0xFFFF)
    PRG_ROM_START = (auto(), 0xC000, 0xFFFF)  # the PRG ROM size is defined by the iNES header, so `end` is a dummy value
    NMI = (auto(), 0xFFFA, 0xFFFB)
    RESET = (auto(), 0xFFFC, 0xFFFD)
    IRQ = (auto(), 0xFFFE, 0xFFFF)

    def __init__(self, value, start, end):
        self._value_ = value
        self.start = start
        self.end = end

    def contains(self, addr):
        return self.start <= addr <= self.end

    def wrap(self, addr):
        return (self.start + (addr - self.start) % ((self.end + 1) - self.start))


class Memory:
    def __init__(self, rom=None, ram=None, ppu=None):
        def __pad_or_truncate(some_list, target_len):
            return some_list[:target_len] + [0x0] * (target_len - len(some_list))

        if not ram:
            self.ram = [0x0] * Memory.ram_size()
        else:
            self.ram = __pad_or_truncate(ram, Memory.ram_size())
        self.rom = rom
        self.debug_mem = []
        self.ppu = ppu


    def fetch(self, addr):
        if MemoryPositions.ZERO_PAGE.contains(addr) or \
                MemoryPositions.STACK.contains(addr) or \
                MemoryPositions.RAM.contains(addr):
            return self.ram[addr]
        elif MemoryPositions.RAM_MIRROR_1.contains(addr):
            return self.ram[addr - MemoryPositions.RAM_MIRROR_1.start]
        elif MemoryPositions.RAM_MIRROR_2.contains(addr):
            return self.ram[addr - MemoryPositions.RAM_MIRROR_2.start]
        elif MemoryPositions.RAM_MIRROR_3.contains(addr):
            return self.ram[addr - MemoryPositions.RAM_MIRROR_3.start]
        elif MemoryPositions.PPU_REGISTERS.contains(addr) or addr == 0x4014:
            return self.fetch_ppu(addr)
        elif MemoryPositions.PPU_REGISTERS_MIRROR.contains(addr):
            return self.fetch_ppu(addr%8 + 0x2000)
        elif MemoryPositions.APU_IO_REGISTERS.contains(addr):
            # TODO
            return 0xFF
        elif MemoryPositions.CONTROLLERS.contains(addr):
            return self.ppu.screen.read_keys()
        elif MemoryPositions.APU_IO_EXTRAS.contains(addr):
            # TODO
            return 0xFF
        elif MemoryPositions.CARTRIDGE.contains(addr):
            return self.rom[addr - MemoryPositions.PRG_ROM_START.start]
        else:
            raise IndexError("Invalid Address 0x{:04x}".format(addr))

    def store(self, addr, value):
        if MemoryPositions.ZERO_PAGE.contains(addr) or \
                MemoryPositions.STACK.contains(addr) or \
                MemoryPositions.RAM.contains(addr):
            self.ram[addr] = value
        elif MemoryPositions.RAM_MIRROR_1.contains(addr):
            self.ram[addr - MemoryPositions.RAM_MIRROR_1.start] = value
        elif MemoryPositions.RAM_MIRROR_2.contains(addr):
            self.ram[addr - MemoryPositions.RAM_MIRROR_2.start] = value
        elif MemoryPositions.RAM_MIRROR_3.contains(addr):
            self.ram[addr - MemoryPositions.RAM_MIRROR_3.start] = value
        elif MemoryPositions.PPU_REGISTERS.contains(addr) or addr == 0x4014:
            self.store_ppu(addr, value)
        elif MemoryPositions.PPU_REGISTERS_MIRROR.contains(addr):
            self.store_ppu(addr%8 + 0x2000, value)
        elif 0x4000 <= addr <= 0xFFFF:
            # TODO: remove later
            self.debug_mem.append(("%04X" % addr, "%02X" % value))
        elif MemoryPositions.APU_IO_REGISTERS.contains(addr):
            # TODO
            return
        elif MemoryPositions.CONTROLLERS.contains(addr):
            self.ppu.screen.latch_keys(value & 0b00000001)
        elif MemoryPositions.APU_IO_EXTRAS.contains(addr):
            # TODO
            return
        else:
            raise IndexError("Invalid Address 0x{:04x}".format(addr))

    def fetch_ppu(self, addr):
        if addr == 0x2002:
            self.ppu.hi_lo_latch = False
            status = self.ppu.ppustatus
            self.ppu.ppustatus = self.ppu.ppustatus & 0b01111111 # clears vblank flag
            return status
        elif addr == 0x2004:
            return self.ppu.oam[oamaddr]
        elif addr == 0x2007:
            # self.ppu.ppudata = self.ppu.ram[self.ppu.ppuaddr]
            self.ppu.ppudata = self.ppu.fetch(self.ppu.ppuaddr)
            return self.ppu.ppudata
        else:
            return 0x00

    def store_ppu(self, addr, value):
        if addr == 0x2000:
            self.ppu.ppuctrl = value
        elif addr == 0x2001:
            self.ppu.ppumask = value
        elif addr == 0x2003:
            self.ppu.oamaddr = value
        elif addr == 0x2004:
            self.ppu.oam[oamaddr] = value
        elif addr == 0x2005:
            # Write low byte
            if self.ppu.high_latch:
                self.ppu.ppuscroll |= value
            # Write high byte
            else:
                self.ppu.ppuscroll = value << 8
                self.ppu.hi_lo_latch = True
        elif addr == 0x2006:
            # Write low byte
            if self.ppu.hi_lo_latch:
                self.ppu.ppuaddr |= value
            # Write high byte
            else:
                self.ppu.ppuaddr = value << 8
                self.ppu.hi_lo_latch = True
        elif addr == 0x2007:
            self.ppu.store(self.ppu.ppuaddr, value)
            #print(hex(self.ppu.ppuaddr) + " "+str(self.ppu.ram[self.ppu.ppuaddr]))
            if (self.ppu.ppuaddr >= 0x3F00 and self.ppu.ppuaddr <=0x3F31): #palette change
                self.ppu.getPalettes()
        elif addr == 0x4014:
            self.ppu.oamdma = value
            # DMA Transfer
            # TODO: Cycle count
            #     513 or 514 cycles after the $4014 write tick. (1 dummy read cycle
            #     while waiting for writes to complete, +1 if on an odd CPU cycle,
            #     then 256 alternating read/write cycles.)
            for dma_addr in range(0x0,0x100):
                self.ppu.oam[dma_addr] = self.ram[value << 8 | dma_addr]

    def stack_push(self, cpu, value):
        self.store(cpu.sp, value)
        cpu.sp -= 1
        cpu.sp = cpu.sp & 0xff ^ 0x0100

    def stack_pop(self, cpu):
        value = self.fetch(cpu.sp + 1)
        cpu.sp += 1
        cpu.sp = cpu.sp & 0xff ^ 0x0100
        return value

    def get_effective_address(self, addr):
        if MemoryPositions.RAM_MIRROR_1.contains(addr):
            return addr - MemoryPositions.RAM_MIRROR_1.start
        elif MemoryPositions.RAM_MIRROR_2.contains(addr):
            return addr - MemoryPositions.RAM_MIRROR_2.start
        elif MemoryPositions.RAM_MIRROR_3.contains(addr):
            return addr - MemoryPositions.RAM_MIRROR_3.start
        else:
            return addr

    @staticmethod
    def ram_size():
        return 2 * KB
