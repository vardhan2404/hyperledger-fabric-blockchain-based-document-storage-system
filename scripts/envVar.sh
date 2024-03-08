#!/bin/bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#

# This is a collection of bash functions used by different scripts

# imports
. scripts/utils.sh

export CORE_PEER_TLS_ENABLED=true
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/isfcr.com/tlsca/tlsca.isfcr.com-cert.pem
export PEER0_RR_CA=${PWD}/organizations/peerOrganizations/rr.isfcr.com/tlsca/tlsca.rr.isfcr.com-cert.pem
export PEER0_EC_CA=${PWD}/organizations/peerOrganizations/ec.isfcr.com/tlsca/tlsca.ec.isfcr.com-cert.pem
#export PEER0_ORG3_CA=${PWD}/organizations/peerOrganizations/org3.isfcr.com/tlsca/tlsca.org3.isfcr.com-cert.pem
export ORDERER_ADMIN_TLS_SIGN_CERT=${PWD}/organizations/ordererOrganizations/isfcr.com/orderers/orderer.isfcr.com/tls/server.crt
export ORDERER_ADMIN_TLS_PRIVATE_KEY=${PWD}/organizations/ordererOrganizations/isfcr.com/orderers/orderer.isfcr.com/tls/server.key

# Set environment variables for the peer org



setGlobalWithAdminKeys(){
  echo "org $1 peer $2"
  local USING_ORG=$1
  local USING_PEER=$2

  infoln "Using organization ${USING_ORG}"
  if [ $USING_ORG -eq 1 ]; then
    export CORE_PEER_LOCALMSPID="rrMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_RR_CA
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/rr.isfcr.com/users/Admin@rr.isfcr.com/msp
    if [ $USING_PEER -eq 0 ]; then
      export CORE_PEER_ADDRESS=localhost:7051
    elif [ $USING_PEER -eq 1 ]; then
      export CORE_PEER_ADDRESS=localhost:7057
    elif [ $USING_PEER -eq 2 ]; then
      export CORE_PEER_ADDRESS=localhost:7055
    fi
  elif [ $USING_ORG -eq 2 ]; then
    export CORE_PEER_LOCALMSPID="ecMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_EC_CA
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/ec.isfcr.com/users/Admin@ec.isfcr.com/msp
    if [ $USING_PEER -eq 0 ]; then
      export CORE_PEER_ADDRESS=localhost:9051
    elif [ $USING_PEER -eq 1 ]; then
      export CORE_PEER_ADDRESS=localhost:9053
    
    fi

  #if any new organization is used we need to use this
  elif [ $USING_ORG -eq 3 ]; then
    export CORE_PEER_LOCALMSPID="Org3MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG3_CA
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org3.example.com/users/Admin@org3.example.com/msp
    export CORE_PEER_ADDRESS=localhost:11051
  else
    errorln "ORG Unknown"
  fi

  if [ "$VERBOSE" == "true" ]; then
    env | grep CORE
  fi

}

# OLD set Globals
# setGlobals() {
#   local USING_ORG=$1
#   local USING_PEER=$2
#   if [ -z "$OVERRIDE_ORG" ]; then
#     USING_ORG=$1
#   else
#     USING_ORG="${OVERRIDE_ORG}"
#   fi
#   infoln "Using organization ${USING_ORG}"
#   if [ $USING_ORG -eq 1 ]; then
#     export CORE_PEER_LOCALMSPID="rrMSP"
#     export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_RR_CA
#     export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/rr.isfcr.com/tlsca/tlsca.rr.isfcr.com-cert.pem
#     export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/rr.isfcr.com/users/Admin@rr.isfcr.com/msp
#     export CORE_PEER_ADDRESS=localhost:9996
#   elif [ $USING_ORG -eq 2 ]; then
#     export CORE_PEER_LOCALMSPID="ecMSP"
#     export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_EC_CA
#     export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/ec.isfcr.com/users/Admin@ec.isfcr.com/msp
#     export CORE_PEER_ADDRESS=localhost:9990

