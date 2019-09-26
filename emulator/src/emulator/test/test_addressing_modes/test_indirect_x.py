from emulator.adressing import IndirectX
from emulator.cpu import CPU
from emulator.memory import Memory, MemoryPositions


def test_read_IndirectX():
    address_mode = IndirectX
    cpu = CPU()
    memory = Memory(rom=[0x00, 0x02, 0x7F, 0x01, 0xFF, 0x07], ram=list(map(lambda x: x % 256, range(Memory.ram_size()))))

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.inc_pc_by(-cpu.pc + MemoryPositions.PRG_ROM_START.start)
    cpu.x = 4
    address = address_mode.fetch_address(cpu, memory)
    value = address_mode.read_from(cpu, memory, address)
    assert address == (((0x0504 + 1) % 256) << 8 | (0x0504 % 256))
    assert value == address % 256
    assert cpu.cycle == 5

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.x = 0xFF
    address = address_mode.fetch_address(cpu, memory)
    value = address_mode.read_from(cpu, memory, address)
    assert address == (((0x0201 + 1) % 256) << 8 | (0x0201 % 256))
    assert value == address % 256
    assert cpu.cycle == 5


def test_read_IndirectX():
    address_mode = IndirectX
    cpu = CPU()
    memory = Memory(rom=[0x00, 0x02, 0x7F, 0x01, 0xFF, 0x07], ram=list(map(lambda x: x % 256, range(Memory.ram_size()))))

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.inc_pc_by(-cpu.pc + MemoryPositions.PRG_ROM_START.start)
    cpu.x = 4
    address = address_mode.fetch_address(cpu, memory)
    address_mode.write_to(cpu, memory, address, 6)
    assert address == (((0x0504 + 1) % 256) << 8 | (0x0504 % 256))
    assert memory.ram[address] == 6
    assert cpu.cycle == 5

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.x = 0xFF
    address = address_mode.fetch_address(cpu, memory)
    address_mode.write_to(cpu, memory, address, 9)
    assert address == (((0x0201 + 1) % 256) << 8 | (0x0201 % 256))
    assert memory.ram[address] == 9
    assert cpu.cycle == 5
