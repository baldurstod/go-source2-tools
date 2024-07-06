package model

type DataChannel struct {
	channelClass string
	description  string
	flags        uint64
	channelType  int32
	elementsName []string
	variableName string
	//Grouping
}

func newDataChannel() *DataChannel {
	return &DataChannel{
		elementsName: make([]string, 0),
	}
}
