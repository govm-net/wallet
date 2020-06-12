package trans

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lengzhao/govm/wallet"
)

const (
	// HashLen the byte length of Hash
	HashLen = 32
	// AddressLen the byte length of Address
	AddressLen = 24
	// AdminNum admin number
	AdminNum = 23
)

// Hash The KEY of the block of transaction
type Hash [HashLen]byte

// Address the wallet address
type Address [AddressLen]byte

// Empty Check whether Hash is empty
func (h Hash) Empty() bool {
	return h == (Hash{})
}

// MarshalJSON marshal by base64
func (h Hash) MarshalJSON() ([]byte, error) {
	if h.Empty() {
		return json.Marshal(nil)
	}
	return json.Marshal(h[:])
}

// UnmarshalJSON UnmarshalJSON
func (h *Hash) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return nil
	}
	var v []byte
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}
	copy(h[:], v)
	return nil
}

// Empty Check where Address is empty
func (a Address) Empty() bool {
	return a == (Address{})
}

// MarshalJSON marshal by base64
func (a Address) MarshalJSON() ([]byte, error) {
	if a.Empty() {
		return json.Marshal(nil)
	}
	return json.Marshal(a[:])
}

// UnmarshalJSON UnmarshalJSON
func (a *Address) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return nil
	}
	var v []byte
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}
	copy(a[:], v)
	return nil
}

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
	// OpsRegisterAdmin Registered as a admin
	OpsRegisterAdmin
	// OpsVote vote admin
	OpsVote
	// OpsUnvote unvote
	OpsUnvote
	// OpsReportError error block
	OpsReportError
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
func NewTransaction(chain uint64, user []byte, cost uint64) *StTrans {
	out := StTrans{}
	out.Chain = chain
	Decode(user, &out.User)
	out.Time = uint64(time.Now().Unix()) * 1000
	out.Cost = cost
	out.Energy = 1000000
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

// SetEnergy set energy
func (t *StTrans) SetEnergy(energy uint64) bool {
	if energy > t.Energy {
		t.Energy = energy
		return true
	}
	return false
}

// CreateTransfer transfer
func (t *StTrans) CreateTransfer(payee, msg string) error {
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
	t.Ops = OpsTransfer
	t.Data = p
	if msg != "" {
		t.Data = append(t.Data, []byte(msg)...)
	}
	return nil
}

// CreateMove move coin to other chain
func (t *StTrans) CreateMove(dstChain uint64) {
	t.Ops = OpsMove
	t.Data = Encode(dstChain)
}

// RunApp run app
func (t *StTrans) RunApp(app string, data []byte) error {
	p, err := hex.DecodeString(app)
	if err != nil {
		fmt.Println("error app hash:", app)
		return fmt.Errorf("error app hash:%s,err:%s", app, err)
	}
	if len(p) != HashLen {
		fmt.Println("error app length:", app)
		return fmt.Errorf("error app:%d", len(p))
	}
	t.Ops = OpsRunApp
	t.Data = p
	if len(data) != 0 {
		t.Data = append(t.Data, data...)
	}
	energy := 20*uint64(len(t.Data)) + 10000
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
func (t *StTrans) UpdateAppLife(app string, life uint64) error {
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
	return nil
}

// CreateVote vote
func (t *StTrans) CreateVote(payee string) error {
	p, err := hex.DecodeString(payee)
	if err != nil {
		fmt.Println("error peer address:", payee)
		return err
	}
	if len(p) != AddressLen {
		fmt.Println("error peer address length:", payee)
		return fmt.Errorf("error address length:%d", len(p))
	}
	if t.Cost%1000000000 != 0 {
		fmt.Println("error vote.", t.Cost)
		return fmt.Errorf("error vote")
	}
	t.Ops = OpsVote
	t.Data = p
	return nil
}

// Unvote cancel vote
func (t *StTrans) Unvote() error {
	t.Ops = OpsUnvote
	return nil
}

// RegisterMiner RegisterMiner
func (t *StTrans) RegisterMiner(chain uint64, peer string) error {
	t.Ops = OpsRegisterMiner

	if chain != 0 && chain != t.Chain {
		t.Data = Encode(chain)
	}
	if peer != "" {
		p, err := hex.DecodeString(peer)
		if err != nil {
			fmt.Println("error address:", peer)
			return err
		}
		if len(p) != AddressLen {
			fmt.Println("error address length:", peer)
			return fmt.Errorf("error address length:%d", len(p))
		}
		t.Data = Encode(chain)
		t.Data = append(t.Data, p...)
	}
	return nil
}
