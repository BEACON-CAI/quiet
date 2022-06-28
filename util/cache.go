package util

import (
	"encoding/gob"
	"fmt"
	"os"
	"strings"
	"time"
	"log"

	"quiet/crackpwd/models"
	"quiet/vars"

	"github.com/patrickmn/go-cache"
)

func init() {
	gob.Register(models.Service{})
	gob.Register(models.CrackResult{})
}


// svae result
func SaveResult(result models.CrackResult, err error) {
	if err == nil && result.Result {
		var k string
		protocol := strings.ToUpper(result.Service.Protocol)

		if protocol == "REDIS" {
			k = fmt.Sprintf("%v-%v-%v", result.Service.IP, result.Service.Port, result.Service.Protocol)
		} else {
			k = fmt.Sprintf("%v-%v-%v", result.Service.IP, result.Service.Port, result.Service.Username)
		}

		h := MakeTaskHash(k)
		SetTaskHash(h)

		_, found := vars.CacheService.Get(k)
		if !found {
			log.Printf("Ip: %v, Port: %v, Protocol: [%v], Username: %v, Password: %v", result.Service.IP,
				result.Service.Port, result.Service.Protocol, result.Service.Username, result.Service.Password)
		}
		vars.CacheService.Set(k, result, cache.NoExpiration)
	}
}

func CacheStatus() (count int, items map[string]cache.Item) {
	count = vars.CacheService.ItemCount()
	items = vars.CacheService.Items()
	return count, items
}

func ResultTotal() {
	vars.ProgressBarPC.Finish()
	log.Printf("Finshed scan, total result: %v, used time: %v",
		vars.CacheService.ItemCount(),
		time.Since(vars.StartTime))
}

func SaveResultToFile() error {
	return vars.CacheService.SaveFile("password_crack.db")
}

func DumpToFile(filename string) (err error) {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	_, items := CacheStatus()
	for _, v := range items {
		result := v.Object.(models.CrackResult)
		_, _ = file.WriteString(fmt.Sprintf("%v:%v|%v,%v:%v\n",
			result.Service.IP,
			result.Service.Port,
			result.Service.Protocol,
			result.Service.Username,
			result.Service.Password),
		)
	}

	return err
}
