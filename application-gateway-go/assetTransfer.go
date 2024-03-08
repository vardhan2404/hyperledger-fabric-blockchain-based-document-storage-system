/*

Copyright 2021 IBM All Rights Reserved.



SPDX-License-Identifier: Apache-2.0

*/



package main



import (

	"bytes"
	
	"crypto/sha256"

	"context"
	
	"log"
	
	"io"

	"crypto/x509"

	"encoding/json"
	
	"strconv"

	"errors"

	"net/http"

	"fmt"

	"os"

	"path"

	"time"
	
	"html/template"
	
	"github.com/gorilla/mux"
	
	"database/sql"
	
	_ "github.com/mattn/go-sqlite3"



	"github.com/hyperledger/fabric-gateway/pkg/client"

	"github.com/hyperledger/fabric-gateway/pkg/identity"

	"github.com/hyperledger/fabric-protos-go-apiv2/gateway"

	"google.golang.org/grpc"

	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc/status"	

)



const (

	mspID        = "rrMSP"

	cryptoPath   = "../organizations/peerOrganizations/rr.isfcr.com"

	certPath     = cryptoPath + "/users/User1@rr.isfcr.com/msp/signcerts/cert.pem"

	keyPath      = cryptoPath + "/users/User1@rr.isfcr.com/msp/keystore/"

	tlsCertPath  = cryptoPath + "/peers/peer0.rr.isfcr.com/tls/ca.crt"

	peerEndpoint = "localhost:7051"

	gatewayPeer  = "peer0.rr.isfcr.com"
	
	ChaincodeURL = "http://localhost:7050"

)



var now = time.Now()

var assetId = fmt.Sprintf("asset%d", now.Unix()*1e3+int64(now.Nanosecond())/1e6)

var tpl = template.Must(template.ParseFiles("templates/index.html"))

const uploadDir = "uploads"

type Asset struct {
	ID          string `json:"ID"`
	Color       string `json:"Color"`
	Size        string    `json:"Size"`
	Owner       string `json:"Owner"`
	AppraisedBy string `json:"AppraisedBy"`
}

type MyData struct {

    Name  string

    Size   string

    ID string

}



