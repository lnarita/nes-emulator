from enum import unique, Enum, auto
from emulator.window import Window
import pygame
from pygame import Color

COLORS = [Color(124,124,124),Color(0,0,252),Color(0,0,188),Color(68,40,188),
        Color(148,0,132),Color(168,0,32),Color(168,16,0),Color(136,20,0),
        Color(80,48,0),Color(0,120,0),Color(0,104,0),Color(0,88,0),
        Color(0,64,88),Color(0,0,0),Color(0,0,0),Color(0,0,0),
        Color(188,188,188),Color(0,120,248),Color(0,88,248),
        Color(104,68,252),Color(216,0,204),Color(228,0,88),
        Color(248,56,0),Color(228,92,16),Color(172,124,0),
        Color(0,184,0),Color(0,168,0),Color(0,168,68),Color(0,136,136),
        Color(0,0,0),Color(0,0,0),Color(0,0,0),Color(248,248,248),
        Color(60,188,252),Color(104,136,252),Color(152,120,248),Color(248,120,248),
        Color(248,88,152),Color(248,120,88),Color(252,160,68),Color(248,184,0),
        Color(184,248,24),Color(88,216,84),Color(88,248,152),Color(0,232,216),
        Color(120,120,120),Color(0,0,0),Color(0,0,0),Color(252,252,252),
        Color(164,228,252),Color(184,184,248),Color(216,184,248),Color(248,184,248),
        Color(248,164,192),Color(240,208,176),Color(252,224,168),Color(248,216,120),
        Color(216,248,120),Color(184,248,184),Color(184,248,216),Color(0,252,252),
        Color(248,216,248),Color(0,0,0),Color(0,0,0)]

