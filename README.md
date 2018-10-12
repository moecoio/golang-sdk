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
