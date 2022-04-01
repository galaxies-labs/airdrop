package impt

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/go-sql-driver/mysql"
)

type Record struct {
	Address string
	Uglx    uint64
}

func ImportSnapshot(fp string, st string) error {
	if st != "cosmos" && st != "osmo" && st != "osmo-pool" {
		return fmt.Errorf("invalud type : %s", st)
	}

	db, err := sql.Open("mysql", os.Getenv("MYSQL"))

	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(time.Duration(time.Minute * 3))
	db.SetMaxIdleConns(1)
	db.SetMaxIdleConns(1)

	fmt.Println("Database connect success")

	defer db.Close()

	// check exist then delete
	rows, err := db.Query("SELECT address, uglx, type FROM airdrop Where type = ?", st)
	if err != nil {
		return err
	}

	defer rows.Close()

	existCnt := 0
	for rows.Next() {
		existCnt++
	}

	fmt.Printf("Exist rows count : %v\n", existCnt)

	result, err := db.Exec("DELETE FROM airdrop WHERE type = ?", st)
	if err != nil {
		return err
	}

	nRow, err := result.RowsAffected()
	if err != nil {
		return err
	}

	fmt.Printf("Deleted rows count : %v\n", nRow)

	// read file then marshal
	records := []Record{}

	file, err := os.Open(fp)
	if err != nil {
		return err
	}
	defer file.Close()

	bz, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bz, &records)
	if err != nil {
		return err
	}

	// convert bech32 to pub address
	totalUglx := uint64(0)
	for i, record := range records {
		//bech32
		bz, err := sdk.GetFromBech32(record.Address, strings.Split(st, "-")[0])
		if err != nil {
			fmt.Println("error")
			return err
		}
		records[i].Address = hex.EncodeToString(bz)

		totalUglx += record.Uglx
	}
	fmt.Printf("total records : %v, total uglx : %v\n", len(records), totalUglx)

	// divide array if to many rows
	chunkSize := 20_000
	var divided [][]Record
	for i := 0; i < len(records); i += chunkSize {
		end := i + chunkSize

		if end > len(records) {
			end = len(records)
		}

		divided = append(divided, records[i:end])

	}

	// insert to sql
	totalRows := int64(0)
	for _, records := range divided {
		valueStrings := []string{}
		valueArgs := []interface{}{}
		for _, record := range records {
			valueStrings = append(valueStrings, "(?, ?, ?)")
			valueArgs = append(valueArgs, record.Address)
			valueArgs = append(valueArgs, record.Uglx)
			valueArgs = append(valueArgs, st)
		}

		stmt := fmt.Sprintf("INSERT INTO airdrop (address, uglx, type) VALUES %s", strings.Join(valueStrings, ","))

		tx, err := db.Begin()
		if err != nil {
			return err
		}

		result, err = tx.Exec(stmt, valueArgs...)
		if err != nil {
			tx.Rollback()
			return err
		}

		if err = tx.Commit(); err != nil {
			return err
		}

		ir, err := result.RowsAffected()
		if err != nil {
			return err
		}
		totalRows += ir

	}

	fmt.Printf("inserted records : %v\n", totalRows)

	return nil
}
