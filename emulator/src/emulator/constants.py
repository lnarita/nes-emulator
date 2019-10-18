from enum import auto, Enum

CPU_FREQUENCY = 1.79e6
CYCLE_PERIOD = 1. / CPU_FREQUENCY  # seconds
CYCLE_PERIOD_SIZE = 10

KB = 1024

HIGH_BITS_MASK = 0b1111111100000000
LOW_BITS_MASK = 0b0000000011111111

NEGATIVE_BIT = 0b0000000010000000
