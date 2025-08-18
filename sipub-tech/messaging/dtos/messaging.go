package dtos

type Message struct {
	Metadata MessageMetadata
	Data any
} 

type MessageMetadata struct {
	CorrelationId string
}
