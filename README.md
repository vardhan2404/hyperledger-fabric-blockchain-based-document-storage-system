# DSM ISFCR NET (Distributed Filestorage System)

## Documentation

### Port Numbers
<br> Ex - Exposed ports on the host
<br> In - Internal usage ports

- Orderer
  - 7050
  - 7053
  - 9443


- RR
  - peer0
    - 7051 (Ex)
    - 7052 (In)
    - 9444 (Ex)
  
  - peer1
    - 7057 (Ex)
    - 7054 (In)
    - 9445 (Ex)
  
  - peer2
    - 7055 (Ex)
    - 7056 (In)
    - 9446 (Ex)

- EC
  - peer0
    - 9051 (Ex)
    - 9052 (In)
    - 9447 (Ex)
  
  - peer1
    - 9053 (Ex)
    - 9054 (In)
    - 9448 (Ex)

so to change the org1 to peer0.rr and org2 to peer0.ec we need to change envVar.sh file



### To get peer onto the terminal and be able to control using peer command 
run 
- export PATH=${PWD}/../bin:$PATH
- export FABRIC_CFG_PATH=$PWD/../config/

then you should be able to see
- peer version




The chaincode depemdendies are installed in the chaincode go.mod file so any other dependencies that needed to be added must be put there

- The chaincode needs to be installed on every peer that will endorse a transaction.



## Needed functionalities
- [x] need to write register peer function on register enroll.sh inside fabric ca file



## To install chaincode we need to set the global variables on the terminal

### Set package
```
export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/
peer lifecycle chaincode package basic.tar.gz --path ../asset-transfer-basic/chaincode-go/ --lang golang --label basic_1.0
```

### Installing the chaincode 
```
peer lifecycle chaincode install basic.tar.gz
```

### Set Env
- peer0.rr
```
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="rrMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/rr.isfcr.com/peers/peer0.rr.isfcr.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/rr.isfcr.com/users/Admin@rr.isfcr.com/msp
export CORE_PEER_ADDRESS=localhost:7051
```


- peer1.rr
```
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="rrMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/rr.isfcr.com/peers/peer1.rr.isfcr.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/rr.isfcr.com/users/Admin@rr.isfcr.com/msp
export CORE_PEER_ADDRESS=localhost:7057
```


- peer2.rr
~~~
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="rrMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/rr.isfcr.com/peers/peer2.rr.isfcr.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/rr.isfcr.com/users/Admin@rr.isfcr.com/msp
export CORE_PEER_ADDRESS=localhost:7055
~~~


- peer0.ec
```
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="ecMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/ec.isfcr.com/peers/peer0.ec.isfcr.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/ec.isfcr.com/users/Admin@ec.isfcr.com/msp
export CORE_PEER_ADDRESS=localhost:9051
```


- peer1.ec
```
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="ecMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/ec.isfcr.com/peers/peer1.ec.isfcr.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/ec.isfcr.com/users/Admin@ec.isfcr.com/msp
export CORE_PEER_ADDRESS=localhost:9053
```


### Verify install on the peer
~~~
peer lifecycle chaincode queryinstalled
~~~

### Set package ID on export 

use the package ID that you get on the above query to set it in this

```
export CC_PACKAGE_ID=basic_1.0:d9f386982d89197ddffd72e46473c3b056b3068564c5af16c8b60b51a2cff21b
```


### Approve the values for your organisation

```
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.isfcr.com --channelID mychannel --name basic --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile "${PWD}/organizations/ordererOrganizations/isfcr.com/orderers/orderer.isfcr.com/msp/tlscacerts/tlsca.isfcr.com-cert.pem"
```


### Commiting the chaincode into the channel

First we need to check the readiness of the organisation before commiting the code onto the channel

- We need to get true for all orgs before commiting
```
peer lifecycle chaincode checkcommitreadiness --channelID mychannel --name basic --version 1.0 --sequence 1 --tls --cafile "${PWD}/organizations/ordererOrganizations/isfcr.com/orderers/orderer.isfcr.com/msp/tlscacerts/tlsca.isfcr.com-cert.pem" --output json
```

- Then commit the code into the channel 
```
peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.isfcr.com --channelID mychannel --name basic --version 1.0 --sequence 1 --tls --cafile "${PWD}/organizations/ordererOrganizations/isfcr.com/orderers/orderer.isfcr.com/msp/tlscacerts/tlsca.isfcr.com-cert.pem" --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/rr.isfcr.com/peers/peer0.rr.isfcr.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/ec.isfcr.com/peers/peer0.ec.isfcr.com/tls/ca.crt"
```


- You need to get valid status code when the commit is done successfully


### Query the commited chaincodes in the channel

- While using this we need to set the variables in the local terminal and then set it 


```
peer lifecycle chaincode querycommitted --channelID mychannel --name basic --cafile "${PWD}/organizations/ordererOrganizations/isfcr.com/orderers/orderer.isfcr.com/msp/tlscacerts/tlsca.isfcr.com-cert.pem"
```



### Invoking the chaincode 

So once the chaincode is installed we need to check if the chaincode is installed by querying it 

```
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.isfcr.com --tls --cafile "${PWD}/organizations/ordererOrganizations/isfcr.com/orderers/orderer.isfcr.com/msp/tlscacerts/tlsca.isfcr.com-cert.pem" -C mychannel -n basic --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/rr.isfcr.com/peers/peer0.rr.isfcr.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/ec.isfcr.com/peers/peer0.ec.isfcr.com/tls/ca.crt" -c '{"function":"InitLedger","Args":[]}'
```