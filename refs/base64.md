```go
package main

import (
	b64 "encoding/base64"
	"fmt"
)

func main() {
	data := "IyBHZXRQdWxsZXIKCkEgQ0xJIFRvb2wgZm9yIFNpbXBsaWZpZWQgR2l0SHVi\nIENvbnRlbnQgTWFuYWdlbWVudC4KCiMjIEluc3RhbGxhdGlvbjoKClRvIGlu\nc3RhbGwgYGdpdHB1bGxlcmAsIHJ1biB0aGUgZm9sbG93aW5nIGNvbW1hbmQ6\nCgpgYGBiYXNoCmdvIGluc3RhbGwgZ2l0aHViLmNvbS9TeWVkRGV2b3AvZ2l0\ncHVsbGVyCmBgYAoKIyMgRmVhdHVyZXM6CgotIFsgXSBWaWV3IGZpbGVzCi0g\nWyBdIFNlYXJjaCBSZXBvcwogIC0gWyBdIHNlbGVjdCBicmFuY2gKCiMjIENv\nbmZpZ3VyYXRpb246CgogICogKFJlcXVpcmVkKSBTZXQgdGhlIGBHSVRfVE9L\nRU5gIGVudmlyb25tZW50IHZhcmlhYmxlIHRvIGEgZW52IGZvciBhY2Nlc3Mg\ndG8gdGhlIEdpdEh1YiBBUEkuCiAgKiAoT3B0aW9uYWwpIENyZWF0ZSBhIGNv\nbmZpZyBmaWxlIGF0IGAkSE9NRS8uY29uZmlnL2dpdHB1bGxlci55bWxgIGFu\nZCBzZXQgYSBrZXlzLCBsaWtlOgoKYGBgCmVtYWlsOiBleGFtcGxlQGVtYWls\nLmNvbQp1c2VyOiAgU3llZERldm9wCmBgYAo=\n"
	sDec, err := b64.StdEncoding.DecodeString(data)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Print(string(sDec))
}
```
