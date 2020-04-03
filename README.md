# dpc2
DoorParty Connector v2

A small server that connects local RLOGIN clients to a remote RLOGIN server via an SSH tunnel.

This is a replacement for the original [doorparty-connector](https://github.com/echicken/doorparty-connector).  It's smaller and in some cases easier to install, but there is no compelling reason to switch if you already have the original one working.

This program was written in [Go](https://golang.org/), and you can build it from source if you want.  [Prebuilt executables](build/) are available for [Linux X86_64](build/linux_x64/), [Linux ARM6](build/linux_arm6/) (for use on a Raspberry Pi, for example), and [Win32](build/win32/).

## Installation

* Create a directory for doorparty-connector on your system, and place the following two files in it:
	* The appropriate executable for your system from one of the [build directories](build/)
	* A copy of [doorparty-connector.ini](doorparty-connector.ini)
* Edit your copy of [doorparty-connector.ini](doorparty-connector.ini)
	* Fill out `system_tag`, `ssh_username`, and `ssh_password` with the values you were given by the DoorParty administrator (you can omit the [square brackets] from your `system_tag`)
	* Edit `local_interface` and `local_port` if the default values don't suit your system
	* Leave the other values at in their default state.  They have been made customizable just in case the remote server's details change in the future.
	
### Linux

For automatic startup, a [systemd unit file](doorparty-connector.service) is available.  It assumes that `doorparty-connector` and `doorparty-connector.ini` reside within `/opt/doorparty-connector`, so edit the path as needed.

If you use some other init system, you're on your own - but feel free to share your init script / config file / whatever and I can add it to this repo.

### Windows

I dunno, put a shortcut in your Startup menu or some shit like that.
	
## Usage

`doorparty-connector` runs as a server, by default listening on TCP port `9999` of `localhost` for connections from RLOGIN clients.  You must start it up and leave it running.  A typical installation has `doorparty-connector` running on the same machine that hosts your BBS - but this isn't a requirement.

When a user on your BBS chooses to connect to DoorParty:

* Your BBS should open an RLOGIN connection to `doorparty-connector` (eg. port `9999` on `localhost`)
	* The RLOGIN "server-user-name" is the user's alias prefixed with your DoorParty "system tag", eg. `[ec]echicken`
	* The RLOGIN "client-user-name" is the user's DoorParty password
		* You can use whatever you like for this value, and even use the same value for all users if you must, but you must always send the same password for the same user every time they connect
		* The user does not need to know this value
	* `doorparty-connector` then connects the user to DoorParty's RLOGIN server via an SSH tunnel, passing the necessary user details along
		* An account is automatically created for the user on the DoorParty server if it doesn't already exist
			* This account is linked to this user on your BBS
			* The user can only access it from your BBS
			* The user does not need to know the password
			* The user may have several DoorParty accounts, one for each BBS they connect from; they cannot use the same DoorParty account from multiple BBSs

On Mystic, for example, this is menu command `IR`, with a `DATA` field like:
* `/ADDR=localhost:9999 /USER=[system_tag]@USER@ /PASS=some_password`
	* Mind that `system_tag` and `some_password` must be replaced with your own values

### Synchronet

A [script](https://raw.githubusercontent.com/echicken/dpc2/master/synchronet/doorparty.js) is available for Synchronet BBS which makes it simple to set up the connection between your BBS and DoorParty Connector.

* If you don't already have a `mods` directory, create one at the top level of your Synchronet BBS installation (alongside `ctrl`, `data`, `exec`, etc.)
* Place a copy of [doorparty.js](https://raw.githubusercontent.com/echicken/dpc2/master/synchronet/doorparty.js) in your Synchronet `mods` directory
* In SCFG, create an external program:

```
	Name: DoorParty
	Internal Code: DOORPRTY
	Command Line: ?doorparty.js
	Multiple Concurrent Users: Yes
```

* All other settings can be left at their default values.

You can optionally tell this script to bring the user directly into a particular door like so:

```
	Name: DoorParty LORD
	Internal Code: DPLORD
	Command Line: ?doorparty.js lord
	Multiple Concurrent Users: Yes
```

In the `Command Line` of the above example, `lord` is a "door code" that tells DoorParty which game to launch upon connect. [See here](http://wiki.throwbackbbs.com/doku.php?id=doorcode) for a list of valid door codes. You can make several external program entries like this if you want to direct link from your doors menu into several DoorParty games.

#### Advanced

##### DoorParty Connector running on an alternate host/interface or port

* Edit `ctrl/modopts.ini`
	* Create a section with the heading `[doorparty]` if it doesn't already exist
	* If DoorParty Connector is listening on a host or interface other than `localhost`
		* Under the `[doorparty]` section, create a key named `tunnel_host`; set this to the hostname or IP address that DoorParty Connector is listening on
	* If DoorParty Connector is listening on a port other than `9999`
		* Under the `[doorparty]` section, create a key named `tunnel_port`; set this to the port that DoorParty Connector is listening on	

##### Migration from sbbs-dp-rlogin.js

Follow these steps only if you were using the older `sbbs-dp-rlogin.js` script that came with the previous version of DoorParty Connector.

* Edit `ctrl/modopts.ini`
	* Create a section with the heading `[doorparty]` if it doesn't already exist
	* Under the `[doorparty]` section, create a key named `password`; the value should be the same as what you had in `ctrl/sbbs-dp-rlogin.ini`
		* If you didn't have `ctrl/sbbs-dp-rlogin.ini`, then `password` is whatever you were using as the first argument to `sbbs-dp-rlogin.js` (check your old `Command Line` in SCFG, or if you've already deleted it, ask the DoorParty administrator for help; if he's not sure what to do, ask him to ask me)
	* If everything's working now, you can safely delete `ctrl/sbbs-dp-rlogin.ini` (if it exists)

###### Note

The previous version of `sbbs-dp-rlogin.js` was written under the impression that Synchronet would only connect to port `513` of a remote system when creating an outbound RLOGIN connection. This is not the case in recent versions of Synchronet and quite possibly never was. However, this meant that DoorParty Connector had to bind to port `513` to accept RLOGIN connections from the local system, and it often meant that Synchronet's own RLOGIN server had to be moved to a different port. Sorry about that.

There is no longer any need for contention between DoorParty Connector and Synchronet's RLOGIN server. You can bind Synchronet's RLOGIN server to port `513` if you like, or have both servers using the same interface, etc.
