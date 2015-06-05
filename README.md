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
The configuration files can be placed anywhere on you system and then referenced when running the program.

#### Client Config

 * connections - list of IP:PORT combinations to check
 * port - port for the client to run on

#### Server Config

 * connections - list of IP:PORT combinations of clients to call
 * port - port for the server to run on