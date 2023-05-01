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
| ip2region |      加拿大       | 魁北克 | 魁北克  |     |            |
+-----------+------------------+--------+--------+-----+------------+
|   qqwry   | 加拿大 Videotron  |        |        |     |            |
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
λ lip packet google.com:80 "GET /"
HTTP/1.0 200 OK
Content-Type: text/html; charset=ISO-8859-1
Server: gws
...
Vary: Accept-Encoding

<!doctype html>
...
```

### WHOIS Lookup

```bash
λ lip whois google.com
Domain:
+------------+------------+-----------------------+--------------------------------+--------------------------------+--------+----------------------+----------------------+----------------------+
|   DOMAIN   |  PUNYCODE  |      WHOISSERVER      |             STATUS             |          NAMESERVERS           | DNSSEC |     CREATEDDATE      |     UPDATEDDATE      |    EXPIRATIONDATE    |
+------------+------------+-----------------------+--------------------------------+--------------------------------+--------+----------------------+----------------------+----------------------+
| google.com | google.com | whois.markmonitor.com |    clientdeleteprohibited,     |        ns1.google.com,         | false  | 1997-09-15T04:00:00Z | 2019-09-09T15:39:04Z | 2028-09-14T04:00:00Z |
|            |            |                       |   clienttransferprohibited,    |        ns2.google.com,         |        |                      |                      |                      |
|            |            |                       |    clientupdateprohibited,     | ns3.google.com, ns4.google.com |        |                      |                      |                      |
|            |            |                       |    serverdeleteprohibited,     |                                |        |                      |                      |                      |
|            |            |                       |   servertransferprohibited,    |                                |        |                      |                      |                      |
|            |            |                       |     serverupdateprohibited     |                                |        |                      |                      |                      |
+------------+------------+-----------------------+--------------------------------+--------------------------------+--------+----------------------+----------------------+----------------------+

Contacts:
+----------------+------------------+--------------+--------+------+----------+------------+---------+---------------+--------------------------------------------------+----------------------------+
|     SOURCE     |       NAME       | ORGANIZATION | STREET | CITY | PROVINCE | POSTALCODE | COUNTRY |     PHONE     |                      EMAIL                       |        REFERRALURL         |
+----------------+------------------+--------------+--------+------+----------+------------+---------+---------------+--------------------------------------------------+----------------------------+
|   Registrar    | MarkMonitor Inc. |              |        |      |          |            |         | +1.2086851750 |         abusecomplaints@markmonitor.com          | http://www.markmonitor.com |
+----------------+------------------+--------------+--------+------+----------+------------+---------+---------------+--------------------------------------------------+----------------------------+
|   Registrant   |                  |  Google LLC  |        |      |    CA    |            |   US    |               |           select request email form at           |                            |
|                |                  |              |        |      |          |            |         |               | https://domains.markmonitor.com/whois/google.com |                            |
+----------------+------------------+              +--------+------+          +------------+         +---------------+                                                  +----------------------------+
| Administrative |                  |              |        |      |          |            |         |               |                                                  |                            |
|                |                  |              |        |      |          |            |         |               |                                                  |                            |
+----------------+------------------+              +--------+------+          +------------+         +---------------+                                                  +----------------------------+
|   Technical    |                  |              |        |      |          |            |         |               |                                                  |                            |
|                |                  |              |        |      |          |            |         |               |                                                  |                            |
+----------------+------------------+--------------+--------+------+----------+------------+---------+---------------+--------------------------------------------------+----------------------------+
```


## Todos

- [x] IP Lookup
- [x] TCPing
- [x] Packet Sender
- [x] WHOIS Lookup
- [ ] Reverse DNS Lookup
- [ ] SSL Certificate Check
- [ ] Bandwidth Test