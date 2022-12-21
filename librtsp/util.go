package librtsp

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
)

const DefaultAlgorithm = "MD5"

func parseAuth(str string, iterator func(k, v string) error) error {
	isQuotes, offset, parseOnce := false, 0, false

	for i, char := range str {
		if char == '"' {
			isQuotes = !isQuotes
			//last params
			parseOnce = !isQuotes && i == len(str)-1
		} else if (char == ',' || i == len(str)-1) && !isQuotes {
			parseOnce = true
		}

		if !parseOnce {
			continue
		}
		parseOnce = false

		var params string
		if char == ',' {
			params = str[offset:i]
		} else {
			params = str[offset:]
		}
		offset = i + 1

		split := strings.Split(params, "=")
		if len(split) != 2 {
			return fmt.Errorf("bad auth params :%s", params)
		}

		k := strings.TrimSpace(split[0])
		v := strings.TrimSpace(split[1])
		if strings.HasPrefix(v, "\"") || strings.HasSuffix(v, "\"") {
			v = strings.Trim(v, "\"")
		}

		if err := iterator(k, v); err != nil {
			return err
		}

	}

	return nil
}

func parseWWWAuthenticateHeader(str string) (map[string]string, error) {
	index := strings.Index(str, " ")
	if index < 0 {
		return nil, fmt.Errorf("not Find digest in WWWAuthenticate %s", str)
	}

	schema := str[:index]
	if schema != "Digest" {

	}

	params := make(map[string]string, 10)
	if err := parseAuth(str[index+1:], func(k, v string) error {
		//fmt.Printf("name: %s , value: %s", n, v)
		params[strings.ToLower(k)] = v
		return nil
	}); err != nil {
		return nil, err
	}
	return params, nil
}

func h(data string) string {
	hash := md5.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

func calculateResponse(username, realm, nonce, uri, password string) string {
	//H(data) = MD5(data)
	//KD(secret, data) = H(concat(secret, ":", data))
	//request-digest  = <"> < KD ( H(A1), unq(nonce-value) ":" H(A2) ) > <">
	A1 := fmt.Sprintf("%s:%s:%s", username, realm, password)
	A2 := fmt.Sprintf("%s:%s", "DESCRIBE", uri)

	return h(h(A1) + ":" + nonce + ":" + h(A2))
}

func generateCredentials(params map[string]string, password string) (string, error) {

	realm := params["realm"]
	nonce := params["nonce"]
	uri := params["uri"]
	username := params["username"]

	if realm == "" || nonce == "" {
		return "", fmt.Errorf("authorization is missing required fields")
	}

	//if "auth" == wwwAuthenticateHeader.Qop() || "auth-int" == wwwAuthenticateHeader.Qop() {
	//	return false
	//}

	response := calculateResponse(username, realm, nonce, uri, password)
	return fmt.Sprintf("Digest username=\"%s\", realm=\"%s\", nonce=\"%s\", uri=\"%s\", response=\"%s\", algorithm=%s", username, realm, nonce, uri, response, DefaultAlgorithm), nil
}
