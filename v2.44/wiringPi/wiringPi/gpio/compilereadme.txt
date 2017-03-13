echo [Compile] gpio.c
gcc -c -O2 -Wall -Wextra -I/usr/local/include -Winline -pipe gpio.c -o gpio.o
echo [Compile] readall.c
gcc -c -O2 -Wall -Wextra -I/usr/local/include -Winline -pipe readall.c -o readall.o
echo [Compile] pins.c
gcc -c -O2 -Wall -Wextra -I/usr/local/include -Winline -pipe pins.c -o pins.o
echo [Link]
gcc -o gpio gpio.o readall.o pins.o -L/usr/local/lib -lwiringPi -lwiringPiDev -lpthread -lrt -lm -lcrypt
