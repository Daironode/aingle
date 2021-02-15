
package ontfs

import (
	"fmt"
	"testing"
)

func TestErrors_Serialization(t *testing.T) {
	var e Errors
	e.AddObjectError("file1", "transfer error1")
	e.AddObjectError("file2", "transfer error2")
	e.AddObjectError("file3", "transfer error3")

	data := e.ToString()
	fmt.Printf("%v\n", data)

	var f Errors
	f.FromString(data)
	for obj, err := range f.ObjectErrors {
		fmt.Printf("obj:%s   error: %s\n", obj, err)
	}
}
