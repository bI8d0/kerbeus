package asreq

import (
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jcmturner/gokrb5/v8/messages"
	"github.com/jcmturner/gokrb5/v8/types"

	"kerbeus/src/models"
)

func Process(buffer []byte, filename string, seen map[string]bool,
	pendingReqs map[string]*models.PendingASREQ, srcIP, dstIP string) {
	var asReq messages.ASReq
	err := asReq.Unmarshal(buffer)
	if err != nil {
		log.Printf("\r\033[K⚠️  Error parsing AS-REQ: %v", err)
		return
	}

	usuario := asReq.ReqBody.CName.PrincipalNameString()
	realm := asReq.ReqBody.Realm

	if usuario == "" || realm == "" {
		return
	}

	userKey := fmt.Sprintf("\r\033[KAS-REQ:%s", usuario)
	if seen[userKey] {
		return
	}

	for _, padata := range asReq.PAData {
		if padata.PADataType != 2 {
			continue
		}

		var encData types.EncryptedData
		err := encData.Unmarshal(padata.PADataValue)
		if err != nil {
			continue
		}

		etype := encData.EType

		if etype != 17 && etype != 18 && etype != 23 {
			continue
		}

		cipher := strings.ToLower(hex.EncodeToString(encData.Cipher))

		hashKey := fmt.Sprintf("HASH:%s", usuario)
		pendingReqs[hashKey] = &models.PendingASREQ{
			Usuario:   usuario,
			Realm:     realm,
			EType:     etype,
			Salt:      cipher,
			Timestamp: time.Now(),
			SrcIP:     srcIP,
			DstIP:     dstIP,
		}
		return
	}
}
