from emulator.adressing import IndirectX, IndirectY
from emulator.cpu import CPU
from emulator.memory import Memory, MemoryPositions


def test_read_IndirectY():
    address_mode = IndirectY
    cpu = CPU()
    memory = Memory(rom=[0x00, 0x02, 0x7F, 0x01, 0xFF, 0x07], ram=list(map(lambda x: x % 256, range(Memory.ram_size()))))

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.inc_pc_by(-cpu.pc + MemoryPositions.PRG_ROM_START.start)
    cpu.y = 4
    address = address_mode.fetch_address(cpu, memory)
    value = address_mode.read_from(cpu, memory, address)
    assert address == 0x0104
    assert value == address % 256
    assert cpu.cycle == 5

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.y = 0xFF
    address = address_mode.fetch_address(cpu, memory)
    value = address_mode.read_from(cpu, memory, address)
    assert address == 0x0401
    assert value == address % 256
    assert cpu.cycle == 6


def test_write_IndirectY():
    address_mode = IndirectY
    cpu = CPU()
    memory = Memory(rom=[0x00, 0x02, 0x7F, 0x01, 0xFF, 0x07], ram=list(map(lambda x: x % 256, range(Memory.ram_size()))))

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.inc_pc_by(-cpu.pc + MemoryPositions.PRG_ROM_START.start)
    cpu.y = 4
    address = address_mode.fetch_address(cpu, memory)
    address_mode.write_to(cpu, memory, address, 9)
    assert address == 0x0104
    assert memory.ram[address] == 9
    assert cpu.cycle == 5

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.y = 0xFF
    address = address_mode.fetch_address(cpu, memory)
    address_mode.write_to(cpu, memory, address, 50)
    assert address == 0x0401
    assert memory.ram[address] == 50
    assert cpu.cycle == 6
