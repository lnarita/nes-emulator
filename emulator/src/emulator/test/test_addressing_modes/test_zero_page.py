from emulator.adressing import ZeroPage
from emulator.cpu import CPU
from emulator.memory import Memory, MemoryPositions


def test_read_ZeroPage():
    address_mode = ZeroPage
    cpu = CPU()
    memory = Memory(rom=list(range(256)), ram=list(reversed(range(256))))

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.inc_pc_by(-cpu.pc + MemoryPositions.PRG_ROM_START.start)
    address = address_mode.fetch_address(cpu, memory)
    value = address_mode.read_from(cpu, memory, address)
    assert address == 0
    assert value == 255
    assert cpu.cycle == 2

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.inc_pc_by(-cpu.pc + MemoryPositions.PRG_ROM_START.start + 20)
    address = address_mode.fetch_address(cpu, memory)
    value = address_mode.read_from(cpu, memory, address)

    assert address == 20
    assert value == 235
    assert cpu.cycle == 2

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.inc_pc_by(-cpu.pc + MemoryPositions.PRG_ROM_START.start + 100)
    address = address_mode.fetch_address(cpu, memory)
    value = address_mode.read_from(cpu, memory, address)

    assert address == 100
    assert value == 155
    assert cpu.cycle == 2

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.inc_pc_by(-cpu.pc + MemoryPositions.PRG_ROM_START.start + 255)
    address = address_mode.fetch_address(cpu, memory)
    value = address_mode.read_from(cpu, memory, address)

    assert address == 255
    assert value == 0
    assert cpu.cycle == 2


def test_write_ZeroPage():
    address_mode = ZeroPage
    cpu = CPU()
    memory = Memory(rom=list(range(256)))

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.inc_pc_by(-cpu.pc + MemoryPositions.PRG_ROM_START.start)
    address = address_mode.fetch_address(cpu, memory)
    address_mode.write_to(cpu, memory, address, 6)
    assert address == 0
    assert memory.ram[address] == 6
    assert cpu.cycle == 2

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.inc_pc_by(-cpu.pc + MemoryPositions.PRG_ROM_START.start + 20)
    address = address_mode.fetch_address(cpu, memory)
    address_mode.write_to(cpu, memory, address, 105)

    assert address == 20
    assert memory.ram[address] == 105
    assert cpu.cycle == 2

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.inc_pc_by(-cpu.pc + MemoryPositions.PRG_ROM_START.start + 100)
    address = address_mode.fetch_address(cpu, memory)
    address_mode.write_to(cpu, memory, address, 143)

    assert address == 100
    assert memory.ram[address] == 143
    assert cpu.cycle == 2

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.inc_pc_by(-cpu.pc + MemoryPositions.PRG_ROM_START.start + 255)
    address = address_mode.fetch_address(cpu, memory)
    address_mode.write_to(cpu, memory, address, 255)

    assert address == 255
    assert memory.ram[address] == 255
    assert cpu.cycle == 2