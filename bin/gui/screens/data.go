package screens

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Decode binary.BigEndian decode
func Decode(in []byte, out interface{}) int {
	buf := bytes.NewReader(in)
	err := binary.Read(buf, binary.BigEndian, out)
	if err != nil {
		log.Println("fail to decode interface:", in[:20], len(in))
		log.Printf("type:%T\n", out)
		return 0
	}
	return len(in) - buf.Len()
}

func getDataFromServer(chain uint64, server, app, structName, key string) []byte {
	if app == "" {
		app = "ff0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"
	}
	urlStr := fmt.Sprintf("%s/api/v1/%d/data?app_name=%s&is_db_data=true&raw=true&key=%s&struct_name=%s",
		server, chain, app, key, structName)
	resp, err := http.Get(urlStr)
	if err != nil {
		log.Println("fail to get data:", urlStr, err)
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Println("error response status:", resp.Status, urlStr)
		return nil
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("fail to read response body:", err)
		return nil
	}

	return data
}
