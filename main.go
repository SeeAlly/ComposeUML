package main

import (
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"reflect"
)

// type service struct {
//   Depends []string `yaml:"depends"`
// }

// type composeFile struct {
//   Services map[string]service `yaml:"services"`
//   networks map[string]interface{}
//   volumes  map[string]interface{}
//   version  string
// }

func parseDependencies(file map[interface{}]interface{}) map[string][]string {
	servicesWithDependencies := make(map[string][]string)

	services := file["services"]
	servicesMapIterator := reflect.ValueOf(services).MapRange()
	for servicesMapIterator.Next() {
		service := servicesMapIterator.Key().Interface().(string)
		fields := servicesMapIterator.Value().Interface()
		fieldsMapIterator := reflect.ValueOf(fields).MapRange()
		var serviceDependencies []string
		for fieldsMapIterator.Next() {
			key := fieldsMapIterator.Key().Interface().(string)
			if key == "depends_on" {
				deps := fieldsMapIterator.Value().Interface().([]interface{})
				for _, v := range deps {
					serviceDependencies = append(serviceDependencies, v.(string))
				}
			}
		}
		if len(serviceDependencies) == 0 {
			log.Printf("Service %s does not got dependencies", service)
			continue
		}
		log.Printf("%s depends on %v", service, serviceDependencies)
		servicesWithDependencies[service] = serviceDependencies
	}

	return servicesWithDependencies
}

func main() {
	filePath := flag.String("f", "./docker-compose.yml", "Path to the compose file")
	log.Printf("Gonna parse %v", *filePath)

	log.Print("Opening file...")
	yamlFile, err := ioutil.ReadFile(*filePath)
	if err != nil {
		log.Panicf("yamlFile.Get err   #%v ", err)
	}
	log.Print("Opened file")

	log.Print("Parsing file...")
	parsedFile := make(map[interface{}]interface{})
	err = yaml.Unmarshal(yamlFile, parsedFile)
	if err != nil {
		log.Panicf("Unmarshal: %v", err)
	}
	log.Print("Parsed file")

	log.Print("Getting services dependencies list...")
	servicesWithDependencies := parseDependencies(parsedFile)
	log.Print("Got services dependencies list")

	log.Printf("%v", servicesWithDependencies)
	log.Print("Done")
}
