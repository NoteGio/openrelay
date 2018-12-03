package main

import (
	dbModule "github.com/notegio/openrelay/db"
	"github.com/notegio/openrelay/common"
	"log"
	"os"
	"fmt"
	"strings"
)

func main() {
	db, err := dbModule.GetDB(os.Args[1], os.Args[2])
	if err != nil {
		log.Fatalf("Could not open database connection: %v", err.Error())
	}
	if err := db.AutoMigrate(&dbModule.Order{}).Error; err != nil {
		log.Fatalf("Error migrating order table: %v", err.Error())
	}
	if err := db.AutoMigrate(&dbModule.Cancellation{}).Error; err != nil {
		log.Fatalf("Error migrating cancellation table: %v", err.Error())
	}
	if err := db.AutoMigrate(&dbModule.Exchange{}).Error; err != nil {
		log.Fatalf("Error migrating exchange table: %v", err.Error())
	}
	if err := db.AutoMigrate(&dbModule.Terms{}).Error; err != nil {
		log.Fatalf("Error migrating terms table: %v", err.Error())
	}
	if err := db.AutoMigrate(&dbModule.TermsSig{}).Error; err != nil {
		log.Fatalf("Error migrating term_sigs table: %v", err.Error())
	}
	if err := db.AutoMigrate(&dbModule.HashMask{}).Error; err != nil {
		log.Fatalf("Error migrating hash_masks table: %v", err.Error())
	}
	kovanAddress, _ := common.HexToAddress("0x35dd2932454449b14cee11a94d3674a936d5d7b2")
	db.Where(
		&dbModule.Exchange{Network: 42},
	).FirstOrCreate(&dbModule.Exchange{Network: 42, Address: kovanAddress })
	ganacheAddress, _ := common.HexToAddress("0x48bacb9266a570d521063ef5dd96e61686dbe788")
	db.Where(
		&dbModule.Exchange{Network: 50},
	).FirstOrCreate(&dbModule.Exchange{Network: 50, Address: ganacheAddress })
	mainnetAddress, _ := common.HexToAddress("0x4f833a24e1f95d70f028921e27040ca56e09ab0b")
	db.Where(
		&dbModule.Exchange{Network: 1},
	).FirstOrCreate(&dbModule.Exchange{Network: 1, Address: mainnetAddress })
	if db.Model(&Terms{}).First(&Terms{}).RecordNotFound() {
		if err := dbModule.NewTermsManager(db).UpdateTerms("en", "don't break the law"); err != nil {
			log.Fatalf("Error setting terms: %v", err.Error())
		}
	}
	if err := db.Model(&dbModule.Order{}).AddIndex("idx_order_maker_asset_taker_asset_data", "maker_asset_data", "taker_asset_data").Error; err != nil {
		log.Fatalf("Error adding token pair index: %v", err.Error())
	}
	for _, credString := range(os.Args[3:]) {
		creds := strings.Split(credString, ";")
		if len(creds) != 3 {
			log.Printf("Malformed credential string: %v", credString)
			continue
		}
		username, passwordURI, permissions := creds[0], creds[1], creds[2]
		password := common.GetSecret(passwordURI)
		if dialect := db.Dialect().GetName(); dialect == "postgres" {
			// I don't like using string formatting instead of paramterization, but I
			// don't know of a way to parameterize the username in this statement. It
			// should still be fairly safe, because if you're able to execute this
			// command you already have administrative database access.
			if err = db.Exec(fmt.Sprintf("CREATE USER %v WITH PASSWORD '%v'", username, password)).Error; err != nil {
				log.Printf(err.Error())
			}
			for _, permission := range(strings.Split(permissions, ",")) {
				permArray := strings.Split(permission, ".")
				if len(permArray) != 2 {
					log.Printf("Malformed permission string '$v'", permission)
					continue
				}
				table, permission := permArray[0], permArray[1]
				// I don't like using string formatting instead of paramterization, but I
				// don't know of a way to parameterize the elements in this statement. It
				// should still be fairly safe, because if you're able to execute this
				// command you already have administrative database access.
				if err = db.Exec(fmt.Sprintf("GRANT %v ON TABLE %v TO %v", permission, table, username)).Error; err != nil {
					log.Printf(err.Error())
				}
			}
			if err = db.Exec(fmt.Sprintf("GRANT USAGE, SELECT on ALL SEQUENCES in SCHEMA public to %v", username)).Error; err != nil {
				log.Printf(err.Error())
			}
		} else if dialect == "mysql" {
			if err := db.Exec(fmt.Sprintf("CREATE USER '%v' IDENTIFIED BY '%v'", username, password)).Error; err != nil {
				log.Printf(err.Error())
			}
			result := make(map[string]string)
			if err := db.Exec("SELECT DATABASE()").Row().Scan(result); err != nil {
				log.Printf(err.Error())
			}
			log.Printf("'%v'", result)
			databaseName := result["DATABASE()"]
			log.Printf("Database name: %v", databaseName)
			for _, permission := range(strings.Split(permissions, ",")) {
				permArray := strings.Split(permission, ".")
				if len(permArray) != 2 {
					log.Printf("Malformed permission string '$v'", permission)
					continue
				}
				table, permission := permArray[0], permArray[1]
				// I don't like using string formatting instead of paramterization, but I
				// don't know of a way to parameterize the elements in this statement. It
				// should still be fairly safe, because if you're able to execute this
				// command you already have administrative database access.
				if err = db.Exec(fmt.Sprintf("GRANT %v ON %v.%v TO '%v'", permission, databaseName, table, username)).Error; err != nil {
					log.Printf(err.Error())
				}
			}
			if err := db.Exec("FLUSH PRIVILEGES;").Error; err != nil {
				log.Printf(err.Error());
			}
		}
		log.Printf("Created '%v'", credString)
	}
}
