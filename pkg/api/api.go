package api

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/rs/cors"
)

type Type string

const (
	OsmoPool = Type("osmo-pool")
	Cosmos   = Type("cosmos")
	Osmo     = Type("osmo")
)

type Record struct {
	Address string
	Uglx    float64
	Type    Type
}

type Response struct {
	StatusCode int
	Message    string
	Data       interface{}
}

func makeError(response *Response, message string) {
	response.StatusCode = 500
	response.Message = message
}

func AddressClaimableHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{}
	defer func() {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(response.StatusCode)
		json.NewEncoder(w).Encode(response)
	}()
	vars := mux.Vars(r)
	address := vars["address"]

	if len(address) == 0 {
		makeError(&response, "Please enter a valid address.")
		return
	}

	var prefix string

	if strings.HasPrefix(address, "cosmos") {
		prefix = "cosmos"
	} else if strings.HasPrefix(address, "osmo") {
		prefix = "osmo"
	} else {
		makeError(&response, "Only Cosmos and Osmosis address can be searched.")
		return
	}

	bz, err := sdk.GetFromBech32(address, prefix)
	if err != nil {
		response.StatusCode = 500
		response.Message = err.Error()
		return
	}

	pubAddr := hex.EncodeToString(bz)

	fmt.Printf("%s - (%s) has been searched\n", address, pubAddr)

	db, err := sql.Open("mysql", os.Getenv("MYSQL"))
	if err != nil {
		makeError(&response, fmt.Sprintf("Database connect error, err : %s", err.Error()))
		return
	}

	db.SetConnMaxLifetime(time.Duration(time.Minute * 3))
	db.SetMaxIdleConns(10)
	db.SetMaxIdleConns(10)
	defer db.Close()

	rows, err := db.Query("SELECT address, type, uglx FROM airdrop WHERE address = ? ", pubAddr)
	if err != nil {
		makeError(&response, fmt.Sprintf("Database query error, err : %s", err.Error()))
		return
	}
	defer rows.Close()

	records := []Record{}

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var record Record
		if err := rows.Scan(&record.Address, &record.Type, &record.Uglx); err != nil {
			makeError(&response, fmt.Sprintf("Database select error, err : %s", err.Error()))
			return

		}
		records = append(records, record)
	}

	response.StatusCode = 200
	response.Data = records

}

func NewApiServer() error {
	if err := godotenv.Load(".env"); err != nil {
		return err
	}

	fmt.Printf("Starting server at port : %s\n", os.Getenv("PORT"))

	r := mux.NewRouter()

	r.HandleFunc("/addresses/{address}/claimable", AddressClaimableHandler)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://galaxychain.zone"},
		AllowCredentials: true,
	})

	http.Handle("/", c.Handler(r))

	if err := http.ListenAndServe(":"+os.Getenv("PORT"), r); err != nil {
		return err
	}
	return nil
}
