package baiducloud

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/thoas/go-funk"
	"net/http"
	"strings"
)

func (b *BaiduCloud) Request(data map[string]interface{}, result interface{}) error {
	methods := []string{"GET", "PUT", "POST", "DELETE"}
	if funk.Contains(methods, strings.ToUpper(b.Method)) {

		bytesData, err := json.Marshal(data)
		if err != nil {
			logrus.Println("json error", err)
			return err
		} else {
			request, err := http.NewRequest(strings.ToUpper(b.Method), b.Url, bytes.NewReader(bytesData))
			if err != nil {
				logrus.Println("json error", err)
				return err
			}
			for k, v := range b.Headers {
				request.Header.Set(k, v)
			}
			request.Header.Set("authorization", b.GetAuthorization())
			client := http.Client{}
			resp, err := client.Do(request)
			if err != nil {
				fmt.Println(err.Error())
				return err
			}
			err1 := json.NewDecoder(resp.Body).Decode(result)
			if err1 != nil {
				return err1
			}
			return nil
		}
	} else {
		return errors.New(fmt.Sprintf("%s not allow", b.Method))
	}
}
