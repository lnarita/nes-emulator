from emulator.constants import HIGH_BITS_MASK, LOW_BITS_MASK
from emulator.memory import MemoryPositions


class AddressMode:
    @classmethod
    def write_to(cls, cpu, memory, addr, value):
        result = cpu.exec_in_cycle(memory.store, addr, value)
        if addr == 0x2007:
            if memory.ppu.ppuctrl & 0b0000100:
                memory.ppu.ppuaddr += 32
            else:
                memory.ppu.ppuaddr += 1
        elif addr == 0x2004:
            memory.ppu.oamaddr += 1
        elif addr == 0x4016:
            memory.controller.read_count = 0

        return result

    @classmethod
    def read_from(cls, cpu, memory, addr):
        result = cpu.exec_in_cycle(memory.fetch, addr)
        if addr == 0x2007:
            if memory.ppu.ppuctrl & 0b0000100:
                memory.ppu.ppuaddr += 32
            else:
                memory.ppu.ppuaddr += 1
        elif addr == 0x4016:
            memory.controller.read_count += 1
        return result

    @classmethod
    def fetch_address(cls, cpu, memory):
        """
        This is the method that should fetch the parameter(s) from the ROM binary,
        for instructions that deal with memory, call this method, save the return as the address it want to read / modify
        and use it in `read_from` or `write_to`
        """
        pass

    @classmethod
    def read_16_bits_low(cls, memory, addr):
        addr_low = memory.fetch(addr)
        return addr_low & LOW_BITS_MASK

    @classmethod
    def read_16_bits_high(cls, memory, addr):
        addr_high = memory.fetch(addr)
        return (addr_high << 8) & HIGH_BITS_MASK

    @classmethod
    def get_16_bits_addr_from_high_low(cls, high, low):
        """
        Combine high and low to create a single 16 bits address
        """
        return high | low


class Indirect(AddressMode):
    """
        #   address  R/W description
       --- --------- --- ------------------------------------------
        1     PC      R  fetch opcode, increment PC
        2     PC      R  fetch pointer address low, increment PC
        3     PC      R  fetch pointer address high, increment PC
        4   pointer   R  fetch low address to latch
        5  pointer+1* R  fetch PCH, copy latch to PCL

       Note: * The PCH will always be fetched from the same page
               than PCL, i.e. page boundary crossing is not handled.
    """

    @classmethod
    def write_to(cls, cpu, memory, addr, value):
        pass

    @classmethod
    def read_from(cls, cpu, memory, addr):
        pass

    @classmethod
    def fetch_address(cls, cpu, memory):
        def _read_ptr_low():
            low = cls.read_16_bits_low(memory, cpu.pc)
            cls.low = low
            cpu.inc_pc_by(1)
            return low

        def _read_ptr_high():
            high = cls.read_16_bits_high(memory, cpu.pc)
            cls.high = high
            cpu.inc_pc_by(1)
            return high

        def _read_addr_low(addr):
            low = cls.read_16_bits_low(memory, addr)
            return low

        def _read_addr_high_and_return_address(addr, l):
            if (addr & LOW_BITS_MASK) == 0xFF:
                h = cls.read_16_bits_high(memory, addr & HIGH_BITS_MASK)
            else:
                h = cls.read_16_bits_high(memory, addr + 1)
            value = cls.get_16_bits_addr_from_high_low(l, h)
            return value

        cls.low = None
        cls.high = None
        cls.addr = None
        cls.data = None
        pointer_low = cpu.exec_in_cycle(_read_ptr_low)  # 2
        pointer_high = cpu.exec_in_cycle(_read_ptr_high)  # 3
        pointer = cls.get_16_bits_addr_from_high_low(pointer_low, pointer_high)
        address_low = cpu.exec_in_cycle(_read_addr_low, pointer)  # 4
        address = cpu.exec_in_cycle(_read_addr_high_and_return_address, pointer, address_low)  # 5
        cls.addr = "($%04X)" % pointer
        cls.data = "= %04X" % address
        return address