class PPU:
    #shifter
    class TileDataShifter():#shifts right 16 bits
        def __init__(self):
            self.attribute = 0
            self.patternlo = 0
            self.patternhi = 0

        def shift(self):
            #gets rightmost bit and shifts
            atBit = (self.attribute & 1) 
            loBit = (self.patternlo & 1)
            hiBit = (self.patternhi & 1)
            self.attribute >>= 1
            self.patternlo >>= 1
            self.patternhi >>= 1
            return (atBit,loBit,hiBit)

        def shift8(self):
            self.attribute >>= 8
            self.patternlo >>= 8
            self.patternhi >>= 8

        #DO NOT RELOAD WITHOUT USING THE FIRST 8 BITS FIRST
        def reload(self, attribute, patternlo, patternhi):
            #new byte is high
            #old data is low
            self.attribute = ((attribute << 8) | (self.attribute & 0b11111111))
            self.patternlo = ((patternlo << 8) | (self.patternlo & 0b11111111))
            self.patternhi = ((patternhi << 8) | (self.patternhi & 0b11111111))

    def __init__(self, ppuctrl=0x0, ppumask=0x0, ppustatus=0x0, oamaddr=0x0, oamdata=0x0, ppuscroll=0x0, ppuaddr=0x0, ppudata=0x0, oamdma=0x0, hi_lo_latch=False, ‭mirroring=True):
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

    #palettes
    self.bgPalettes = [[0 for x in range(4)] for y in range(4)]
    self.sprPalettes = [[0 for x in range(4)] for y in range(4)]

    #clock sync
    self.scnLine = 261
    self.clock = 0
    self.evenFrame = False

    #screen
    self.screen = Window()


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
            self.nametables[nametable_addr(addr)] = value
        elif MemoryPositions.NAMETABLES_MIRROR.contains(addr):
            self.store(addr % 0x1000 + MemoryPositions.NAMETABLES.start)
        elif MemoryPositions.PALLETES.contains(addr):
            self.palletes[addr - MemoryPositions.PALLETES.start] = value
        elif MemoryPositions.PALLETE_MIRROR.contains(addr):
            self.store(addr % 0x20 + MemoryPositions.PALLETES.start)
        elif MemoryPositions.MIRROR.contains(addr):
            self.store(addr % 0x4000)
        else:
            raise IndexError("Invalid Address 0x{:04x}".format(addr))

    def setNMI(self,cpu=None,memory=None, nmi=None):
        self.cpu = cpu
        self.memory = memory
        self.nmi = nmi

    def getPalettes(self):
        bgAddr = 0x3F00
        sprAddr = bgAddr+16

        for i in range(16):
            k = self.ram[bgAddr + i]
            c = COLORS[k]
            self.bgPalettes[i//4][i%4] = c

            k = self.ram[sprAddr + i]
            c = COLORS[k]
            self.sprPalettes[i//4][i%4]= c


    def tick(self):

        if (self.scnLine >= 0 and self.scnLine <= 239): #TODO: visible scanlines
            #cycle 0: idle
            if (self.cycle >=1 and self.cycle<=256):
                pass
                #every 8 cycles:
                #load nametable byte 
                #attribute table byte
                #pattern table low
                #pattern table hi
                #2 cycles each

                #every cycle:
                #get pixel from latch
                #print pixel

            if (self.cycle >=9 and self.cycle<= 257 and self.cycle%8 == 1):
                #reload shifters
                pass

            if (self.cycle >=257 and self.cycle <=320):
                pass
            
        elif (self.scnLine == 240): #TODO: Post render scanline
            pass
        elif (self.scnLine >=241 and self.scnLine <= 260):#TODO: vblank
            pass
        else:# TODO: pre-render scanline
            #gets ppu base nameTable address
            nameTable = 0xF100
            #gets ppu base patternTable address
            patternTable = 0xD000
            #gets ppu base attributeTable address
            attributeTable = 0xF4C0

            tileType = self.memory.fetch(nameTable)

            pass

        #odd frames have less clocks
        if (self.evenFrame == False and self.scnLine == 261 and self.clock == 339):
            self.evenFrame = True
            self.scnLine = 0
            self.clock = 0
        elif (self.scnLine == 261 and self.clock == 340): # last clock of a frame
            self.evenFrame = not self.evenFrame
            self.scnLine = 0
            self.clock = 0
        elif (self.clock == 340): # last clock of a scanline
            self.scnLine += 1
            self.clock = 0
        else:
            self.clock +=1 # a normal clock



    def vBlank(self,vb): #TODO: NMI
        if(vb):#vblank starts
            self.ppustatus = self.ppustatus | 0b10000000 #updates vblank flag
            
            generateNMI = (self.ppuctrl & 0b10000000) > 0
            if generateNMI:
               pass
               #self.nmi(self.cpu,self.memory)

        else:#vblank ends
            self.ppustatus = self.ppustatus & 0b01111111 # updates vblank flag

    def scanLine(self,line):# TODO: 8x16 mode
        nameTable = 0
        attributeTable = 0
        nameTableCtrlBits = (self.ppuctrl & 0b11)

        if (nameTableCtrlBits == 0):
            nameTable = 0x2000
            attributeTable = 0x23C0
        elif (nameTableCtrlBits==1):
            nameTable = 0x2400
            attributeTable = 0x27C0
        elif (nameTableCtrlBits==2):
            nameTable = 0x2800
            attributeTable = 0x2BC0
        else:
            nameTable = 0x2C00
            attributeTable = 0x2FC0

        patternTable = 0
        patternTableCtrlBits = (self.ppuctrl & 0b10000)
        if (patternTableCtrlBits==0):
            patternTable = 0x0
        else:
            #patternTable = 0x4200
            patternTable = 0x1000

        if line < 240: #visible scanlines
            for i in range(256):

                #calculates the tile position
                tileVer = line//8
                tileHor = i//8
                tileNo = 32*tileVer+tileHor

                #gets the tile type by looking at name table
                tileType = self.ram[nameTable+tileNo] #1 byte tile
                tileType = 255
                patternVer = line%8
                patternHor = i%8

                #gets the pattern
                patternlo = self.memory.fetch(patternTable+tileType*16+patternVer)
                patternhi = self.memory.fetch(patternTable+tileType*16+8+patternVer)
                #print(hex(patternTable+tileType*16))
                #print (bin(patternhi) + " " + bin(patternlo))
                
                #gets the colorCode
                hiBit = (patternhi&(1<<7-patternHor)) > 0
                loBit = (patternlo&(1<<7-patternHor)) > 0
                colorCode = hiBit*2 + loBit

                #which block
                blockVer = line//32
                blockHor = i//32
                blockNo = 8*blockVer + blockHor
                attribute = self.ram[attributeTable+blockNo]

                #which tile in block
                whichTile = 0
                offVer = line%32
                offHor = i%32
                if (offVer<16 and offHor<16):
                    whichTile=0
                elif (offVer<16 and offHor>=16):
                    whichTile=1
                elif (offVer>=16 and offHor<16):
                    whichTile=2
                else:
                    whichTile=3

                paletteNo = ((attribute & (0b11 << whichTile*2)) >> whichTile*2)
                color = self.bgPalettes[paletteNo][colorCode]
                self.screen.setPixel(i,line,color)

            if line==239:
                self.screen.flip()     

        if line==241: # vblank
            pass
            self.vBlank(True)
        if line== 261:
            pass
            #vBlank(False)

