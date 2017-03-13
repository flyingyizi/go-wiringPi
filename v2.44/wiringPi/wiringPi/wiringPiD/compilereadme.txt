echo [Compile] wiringpid.c
gcc -c -O2 -Wall -Wextra -I/usr/local/include -Winline -pipe wiringpid.c -o wiringpid.o
echo [Compile] network.c
gcc -c -O2 -Wall -Wextra -I/usr/local/include -Winline -pipe network.c -o network.o
echo [Compile] runRemote.c
gcc -c -O2 -Wall -Wextra -I/usr/local/include -Winline -pipe runRemote.c -o runRemote.o
echo [Compile] daemonise.c
gcc -c -O2 -Wall -Wextra -I/usr/local/include -Winline -pipe daemonise.c -o daemonise.o
echo [Link]
gcc -o wiringpid wiringpid.o network.o runRemote.o daemonise.o -L/usr/local/lib -lwiringPi -lwiringPiDev -lpthread -lrt -lm -lcrypt
