from emulator.adressing import Absolute
from emulator.cpu import CPU
from emulator.memory import Memory, MemoryPositions


def test_read_Absolute():
    address_mode = Absolute
    cpu = CPU()
    memory = Memory(rom=[0x32, 0x02, 0x7F, 0x01, 0xFF, 0x07], ram=list(map(lambda x: x % 256, range(Memory.ram_size()))))

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.inc_pc_by(-cpu.pc + MemoryPositions.PRG_ROM_START.start)
    address = address_mode.fetch_address(cpu, memory)
    value = address_mode.read_from(cpu, memory, address)
    assert address == 0x0232
    assert value == (0x0232 % 256)
    assert cpu.cycle == 3

    cpu.inc_cycle_by(-cpu.cycle)
    address = address_mode.fetch_address(cpu, memory)
    value = address_mode.read_from(cpu, memory, address)

    assert address == 0x017F
    assert value == (0x017F % 256)
    assert cpu.cycle == 3

    cpu.inc_cycle_by(-cpu.cycle)
    address = address_mode.fetch_address(cpu, memory)
    value = address_mode.read_from(cpu, memory, address)

    assert address == 0x07FF
    assert value == (0x07FF % 256)
    assert cpu.cycle == 3


def test_write_Absolute():
    address_mode = Absolute
    cpu = CPU()
    memory = Memory(rom=[0x30, 0x01, 0xFF, 0x02, 0x0A, 0x05], ram=list(map(lambda x: x % 256, range(Memory.ram_size()))))

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.inc_pc_by(-cpu.pc + MemoryPositions.PRG_ROM_START.start)
    address = address_mode.fetch_address(cpu, memory)
    address_mode.write_to(cpu, memory, address, 50)
    assert address == 0x0130
    assert memory.ram[address] == 50
    assert cpu.cycle == 3

    cpu.inc_cycle_by(-cpu.cycle)
    address = address_mode.fetch_address(cpu, memory)
    address_mode.write_to(cpu, memory, address, 100)
    assert address == 0x02FF
    assert memory.ram[address] == 100
    assert cpu.cycle == 3

    cpu.inc_cycle_by(-cpu.cycle)
    address = address_mode.fetch_address(cpu, memory)
    address_mode.write_to(cpu, memory, address, 67)
    assert address == 0x050A
    assert memory.ram[address] == 67
    assert cpu.cycle == 3