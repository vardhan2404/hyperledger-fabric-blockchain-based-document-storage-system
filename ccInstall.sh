#!/bin/bash

. scripts/utils.sh



function installCC(){
    export PATH=${PWD}/../bin:$PATH
    export FABRIC_CFG_PATH=$PWD/../config/
    peer version

    #set the address accordingly
    peer lifecycle chaincode package basic.tar.gz --path ./chaincode_go/ --lang golang --label basic_1.0
     
    #peer0.rr
    infoln "peer0 rr installing..."
    export CORE_PEER_TLS_ENABLED=true
    export CORE_PEER_LOCALMSPID="rrMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/rr.isfcr.com/peers/peer0.rr.isfcr.com/tls/ca.crt
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/rr.isfcr.com/users/Admin@rr.isfcr.com/msp
    export CORE_PEER_ADDRESS=localhost:7051

    peer lifecycle chaincode install basic.tar.gz
    infoln "peer0 rr installed"

    #peer1.rr
    infoln "peer1 rr installing..."
    export CORE_PEER_TLS_ENABLED=true
    export CORE_PEER_LOCALMSPID="rrMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/rr.isfcr.com/peers/peer1.rr.isfcr.com/tls/ca.crt
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/rr.isfcr.com/users/Admin@rr.isfcr.com/msp
    export CORE_PEER_ADDRESS=localhost:7057

    peer lifecycle chaincode install basic.tar.gz
    infoln "peer1 rr installed"


    #peer2.rr
    infoln "peer2 rr installing..."
    export CORE_PEER_TLS_ENABLED=true
    export CORE_PEER_LOCALMSPID="rrMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/rr.isfcr.com/peers/peer2.rr.isfcr.com/tls/ca.crt
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/rr.isfcr.com/users/Admin@rr.isfcr.com/msp
    export CORE_PEER_ADDRESS=localhost:7055

    peer lifecycle chaincode install basic.tar.gz
    infoln "peer2 rr installed"

    
    #peer0.ec
    infoln "peer0 ec installing..."
    export CORE_PEER_TLS_ENABLED=true
    export CORE_PEER_LOCALMSPID="ecMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/ec.isfcr.com/peers/peer0.ec.isfcr.com/tls/ca.crt
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/ec.isfcr.com/users/Admin@ec.isfcr.com/msp
    export CORE_PEER_ADDRESS=localhost:9051 

    peer lifecycle chaincode install basic.tar.gz
    infoln "peer0 ec installed"
    

    #peer1.ec
    infoln "peer1 ec installing..."
    export CORE_PEER_TLS_ENABLED=true
    export CORE_PEER_LOCALMSPID="ecMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/ec.isfcr.com/peers/peer1.ec.isfcr.com/tls/ca.crt
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/ec.isfcr.com/users/Admin@ec.isfcr.com/msp
    export CORE_PEER_ADDRESS=localhost:9053

    peer lifecycle chaincode install basic.tar.gz
    infoln "peer1 ec installed"


    infoln "Package installed ID present"
    peer lifecycle chaincode queryinstalled

}

function queryInstalled(){
    peer lifecycle chaincode queryinstalled
}

function approveCC(){
    #rr approve
    infoln "rr approve"
    export CORE_PEER_TLS_ENABLED=true
    export CORE_PEER_LOCALMSPID="rrMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/rr.isfcr.com/peers/peer0.rr.isfcr.com/tls/ca.crt
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/rr.isfcr.com/users/Admin@rr.isfcr.com/msp
    export CORE_PEER_ADDRESS=localhost:7051

    peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.isfcr.com --channelID mychannel --name basic --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile "${PWD}/organizations/ordererOrganizations/isfcr.com/orderers/orderer.isfcr.com/msp/tlscacerts/tlsca.isfcr.com-cert.pem"
    infoln "rr approved"

    #ec approve
    infoln "ec approve"
    export CORE_PEER_TLS_ENABLED=true
    export CORE_PEER_LOCALMSPID="ecMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/ec.isfcr.com/peers/peer0.ec.isfcr.com/tls/ca.crt
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/ec.isfcr.com/users/Admin@ec.isfcr.com/msp
    export CORE_PEER_ADDRESS=localhost:9051 

    peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.isfcr.com --channelID mychannel --name basic --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile "${PWD}/organizations/ordererOrganizations/isfcr.com/orderers/orderer.isfcr.com/msp/tlscacerts/tlsca.isfcr.com-cert.pem"
    infoln "ec approved"

}

function commitCC(){
    peer lifecycle chaincode checkcommitreadiness --channelID mychannel --name basic --version 1.0 --sequence 1 --tls --cafile "${PWD}/organizations/ordererOrganizations/isfcr.com/orderers/orderer.isfcr.com/msp/tlscacerts/tlsca.isfcr.com-cert.pem" --output json

    infoln "Commiting"

    peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.isfcr.com --channelID mychannel --name basic --version 1.0 --sequence 1 --tls --cafile "${PWD}/organizations/ordererOrganizations/isfcr.com/orderers/orderer.isfcr.com/msp/tlscacerts/tlsca.isfcr.com-cert.pem" --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/rr.isfcr.com/peers/peer0.rr.isfcr.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/ec.isfcr.com/peers/peer0.ec.isfcr.com/tls/ca.crt"
    infoln "Commited"


    peer lifecycle chaincode querycommitted --channelID mychannel --name basic --cafile "${PWD}/organizations/ordererOrganizations/isfcr.com/orderers/orderer.isfcr.com/msp/tlscacerts/tlsca.isfcr.com-cert.pem"
}


function invokeCC(){
    peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.isfcr.com --tls --cafile "${PWD}/organizations/ordererOrganizations/isfcr.com/orderers/orderer.isfcr.com/msp/tlscacerts/tlsca.isfcr.com-cert.pem" -C mychannel -n basic --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/rr.isfcr.com/peers/peer0.rr.isfcr.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/ec.isfcr.com/peers/peer0.ec.isfcr.com/tls/ca.crt" -c '{"function":"InitLedger","Args":[]}'
}



if [[ $# -lt 1 ]] ; then
    errorln "Invalid Flag"
    exit 0
else
    MODE=$1
    infoln "Getting Called 1"
    shift
    infoln "Getting Called 2"
fi

# Parsing the the falg
while [[ $# -ge 1 ]] ; do
    infoln "Getting Called 3"
    case $key in
    -pid )
        export CC_PACKAGE_ID="$2"
        shift
        ;;
    * )
        errorln "Unknown flag: $key"
        exit 1
        ;;
    esac
    shift 
done

if [ "$MODE" == "full" ] ; then 
    installCC
elif [ "$MODE" == "query" ] ; then 
    queryInstalled
elif [ "$MODE" == "package" ] ; then 
    approveCC
    commitCC
    invokeCC
else
    errorln "Error"
    exit 1
fi

