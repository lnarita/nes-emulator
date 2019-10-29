# pygame silent import
import os, sys
with open(os.devnull, 'w') as f:
    oldstdout = sys.stdout
    sys.stdout = f

    import pygame

    sys.stdout = oldstdout

from pygame import Color
import time

W_WIDTH = 256*2
W_HEIGHT = 240*2




class Window():
	def __init__(self):
		check_errors = pygame.init()
		if check_errors[1] > 0:
		    print("pygame error")
		    sys.exit(-1)

		size = (W_WIDTH, W_HEIGHT)

		self.surface = pygame.display.set_mode(size)
		pygame.display.set_caption("NES emulator")

	def setPixel(self,i,j,color):
		pygame.draw.rect(self.surface, color, pygame.Rect(i*2,j*2,2,2))


	def flip(self):
		pygame.display.flip()


def main():
	screen = Window()
	for i in range(262):
		screen.scanLine(i)
	time.sleep(10)

if __name__ == "__main__":
	main()
