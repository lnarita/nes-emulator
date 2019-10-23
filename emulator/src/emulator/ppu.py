class PPU:
    def __init__(self, ppuctrl=0x0, ppumask=0x0, ppustatus=0x0, oamaddr=0x0, oamdata=0x0, ppuscroll=0x0, ppuaddr=0x0, ppudata=0x0, oamdma=0x0, hi_lo_latch=False):
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
