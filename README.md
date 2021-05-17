# pscan
scan a target machine for all open ports with go
<br>
works on windows or linux (not tested on mac). I learned how to do this from reading https://nostarch.com/blackhatgo

```bash
usage: pscan ADDR [ARGS]

optional args:
--workers   how many workers to dispatch (max is 1000)
--wait      how long to wait in ms before we fail the port (default is 90)
--range     range of ports to scan (42-6666)

examples: (on windows command is pscan.exe)

$ pscan 192.168.1.87 (defaults applied are ports 0-65535 with 250 workers
                      waiting 90ms)

$ pscan 192.168.1.87 --workers 666 --wait 100

$ pscan 192.168.1.87 --workers 666 --wait 100 --range 40000-42000
```
