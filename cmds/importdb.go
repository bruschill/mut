package cmds

import (
	"log"
	"os/exec"
	"regexp"
	"strings"
)

const tmpDir = "$HOME/.mbc_db_dumps"

//ImportDB imports the database from the database that corresponds to the appName parameter
func ImportDB(appName string) {
	//verifying that the heroku cli utility is installed
	_, err := exec.Command("which", "heroku").Output()

	if err != nil {
		log.Fatalln("The Heroku cli utility must be installed.")
	}

	//check that user is logged into heroku as correct user
	out, err := exec.Command("heroku", "orgs").Output()

	//convert out []byte to string
	outStr := string(out[:])

	if err != nil || !strings.Contains(outStr, "merchantsbonding") {
		log.Fatalln("You must be logged into heroku with a user that is part of the merchantsbonding organization.")
	}
}

func dumpDB(appName string) {
	fileName := appName + ".dump"
	herokuCmdStr := "$(heroku pg:backups public-url -a" + appName + ")"
	//make tmpDir
	//exclude dir from Time Machine
	//curl -o filename to tmpDir from $(heroku pg:backups public-url -a appName)
	err := exec.Command("curl", "-o", tmpDir, fileName, herokuCmdStr).Run()

	if err != nil {
		log.Fatalln("There was an error when downloading the database backup. Try again later.")
	}

}

func loadDB(appName string) {
	//dropdb dbName --if-exists 2> /dev/null
	//createdb dbName --if-exists 2> /dev/null
	//pg_restore
	//  -j runtime.GOMAXPROCS(0)
	//  -c
	//  -x
	//  -O
	//  -h localhost
	//  -U $USER
	//  -d dbName
	//  tmpDir + filename
}

func cleanup(appName string) {
	//if appName ends with "-mbc"
	//  run "psql -U $USER -d #{db_name} -c 'UPDATE efile_applications SET social_security_number = NULL, driver_license_number = NULL'"
	//  run "psql -U $USER -d #{db_name} -c 'UPDATE agencies SET tax_id_number = NULL'"
	//  run "psql -U $USER -d #{db_name} -c 'UPDATE agents SET ssn = NULL'"
}

func toDBName(appName string) string {
	_, err := regexp.Compile("/mbc-(prod|test)-/")

	if err != nil {
		log.Fatalln("Attempted to compile invalid regular expression when converting appName to dbName")
	}

	return ""
}
