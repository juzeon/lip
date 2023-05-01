# lip - Look up IP

`lip` is a versatile command-line tool designed for network administrators and enthusiasts, offering a suite of robust functions that enable users to perform a variety of network-related tasks efficiently and conveniently. With `lip`, you can execute IP address lookups, bandwidth tests, WHOIS lookups, reverse DNS lookups, SSL certificate checks, and TCPing and Telnet connections from the comfort of your terminal.

*This project is under development currently.*

## Install

Install `lip` by cloning this repository and compiling it:

```bash
git clone https://github.com/juzeon/lip
cd lip/
go install
```

Or, install `lip` with the go command below:

```bash
go install github.com/juzeon/lip@latest
```

## Usage

Type `lip -h` for help.

### IP Lookup

```bash
λ lip 24.48.0.1
Offline lookup result of 24.48.0.1:
+-----------+------------------+--------+--------+-----+------------+
|  Source   |     Country      | Region |  City  | ISP | Additional |
+-----------+------------------+--------+--------+-----+------------+
| ip2region |      加拿大      | 魁北克 | 魁北克 |     |            |
+-----------+------------------+--------+--------+-----+------------+
|   qqwry   | 加拿大 Videotron |        |        |     |            |
+-----------+------------------+--------+--------+-----+------------+
Fetching results from online sources...
Online lookup result of 24.48.0.1:
+--------+---------+--------+----------+--------------------------------+---------------------------------------+
| Source | Country | Region |   City   |              ISP               |              Additional               |
+--------+---------+--------+----------+--------------------------------+---------------------------------------+
| ip-api | Canada  | Quebec | Montreal |   Le Groupe Videotron Ltee,    |                                       |
|        |         |        |          |     Videotron Ltee, AS5769     |                                       |
|        |         |        |          |     Videotron Telecom Ltee     |                                       |
+--------+---------+        +----------+--------------------------------+---------------------------------------+
| ipinfo |   CA    |        |  Dorval  |  AS5769 Videotron Telecom Ltee |               hostname:               |
|        |         |        |          | videotron.com, Videotron Ltee  | modemcable001.0-48-24.mc.videotron.ca |
|        |         |        |          |                                |                                       |
+--------+---------+--------+----------+--------------------------------+---------------------------------------+
```

### TCPing

```bash
λ lip tcping 10.10.10.25 443
Probing 10.10.10.25:443/tcp - Port is open (open) - time=5.8126ms
Probing 10.10.10.25:443/tcp - Port is open (open) - time=2.6021ms
Probing 10.10.10.25:443/tcp - Port is open (open) - time=2.5113ms
Probing 10.10.10.25:443/tcp - Port is open (open) - time=3.6459ms
```

### Packet Sender
```bash
lip packet 127.0.0.1:9999 "hello"
hello
```

## Todos

- [x] IP Lookup
- [x] TCPing
- [x] Packet Sender
- [ ] WHOIS lookup
- [ ] Reverse DNS lookup
- [ ] SSL Certificate Check
- [ ] Bandwidth Test