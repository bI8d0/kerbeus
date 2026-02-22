package utils

import (
	"bufio"
	"os"
	"strings"
)

func UserExistsInFile(filename, username string) bool {
	file, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 && parts[0] == username {
				return true
			}
		}
		if strings.HasPrefix(line, "$krb5pa$") || strings.HasPrefix(line, "$krb5asrep$") {
			hashParts := strings.Split(line, "$")
			if len(hashParts) >= 4 && hashParts[3] == username {
				return true
			}
		}
	}
	return false
}

func UpdateHashForUser(filename, username, salt, newHash string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var lines []string
	updated := false
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if updated {
			lines = append(lines, line)
			continue
		}

		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 && parts[0] == username {
				hashParts := strings.Split(parts[1], "$")
				if len(hashParts) >= 7 {
					reconstructed := username + ":$" + hashParts[1] + "$" + hashParts[2] + "$" +
						hashParts[3] + "$" + hashParts[4] + "$" + salt + "$" + newHash
					lines = append(lines, reconstructed)
					updated = true
					continue
				}
			}
		}

		if (strings.HasPrefix(line, "$krb5pa$") || strings.HasPrefix(line, "$krb5asrep$")) && !updated {
			hashParts := strings.Split(line, "$")
			if len(hashParts) >= 7 && hashParts[3] == username {
				reconstructed := username + ":$" + hashParts[1] + "$" + hashParts[2] + "$" +
					hashParts[3] + "$" + hashParts[4] + "$" + salt + "$" + newHash
				lines = append(lines, reconstructed)
				updated = true
				continue
			}
		}

		lines = append(lines, line)
	}

	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outFile.Close()

	writer := bufio.NewWriter(outFile)
	for _, line := range lines {
		writer.WriteString(line + "\n")
	}
	return writer.Flush()
}
