package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lengzhao/wallet/trans"
)

func proxyHTTP(w http.ResponseWriter, r *http.Request) {
	remote, err := url.Parse(conf.APIServer)
	if err != nil {
		log.Println("fail to proxy:", err)
		return
	}
	// log.Println("proxy:", r.URL.String())
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(w, r)
}

func postTrans(chain uint64, data []byte) error {
	buf := bytes.NewBuffer(data)
	urlStr := fmt.Sprintf("%s/api/v1/%d/transaction/new", conf.APIServer, chain)
	req, err := http.NewRequest(http.MethodPost, urlStr, buf)
	if err != nil {
		log.Println("fail to new http request:", urlStr, err)
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("fail to do request:", urlStr, err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Println("error status:", resp.Status)
		data, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("error response status:%s,msg:%s", resp.Status, data)
	}
	return nil
}

// AccountGet get account of the address on the chain
func AccountGet(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	q := r.URL.Query()
	q.Set("address", conf.AddressStr)
	r.URL.RawQuery = q.Encode()
	proxyHTTP(w, r)
}

// TransMoveInfo move info
type TransMoveInfo struct {
	DstChain uint64 `json:"dst_chain,omitempty"`
	Cost     uint64 `json:"cost,omitempty"`
	Energy   uint64 `json:"energy,omitempty"`
	TransKey string `json:"trans_key,omitempty"`
}

// TransactionMovePost move cost to other chain
func TransactionMovePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chainStr := vars["chain"]
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "fail to read body of request,", err, chainStr)
		return
	}
	chain, err := strconv.ParseUint(chainStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error chain"))
		return
	}
	info := TransMoveInfo{}
	err = json.Unmarshal(data, &info)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "fail to Unmarshal body of request,", err)
		return
	}
	trans := trans.NewTransaction(chain, conf.Address)
	trans.CreateMove(info.DstChain, info.Cost, info.Energy)

	td := trans.GetSignData()
	sign := conf.Wallet.Sign(td)
	trans.SetSign(sign)
	td = trans.Output()
	key := trans.Key[:]

	err = postTrans(chain, td)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error:%s", err)
		return
	}

	info.Energy = trans.Energy
	info.TransKey = hex.EncodeToString(key)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	enc.Encode(info)
}

// TransferInfo transfer info
type TransferInfo struct {
	Peer     string `json:"peer,omitempty"`
	Cost     uint64 `json:"cost,omitempty"`
	Energy   uint64 `json:"energy,omitempty"`
	TransKey string `json:"trans_key,omitempty"`
}

// TransactionTransferPost transfer
func TransactionTransferPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chainStr := vars["chain"]
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "fail to read body of request,", err, chainStr)
		return
	}
	chain, err := strconv.ParseUint(chainStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error chain"))
		return
	}
	info := TransferInfo{}
	err = json.Unmarshal(data, &info)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "fail to Unmarshal body of request,", err)
		return
	}
	if info.Cost == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "error cost value,", info.Cost)
		return
	}

	trans := trans.NewTransaction(chain, conf.Address)
	trans.CreateTransfer(info.Peer, "", info.Cost, info.Energy)
	td := trans.GetSignData()
	sign := conf.Wallet.Sign(td)
	trans.SetSign(sign)
	td = trans.Output()
	key := trans.Key[:]

	err = postTrans(chain, td)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error:%s", err)
		return
	}

	info.Energy = trans.Energy
	info.TransKey = hex.EncodeToString(key)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	enc.Encode(info)
}

// RunApp run app
type RunApp struct {
	Cost      uint64      `json:"cost,omitempty"`
	Energy    uint64      `json:"energy,omitempty"`
	AppName   string      `json:"app_name,omitempty"`
	Param     string      `json:"param,omitempty"`
	ParamType string      `json:"param_type,omitempty"`
	JSONParam interface{} `json:"json_param,omitempty"`
}

// RespOfNewTrans the response of New Transaction
type RespOfNewTrans struct {
	Chain    uint64 `json:"chain,omitempty"`
	TransKey string `json:"trans_key,omitempty"`
}

// TransactionRunAppPost run app
func TransactionRunAppPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chainStr := vars["chain"]
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "fail to read body of request,", err, chainStr)
		return
	}
	chain, err := strconv.ParseUint(chainStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error chain"))
		return
	}
	info := RunApp{}
	err = json.Unmarshal(data, &info)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "fail to Unmarshal body of request,", err)
		return
	}
	log.Println("run app:", info)
	var param []byte
	switch info.ParamType {
	case "json":
		if len(info.Param) > 0 {
			prefix, err := hex.DecodeString(info.Param)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, "error param, hope hex string,", err)
				return
			}
			param = prefix
		}
		jData, _ := json.Marshal(info.JSONParam)
		param = append(param, jData...)
	case "string":
		param = []byte(info.Param)
	default:
		if len(info.Param) > 0 {
			param, err = hex.DecodeString(info.Param)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, "error param, hope hex string,", err)
				return
			}
		}
	}

	trans := trans.NewTransaction(chain, conf.Address)
	err = trans.RunApp(info.AppName, info.Cost, info.Energy, param)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error:%s", err)
		return
	}

	td := trans.GetSignData()
	sign := conf.Wallet.Sign(td)
	trans.SetSign(sign)
	td = trans.Output()

	err = postTrans(chain, td)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error:%s", err)
		return
	}

	resp := RespOfNewTrans{chain, hex.EncodeToString(trans.Key[:])}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	enc.Encode(resp)
}

// AppLife app life
type AppLife struct {
	Energy   uint64 `json:"energy,omitempty"`
	AppName  string `json:"app_name,omitempty"`
	Life     uint64 `json:"life,omitempty"`
	TransKey string `json:"trans_key,omitempty"`
}

// TransactionAppLifePost update app life
func TransactionAppLifePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chainStr := vars["chain"]
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "fail to read body of request,", err, chainStr)
		return
	}
	chain, err := strconv.ParseUint(chainStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error chain"))
		return
	}
	info := AppLife{}
	err = json.Unmarshal(data, &info)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "fail to Unmarshal body of request,", err)
		return
	}

	trans := trans.NewTransaction(chain, conf.Address)
	trans.UpdateAppLife(info.AppName, info.Life, info.Energy)

	td := trans.GetSignData()
	sign := conf.Wallet.Sign(td)
	trans.SetSign(sign)
	td = trans.Output()

	err = postTrans(chain, td)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error:%s", err)
		return
	}

	info.Energy = trans.Energy
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	enc.Encode(info)
}
