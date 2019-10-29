from emulator.constants import KB


class INesHeader:
    def __init__(self, prg_rom_size, chr_rom_size, flags_6, flags_7, prg_ram_size, flags_9, flags_10):
        self.prg_rom_size = prg_rom_size * 16 * KB
        self.chr_rom_size = chr_rom_size * 8 * KB
        self.flags_6 = flags_6
        self.flags_7 = flags_7
        self.prg_ram_size = prg_ram_size * 8 * KB
        self.flags_9 = flags_9
        self.flags_10 = flags_10


class Cartridge:
    def __init__(self, header, prg_rom, chr_rom):
        self.header = header
        self.prg_rom = prg_rom
        self.chr_rom = chr_rom

    @classmethod
    def from_bytes(cls, cartridge_content):
        if cartridge_content and len(cartridge_content) > 16:
            # 0-3: Constant $4E $45 $53 $1A (“NES” followed by MS-DOS end-of-file)
            if "NES\x1A" == cartridge_content[:4].decode("utf-8"):
                # 4: Size of PRG ROM in 16 KB units
                prg_rom_size = cartridge_content[4]
                # 5: Size of CHR ROM in 8 KB units (Value 0 means the board uses CHR RAM)
                chr_rom_size = cartridge_content[5]
                # 6: Flags 6
                flags_6 = cartridge_content[6]
                # 7: Flags 7
                flags_7 = cartridge_content[7]
                # 8: Size of PRG RAM in 8 KB units (Value 0 infers 8 KB for compatibility; see PRG RAM circuit)
                prg_ram_size = cartridge_content[8]
                # 9: Flags 9
                flags_9 = cartridge_content[9]
                # 10: Flags 10 (unofficial)
                flags_10 = cartridge_content[10]
                # 11-15: Zero filled

                header = INesHeader(prg_rom_size, chr_rom_size, flags_6, flags_7, prg_ram_size, flags_9, flags_10)
                prg_rom = cartridge_content[16:(header.prg_rom_size + 16)]
                if len(cartridge_content) < (header.prg_rom_size + 16 + header.chr_rom_size):
                    cartridge_content = bytearray(cartridge_content)
                    for i in range((header.prg_rom_size + 16 + header.chr_rom_size) - len(cartridge_content)):
                        cartridge_content.append(0)
                chr_rom = cartridge_content[(header.prg_rom_size + 16):(header.prg_rom_size + 16 + header.chr_rom_size)]

                return cls(header, prg_rom, chr_rom)
        raise ValueError("Invalid iNES Header!")
