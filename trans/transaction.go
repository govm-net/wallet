package trans

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/lengzhao/govm/wallet"
)

const (
	// HashLen the byte length of Hash
	HashLen = 32
	// AddressLen the byte length of Address
	AddressLen = 24
)

// Hash The KEY of the block of transaction
type Hash [HashLen]byte

// Address the wallet address
type Address [AddressLen]byte

// TransactionHead transaction = sign + head + data
type TransactionHead struct {
	//signLen uint8
	//sing  []byte
	Time   uint64
	User   Address
	Chain  uint64
	Energy uint64
	Cost   uint64
	Ops    uint8
}

const (
	// OpsTransfer pTransfer
	OpsTransfer = uint8(iota)
	// OpsMove Move out of coin, move from this chain to adjacent chains
	OpsMove
	// OpsNewChain create new chain
	OpsNewChain
	// OpsNewApp create new app
	OpsNewApp
	// OpsRunApp run app
	OpsRunApp
	// OpsUpdateAppLife update app life
	OpsUpdateAppLife
	// OpsRegisterMiner Registered as a miner
	OpsRegisterMiner
	// OpsDisableAdmin disable admin
	OpsDisableAdmin
)

// time
const (
	TimeMillisecond = 1
	TimeSecond      = 1000 * TimeMillisecond
	TimeMinute      = 60 * TimeSecond
	TimeHour        = 60 * TimeMinute
	TimeDay         = 24 * TimeHour
	TimeYear        = 31558150 * TimeSecond
	TimeMonth       = TimeYear / 12
)

// StTrans transaction define
type StTrans struct {
	TransactionHead
	Sign []byte
	Data []byte
	Key  []byte
}

// Encode binary encode
func Encode(in interface{}) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, in)
	if err != nil {
		fmt.Println("fail to encode interface:", in)
		return nil
		// return nil
	}
	return buf.Bytes()
}

// Decode binary decode
func Decode(in []byte, out interface{}) int {
	buf := bytes.NewReader(in)
	err := binary.Read(buf, binary.BigEndian, out)
	if err != nil {
		fmt.Println("fail to decode interface:", in[:20])
		return 0
	}
	return len(in) - buf.Len()
}

// NewTransaction create a new transaction struct
/* 1. NewTransaction
   2. Create*:CreateTransfer,CreateMove...
   3. update energy,time...
   4. GetSignData
   5. SetSign
   6. Output
*/
func NewTransaction(chain uint64, user []byte) *StTrans {
	out := StTrans{}
	out.Chain = chain
	Decode(user, &out.User)
	out.Time = uint64(time.Now().Add(-time.Hour).Unix()) * 1000
	out.Energy = 10000000
	return &out
}

// GetSignData get data of sign
func (t *StTrans) GetSignData() []byte {
	data := Encode(t.TransactionHead)
	data = append(data, t.Data...)
	return data
}

// SetTheSign set the sign
func (t *StTrans) SetTheSign(in []byte) error {
	t.Sign = in
	return nil
}

// Output output the data, the full data of transaction
func (t *StTrans) Output() []byte {
	out := make([]byte, 1, 1000)
	out[0] = uint8(len(t.Sign))
	out = append(out, t.Sign...)
	out = append(out, Encode(t.TransactionHead)...)
	out = append(out, t.Data...)
	t.Key = wallet.GetHash(out)
	return out
}

// CreateTransfer transfer
func (t *StTrans) CreateTransfer(payee, msg string, value, energy uint64) error {
	p, err := hex.DecodeString(payee)
	if err != nil {
		fmt.Println("error peer address:", payee)
		return err
	}
	if len(p) != AddressLen {
		fmt.Println("error peer address length:", payee)
		return fmt.Errorf("error address length:%d", len(p))
	}
	if len(msg) > 100 {
		return fmt.Errorf("data too long:%d", len(msg))
	}
	t.Cost = value
	t.Ops = OpsTransfer
	t.Data = p
	if energy > t.Energy {
		t.Energy = energy
	}
	if msg != "" {
		t.Data = append(t.Data, []byte(msg)...)
	}
	return nil
}

// CreateMove move coin to other chain
func (t *StTrans) CreateMove(dstChain, value, energy uint64) {
	t.Cost = value
	t.Ops = OpsMove
	t.Data = Encode(dstChain)
	if energy > t.Energy {
		t.Energy = energy
	}
}

// RunApp run app
func (t *StTrans) RunApp(app string, cost, energy uint64, data []byte) error {
	p, err := hex.DecodeString(app)
	if err != nil {
		fmt.Println("error app hash:", app)
		return fmt.Errorf("error app hash:%s,err:%s", app, err)
	}
	if len(p) != HashLen {
		fmt.Println("error app length:", app)
		return fmt.Errorf("error app:%d", len(p))
	}
	t.Cost = cost
	t.Ops = OpsRunApp
	t.Data = p
	if len(data) != 0 {
		t.Data = append(t.Data, data...)
	}
	t.Energy = 20*uint64(len(t.Data)) + 10000
	if energy > t.Energy {
		t.Energy = energy
	}
	return nil
}

// UpdateInfo Information of update app life
type UpdateInfo struct {
	Name Hash
	Life uint64
}

// UpdateAppLife update app life
func (t *StTrans) UpdateAppLife(app string, life, energy uint64) error {
	p, err := hex.DecodeString(app)
	if err != nil {
		fmt.Println("error app hash:", app)
		return err
	}
	t.Cost = (life/TimeHour + 1) * 2000
	t.Ops = OpsUpdateAppLife
	info := UpdateInfo{}
	n := Decode(p, &info.Name)
	if n == 0 {
		return fmt.Errorf("error app name")
	}
	info.Life = life
	t.Data = Encode(info)
	if energy > t.Energy {
		t.Energy = energy
	}
	return nil
}
