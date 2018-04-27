package db

import (
	"strings"

	//required by mattes/migrate
	_ "github.com/mattes/migrate/driver/postgres"
	"github.com/mattes/migrate/migrate"
)

//RunMigrateScripts migrates all the schemas from dbMigrationScriptsPath
//into the dbURL specified
func RunMigrateScripts(dbURL, dbMigrationScriptsPath string) {
	allErrors, ok := migrate.UpSync(dbURL, dbMigrationScriptsPath)
	if !ok {
		errString := []string{"Error running migration scripts : \n"}
		for _, err := range allErrors {
			errString = append(errString, err.Error())
		}
		panic(strings.Join(errString, "\n"))
	}
}
