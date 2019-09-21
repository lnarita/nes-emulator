from enum import auto, Enum

CPU_FREQUENCY = 1.79e6
CYCLE_PERIOD = 1. / CPU_FREQUENCY  # seconds

KB = 1024

HIGH_BITS_MASK = 0b1111111100000000
LOW_BITS_MASK = 0b0000000011111111

NEGATIVE_BIT = 0b0000000010000000


# TODO: this is a WIP, it should probably not be designed like this


class DO_NO_USE_AddressingMode(Enum):
    # Indirect
    INDIRECT = (auto(),)
    # (Indirect, X)
    INDIRECT_X = (auto(),)
    # (Indirect), Y
    INDIRECT_Y = (auto(),)
    # Zero Page
    ZERO_PAGE = (auto(),)
    # Zero Page, X
    ZERO_PAGE_X = (auto(),)
    # Zero Page, Y
    ZERO_PAGE_Y = (auto(),)
    # Absolute
    ABSOLUTE = (auto(),)
    # Absolute, Y
    ABSOLUTE_Y = (auto(),)
    # Absolute, X
    ABSOLUTE_X = (auto(),)
    # Immediate
    IMMEDIATE = (auto(),)
    # Accumulator
    ACCUMULATOR = (auto(),)
    # Relative (Branching)
    RELATIVE = (auto(),)


