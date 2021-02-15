 
package nodeinfo

import "strings"

type NgbNodeInfo struct {
	NgbId        string //neighbor node id
	NgbType      string
	NgbAddr      string
	HttpInfoAddr string
	HttpInfoPort uint16
	NgbVersion   string
}

type NgbNodeInfoSlice []NgbNodeInfo

func (n NgbNodeInfoSlice) Len() int {
	return len(n)
}

func (n NgbNodeInfoSlice) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

func (n NgbNodeInfoSlice) Less(i, j int) bool {
	if 0 <= strings.Compare(n[i].HttpInfoAddr, n[j].HttpInfoAddr) {
		return false
	} else {
		return true
	}
}
