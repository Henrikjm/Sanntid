# Python 3.3.3 and 2.7.6
# python helloworld_python.py

from threading import Thread

i = 0

def adder():
# In Python you "import" a global variable, instead of "export"ing it when you declare it
    global i

    
# In Python 2 this generates a list of integers (which takes time),
# while in Python 3 this is an iterable (which is much faster to generate).
# In python 2, an iterable is created with xrange()
    for x in range(0, 999999):
        	i += 1

def adder2():
	global i
	for x in range(0, 1000000):
		i -= 1

def main():
	
	adder_thr = Thread(target = adder)
	adder_thr2 = Thread(target = adder2)
	adder_thr.start()
	adder_thr2.start()
	adder_thr.join()
	adder_thr2.join()
	print("Done: " + str(i))


main()