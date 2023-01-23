package database

import (
	"chimera/network"
	"chimera/utility/configuration"
	"chimera/utility/logging"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var database *sql.DB

func Initialize() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s@tcp(127.0.0.1:3306)/chimera", configuration.GetConfiguration().Connection.DBLogin))

	if err != nil {
		logging.Fatal("Failed to contact database! (%s)", err.Error())
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(50)
	database = db
	logging.Info("Database Service", "Connected to Database")
}

func Query(data string, args ...any) (*sql.Rows, error) {
	return database.Query(data, args...)
}

func SetLastLoginDate(uin int) error {
	row, err := Query("UPDATE userdetails SET LastLogin= ? WHERE UIN= ?", time.Now().UnixNano(), uin)
	if err != nil {
		logging.Error("Database/GetUserData", "Failed to get userdata: %s", err.Error())
		return err
	}
	row.Close()
	return err
}

func GetUserMetaDetailsDataByUIN(uin int) (network.Meta, error) {

	var meta network.Meta

	row, err := Query("SELECT * from meta WHERE UIN= ?", uin)

	if err != nil {
		logging.Error("Database/GetMetaDetailsData", "Failed to get userdata: %s", err.Error())
		return meta, err
	}

	row.Next()
	row.Scan(&meta.UIN, &meta.UsageFlag, &meta.AccountFlag)
	row.Close()

	return meta, err
}

func GetUserDetailsDataByUIN(uin int) (network.User, error) {

	var user network.User

	row, err := Query("SELECT * from userdetails WHERE UIN= ?", uin)

	if err != nil {
		logging.Error("Database/GetUserData", "Failed to get userdata: %s", err.Error())
		return user, err
	}

	row.Next()
	row.Scan(&user.UIN, &user.AvatarBlob, &user.AvatarImageType, &user.StatusCode, &user.StatusMessage, &user.LastLogin, &user.SignupDate)
	row.Close()

	return user, err
}

func GetAccountDataByEmail(email string) (network.Account, error) {

	var acc network.Account

	row, err := Query("SELECT * from accounts WHERE Mail= ?", email)

	if err != nil {
		logging.Error("Database/GetUserData", "Failed to get userdata: %s", err.Error())
		return acc, err
	}

	row.Next()
	row.Scan(&acc.UIN, &acc.DisplayName, &acc.Mail, &acc.Password)
	row.Close()

	return acc, err
}

func GetAccountDataByDisplayName(displayName string) (network.Account, error) {

	var acc network.Account

	row, err := Query("SELECT * from accounts WHERE DisplayName= ?", displayName)

	if err != nil {
		logging.Error("Database/GetUserData", "Failed to get userdata: %s", err.Error())
		return acc, err
	}

	row.Next()
	row.Scan(&acc.UIN, &acc.DisplayName, &acc.Mail, &acc.Password)
	row.Close()

	return acc, err
}

func GetAccountDataByUIN(uin int) (network.Account, error) {

	var acc network.Account

	row, err := Query("SELECT * from accounts WHERE UIN= ?", uin)

	if err != nil {
		logging.Error("Database/GetUserData", "Failed to get userdata: %s", err.Error())
		return acc, err
	}

	row.Next()
	row.Scan(&acc.UIN, &acc.DisplayName, &acc.Mail, &acc.Password)
	row.Close()

	return acc, err
}

func GetEncryptionAttributes(uin int) (network.EncryptionAttributes, error) {

	var enc network.EncryptionAttributes

	row, err := database.Query("SELECT InviteToken FROM invites WHERE NewUser= ?", uin)

	if err != nil {
		logging.Error("MySpace/Authentication", "Failed to fetch invite token! (%s)", err.Error())
		return enc, err
	}

	row.Next()
	row.Scan(&enc.IToken)
	row.Close()

	row, err = database.Query("SELECT SignupDate FROM userdetails WHERE UIN= ?", uin)

	if err != nil {
		logging.Error("MySpace/Authentication", "Failed to fetch gift owner! (%s)", err.Error())
		return enc, err
	}

	row.Next()
	row.Scan(&enc.SDate)
	row.Close()

	row, err = database.Query("SELECT RandomSeed FROM userdetails WHERE NewUser= ?", uin)

	if err != nil {
		logging.Error("MySpace/Authentication", "Failed to fetch gift owner! (%s)", err.Error())
		return enc, err
	}

	row.Next()
	row.Scan(&enc.RSeed)
	row.Close()

	row, err = database.Query("SELECT GiftOwner FROM invites WHERE NewUser= ?", uin)

	if err != nil {
		logging.Error("MySpace/Authentication", "Failed to fetch gift owner! (%s)", err.Error())
		return enc, err
	}

	row.Next()
	row.Scan(&enc.GOwner)
	row.Close()

	enc.NUser = strconv.Itoa(uin)

	return enc, err
}
