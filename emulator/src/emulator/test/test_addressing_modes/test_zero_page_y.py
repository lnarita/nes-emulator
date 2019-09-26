from emulator.adressing import ZeroPageY
from emulator.cpu import CPU
from emulator.memory import Memory, MemoryPositions


def test_read_ZeroPageY():
    address_mode = ZeroPageY
    cpu = CPU()
    memory = Memory(rom=list(range(256)), ram=list(reversed(range(256))))

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.inc_pc_by(-cpu.pc + MemoryPositions.PRG_ROM_START.start)
    cpu.y = 1
    address = address_mode.fetch_address(cpu, memory)
    value = address_mode.read_from(cpu, memory, address)
    assert address == 1
    assert value == 254
    assert cpu.cycle == 3

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.inc_pc_by(-cpu.pc + MemoryPositions.PRG_ROM_START.start + 20)
    cpu.y = 10
    address = address_mode.fetch_address(cpu, memory)
    value = address_mode.read_from(cpu, memory, address)

    assert address == 30
    assert value == 225
    assert cpu.cycle == 3

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.inc_pc_by(-cpu.pc + MemoryPositions.PRG_ROM_START.start + 100)
    cpu.y = 5
    address = address_mode.fetch_address(cpu, memory)
    value = address_mode.read_from(cpu, memory, address)

    assert address == 105
    assert value == 150
    assert cpu.cycle == 3

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.inc_pc_by(-cpu.pc + MemoryPositions.PRG_ROM_START.start + 255)
    cpu.y = 0
    address = address_mode.fetch_address(cpu, memory)
    value = address_mode.read_from(cpu, memory, address)

    assert address == 255
    assert value == 0
    assert cpu.cycle == 3

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.inc_pc_by(-cpu.pc + MemoryPositions.PRG_ROM_START.start + 255)
    cpu.y = 5
    address = address_mode.fetch_address(cpu, memory)
    value = address_mode.read_from(cpu, memory, address)

    assert address == 4
    assert value == 251
    assert cpu.cycle == 3


def test_write_ZeroPageY():
    address_mode = ZeroPageY
    cpu = CPU()
    memory = Memory(rom=list(range(256)))

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.inc_pc_by(-cpu.pc + MemoryPositions.PRG_ROM_START.start)
    cpu.y = 5
    address = address_mode.fetch_address(cpu, memory)
    address_mode.write_to(cpu, memory, address, 6)
    assert address == 5
    assert memory.ram[address] == 6
    assert cpu.cycle == 3

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.inc_pc_by(-cpu.pc + MemoryPositions.PRG_ROM_START.start + 20)
    cpu.y = 15
    address = address_mode.fetch_address(cpu, memory)
    address_mode.write_to(cpu, memory, address, 105)

    assert address == 35
    assert memory.ram[address] == 105
    assert cpu.cycle == 3

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.inc_pc_by(-cpu.pc + MemoryPositions.PRG_ROM_START.start + 100)
    cpu.y = 155
    address = address_mode.fetch_address(cpu, memory)
    address_mode.write_to(cpu, memory, address, 143)

    assert address == 255
    assert memory.ram[address] == 143
    assert cpu.cycle == 3

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.inc_pc_by(-cpu.pc + MemoryPositions.PRG_ROM_START.start + 255)
    cpu.y = 0
    address = address_mode.fetch_address(cpu, memory)
    address_mode.write_to(cpu, memory, address, 255)

    assert address == 255
    assert memory.ram[address] == 255
    assert cpu.cycle == 3

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.inc_pc_by(-cpu.pc + MemoryPositions.PRG_ROM_START.start + 255)
    cpu.y = 75
    address = address_mode.fetch_address(cpu, memory)
    address_mode.write_to(cpu, memory, address, 255)

    assert address == 74
    assert memory.ram[address] == 255
    assert cpu.cycle == 3
