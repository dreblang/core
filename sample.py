import time

st = time.time()

def loopfunc(n):
	a = 0
	sum = 0
	while (a < n):
		sum = a + sum
		a = a + 1

	return sum
print(loopfunc(10000000))
print(time.time() - st)
