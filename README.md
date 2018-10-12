Moeco Golang SDK it is code for a gateway that written on Go.

The repository contains the next folders and files:
  *  install.sh - bash script for build all dependencies;
  *  run.sh - bash script for start Moeco SDK demon;
  *  cmd - containing the file for import Moeco SDK Golang library;
  *  src - the source code of Moeco SDK Golang library;
  *  ble - all about Bluetooth;
  *  clients/prot - HTTP path (for gate registration, sending request, etc);
  *  db - SQLite path;
  *  sdk - main Moeco SDK module;
  *  typeutil - type conversion functions.

To change Masternode data you should open the file moeco-golang-sdk/cmd/moecosdk.go and find MoecoSDK:
  *  first field - Masternode address;
  *  second field - gate owner API key;
  *  third field - gate id.


To install and test Moeco Golang SDK you need PC with OS Linux (recommend use Debian or Ubuntu).
1. Unzip the archive.
2. Run install.sh.
3. Run run.sh.

Moeco Golang SDK will start working in the background mode.

In the folder where you ran install.sh will be created nohup.out log-file. If everything working well this file will contain:
 * timestamps;
 * information about connection with Masternode;
 * information about all transactions;
 * information about found devices.

You need to use own Moeco key, to do so you should change the fields in the file moeco-golang-sdk/cmd/moecosdk.go.
 * the first field (look at the screenshot below) - Masternode address;
 * the second field - gate owner API key;
 * the third field - gate id.

And run run.sh again.
