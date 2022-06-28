package plugins

import (
	"quiet/crackpwd/models"
)

type ScanFunc func(service models.Service) (result models.CrackResult, err error)

var (
	ScanFuncMap map[string]ScanFunc
)

func init() {
	ScanFuncMap = make(map[string]ScanFunc)
	ScanFuncMap["SSH"] = CrackSsh
}
