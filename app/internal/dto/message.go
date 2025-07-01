package dto

import (
	"encoding/json"
	"log"

	"github.com/go-playground/validator/v10"
)

var validMessageTypes = map[string]bool{
	"connect":    true,
	"disconnect": true,
	"chat":       true,
	"restore":    true,
	"empty":      true,
}

func validateMessage(msg *Message) bool {
	validate := validator.New()

	err := validate.Struct(msg)
	if err != nil {
		log.Println("Invalid message:", err)
		return false
	}

	_, valid := validMessageTypes[msg.Mode]
	if !valid {
		log.Println("Invalid message type:", msg.Mode)
		return false
	}

	return true
}

func DeserializeMessage(jsonData []byte) *Message {
	var msg Message

	err := json.Unmarshal(jsonData, &msg)
	if err != nil {
		log.Println("Error deserializing message:", string(jsonData))
		log.Println(err)
		return nil
	}

	if !validateMessage(&msg) {
		log.Println("Error validating message:", string(jsonData))
		return nil
	}

	return &msg
}

func SerializeMessage(m *Message) []byte {
	data, err := json.Marshal(m)
	if err != nil {
		log.Println("Error serializing message:", err)
		return nil
	}

	return data
}

func EncodeRaw(arr []Message) (json.RawMessage, error) {
	b, err := json.Marshal(arr)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(b), nil
}

// func PrintMessage(m *MessageDTO) {
// 	res, err := json.MarshalIndent(m, "", "  ")
// 	if err != nil {
// 		log.Println("Error printing message:", err)
// 		return
// 	}
// 	log.Println(res)
// }

// func PrintJson(j []byte) {
// 	var buffer bytes.Buffer
// 	err := json.Indent(&buffer, j, "", "\t")
// 	if err != nil {
// 		log.Println("Error prettifying JSON:", j)
// 		return
// 	}
// 	log.Println(buffer.String())
// }
