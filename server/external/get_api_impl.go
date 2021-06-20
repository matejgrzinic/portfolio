package external

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type getImpl struct{}

func (e *getImpl) getCryptocurrency() ([]byte, error) {
	resp, err := http.Get(os.Getenv("CRYPTOCURRENCY_URL"))
	if err != nil {
		return nil, fmt.Errorf("get cryptocurrency api: %v", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read cryptocurrency api response: %v", err)
	}

	return data, nil
}

func (e *getImpl) getFiat() ([]byte, error) {
	resp, err := http.Get(os.Getenv("FIAT_URL"))
	if err != nil {
		return nil, fmt.Errorf("get fiat api: %v", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read fiat api response: %v", err)
	}

	return data, nil
}
