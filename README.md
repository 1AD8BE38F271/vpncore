waiting...


# DESC
This project is originally a fold of [water](https://github.com/songgao/water), that is a tun/tap util implemented with Go. I extend it for supporting Mac OS X, and DNS/routing feature. Now I recreate a project for better portability, implements basic toolset for VPN Server/Client, including but not limited to tun/tap, DNS, routing management, transport.


# PROJECT
* *vpncore/tuntap* is a native Go library for [TUN/TAP](http://en.wikipedia.org/wiki/TUN/TAP) interfaces.

* *vpncore/tcpip* has some useful functions to interpret MAC frame headers and IP packet headers. It also contains some constants such as protocol numbers and ethernet frame types.

# TODO
* kcp/obfs protocol


# tun/tap
## Mac OS X tips for Tun/Tap
### < Mac OS X 10.10

Since the Mac OS kernel has no tun/tap device(but FreeBSD have), you need to use the third-party opensource Tun/Tap device
kernel extension(kext), for example, [tuntaposx](http://tuntaposx.sourceforge.net).

```
sudo kextload /opt/local/Library/Extensions/tap.kext
sudo kextload /opt/local/Library/Extensions/tun.kext
```

If you had used Tunnelblick/Viscosity before, you may need to check if some kinds of Tun/Tap kext have already loaded,
if so, close them before loading new Tun/Tap kext:

```
kextstat |grep tun
kextstat |grep tap

sudo kextunload /Library/Extensions/tun.kext
sudo kextunload /Library/Extensions/tap.kext
```

### >= Mac OS X 10.10

Mac 10.10 or later OS version deny from loading unsigned kext, so:

1. [Disable it](https://github.com/sergeybratus/netfluke/blob/master/howto-disable-kext-signing.txt)

2. [Use the signed kext](https://github.com/sergeybratus/netfluke/blob/master/howto-load-signed-tunnelblick-drivers.txt).You can extract tap-signed.kext/tun-signed.kext file from [Tunnelblick.app](https://tunnelblick.net) App Bundle, then loads it using command line:

    ```
    sudo chmod -R 644  /Applications/Tunnelblick.app/Contents/Resources/tap-signed.kext
    sudo chown -R root:wheel  /Applications/Tunnelblick.app/Contents/Resources/tap-signed.kext
    
    kextstat |grep tap 
    sudo kextutil -d /Applications/Tunnelblick.app/Contents/Resources/tap-signed.kext -b net.tunnelblick.tap
    sudo kextunload -b  net.tunnelblick.tap
    ```
3. Use [utun](https://github.com/songgao/water/issues/3) to implements Tun without loading third-party tun/tap kext. Reference to implementation of [vpn-ws/tun.c](https://github.com/unbit/vpn-ws/blob/master/src/tuntap.c). 



