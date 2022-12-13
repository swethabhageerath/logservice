package models

import "fmt"

type LogDetail struct {
	PackageName string
	FileName    string
	Operation   string
}

func (l LogDetail) String() string {
	return fmt.Sprintf("{\"PackageName\":%s, \"FileName\":%s, \"Operation\":%s}", l.PackageName, l.FileName, l.Operation)
}
