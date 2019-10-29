# pygame silent import
import os, sys
with open(os.devnull, 'w') as f:
    oldstdout = sys.stdout
    sys.stdout = f

    import pygame

    sys.stdout = oldstdout

class Controller:
    def __init__(self, read_count=0, strobe=False):
    	self.keys = [False] * 8
    	self.read_count = read_count
    	self.strobe = strobe

    def read_keys(self):
        # print(self.read_count)
        if self.read_count >= 8:
            return 1

        if not self.strobe:
            key = self.keys[self.read_count]
            # print("keys", self.keys)
            return int(key)
        else:
            pygame.event.pump()
            pressed = pygame.key.get_pressed()
            self.keys[0] = pressed[pygame.K_3]      # A
            self.keys[1] = pressed[pygame.K_2]      # B
            self.keys[2] = pressed[pygame.K_TAB]    # Select
            self.keys[3] = pressed[pygame.K_RETURN] # Start
            self.keys[4] = pressed[pygame.K_w]      # Up
            self.keys[5] = pressed[pygame.K_s]      # Down
            self.keys[6] = pressed[pygame.K_a]      # Left
            self.keys[7] = pressed[pygame.K_d]      # Right

            key = self.keys[self.read_count]
            # print("keys", self.keys)
            return int(key)


    def latch_keys(self, value):
        self.strobe = value
        pygame.event.pump()
        pressed = pygame.key.get_pressed()
        self.keys[0] = pressed[pygame.K_3]      # A
        self.keys[1] = pressed[pygame.K_2]      # B
        self.keys[2] = pressed[pygame.K_TAB]    # Select
        self.keys[3] = pressed[pygame.K_RETURN] # Start
        self.keys[4] = pressed[pygame.K_w]      # Up
        self.keys[5] = pressed[pygame.K_s]      # Down
        self.keys[6] = pressed[pygame.K_a]      # Left
        self.keys[7] = pressed[pygame.K_d]      # Right