class DO_NOT_USE_Opcodes(Enum):
    ADC_1 = (0x61, DO_NO_USE_AddressingMode.INDIRECT_X, 6,)
    ADC_2 = (0x65, DO_NO_USE_AddressingMode.ZERO_PAGE, 3,)
    ADC_3 = (0x69, DO_NO_USE_AddressingMode.IMMEDIATE, 2,)
    ADC_4 = (0x6D, DO_NO_USE_AddressingMode.ABSOLUTE, 4,)
    ADC_5 = (0x71, DO_NO_USE_AddressingMode.INDIRECT_Y, 5,)  # Add one cycle if indexing across page boundary
    ADC_6 = (0x75, DO_NO_USE_AddressingMode.ZERO_PAGE_X, 4,)
    ADC_7 = (0x79, DO_NO_USE_AddressingMode.ABSOLUTE_Y, 4,)  # Add one cycle if indexing across page boundary
    ADC_8 = (0x7D, DO_NO_USE_AddressingMode.ABSOLUTE_X, 4,)  # Add one cycle if indexing across page boundary
    AND_1 = (0x21, DO_NO_USE_AddressingMode.INDIRECT_X, 6,)
    AND_2 = (0x25, DO_NO_USE_AddressingMode.ZERO_PAGE, 3,)
    AND_3 = (0x29, DO_NO_USE_AddressingMode.IMMEDIATE, 2,)
    AND_4 = (0x2D, DO_NO_USE_AddressingMode.ABSOLUTE, 4,)
    AND_5 = (0x31, DO_NO_USE_AddressingMode.INDIRECT_Y, 5,)  # Add one cycle if indexing across page boundary
    AND_6 = (0x35, DO_NO_USE_AddressingMode.ZERO_PAGE_X, 4,)
    AND_7 = (0x39, DO_NO_USE_AddressingMode.ABSOLUTE_Y, 4,)  # Add one cycle if indexing across page boundary
    AND_8 = (0x3D, DO_NO_USE_AddressingMode.ABSOLUTE_X, 4,)  # Add one cycle if indexing across page boundary
    ASL_1 = (0x06, DO_NO_USE_AddressingMode.ZERO_PAGE, 5,)
    ASL_2 = (0x0A, DO_NO_USE_AddressingMode.ACCUMULATOR, 2,)
    ASL_3 = (0x0E, DO_NO_USE_AddressingMode.ABSOLUTE, 6,)
    ASL_4 = (0x16, DO_NO_USE_AddressingMode.ZERO_PAGE_X, 6,)
    ASL_5 = (0x1E, DO_NO_USE_AddressingMode.ABSOLUTE_X, 7,)
    BCC = (0x90, DO_NO_USE_AddressingMode.RELATIVE, 2,)  # Add one cycle if branch is taken, and one additional if branching operation crosses page boundary
    BCS = (0xB0, DO_NO_USE_AddressingMode.RELATIVE, 2,)  # Add one cycle if branch is taken, and one additional if branching operation crosses page boundary
    BEQ = (0xF0, DO_NO_USE_AddressingMode.RELATIVE, 2,)  # Add one cycle if branch is taken, and one additional if branching operation crosses page boundary
    BIT_1 = (0x24, DO_NO_USE_AddressingMode.ZERO_PAGE, 3,)
    BIT_2 = (0x2C, DO_NO_USE_AddressingMode.ABSOLUTE, 4,)
    BMI = (0x30, DO_NO_USE_AddressingMode.RELATIVE, 2,)  # Add one cycle if branch is taken, and one additional if branching operation crosses page boundary
    BNE = (0xD0, DO_NO_USE_AddressingMode.RELATIVE, 2,)  # Add one cycle if branch is taken, and one additional if branching operation crosses page boundary
    BPL = (0x10, DO_NO_USE_AddressingMode.RELATIVE, 2,)  # Add one cycle if branch is taken, and one additional if branching operation crosses page boundary
    BRK = (0x00, None, 7,)
    BVC = (0x50, DO_NO_USE_AddressingMode.RELATIVE, 2,)  # Add one cycle if branch is taken, and one additional if branching operation crosses page boundary
    BVS = (0x70, DO_NO_USE_AddressingMode.RELATIVE, 2,)  # Add one cycle if branch is taken, and one additional if branching operation crosses page boundary
    CLC = (0x18, None, 2,)
    CLD = (0xD8, None, 2,)
    CLI = (0x58, None, 2,)
    CLV = (0xB8, None, 2,)
    CMP_1 = (0xC1, DO_NO_USE_AddressingMode.INDIRECT_X, 6,)
    CMP_2 = (0xC5, DO_NO_USE_AddressingMode.ZERO_PAGE, 3,)
    CMP_3 = (0xC9, DO_NO_USE_AddressingMode.IMMEDIATE, 2,)
    CMP_4 = (0xCD, DO_NO_USE_AddressingMode.ABSOLUTE, 4,)
    CMP_5 = (0xD1, DO_NO_USE_AddressingMode.INDIRECT_Y, 5,)  # Add one cycle if indexing across page boundary
    CMP_6 = (0xD5, DO_NO_USE_AddressingMode.ZERO_PAGE_X, 4,)
    CMP_7 = (0xD9, DO_NO_USE_AddressingMode.ABSOLUTE_Y, 4,)  # Add one cycle if indexing across page boundary
    CMP_8 = (0xDD, DO_NO_USE_AddressingMode.ABSOLUTE_X, 4,)  # Add one cycle if indexing across page boundary
    CPX_1 = (0xE0, DO_NO_USE_AddressingMode.IMMEDIATE, 2,)
    CPX_2 = (0xE4, DO_NO_USE_AddressingMode.ZERO_PAGE, 3,)
    CPX_3 = (0xEC, DO_NO_USE_AddressingMode.ABSOLUTE, 4,)
    CPY_1 = (0xC0, DO_NO_USE_AddressingMode.IMMEDIATE, 2,)
    CPY_2 = (0xC4, DO_NO_USE_AddressingMode.ZERO_PAGE, 3,)
    CPY_3 = (0xCC, DO_NO_USE_AddressingMode.ABSOLUTE, 4,)
    DEC_1 = (0xC6, DO_NO_USE_AddressingMode.ZERO_PAGE, 5,)
    DEC_2 = (0xCE, DO_NO_USE_AddressingMode.ABSOLUTE, 6,)
    DEC_3 = (0xD6, DO_NO_USE_AddressingMode.ZERO_PAGE_X, 6,)
    DEC_4 = (0xDE, DO_NO_USE_AddressingMode.ABSOLUTE_X, 7,)
    DEX = (0xCA, None, 2,)
    DEY = (0x88, None, 2,)
    EOR_1 = (0x41, DO_NO_USE_AddressingMode.INDIRECT_X, 6,)
    EOR_2 = (0x45, DO_NO_USE_AddressingMode.ZERO_PAGE, 3,)
    EOR_3 = (0x49, DO_NO_USE_AddressingMode.IMMEDIATE, 2,)
    EOR_4 = (0x4D, DO_NO_USE_AddressingMode.ABSOLUTE, 4,)
    EOR_5 = (0x51, DO_NO_USE_AddressingMode.INDIRECT_Y, 5,)  # Add one cycle if indexing across page boundary
    EOR_6 = (0x55, DO_NO_USE_AddressingMode.ZERO_PAGE_X, 4,)
    EOR_7 = (0x59, DO_NO_USE_AddressingMode.ABSOLUTE_Y, 4,)  # Add one cycle if indexing across page boundary
    EOR_8 = (0x5D, DO_NO_USE_AddressingMode.ABSOLUTE_X, 4,)  # Add one cycle if indexing across page boundary
    INC_1 = (0xE6, DO_NO_USE_AddressingMode.ZERO_PAGE, 5,)
    INC_2 = (0xEE, DO_NO_USE_AddressingMode.ABSOLUTE, 6,)
    INC_3 = (0xF6, DO_NO_USE_AddressingMode.ZERO_PAGE_X, 6,)
    INC_4 = (0xFE, DO_NO_USE_AddressingMode.ABSOLUTE_X, 7,)
    INX = (0xE8, None, 2,)
    INY = (0xC8, None, 2,)
    JMP_1 = (0x4C, DO_NO_USE_AddressingMode.ABSOLUTE, 3,)
    JMP_2 = (0x6C, DO_NO_USE_AddressingMode.INDIRECT, 5,)
    JSR = (0x20, DO_NO_USE_AddressingMode.ABSOLUTE, 6,)
    LDA_1 = (0xA1, DO_NO_USE_AddressingMode.INDIRECT_X, 6)
    LDA_2 = (0xA5, DO_NO_USE_AddressingMode.ZERO_PAGE, 3,)
    LDA_3 = (0xA9, DO_NO_USE_AddressingMode.IMMEDIATE, 2,)
    LDA_4 = (0xAD, DO_NO_USE_AddressingMode.ABSOLUTE, 4,)
    LDA_5 = (0xB1, DO_NO_USE_AddressingMode.INDIRECT_Y, 5)  # Add one cycle if indexing across page boundary
    LDA_6 = (0xB5, DO_NO_USE_AddressingMode.ZERO_PAGE_X, 4,)
    LDA_7 = (0xB9, DO_NO_USE_AddressingMode.ABSOLUTE_Y, 4,)  # Add one cycle if indexing across page boundary
    LDA_8 = (0xBD, DO_NO_USE_AddressingMode.ABSOLUTE_X, 4,)  # Add one cycle if indexing across page boundary
    LDX_1 = (0xA2, DO_NO_USE_AddressingMode.IMMEDIATE, 2,)
    LDX_2 = (0xA6, DO_NO_USE_AddressingMode.ZERO_PAGE, 3,)
    LDX_3 = (0xAE, DO_NO_USE_AddressingMode.ABSOLUTE, 4,)
    LDX_4 = (0xB6, DO_NO_USE_AddressingMode.ZERO_PAGE_Y, 4,)
    LDX_5 = (0xBE, DO_NO_USE_AddressingMode.ABSOLUTE_Y, 4,)  # Add one cycle if indexing across page boundary
    LDY_1 = (0xA0, DO_NO_USE_AddressingMode.IMMEDIATE, 2,)
    LDY_2 = (0xA4, DO_NO_USE_AddressingMode.ZERO_PAGE, 3,)
    LDY_3 = (0xAC, DO_NO_USE_AddressingMode.ABSOLUTE, 4,)
    LDY_4 = (0xB4, DO_NO_USE_AddressingMode.ZERO_PAGE_X, 4,)
    LDY_5 = (0xBC, DO_NO_USE_AddressingMode.ABSOLUTE_X, 4,)  # Add one cycle if indexing across page boundary
    LSR_1 = (0x46, DO_NO_USE_AddressingMode.ZERO_PAGE, 5,)
    LSR_2 = (0x4A, DO_NO_USE_AddressingMode.ACCUMULATOR, 2,)
    LSR_3 = (0x4E, DO_NO_USE_AddressingMode.ABSOLUTE, 6,)
    LSR_4 = (0x56, DO_NO_USE_AddressingMode.ZERO_PAGE_X, 6,)
    LSR_5 = (0x5E, DO_NO_USE_AddressingMode.ABSOLUTE_X, 7,)
    NOP = (0xEA, None, 2)
    ORA_1 = (0x01, DO_NO_USE_AddressingMode.INDIRECT_X, 6,)
    ORA_2 = (0x05, DO_NO_USE_AddressingMode.ZERO_PAGE, 3,)
    ORA_3 = (0x09, DO_NO_USE_AddressingMode.IMMEDIATE, 2,)
    ORA_4 = (0x0D, DO_NO_USE_AddressingMode.ABSOLUTE, 4,)
    ORA_5 = (0x11, DO_NO_USE_AddressingMode.INDIRECT_Y, 5,)  # Add one cycle if indexing across page boundary
    ORA_6 = (0x15, DO_NO_USE_AddressingMode.ZERO_PAGE_X, 4,)
    ORA_7 = (0x19, DO_NO_USE_AddressingMode.ABSOLUTE_Y, 4,)  # Add one cycle if indexing across page boundary
    ORA_8 = (0x1D, DO_NO_USE_AddressingMode.ABSOLUTE_X, 4,)  # Add one cycle if indexing across page boundary
    PHA = (0x48, None, 3,)
    PHP = (0x08, None, 3,)
    PLA = (0x68, None, 4,)
    PLP = (0x28, None, 4,)
    ROL_1 = (0x26, DO_NO_USE_AddressingMode.ZERO_PAGE, 5,)
    ROL_2 = (0x2A, DO_NO_USE_AddressingMode.ACCUMULATOR, 2,)
    ROL_3 = (0x2E, DO_NO_USE_AddressingMode.ABSOLUTE, 6,)
    ROL_4 = (0x36, DO_NO_USE_AddressingMode.ZERO_PAGE_X, 6,)
    ROL_5 = (0x3E, DO_NO_USE_AddressingMode.ABSOLUTE_X, 7,)
    ROR_1 = (0x66, DO_NO_USE_AddressingMode.ZERO_PAGE, 5,)
    ROR_2 = (0x6A, DO_NO_USE_AddressingMode.ACCUMULATOR, 2,)
    ROR_3 = (0x6E, DO_NO_USE_AddressingMode.ABSOLUTE, 6,)
    ROR_4 = (0x76, DO_NO_USE_AddressingMode.ZERO_PAGE_X, 6,)
    ROR_5 = (0x7E, DO_NO_USE_AddressingMode.ABSOLUTE_X, 7,)
    RTI = (0x40, None, 6,)
    RTS = (0x60, None, 6,)
    SBC_1 = (0xE1, DO_NO_USE_AddressingMode.INDIRECT_X, 6)
    SBC_2 = (0xE5, DO_NO_USE_AddressingMode.ZERO_PAGE, 3,)
    SBC_3 = (0xE9, DO_NO_USE_AddressingMode.IMMEDIATE, 2,)
    SBC_4 = (0xED, DO_NO_USE_AddressingMode.ABSOLUTE, 4,)
    SBC_5 = (0xF1, DO_NO_USE_AddressingMode.INDIRECT_Y, 5,)  # Add one cycle if indexing across page boundary
    SBC_6 = (0xF5, DO_NO_USE_AddressingMode.ZERO_PAGE_X, 4,)
    SBC_7 = (0xF9, DO_NO_USE_AddressingMode.ABSOLUTE_Y, 4,)  # Add one cycle if indexing across page boundary
    SBC_8 = (0xFD, DO_NO_USE_AddressingMode.ABSOLUTE_X, 4,)  # Add one cycle if indexing across page boundary
    SEC = (0x38, None, 2,)
    SED = (0xF8, None, 2,)
    SEI = (0x78, None, 2,)
    STA_1 = (0x81, DO_NO_USE_AddressingMode.INDIRECT_X, 6,)
    STA_2 = (0x85, DO_NO_USE_AddressingMode.ZERO_PAGE, 3,)
    STA_3 = (0x8D, DO_NO_USE_AddressingMode.ABSOLUTE, 4,)
    STA_4 = (0x91, DO_NO_USE_AddressingMode.INDIRECT_Y, 6,)
    STA_5 = (0x95, DO_NO_USE_AddressingMode.ZERO_PAGE_X, 4,)
    STA_6 = (0x99, DO_NO_USE_AddressingMode.ABSOLUTE_Y, 5,)
    STA_7 = (0x9D, DO_NO_USE_AddressingMode.ABSOLUTE_X, 5,)
    STX_1 = (0x86, DO_NO_USE_AddressingMode.ZERO_PAGE, 3,)
    STX_2 = (0x8E, DO_NO_USE_AddressingMode.ABSOLUTE, 4,)
    STX_3 = (0x96, DO_NO_USE_AddressingMode.ZERO_PAGE_Y, 4,)
    STY_1 = (0x84, DO_NO_USE_AddressingMode.ZERO_PAGE, 3,)
    STY_2 = (0x8C, DO_NO_USE_AddressingMode.ABSOLUTE, 4,)
    STY_3 = (0x94, DO_NO_USE_AddressingMode.ZERO_PAGE_X, 4,)
    TAX = (0xAA, None, 2,)
    TAY = (0xA8, None, 2,)
    TSX = (0xBA, None, 2,)
    TXA = (0x8A, None, 2,)
    TXS = (0x9A, None, 2,)
    TYA = (0x98, None, 2,)

    def __init__(self, value, addressing_mode=None, cycles=None):
        super().__init__()
        self._value_ = value
        self.addressing_mode = addressing_mode
        self.cycles = cycles

    @staticmethod
    def from_hex(code):
        for op in DO_NOT_USE_Opcodes:
            if op.value == code:
                return op
