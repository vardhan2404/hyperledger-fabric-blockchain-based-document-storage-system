#!/bin/bash
#
# SPDX-License-Identifier: Apache-2.0




# default to using RR ie RR campus
ORG=${1:-RR}

# Exit on first error, print all commands.
set -e
set -o pipefail

# Where am I?
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"

ORDERER_CA=${DIR}/network-prototype/organizations/ordererOrganizations/isfcr.com/tlsca/tlsca.isfcr.com-cert.pem
PEER0_RR_CA=${DIR}/network-prototype/organizations/peerOrganizations/rr.isfcr.com/tlsca/tlsca.rr.isfcr.com-cert.pem
PEER0_EC_CA=${DIR}/network-prototype/organizations/peerOrganizations/ec.isfcr.com/tlsca/tlsca.ec.isfcr.com-cert.pem
PEER0_ORG3_CA=${DIR}/network-prototype/organizations/peerOrganizations/org3.isfcr.com/tlsca/tlsca.org3.isfcr.com-cert.pem


if [[ ${ORG,,} == "RR" || ${ORG,,} == "digibank" ]]; then

   CORE_PEER_LOCALMSPID=rrMSP
   CORE_PEER_MSPCONFIGPATH=${DIR}/network-prototype/organizations/peerOrganizations/rr.isfcr.com/users/Admin@rr.isfcr.com/msp
   #updated values according to the local host
   CORE_PEER_ADDRESS=localhost:9996
   CORE_PEER_TLS_ROOTCERT_FILE=${DIR}/network-prototype/organizations/peerOrganizations/rr.isfcr.com/tlsca/tlsca.rr.isfcr.com-cert.pem

elif [[ ${ORG,,} == "EC" || ${ORG,,} == "magnetocorp" ]]; then

   CORE_PEER_LOCALMSPID=ecMSP
   CORE_PEER_MSPCONFIGPATH=${DIR}/network-prototype/organizations/peerOrganizations/ec.isfcr.com/users/Admin@ec.isfcr.com/msp
   #updated values according to the local host
   CORE_PEER_ADDRESS=localhost:9990
   CORE_PEER_TLS_ROOTCERT_FILE=${DIR}/network-prototype/organizations/peerOrganizations/ec.isfcr.com/tlsca/tlsca.ec.isfcr.com-cert.pem

else
   echo "Unknown \"$ORG\", please choose RR/Digibank or EC/Magnetocorp"
   echo "For isfcr to get the environment variables to set upa EC shell environment run:  ./setOrgEnv.sh EC"
   echo
   echo "This can be automated to set them as well with:"
   echo
   echo 'export $(./setOrgEnv.sh EC | xargs)'
   exit 1
fi

# output the variables that need to be set
echo "CORE_PEER_TLS_ENABLED=true"
echo "ORDERER_CA=${ORDERER_CA}"
echo "PEER0_RR_CA=${PEER0_RR_CA}"
echo "PEER0_EC_CA=${PEER0_EC_CA}"
echo "PEER0_ORG3_CA=${PEER0_ORG3_CA}"

echo "CORE_PEER_MSPCONFIGPATH=${CORE_PEER_MSPCONFIGPATH}"
echo "CORE_PEER_ADDRESS=${CORE_PEER_ADDRESS}"
echo "CORE_PEER_TLS_ROOTCERT_FILE=${CORE_PEER_TLS_ROOTCERT_FILE}"

echo "CORE_PEER_LOCALMSPID=${CORE_PEER_LOCALMSPID}"