class IndirectX(AddressMode):
    """
        #    address   R/W description
       --- ----------- --- ------------------------------------------
        1      PC       R  fetch opcode, increment PC
        2      PC       R  fetch pointer address, increment PC
        3    pointer    R  read from the address, add X to it
        4   pointer+X   R  fetch effective address low
        5  pointer+X+1  R  fetch effective address high
    """

    @classmethod
    def fetch_address(cls, cpu, memory):
        def _read_ptr():
            ptr = memory.fetch(cpu.pc)
            cls.low = ptr
            cpu.inc_pc_by(1)
            return ptr

        def _calc_real_addr(addr):
            # The effective address is wrapped around to always land on the zero page
            real_addr = MemoryPositions.ZERO_PAGE.wrap(addr + cpu.x)
            return real_addr

        def _read_addr_low(addr):
            low = cls.read_16_bits_low(memory, addr)
            return low

        def _read_addr_high(addr, l):
            if (addr & LOW_BITS_MASK) == 0xFF:
                h = cls.read_16_bits_high(memory, addr & HIGH_BITS_MASK)
            else:
                h = cls.read_16_bits_high(memory, addr + 1)
            value = cls.get_16_bits_addr_from_high_low(l, h)
            return value

        cls.low = None
        cls.high = None
        cls.addr = None
        cls.data = None
        pointer = cpu.exec_in_cycle(_read_ptr)  # 2
        real_pointer = cpu.exec_in_cycle(_calc_real_addr, pointer)  # 3
        addr_low = cpu.exec_in_cycle(_read_addr_low, real_pointer)  # 4
        effective_addr = cpu.exec_in_cycle(_read_addr_high, real_pointer, addr_low)  # 5
        cls.addr = "(${0:02X},X) @ {1:02X} = {2:04X}".format(pointer, real_pointer, effective_addr)
        return effective_addr


class IndirectY(AddressMode):
    """
        #    address   R/W description
       --- ----------- --- ------------------------------------------
        1      PC       R  fetch opcode, increment PC
        2      PC       R  fetch pointer address, increment PC
        3    pointer    R  fetch effective address low
        4   pointer+1   R  fetch effective address high,
                           add Y to low byte of effective address
        5   address+Y*  R  read from effective address,
                           fix high byte of effective address

        The 6502 has one 8-bit ALU and one 16-bit upcounter (for PC).
        To calculate a,x or a,y addressing in an instruction other than STA, STX, or STY,
            it uses the 8-bit ALU to first calculate the low byte while it fetches the high byte.
        If there's a carry out, it goes "oops", applies the carry using the ALU, and repeats the read at the correct address.
        Store instructions always have this "oops" cycle: the CPU first reads from the partially added address and then writes to the correct address.
        The same thing happens on (d),y indirect addressing.
    """

    @classmethod
    def fetch_address(cls, cpu, memory):
        def _read_ptr():
            ptr = memory.fetch(cpu.pc)
            cls.low = ptr
            cpu.inc_pc_by(1)
            return ptr

        def _read_addr_low_from_pointer(ptr):
            low = cls.read_16_bits_low(memory, ptr)
            return low

        def _read_addr_high_from_pointer(ptr, low):
            high = cls.read_16_bits_high(memory, MemoryPositions.ZERO_PAGE.wrap(ptr + 1))
            real_low = low + cpu.y
            return high, real_low

        def _read_from_real_addr(high, low):
            def handle_overflow():
                h = (high + overflow) & HIGH_BITS_MASK
                l = low & LOW_BITS_MASK
                return h, l

            overflow = low & HIGH_BITS_MASK
            if overflow > 0:
                real_high, real_low = cpu.exec_in_cycle(handle_overflow)  # stall
            else:
                real_high, real_low = high, low
            real_addr = cls.get_16_bits_addr_from_high_low(real_high, real_low)
            return real_addr

        cls.low = None
        cls.high = None
        cls.addr = None
        cls.data = None

        pointer = cpu.exec_in_cycle(_read_ptr)  # 2
        real_pointer = cpu.exec_in_cycle(_read_addr_low_from_pointer, pointer)  # 3
        addr_high, addr_low = cpu.exec_in_cycle(_read_addr_high_from_pointer, pointer, real_pointer)  # 4
        effective_addr = _read_from_real_addr(addr_high, addr_low)
        cls.addr = "($%02X),Y = %04X @ %04X" % (pointer, (addr_high & HIGH_BITS_MASK) | (real_pointer & LOW_BITS_MASK), effective_addr)
        return effective_addr


class ZeroPage(AddressMode):
    """
        #  address R/W description
       --- ------- --- ------------------------------------------
        1    PC     R  fetch opcode, increment PC
        2    PC     R  fetch address, increment PC
    """

    @classmethod
    def fetch_address(cls, cpu, memory):
        def _read_addr():
            low = cls.read_16_bits_low(memory, cpu.pc)
            cls.low = low
            cpu.inc_pc_by(1)
            return low

        cls.low = None
        cls.high = None
        cls.addr = None
        cls.data = None
        addr = cpu.exec_in_cycle(_read_addr)
        cls.addr = "$%02X" % addr
        return addr


