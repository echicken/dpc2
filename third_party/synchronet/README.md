# Connecting to DoorParty from Synchronet

A [script](https://raw.githubusercontent.com/echicken/dpc2/master/third_party/synchronet/doorparty.js) is available for Synchronet BBS which makes it simple to set up the connection between your BBS and DoorParty Connector.

* If you don't already have a `mods` directory, create one at the top level of your Synchronet BBS installation (alongside `ctrl`, `data`, `exec`, etc.)
* Place a copy of [doorparty.js](https://raw.githubusercontent.com/echicken/dpc2/master/third_party/synchronet/doorparty.js) in your Synchronet `mods` directory
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

## Advanced

### DoorParty Connector running on an alternate host/interface or port

* Edit `ctrl/modopts.ini`
	* Create a section with the heading `[doorparty]` if it doesn't already exist
	* If DoorParty Connector is listening on a host or interface other than `localhost`
		* Under the `[doorparty]` section, create a key named `tunnel_host`; set this to the hostname or IP address that DoorParty Connector is listening on
	* If DoorParty Connector is listening on a port other than `9999`
		* Under the `[doorparty]` section, create a key named `tunnel_port`; set this to the port that DoorParty Connector is listening on	

### Migration from sbbs-dp-rlogin.js

Follow these steps only if you were using the older `sbbs-dp-rlogin.js` script that came with the previous version of DoorParty Connector.

* Edit `ctrl/modopts.ini`
	* Create a section with the heading `[doorparty]` if it doesn't already exist
	* Under the `[doorparty]` section, create a key named `password`; the value should be the same as what you had in `ctrl/sbbs-dp-rlogin.ini`
		* If you didn't have `ctrl/sbbs-dp-rlogin.ini`, then `password` is whatever you were using as the first argument to `sbbs-dp-rlogin.js` (check your old `Command Line` in SCFG, or if you've already deleted it, ask the DoorParty administrator for help; if he's not sure what to do, ask him to ask me)
	* If everything's working now, you can safely delete `ctrl/sbbs-dp-rlogin.ini` (if it exists)

#### Note

The previous version of `sbbs-dp-rlogin.js` was written under the impression that Synchronet would only connect to port `513` of a remote system when creating an outbound RLOGIN connection. This is not the case in recent versions of Synchronet and quite possibly never was. However, this meant that DoorParty Connector had to bind to port `513` to accept RLOGIN connections from the local system, and it often meant that Synchronet's own RLOGIN server had to be moved to a different port. Sorry about that.

There is no longer any need for contention between DoorParty Connector and Synchronet's RLOGIN server. You can bind Synchronet's RLOGIN server to port `513` if you like, or have both servers using the same interface, etc.
