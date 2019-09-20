class OpCode(object):
    def __init__(self, code, addressing_mode, cycles):
        self.id = code
        self.addressing_mode = addressing_mode
        self.cycles = cycles

    def exec(self, cpu, memory):
        raise NotImplementedError("Class {} not implemented yet! CPU: {}, MEM: {}".format(cls, cpu, memory))

    @classmethod
    def create_variations(cls):
        pass

    @classmethod
    def create_dict_entry(cls, x):
        return tuple((x[0], cls(*x)))

    def __str__(self):
        return "{}(code={:02x}, addr_mode={}, cycles={})".format(type(self).__name__, self.id, self.addressing_mode, self.cycles)