func main() {

	// The gRPC client connection should be shared by all Gateway connections to this endpoint

	clientConnection := newGrpcConnection()

	defer clientConnection.Close()



	id := newIdentity()

	sign := newSign()



	// Create a Gateway connection for a specific client identity

	gw, err := client.Connect(

		id,

		client.WithSign(sign),

		client.WithClientConnection(clientConnection),

		// Default timeouts for different gRPC calls

		client.WithEvaluateTimeout(5*time.Second),

		client.WithEndorseTimeout(15*time.Second),

		client.WithSubmitTimeout(5*time.Second),

		client.WithCommitStatusTimeout(1*time.Minute),

	)

	if err != nil {

		panic(err)

	}

	defer gw.Close()



	// Override default values for chaincode and channel name as they may differ in testing contexts.

	chaincodeName := "basic"

	if ccname := os.Getenv("CHAINCODE_NAME"); ccname != "" {

		chaincodeName = ccname

	}



	channelName := "mychannel"

	if cname := os.Getenv("CHANNEL_NAME"); cname != "" {

		channelName = cname

	}



	network := gw.GetNetwork(channelName)

	contract := network.GetContract(chaincodeName)



	//initLedger(contract)

	//getAllAssets(contract)

	//createAsset(contract)

	//readAssetByID(contract)

	//transferAssetAsync(contract)

	//exampleErrorHandling(contract)
	
	db, err := createDatabase()
	if err != nil {
		panic(fmt.Errorf("failed to create database: %w", err))
	}
	defer db.Close()
	
	router := mux.NewRouter()
    router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // Retrieve filenames from the database
        db, err := createDatabase()
        if err != nil {
            http.Error(w, fmt.Sprintf("Failed to connect to the database: %s", err), http.StatusInternalServerError)
            return
        }
        defer db.Close()

        rows, err := db.Query("SELECT filename FROM assets")
        if err != nil {
            http.Error(w, fmt.Sprintf("Failed to fetch filenames from the database: %s", err), http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        filenames := []string{}
        for rows.Next() {
            var filename string
            if err := rows.Scan(&filename); err != nil {
                http.Error(w, fmt.Sprintf("Failed to scan filename from the database: %s", err), http.StatusInternalServerError)
                return
            }
            filenames = append(filenames, filename)
        }
        if err := rows.Err(); err != nil {
            http.Error(w, fmt.Sprintf("Error while iterating over filenames: %s", err), http.StatusInternalServerError)
            return
        }

        // Pass the filenames to the template
        tpl.Execute(w, filenames)
        createAsset(contract, w, r, db)
    }).Methods("GET", "POST")
    router.HandleFunc("/read", func(w http.ResponseWriter, r *http.Request) {
        readAsset(contract, w, r)
    }).Methods("GET", "POST")
    
    fs := http.FileServer(http.Dir("static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	log.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))



}



// newGrpcConnection creates a gRPC connection to the Gateway server.

func newGrpcConnection() *grpc.ClientConn {

	certificate, err := loadCertificate(tlsCertPath)

	if err != nil {

		panic(err)

	}



	certPool := x509.NewCertPool()

	certPool.AddCert(certificate)

	transportCredentials := credentials.NewClientTLSFromCert(certPool, gatewayPeer)



	connection, err := grpc.Dial(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))

	if err != nil {

		panic(fmt.Errorf("failed to create gRPC connection: %w", err))

	}



	return connection

}



// newIdentity creates a client identity for this Gateway connection using an X.509 certificate.

func newIdentity() *identity.X509Identity {

	certificate, err := loadCertificate(certPath)

	if err != nil {

		panic(err)

	}



	id, err := identity.NewX509Identity(mspID, certificate)

	if err != nil {

		panic(err)

	}



	return id

}



func loadCertificate(filename string) (*x509.Certificate, error) {

	certificatePEM, err := os.ReadFile(filename)

	if err != nil {

		return nil, fmt.Errorf("failed to read certificate file: %w", err)

	}

	return identity.CertificateFromPEM(certificatePEM)

}



// newSign creates a function that generates a digital signature from a message digest using a private key.

func newSign() identity.Sign {

	files, err := os.ReadDir(keyPath)

	if err != nil {

		panic(fmt.Errorf("failed to read private key directory: %w", err))

	}

	privateKeyPEM, err := os.ReadFile(path.Join(keyPath, files[0].Name()))



	if err != nil {

		panic(fmt.Errorf("failed to read private key file: %w", err))

	}



	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)

	if err != nil {

		panic(err)

	}



	sign, err := identity.NewPrivateKeySign(privateKey)

	if err != nil {

		panic(err)

	}



	return sign

}

func createDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "assets.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Create the assets table if it does not exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS assets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			filename TEXT
		)`)
	if err != nil {
		return nil, fmt.Errorf("failed to create assets table: %w", err)
	}

	return db, nil
}

func insertFilenameToDB(db *sql.DB, filename string) error {
	stmt, err := db.Prepare("INSERT INTO assets (filename) VALUES (?)")
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(filename)
	if err != nil {
		return fmt.Errorf("failed to execute insert statement: %w", err)
	}

	return nil
}

func createAsset(contract *client.Contract, w http.ResponseWriter, r *http.Request, db *sql.DB) {
    if r.Method == http.MethodPost {
        err := r.ParseMultipartForm(10 << 20) // Limit the maximum file size to 10 MB
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        file, handler, err := r.FormFile("File")
        if err != nil {
            http.Error(w, "File not found in the request", http.StatusBadRequest)
            return
        }
        defer file.Close()

        // Get the filename from the file handler
        filename := handler.Filename
        fileSize := handler.Size
        fileSizeStr := strconv.FormatInt(fileSize, 10)

        // Create the uploads directory if it doesn't exist
        err = os.MkdirAll(uploadDir, 0755)
        if err != nil {
            http.Error(w, "Failed to create the uploads directory", http.StatusInternalServerError)
            return
        }

        // Create a new file in the uploads directory
        filePath := path.Join(uploadDir, filename)
        f, err := os.Create(filePath)
        if err != nil {
            http.Error(w, "Failed to create the file on the server", http.StatusInternalServerError)
            return
        }
        defer f.Close()

        // Copy the uploaded file's content to the newly created file
        _, err = io.Copy(f, file)
        if err != nil {
            http.Error(w, "Failed to save the file on the server", http.StatusInternalServerError)
            return
        }

        asset := Asset{
            ID:    fmt.Sprintf("%s-%s",r.FormValue("Owner"),filename),
            Size:  fileSizeStr,
            Owner: r.FormValue("Owner"),
        }

        // Invoke the chaincode method to create the asset
        fmt.Printf("\n--> Submit Transaction: CreateAsset, creates new asset with ID, Color, Size, Owner, and AppraisedValue arguments\n")
        
        pathtofile := fmt.Sprintf("uploads/%s", filename)
        f, _ = os.Open(pathtofile)
        if err != nil {
        	fmt.Println("Error opening the file:", err)
        	return
    	}	
    	defer f.Close()

    // Method 1: Reading the contents and storing it in a variable
    	contentBytes, err := io.ReadAll(f)
    	if err != nil {
        	fmt.Println("Error reading the file:", err)
        	return
    	}
    	content := string(contentBytes)

    // Method 2: Printing the contents directly to the console
    	_, err = io.Copy(os.Stdout, f)
    	if err != nil {
        	fmt.Println("Error copying file contents to os.Stdout:", err)
        	return
    	}
	defer f.Close()
	
        _, err = contract.SubmitTransaction("CreateAsset", asset.ID, asset.Size, asset.Owner, filename, content)
        if err != nil {
            http.Error(w, fmt.Sprintf("Failed to submit transaction: %s", err), http.StatusInternalServerError)
            return
        }
        
        if err := insertFilenameToDB(db, filename); err != nil {
		http.Error(w, fmt.Sprintf("Failed to insert filename to database: %s", err), http.StatusInternalServerError)
		return
	}
        
        fmt.Printf("*** Transaction committed successfully\n")
        fmt.Fprintf(w, "Transaction committed successfully\n")
        
        
        
        fmt.Println(pathtofile)
        f, err = os.Open(pathtofile)
  	if err != nil {
    		log.Fatal(err)
  	}
  	defer f.Close()

  	h := sha256.New()
  	if _, err := io.Copy(h, f); err != nil {
    		log.Fatal(err)
  	}

  	fmt.Printf("\n%x\n", h.Sum(nil))

        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    //tpl.Execute(w, nil)
}


func readAsset(contract *client.Contract, w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		asset := Asset{
			ID:          r.FormValue("ID"),
		}
		// Invoke the chaincode method to transfer the asset
		fmt.Printf("\n--> Evaluate Transaction: ReadAsset, function returns asset attributes\n")
		evaluateResult, err := contract.EvaluateTransaction("ReadAssetchaincode", asset.ID)
		if err != nil {
			panic(fmt.Errorf("failed to evaluate transaction: %w", err))
		}
	result := formatJSON(evaluateResult)
	//fmt.Fprintf(w, "found asset with id %s,\n deta
	fmt.Printf("*** Result:%s\n", result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//http.Redirect(w, r, "/", http.StatusSeeOther)

		return

	}

/*
    // Retrieve filenames from the database
    db, err := createDatabase()
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to connect to the database: %s", err), http.StatusInternalServerError)
        return
    }
    defer db.Close()

    rows, err := db.Query("SELECT filename FROM assets")
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to fetch filenames from the database: %s", err), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    filenames := []string{}
    for rows.Next() {
        var filename string
        if err := rows.Scan(&filename); err != nil {
            http.Error(w, fmt.Sprintf("Failed to scan filename from the database: %s", err), http.StatusInternalServerError)
            return
        }
        filenames = append(filenames, filename)
    }
    if err := rows.Err(); err != nil {
        http.Error(w, fmt.Sprintf("Error while iterating over filenames: %s", err), http.StatusInternalServerError)
        return
    }

    // Pass the filenames to the template
    tpl.Execute(w, filenames)
    */
}


// This type of transaction would typically only be run once by an application the first time it was started after its

// initial deployment. A new version of the chaincode deployed later would likely not need to run an "init" function.

func initLedger(contract *client.Contract) {

	fmt.Printf("\n--> Submit Transaction: InitLedger, function creates the initial set of assets on the ledger \n")



	_, err := contract.SubmitTransaction("InitLedger")

	if err != nil {

		panic(fmt.Errorf("failed to submit transaction: %w", err))

	}



	fmt.Printf("*** Transaction committed successfully\n")

}



// Evaluate a transaction to query ledger state.

func getAllAssets(contract *client.Contract) {

	fmt.Println("\n--> Evaluate Transaction: GetAllAssets, function returns all the current assets on the ledger")



	evaluateResult, err := contract.EvaluateTransaction("GetAllAssets")

	if err != nil {

		panic(fmt.Errorf("failed to evaluate transaction: %w", err))

	}

	result := formatJSON(evaluateResult)



	fmt.Printf("*** Result:%s\n", result)

}



/*// Submit a transaction synchronously, blocking until it has been committed to the ledger.



func createAsset(contract *client.Contract) {

	fmt.Printf("\n--> Submit Transaction: CreateAsset, creates new asset with ID, Color, Size, Owner and AppraisedValue arguments \n")



	_, err := contract.SubmitTransaction("CreateAsset", assetId, "5", "Tom")

	if err != nil {

		panic(fmt.Errorf("failed to submit transaction: %w", err))

	}



	fmt.Printf("*** Transaction committed successfully\n")

}

*/



// Evaluate a transaction by assetID to query ledger state.

func readAssetByID(contract *client.Contract) {

	fmt.Printf("\n--> Evaluate Transaction: ReadAsset, function returns asset attributes\n")



	evaluateResult, err := contract.EvaluateTransaction("ReadAsset", assetId)

	if err != nil {

		panic(fmt.Errorf("failed to evaluate transaction: %w", err))

	}

	result := formatJSON(evaluateResult)



	fmt.Printf("*** Result:%s\n", result)

}



// Submit transaction asynchronously, blocking until the transaction has been sent to the orderer, and allowing

// this thread to process the chaincode response (e.g. update a UI) without waiting for the commit notification

func transferAssetAsync(contract *client.Contract) {

	fmt.Printf("\n--> Async Submit Transaction: TransferAsset, updates existing asset owner")



	submitResult, commit, err := contract.SubmitAsync("TransferAsset", client.WithArguments(assetId, "Mark"))

	if err != nil {

		panic(fmt.Errorf("failed to submit transaction asynchronously: %w", err))

	}



	fmt.Printf("\n*** Successfully submitted transaction to transfer ownership from %s to Mark. \n", string(submitResult))

	fmt.Println("*** Waiting for transaction commit.")



	if commitStatus, err := commit.Status(); err != nil {

		panic(fmt.Errorf("failed to get commit status: %w", err))

	} else if !commitStatus.Successful {

		panic(fmt.Errorf("transaction %s failed to commit with status: %d", commitStatus.TransactionID, int32(commitStatus.Code)))

	}



	fmt.Printf("*** Transaction committed successfully\n")

}



// Submit transaction, passing in the wrong number of arguments ,expected to throw an error containing details of any error responses from the smart contract.

func exampleErrorHandling(contract *client.Contract) {

	fmt.Println("\n--> Submit Transaction: UpdateAsset asset70, asset70 does not exist and should return an error")



	_, err := contract.SubmitTransaction("UpdateAsset", "asset70", "blue", "5", "Tomoko", "300")

	if err == nil {

		panic("******** FAILED to return an error")

	}



	fmt.Println("*** Successfully caught the error:")



	switch err := err.(type) {

	case *client.EndorseError:

		fmt.Printf("Endorse error for transaction %s with gRPC status %v: %s\n", err.TransactionID, status.Code(err), err)

	case *client.SubmitError:

		fmt.Printf("Submit error for transaction %s with gRPC status %v: %s\n", err.TransactionID, status.Code(err), err)

	case *client.CommitStatusError:

		if errors.Is(err, context.DeadlineExceeded) {

			fmt.Printf("Timeout waiting for transaction %s commit status: %s", err.TransactionID, err)

		} else {

			fmt.Printf("Error obtaining commit status for transaction %s with gRPC status %v: %s\n", err.TransactionID, status.Code(err), err)

		}

	case *client.CommitError:

		fmt.Printf("Transaction %s failed to commit with status %d: %s\n", err.TransactionID, int32(err.Code), err)

	default:

		panic(fmt.Errorf("unexpected error type %T: %w", err, err))

	}



	// Any error that originates from a peer or orderer node external to the gateway will have its details

	// embedded within the gRPC status error. The following code shows how to extract that.

	statusErr := status.Convert(err)



	details := statusErr.Details()

	if len(details) > 0 {

		fmt.Println("Error Details:")



		for _, detail := range details {

			switch detail := detail.(type) {

			case *gateway.ErrorDetail:

				fmt.Printf("- address: %s, mspId: %s, message: %s\n", detail.Address, detail.MspId, detail.Message)

			}

		}

	}

}



// Format JSON data

func formatJSON(data []byte) string {

	var prettyJSON bytes.Buffer

	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {

		panic(fmt.Errorf("failed to parse JSON: %w", err))

	}

	return prettyJSON.String()

}
