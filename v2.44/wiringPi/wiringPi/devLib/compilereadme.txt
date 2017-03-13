echo [Compile] ds1302.c
gcc -c -O2 -D_GNU_SOURCE -Wformat=2 -Wall -Winline -I. -pipe -fPIC ds1302.c -o ds1302.o
echo [Compile] maxdetect.c
gcc -c -O2 -D_GNU_SOURCE -Wformat=2 -Wall -Winline -I. -pipe -fPIC maxdetect.c -o maxdetect.o
echo [Compile] piNes.c
gcc -c -O2 -D_GNU_SOURCE -Wformat=2 -Wall -Winline -I. -pipe -fPIC piNes.c -o piNes.o
echo [Compile] gertboard.c
gcc -c -O2 -D_GNU_SOURCE -Wformat=2 -Wall -Winline -I. -pipe -fPIC gertboard.c -o gertboard.o
echo [Compile] piFace.c
gcc -c -O2 -D_GNU_SOURCE -Wformat=2 -Wall -Winline -I. -pipe -fPIC piFace.c -o piFace.o
echo [Compile] lcd128x64.c
gcc -c -O2 -D_GNU_SOURCE -Wformat=2 -Wall -Winline -I. -pipe -fPIC lcd128x64.c -o lcd128x64.o
echo [Compile] lcd.c
gcc -c -O2 -D_GNU_SOURCE -Wformat=2 -Wall -Winline -I. -pipe -fPIC lcd.c -o lcd.o
echo [Compile] scrollPhat.c
gcc -c -O2 -D_GNU_SOURCE -Wformat=2 -Wall -Winline -I. -pipe -fPIC scrollPhat.c -o scrollPhat.o
echo [Compile] piGlow.c
gcc -c -O2 -D_GNU_SOURCE -Wformat=2 -Wall -Winline -I. -pipe -fPIC piGlow.c -o piGlow.o
echo "[Link (Dynamic)]"
gcc -shared -Wl,-soname,libwiringPiDev.so -o libwiringPiDev.so.2.44 -lpthread ds1302.o maxdetect.o piNes.o gertboard.o piFace.o lcd128x64.o lcd.o scrollPhat.o piGlow.o