#   elif [ $USING_ORG -eq 3 ]; then
#     export CORE_PEER_LOCALMSPID="Org3MSP"
#     export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG3_CA
#     export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org3.isfcr.com/users/Admin@org3.isfcr.com/msp
#     export CORE_PEER_ADDRESS=localhost:11051
#   else
#     errorln "ORG Unknown"
#   fi

#   if [ "$VERBOSE" == "true" ]; then
#     env | grep CORE
#   fi
# }


setGlobals() {
  echo "org $1 peer $2"
  local USING_ORG=$1
  local USING_PEER=$2
  
  infoln "Using organization ${USING_ORG}"
  if [ $USING_ORG -eq 1 ]; then
    export CORE_PEER_LOCALMSPID="rrMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_RR_CA
    if [ $USING_PEER -eq 0 ]; then
      export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/rr.isfcr.com/users/Admin@rr.isfcr.com/msp
      export CORE_PEER_ADDRESS=localhost:7051
    elif [ $USING_PEER -eq 1 ]; then
      export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/rr.isfcr.com/peers/peer1.rr.isfcr.com/msp
      export CORE_PEER_ADDRESS=localhost:7057
    elif [ $USING_PEER -eq 2 ]; then
      export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/rr.isfcr.com/peers/peer2.rr.isfcr.com/msp
      export CORE_PEER_ADDRESS=localhost:7055
    fi
  elif [ $USING_ORG -eq 2 ]; then
    export CORE_PEER_LOCALMSPID="ecMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_EC_CA
    if [ $USING_PEER -eq 0 ]; then
      export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/ec.isfcr.com/users/Admin@ec.isfcr.com/msp
      export CORE_PEER_ADDRESS=localhost:9051
    elif [ $USING_PEER -eq 1 ]; then
      export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/ec.isfcr.com/peers/peer1.ec.isfcr.com/msp
      export CORE_PEER_ADDRESS=localhost:9053
     # Add more peers here if needed 
    fi
  else
    errorln "ORG Unknown"
  fi

  if [ "$VERBOSE" == "true" ]; then
    env | grep CORE
  fi
}




# Set environment variables for use in the CLI container
setGlobalsCLI() {
  #setGlobalWithAdminKeys $1 0
  setGlobals $1 0
  warnln "We have reached here"
  local USING_ORG=""
  if [ -z "$OVERRIDE_ORG" ]; then
    USING_ORG=$1
  else
    USING_ORG="${OVERRIDE_ORG}"
  fi
  if [ $USING_ORG -eq 1 ]; then
    export CORE_PEER_ADDRESS=peer0.rr.isfcr.com:7051
  elif [ $USING_ORG -eq 2 ]; then
    export CORE_PEER_ADDRESS=peer0.ec.isfcr.com:9051
  elif [ $USING_ORG -eq 3 ]; then
    export CORE_PEER_ADDRESS=peer0.org3.isfcr.com:11051
  else
    errorln "ORG Unknown"
  fi

  warnln "We exit the setGlobal"
}

# parsePeerConnectionParameters $@
# Helper function that sets the peer connection parameters for a chaincode
# operation
parsePeerConnectionParameters() {
  PEER_CONN_PARMS=()
  PEERS=""
  while [ "$#" -gt 0 ]; do
    setGlobals $1
    PEER="peer0.$1"
    ## Set peer addresses
    if [ -z "$PEERS" ]
    then
	PEERS="$PEER"
    else
	PEERS="$PEERS $PEER"
    fi
    PEER_CONN_PARMS=("${PEER_CONN_PARMS[@]}" --peerAddresses $CORE_PEER_ADDRESS)
    ## Set path to TLS certificate
    CA=PEER0_ORG$1_CA
    TLSINFO=(--tlsRootCertFiles "${!CA}")
    PEER_CONN_PARMS=("${PEER_CONN_PARMS[@]}" "${TLSINFO[@]}")
    # shift by one to get to the next organization
    shift
  done
}

verifyResult() {
  if [ $1 -ne 0 ]; then
    fatalln "$2"
  fi
}
