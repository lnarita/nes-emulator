from emulator.adressing import AbsoluteY
from emulator.cpu import CPU
from emulator.memory import Memory, MemoryPositions


def test_read_AbsoluteY():
    address_mode = AbsoluteY
    cpu = CPU()
    memory = Memory(rom=[0x32, 0x02, 0x7F, 0x01, 0x00, 0x20], ram=list(map(lambda x: x % 256, range(Memory.ram_size()))))

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.inc_pc_by(-cpu.pc + MemoryPositions.PRG_ROM_START.start)
    cpu.y = 0x0A
    address = address_mode.fetch_address(cpu, memory)
    value = address_mode.read_from(cpu, memory, address)
    assert address == 0x023C
    assert value == (0x023C % 256)

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.y = 0xA1
    address = address_mode.fetch_address(cpu, memory)
    value = address_mode.read_from(cpu, memory, address)

    assert address == 0x0220
    assert value == (0x0220 % 256)

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.y = 0x07
    address = address_mode.fetch_address(cpu, memory)

    assert address == 0x2007
    # PPU do not assert value for now


def test_write_AbsoluteY():
    address_mode = AbsoluteY
    cpu = CPU()
    memory = Memory(rom=[0x30, 0x01, 0xFF, 0x02, 0x0A, 0x05], ram=list(map(lambda x: x % 256, range(Memory.ram_size()))))

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.inc_pc_by(-cpu.pc + MemoryPositions.PRG_ROM_START.start)
    cpu.y = 0x25
    address = address_mode.fetch_address(cpu, memory)
    address_mode.write_to(cpu, memory, address, 20)
    assert address == 0x0155
    assert memory.ram[address] == 20

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.y = 0x30
    address = address_mode.fetch_address(cpu, memory)
    address_mode.write_to(cpu, memory, address, 100)
    assert address == 0x032F
    assert memory.ram[address] == 100

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.y = 0x20
    address = address_mode.fetch_address(cpu, memory)
    address_mode.write_to(cpu, memory, address, 67)
    assert address == 0x052A
    assert memory.ram[address] == 67