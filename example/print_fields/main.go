// Prints fields of supplied table name along with their Go type and db tag.
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/lebenasa/pqprobe"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

func usage() {
	fmt.Println("Usage: print_fields [database connection string] [table name]")
	fmt.Println("Example: print_fields  postgres://user:pass@host/database musics")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() < 2 {
		flag.Usage()
		return
	}

	connectionString := flag.Arg(0)
	tableName := flag.Arg(1)

	prober, err := pqprobe.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("Unable to connect to database %v: %v", connectionString, err)
	}

	table, err := prober.QueryTable(tableName)
	if err != nil {
		log.Fatalf("Unable to query table fields: %v", errors.Cause(err))
	}

	for _, v := range table.Fields {
		tag := ""
		if v.IsPrimary {
			tag = fmt.Sprintf(" [PrimaryKey] %v", v.IndexDefinition)
		}
		log.Printf("%v %v `db:\"%v\"`%v\n", v.GoName(), v.GoTypeString(), v.Name, tag)
	}

	return
}
