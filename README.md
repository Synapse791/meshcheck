# Meshcheck

Meshcheck enables you to check connectivity throughout your infrastructure and provide information about te connections in an easy to process JSON format. Just set up the client on the machines you want to check connectivity between and use the server to gather all the connection data into a JSON response.

----------

### Client
The client sits on the machines you wish to check. You provide a config file containing a list of IP:PORT combinations to check and the client takes care of the rest!

To run in client mode, use the following command:
```
meshcheck client PATH_TO_CONFIG_DIRECTORY
```

### Server
The server connects to a list of clients and requests their connection status. This will trigger the clients to check all their connections and return the results to the server. At this point, the server will format all of the results into a JSON response.

To run in client mode, use the following command:
```
meshcheck server PATH_TO_CONFIG_DIRECTORY
```

### Configuration
The configuration files can be placed anywhere on you system and then referenced when running the program. An example command for Linux might look like this:
```
meshcheck client /etc/meshcheck-client/
```
In this example, the config files should be placed in the /etc/meshcheck-client/ directory.


#### Client Config

 * connections - list of IP:PORT combinations to check
 * port - port for the client to run on. If not set, port defaults to 6600

**connections example**
```
$ cat /etc/meshcheck/client/connections
10.100.0.10:80
10.100.0.10:22
10.100.0.20:5000
```

**port example**
```
$ cat /etc/meshcheck/client/port
9000
```

#### Server Config

 * connections - list of IP:PORT combinations of clients to call
 * port - port for the server to run on. If not set, port defaults to 6800

**connections example**
```
$ cat /etc/meshcheck/server/connections
10.100.0.10:6600
10.100.0.20:6600
10.100.0.30:6600
```

**port example**
```
$ cat /etc/meshcheck/server/port
5000
```