class ZeroPageX(AddressMode):
    """
        #   address  R/W description
       --- --------- --- ------------------------------------------
        1     PC      R  fetch opcode, increment PC
        2     PC      R  fetch address, increment PC
        3   address   R  read from address, add index register to it
    """

    @classmethod
    def fetch_address(cls, cpu, memory):
        def _read_addr():
            low = cls.read_16_bits_low(memory, cpu.pc)
            cls.low = low
            cpu.inc_pc_by(1)
            return low

        def _calc_real_addr(a):
            return MemoryPositions.ZERO_PAGE.wrap(a + cpu.x)

        cls.low = None
        cls.high = None
        cls.addr = None
        cls.data = None
        addr = cpu.exec_in_cycle(_read_addr)  # 2
        real_addr = cpu.exec_in_cycle(_calc_real_addr, addr)  # 3
        cls.addr = "$%02X,X @ %02X" % (addr, real_addr)
        return real_addr


class ZeroPageY(AddressMode):
    """
        #   address  R/W description
       --- --------- --- ------------------------------------------
        1     PC      R  fetch opcode, increment PC
        2     PC      R  fetch address, increment PC
        3   address   R  read from address, add index register to it
    """

    @classmethod
    def fetch_address(cls, cpu, memory):
        def _read_addr():
            low = cls.read_16_bits_low(memory, cpu.pc)
            cls.low = low
            cpu.inc_pc_by(1)
            return low

        def _calc_real_addr(a):
            return MemoryPositions.ZERO_PAGE.wrap(a + cpu.y)

        cls.low = None
        cls.high = None
        cls.addr = None
        cls.data = None
        addr = cpu.exec_in_cycle(_read_addr)  # 2
        real_addr = cpu.exec_in_cycle(_calc_real_addr, addr)  # 3
        cls.addr = "$%02X,Y @ %02X" % (addr, real_addr)
        return real_addr


class Absolute(AddressMode):
    """
        #  address R/W description
       --- ------- --- ------------------------------------------
        1    PC     R  fetch opcode, increment PC
        2    PC     R  fetch low byte of address, increment PC
        3    PC     R  fetch high byte of address, increment PC
    """

    @classmethod
    def fetch_address(cls, cpu, memory):
        def _read_low():
            low = cls.read_16_bits_low(memory, cpu.pc)
            cls.low = low
            cpu.inc_pc_by(1)
            return low

        def _read_high():
            high = cls.read_16_bits_high(memory, cpu.pc)
            cls.high = high
            cpu.inc_pc_by(1)
            return high

        cls.low = None
        cls.high = None
        cls.addr = None
        cls.data = None
        addr_low = cpu.exec_in_cycle(_read_low)  # 2
        addr_high = cpu.exec_in_cycle(_read_high)  # 3
        addr = cls.get_16_bits_addr_from_high_low(addr_low, addr_high)
        cls.addr = "$%04X" % addr
        return addr


class AbsoluteY(AddressMode):
    """
        #   address  R/W description
       --- --------- --- ------------------------------------------
        1     PC      R  fetch opcode, increment PC
        2     PC      R  fetch low byte of address, increment PC
        3     PC      R  fetch high byte of address,
                         add index register to low address byte,
                         increment PC
        4  address+I* R  read from effective address,
                         fix the high byte of effective address
    """

    @classmethod
    def fetch_address(cls, cpu, memory):
        def _read_addr_low():
            low = cls.read_16_bits_low(memory, cpu.pc)
            cls.low = low
            cpu.inc_pc_by(1)
            return low

        def _read_addr_high(low):
            high = cls.read_16_bits_high(memory, cpu.pc)
            cls.high = high
            real_low = low + cpu.y
            cpu.inc_pc_by(1)
            return high, real_low

        def _read_from_real_addr(high, low):
            def handle_overflow():
                h = (high + overflow) & HIGH_BITS_MASK
                l = low & LOW_BITS_MASK
                return h, l

            overflow = low & HIGH_BITS_MASK
            if overflow > 0:
                real_high, real_low = cpu.exec_in_cycle(handle_overflow)  # stall (4)
            else:
                real_high, real_low = high, low
            real_addr = cls.get_16_bits_addr_from_high_low(real_high, real_low)
            return real_addr

        cls.low = None
        cls.high = None
        cls.addr = None
        cls.data = None

        low_before_inc = cpu.exec_in_cycle(_read_addr_low)  # 2
        high_no_fix, low = cpu.exec_in_cycle(_read_addr_high, low_before_inc)  # 3
        effective_addr = _read_from_real_addr(high_no_fix, low)
        cls.addr = "$%04X,Y @ %04X" % ((high_no_fix | low_before_inc), effective_addr)
        return effective_addr


