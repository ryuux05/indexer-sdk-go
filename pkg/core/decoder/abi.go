package decoder

type ABI []ABIItem

type ABIItem struct {
	Anonymous bool   `json:"anonymous"`
	Type      string `json:"type"`
	Name      string `json:"name"`
	Inputs    []ABIInput `json:"inputs"`
}

type ABIInput struct {
	Indexed      bool   `json:"indexed,omitempty"`
	InternalType string `json:"internalType"`
	Name         string `json:"name"`
	Type         string `json:"type"`
}
