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
	"sync"
	"time"

	"github.com/gorilla/mux"
	core "github.com/govm-net/govm/core"
	"github.com/lengzhao/wallet/trans"
)

// Blance Blance
type Blance struct {
	updateTime int64
	value      map[uint64]uint64
	mu         sync.Mutex
}

func (b *Blance) get(chain uint64) uint64 {
	now := time.Now().UnixNano()
	if b.updateTime > now && b.value[chain] > 0 {
		b.mu.Lock()
		defer b.mu.Unlock()
		return b.value[chain]
	}
	data := getDataFromServer(chain, "", "dbCoin", wallet.AddressStr)
	b.updateTime = time.Now().Add(time.Second * 2).UnixNano()
	if len(data) == 0 {
		return 0
	}
	var cost uint64
	trans.Decode(data, &cost)
	b.mu.Lock()
	defer b.mu.Unlock()
	b.value[chain] = cost
	return cost
}

var balance Blance

func init() {
	balance.value = make(map[uint64]uint64)
}

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
	if q.Get("address") == "" {
		q.Set("address", wallet.AddressStr)
		r.URL.RawQuery = q.Encode()
	}

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
	trans := trans.NewTransaction(chain, wallet.Address, info.Cost)
	trans.SetEnergy(info.Energy)
	if info.Cost+trans.Energy > balance.get(chain) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "not enough cost")
		return
	}
	trans.CreateMove(info.DstChain)

	td := trans.GetSignData()
	sign := wallet.Sign(td)
	trans.SetTheSign(sign)
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

	trans := trans.NewTransaction(chain, wallet.Address, info.Cost)
	trans.SetEnergy(info.Energy)
	if info.Cost+trans.Energy > balance.get(chain) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "not enough cost")
		return
	}
	trans.CreateTransfer(info.Peer, "")
	td := trans.GetSignData()
	sign := wallet.Sign(td)
	trans.SetTheSign(sign)
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

// NewApp new app
type NewApp struct {
	Cost         uint64 `json:"cost,omitempty"`
	Energy       uint64 `json:"energy,omitempty"`
	CodePath     string `json:"code_path,omitempty"`
	IsPrivate    bool   `json:"is_private,omitempty"`
	EnableRun    bool   `json:"enable_run,omitempty"`
	EnableImport bool   `json:"enable_import,omitempty"`
	AppName      string `json:"app_name,omitempty"`
	TransKey     string `json:"trans_key,omitempty"`
}

// TransactionNewAppPost new app
func TransactionNewAppPost(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := recover()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "error code:", err)
		}
	}()
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
	info := NewApp{}
	err = json.Unmarshal(data, &info)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "fail to Unmarshal body of request,", err)
		return
	}

	var flag uint8
	if !info.IsPrivate {
		flag |= core.AppFlagPlublc
		log.Println("1. flag:", flag)
	}
	if info.EnableRun {
		flag |= core.AppFlagRun
		log.Println("2. flag:", flag)
	}
	if info.EnableImport {
		flag |= core.AppFlagImport
		log.Println("3. flag:", flag)
	}
	code, _ := core.CreateAppFromSourceCode(info.CodePath, flag)
	t := trans.NewTransaction(chain, wallet.Address, info.Cost)
	t.SetEnergy(info.Energy)
	if info.Cost+t.Energy > balance.get(chain) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "not enough cost")
		return
	}
	t.Ops = trans.OpsNewApp
	t.Data = code

	td := t.GetSignData()
	sign := wallet.Sign(td)
	t.SetTheSign(sign)
	td = t.Output()
	key := t.Key[:]

	err = postTrans(chain, td)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error:%s", err)
		return
	}

	info.Energy = t.Energy
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
		log.Println("run json param:", info.Param, string(jData))
		param = append(param, jData...)
	case "string":
		param = []byte(info.Param)
		log.Println("run string param:", info.Param)
	default:
		if len(info.Param) > 0 {
			param, err = hex.DecodeString(info.Param)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, "error param, hope hex string,", err)
				return
			}
			log.Println("run param:", info.Param)
		}
	}

	trans := trans.NewTransaction(chain, wallet.Address, info.Cost)
	trans.SetEnergy(info.Energy)
	if info.Cost+trans.Energy > balance.get(chain) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "not enough cost")
		return
	}
	err = trans.RunApp(info.AppName, param)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error:%s", err)
		return
	}

	td := trans.GetSignData()
	sign := wallet.Sign(td)
	trans.SetTheSign(sign)
	td = trans.Output()

	err = postTrans(chain, td)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error:%s", err)
		return
	}
	tk := hex.EncodeToString(trans.Key[:])
	log.Println("run app:", tk, info)
	resp := RespOfNewTrans{chain, tk}
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

	trans := trans.NewTransaction(chain, wallet.Address, 0)
	trans.SetEnergy(info.Energy)
	if trans.Energy > balance.get(chain) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "not enough cost")
		return
	}
	trans.UpdateAppLife(info.AppName, info.Life)

	td := trans.GetSignData()
	sign := wallet.Sign(td)
	trans.SetTheSign(sign)
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