class AbsoluteX(AddressMode):
    """
        #   address  R/W description
       --- --------- --- ------------------------------------------
        1     PC      R  fetch opcode, increment PC
        2     PC      R  fetch low byte of address, increment PC
        3     PC      R  fetch high byte of address,
                         add index register to low address byte,
                         increment PC
        4  address+I* R  read from effective address,
                         fix the high byte of effective address
    """

    @classmethod
    def fetch_address(cls, cpu, memory):
        def _read_addr_low():
            low = cls.read_16_bits_low(memory, cpu.pc)
            cls.low = low
            cpu.inc_pc_by(1)
            return low

        def _read_addr_high(low):
            high = cls.read_16_bits_high(memory, cpu.pc)
            cls.high = high
            real_low = low + cpu.x
            cpu.inc_pc_by(1)
            return high, real_low

        def _read_from_real_addr(high, low):
            def handle_overflow():
                h = (high + overflow) & HIGH_BITS_MASK
                l = low & LOW_BITS_MASK
                return h, l

            overflow = low & HIGH_BITS_MASK
            if overflow > 0:
                real_high, real_low = cpu.exec_in_cycle(handle_overflow)  # stall (4)
            else:
                real_high, real_low = high, low
            real_addr = cls.get_16_bits_addr_from_high_low(real_high, real_low)
            return real_addr

        cls.low = None
        cls.high = None
        cls.addr = None
        cls.data = None

        low_before_inc = cpu.exec_in_cycle(_read_addr_low)  # 2
        high_no_fix, low = cpu.exec_in_cycle(_read_addr_high, low_before_inc)  # 3
        effective_addr = _read_from_real_addr(high_no_fix, low)
        cls.addr = "$%04X,X @ %04X" % ((high_no_fix | low_before_inc), effective_addr)
        return effective_addr


class Immediate(AddressMode):
    """
        #  address R/W description
       --- ------- --- ------------------------------------------
        1    PC     R  fetch opcode, increment PC
        2    PC     R  fetch value, increment PC
    """

    @classmethod
    def write_to(cls, cpu, memory, addr, value):
        pass

    @classmethod
    def read_from(cls, cpu, memory, addr):
        return addr

    @classmethod
    def fetch_address(cls, cpu, memory):
        def _read_immediate():
            low = cls.read_16_bits_low(memory, cpu.pc)
            cls.low = low
            cpu.inc_pc_by(1)
            return low

        cls.low = None
        cls.high = None
        cls.addr = None
        cls.data = None
        value = cpu.exec_in_cycle(_read_immediate)
        cls.addr = "#$%02X" % value
        return value


class Accumulator(AddressMode):
    """
        #  address R/W description
       --- ------- --- -----------------------------------------------
        1    PC     R  fetch opcode, increment PC
        2    PC     R  read next instruction byte (and throw it away)
    """

    @classmethod
    def write_to(cls, cpu, memory, addr, value):
        cpu.a = value

    @classmethod
    def read_from(cls, cpu, memory, addr):
        return cpu.a

    @classmethod
    def fetch_address(cls, cpu, memory):
        def _stall():
            pass

        cls.low = None
        cls.high = None
        cls.addr = None
        cls.data = None
        cls.addr = 'A'

        cpu.exec_in_cycle(_stall)


class Relative(AddressMode):
    """
        #   address  R/W description
       --- --------- --- ---------------------------------------------
        1     PC      R  fetch opcode, increment PC
        2     PC      R  fetch operand, increment PC
        3     PC      R  Fetch opcode of next instruction,
                         If branch is taken, add operand to PCL.
                         Otherwise increment PC.
        4+    PC*     R  Fetch opcode of next instruction.
                         Fix PCH. If it did not change, increment PC.
        5!    PC      R  Fetch opcode of next instruction,
                         increment PC.

       Notes: The opcode fetch of the next instruction is included to
              this diagram for illustration purposes. When determining
              real execution times, remember to subtract the last
              cycle.

              * The high byte of Program Counter (PCH) may be invalid
                at this time, i.e. it may be smaller or bigger by $100.

              + If branch is taken, this cycle will be executed.

              ! If branch occurs to different page, this cycle will be
                executed.
    """

    @classmethod
    def write_to(cls, cpu, memory, addr, value):
        pass

    @classmethod
    def fetch_address(cls, cpu, memory):
        def _read_addr():
            low = cls.read_16_bits_low(memory, cpu.pc)
            cls.low = low
            cpu.inc_pc_by(1)
            return low

        cls.low = None
        cls.high = None
        cls.addr = None
        cls.data = None
        addr = _read_addr()
        cls.addr = "$%02X" % addr
        return addr
