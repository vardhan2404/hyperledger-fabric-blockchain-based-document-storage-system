./startup.sh down
./startup.sh up -s couchdb
./startup.sh createChannel
./startup.sh deployCC -ccn basic -ccp ../asset-transfer-basic/chaincode-go -ccl go
