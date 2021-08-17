package repository

import (
	"github.com/procsy-tech/attorney/entity"
	"time"
	"fmt"
	"encoding/json"
	"crypto/sha256"
	"strings"
)


type (
	DocumentType string
)


const (POADocumentType = "POA"
	)


// Document .
type Document struct {
	Type DocumentType `json:"type"`
}

type POADocument struct{
		Document
		entity.POA
	}
	

func NewPOADocument(e *entity.POA) POADocument{
		timestamp := time.Now().Unix()
		eData, _ := json.Marshal(e)
		h := sha256.Sum256(eData)
		e.BlockchainID = "POA" + fmt.Sprintf("%d", timestamp) + strings.ToUpper(fmt.Sprintf("%x", h[0:4]))
		return POADocument{
			Document{
				Type: POADocumentType,
			},
			*e,
		}
	}
	