// VoteInfo vote info
type VoteInfo struct {
	Peer     string `json:"peer,omitempty"`
	Cost     uint64 `json:"cost,omitempty"`
	Energy   uint64 `json:"energy,omitempty"`
	TransKey string `json:"trans_key,omitempty"`
}

// TransactionVotePost vote
func TransactionVotePost(w http.ResponseWriter, r *http.Request) {
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
	info := VoteInfo{}
	err = json.Unmarshal(data, &info)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "fail to Unmarshal body of request,", err)
		return
	}

	trans := trans.NewTransaction(chain, wallet.Address, info.Cost)
	trans.SetEnergy(info.Energy)
	if info.Cost+trans.Energy > balance.get(chain) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "not enough cost")
		return
	}
	trans.CreateVote(info.Peer)
	td := trans.GetSignData()
	sign := wallet.Sign(td)
	trans.SetTheSign(sign)
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

// TransactionVoteDelete unvote
func TransactionVoteDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chainStr := vars["chain"]
	chain, err := strconv.ParseUint(chainStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error chain"))
		return
	}

	trans := trans.NewTransaction(chain, wallet.Address, 0)
	if trans.Energy > balance.get(chain) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "not enough cost")
		return
	}
	trans.Unvote()
	td := trans.GetSignData()
	sign := wallet.Sign(td)
	trans.SetTheSign(sign)
	td = trans.Output()
	key := trans.Key[:]

	err = postTrans(chain, td)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error:%s", err)
		return
	}

	info := VoteInfo{}
	info.Energy = trans.Energy
	info.TransKey = hex.EncodeToString(key)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	enc.Encode(info)
}

func getDataFromServer(chain uint64, app, structName, key string) []byte {
	if app == "" {
		app = "ff0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"
	}
	urlStr := fmt.Sprintf("%s/api/v1/%d/data?app_name=%s&is_db_data=true&raw=true&key=%s&struct_name=%s",
		conf.APIServer, chain, app, key, structName)
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

// AdminsGet get admin list
func AdminsGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chainStr := vars["chain"]
	chain, err := strconv.ParseUint(chainStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error chain"))
		return
	}
	// key = []byte{core.StatAdmin}
	data := getDataFromServer(chain, "", "dbStat", "0d")
	if len(data) == 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
		return
	}
	var adminList [trans.AdminNum]trans.Address
	trans.Decode(data, &adminList)
	var out []string
	for _, it := range adminList {
		if it.Empty() {
			continue
		}
		val := hex.EncodeToString(it[:])
		out = append(out, val)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	enc.Encode(out)
}

// AdminInfo vote info
type AdminInfo struct {
	Address string `json:"address"`
	Deposit uint64 `json:"deposit"`
	Votes   uint64 `json:"votes"`
}

// AdminInfoGet get admin info
func AdminInfoGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chainStr := vars["chain"]
	r.ParseForm()
	key := r.Form.Get("key")
	chain, err := strconv.ParseUint(chainStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error chain"))
		return
	}
	// key = []byte{core.StatAdmin}
	data := getDataFromServer(chain, "", "dbAdmin", key)
	if len(data) == 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
		return
	}
	var info struct {
		Deposit uint64
		Votes   uint64
	}
	trans.Decode(data, &info)
	var out AdminInfo
	out.Address = key
	out.Deposit = info.Deposit
	out.Votes = info.Votes

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	enc.Encode(out)
}

// Miner miner info
type Miner struct {
	TagetChain uint64 `json:"taget_chain,omitempty"`
	Cost       uint64 `json:"cost,omitempty"`
	Energy     uint64 `json:"energy,omitempty"`
	Miner      string `json:"miner,omitempty"`
	TransKey   string `json:"trans_key,omitempty"`
}

// TransactionMinerPost miner
func TransactionMinerPost(w http.ResponseWriter, r *http.Request) {
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
	info := Miner{}
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

	trans := trans.NewTransaction(chain, wallet.Address, info.Cost)
	trans.SetEnergy(info.Energy)
	if info.Cost+trans.Energy > balance.get(chain) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "not enough cost")
		return
	}
	err = trans.RegisterMiner(info.TagetChain, info.Miner)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err)
		return
	}
	td := trans.GetSignData()
	sign := wallet.Sign(td)
	trans.SetTheSign(sign)
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

// VersionInfo version info
type VersionInfo struct {
	Version string
	// BuildTime string
	// GitHead   string
}

// VersionGet get software version
func VersionGet(w http.ResponseWriter, r *http.Request) {
	info := VersionInfo{Version}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	enc.Encode(info)
}
