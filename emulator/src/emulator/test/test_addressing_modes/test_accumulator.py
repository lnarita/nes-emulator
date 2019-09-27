from emulator.adressing import Accumulator
from emulator.cpu import CPU
from emulator.memory import Memory, MemoryPositions


def test_read_Accumulator():
    address_mode = Accumulator
    cpu = CPU()
    memory = Memory()

    cpu.a = 0x00
    cpu.inc_cycle_by(-cpu.cycle)
    address = address_mode.fetch_address(cpu, memory)
    value = address_mode.read_from(cpu, memory, address)

    assert address is None
    assert value == 0x00

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.inc_pc_by(-cpu.pc)
    cpu.inc_pc_by(MemoryPositions.PRG_ROM_START.start)
    cpu.a = 0x10
    address = address_mode.fetch_address(cpu, memory)
    value = address_mode.read_from(cpu, memory, address)

    assert address is None
    assert value == 0x10

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.inc_pc_by(-cpu.pc)
    cpu.inc_pc_by(MemoryPositions.PRG_ROM_START.start)
    cpu.a = 0xFF
    address = address_mode.fetch_address(cpu, memory)
    value = address_mode.read_from(cpu, memory, address)

    assert address is None
    assert value == 0xFF


def test_read_Accumulator():
    address_mode = Accumulator
    cpu = CPU()
    memory = Memory()

    cpu.inc_cycle_by(-cpu.cycle)
    address = address_mode.fetch_address(cpu, memory)
    address_mode.write_to(cpu, memory, address, 5)

    assert address is None
    assert cpu.a == 5

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.inc_pc_by(-cpu.pc)
    cpu.inc_pc_by(MemoryPositions.PRG_ROM_START.start)
    cpu.a = 0x10
    address = address_mode.fetch_address(cpu, memory)
    address_mode.write_to(cpu, memory, address, 105)

    assert address is None
    assert cpu.a == 105

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.inc_pc_by(-cpu.pc)
    cpu.inc_pc_by(MemoryPositions.PRG_ROM_START.start)
    cpu.a = 0xFF
    address = address_mode.fetch_address(cpu, memory)
    address_mode.write_to(cpu, memory, address, 255)

    assert address is None
    assert cpu.a == 255
