package main

import (
    // "github.com/Sirupsen/logrus"
    "github.com/fsouza/go-dockerclient"
    // "bytes"
    "os"
    "net"
    "io"
    "io/ioutil"
    log "github.com/Sirupsen/logrus"
    "errors"
    "crypto/x509"
    "crypto/tls"
    "path/filepath"
    "net/http"
    "net/url"
    // "time"
    // "strings" // words := strings.Fields(someString)
)

func getDockerClient() (*docker.Client, error) {
    endpoint := getDockerEndpoint()
    if !isDockerMachineInUse() { // using Docker on Linux, Docker for Mac, Docker for Windows
        c, err := docker.NewClient(endpoint)
        if err != nil {
            return nil, errors.New("Could not reach Docker host (" +
                endpoint+  "): " + err.Error())
        } else {
            return c, nil
        }
    } else { // using remote Docker, or Docker Toolbox
        // inspired from https://github.com/fsouza/go-dockerclient/issues/166

        // Still, attachToContainer does not seem to work on Mac with Docker Toolbox.
        // Because of this issue ? https://github.com/fsouza/go-dockerclient/issues/126
        roots := x509.NewCertPool()
        dockerCertPath := os.Getenv("DOCKER_CERT_PATH")
        pemData, err := ioutil.ReadFile(filepath.Join(dockerCertPath, "ca.pem"))
        if err != nil {
            return nil, errors.New("Error while loading ca.pem: " + err.Error())
        }

        //add to pool
        roots.AppendCertsFromPEM(pemData)

       //create certificate
        crt, err := tls.LoadX509KeyPair(
            filepath.Join(dockerCertPath, "cert.pem"),
            filepath.Join(dockerCertPath, "key.pem"))
        if err != nil {
            return nil, errors.New("Error while loading cert.pem and key.pem: " + err.Error())
        }

        //creates the new tls configuration using both the authority and certificate
        conf := &tls.Config{
            RootCAs:      roots,
            Certificates: []tls.Certificate{crt},
        }

        //create our own transport
        tr := &http.Transport{
            TLSClientConfig: conf,
        }

        host, err := url.Parse(endpoint)
        if err != nil {
            return nil, errors.New("Error while parsing endpoint (" + endpoint + "): " +
                err.Error())
        }
        //change tcp to https connection
        host.Scheme = "https"

        //create fsouza/go-dockerclient
        c, err := docker.NewClient(host.String())
        if err != nil {
            return nil, errors.New("Could not reach Docker host (" +
                endpoint+  "): " + err.Error())
        }

        //create a new http client and set on our dockerclient
        c.HTTPClient = &http.Client{Transport: tr}
        return c, nil
    }
}

func isDockerMachineInUse() bool {
    endpoint := os.Getenv("DOCKER_HOST")
    if endpoint == "" {
        return false
    } else {
        return true
    }
}

func getDockerEndpoint() string {
    if isDockerMachineInUse() {
        return os.Getenv("DOCKER_HOST")
    } else {
        return "unix:///var/run/docker.sock"
    }
}

// Bridge (unix socket of current machine) to (port of current machine),
// accepting connection on interface used to connected to the host.
// (for security purposes, it does restrains access from other interfaces).
// Return the IP of the current machine on which the port is bind, or error.
func bridgeUnixSocketToPort(unixSocketName string, port string) (string, error) {
    ip, err := getMyIP()
    // ip = "192.168.64.1"
    if err != nil {
        log.Debug("get IP failed: %v", err)
        return "", err
    } else {
        log.Debug("got ip:", string(ip))

        go func() {
            listener, err := net.Listen("tcp", ip + ":" + port)
            if err != nil {
                log.Error("Failed to setup listener: %v", err)
            } else {
                for {
                    conn, err := listener.Accept()
                    if err != nil {
                        log.Error("ERROR: failed to accept listener: %v", err)
                    }
                    go forward(conn, ip, port, unixSocketName)
                }
            }
        }()

        return ip, nil
    }
}

func getMyIP() (string, error) {
    // log.Debug("getMyIP")
    ifaces, err := net.Interfaces()
    if err != nil {
        return "", err
    }
    // log.Debug("no error to get interfaces: ", len(ifaces))

    // try to find a secure IP on a private interface created by docker toolbox or docker for mac
    for _, i := range ifaces {
        // log.Debug("processing interface ", i.Name)
        addrs, err := i.Addrs()
        // addrs, err := net.InterfaceAddrs()
        usingDockerMachine := isDockerMachineInUse()
        usingDockerForMac := !usingDockerMachine
        if err == nil && ((i.Name == "bridge100" && usingDockerForMac) || (i.Name == "docker0" && usingDockerMachine))  {
            // log.Debug("no error to get Addrs from ", i.Name)
            for _, addr := range addrs {
                var ip net.IP
                switch v := addr.(type) {
                case *net.IPNet:
                        ip = v.IP
                case *net.IPAddr:
                        ip = v.IP
                // default:
                        // log.Debug("no case for ", v)
                }
                // process IP address
                // log.Debug("ip on %s: %s", i.Name, ip.String())
                // return ip.String(), nil
                if finalIP := ip.To4(); finalIP != nil {
                    log.Debug("IPv4 from ", i.Name, ": ", finalIP.String())
                    return finalIP.String(), nil
                }
            }
        }
        // else if err == nil {
        //     log.Warning("Could not find IP on a proper network interface. Opening connection to any IP (security issue)")
        //     return "", nil
        // }
    }

    // in case no secure IP has been found, try to use local network IP
    for _, i := range ifaces {
        // log.Debug("processing interface ", i.Name)
        addrs, err := i.Addrs()
        if err == nil {
            // log.Debug("no error to get Addrs from ", i.Name)
            for _, addr := range addrs {
                if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && i.Name != "lo" && i.Name != "lo0" {
                    var ip net.IP
                    switch v := addr.(type) {
                    case *net.IPNet:
                            ip = v.IP
                    case *net.IPAddr:
                            ip = v.IP
                    }
                    if finalIP := ip.To4(); finalIP != nil {
                        log.Debug("IPv4 from ", i.Name, ": ", finalIP.String())
                        return finalIP.String(), nil
                    }
                    // else if finalIP := ip.To16(); finalIP != nil { // TODO: accept IPv6 as well ?
                    //     log.Debug("IPv6 from ", i.Name, ": ", finalIP.String())
                    //     return finalIP.String(), nil
                    // }
                }
                // process IP address
                // log.Debug("ip on %s: %s", i.Name, ip.String())

            }
        }
        // else if err == nil {
        //     log.Warning("Could not find IP on a proper network interface. Opening connection to any IP (security issue)")
        //     return "", nil
        // }
    }

    // log.Warning("Could not find any network interface. Opening connection to any IP (security issue)")
    return "", errors.New("Could not find any IP address")
}

// inspired from http://blog.evilissimo.net/simple-port-fowarder-in-golang
func forward(conn net.Conn, ip string, port string, unixSocketName string) {
    client, err := net.Dial("unix", unixSocketName)
    if err != nil {
        log.Error("forward: Dial failed: %v", err)
    }
    // log.Printf("forward: Connected to localhost %v\n", conn)
    go func() {
        defer client.Close()
        defer conn.Close()
        io.Copy(client, conn)
    }()
    go func() {
        defer client.Close()
        defer conn.Close()
        io.Copy(conn, client)
    }()
}
