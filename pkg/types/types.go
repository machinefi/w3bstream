package types

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/depends/kit/validator/strfmt"
	"github.com/machinefi/w3bstream/pkg/enums"
)

type UploadConfig struct {
	FilesizeLimitBytes int64 `env:""`
	DiskReserveBytes   int64 `env:""`
}

func (c *UploadConfig) SetDefault() {
	if c.FilesizeLimitBytes == 0 {
		c.FilesizeLimitBytes = 1024 * 1024
	}
	if c.DiskReserveBytes == 0 {
		c.DiskReserveBytes = 20 * 1024 * 1024
	}
}

type FileSystem struct {
	Type enums.FileSystemMode `env:""`
}

func (f *FileSystem) SetDefault() {
	if f.Type > enums.FILE_SYSTEM_MODE__S3 || f.Type <= 0 {
		f.Type = enums.FILE_SYSTEM_MODE__LOCAL
	}
}

type ETHClientConfig struct {
	Endpoints string            `env:""`
	Clients   map[uint32]string `env:"-"`
}

func (c *ETHClientConfig) Init() {
	c.Clients = make(map[uint32]string)
	if !gjson.Valid(c.Endpoints) {
		return
	}
	for k, v := range gjson.Parse(c.Endpoints).Map() {
		chainID, err := strconv.Atoi(k)
		if err != nil {
			continue
		}
		url := v.String()
		c.Clients[uint32(chainID)] = url
	}
}

type Chain struct {
	ChainID                         uint64          `json:"chainID,omitempty"`
	Name                            enums.ChainName `json:"name"`
	Endpoint                        string          `json:"endpoint"`
	AABundlerEndpoint               string          `json:"aaBundlerEndpoint"`
	AAPaymasterEndpoint             string          `json:"aaPaymasterEndpoint"`
	AAEntryPointContractAddress     string          `json:"aaEntryPointContractAddress"`
	AAAccountFactoryContractAddress string          `json:"aaAccountFactoryContractAddress"`
}

func (c *Chain) IsSolana() bool {
	return c.Name == enums.SOLANA_DEVNET || c.Name == enums.SOLANA_TESTNET || c.Name == enums.SOLANA_MAINNET_BETA
}

func (c *Chain) IsEth() bool {
	return c.ChainID != 0
}

func (c *Chain) IsAASupported() bool {
	return c.AABundlerEndpoint != "" && c.AAPaymasterEndpoint != "" && c.AAEntryPointContractAddress != "" && c.AAAccountFactoryContractAddress != ""
}

type ChainConfig struct {
	Configs          string                     `env:""     json:"-"`
	Chains           map[enums.ChainName]*Chain `env:"-"    json:"-"`
	ChainIDs         map[uint64]*Chain          `env:"-"    json:"-"`
	AAUserOpEndpoint string                     `env:""     json:"-"`
}

func (cc *ChainConfig) LivenessCheck() map[string]string {
	m := map[string]string{}

	for _, c := range cc.Chains {
		key := c.Endpoint
		if err := cc.chainLivenessCheck(c); err != nil {
			m[key] = err.Error()
		} else {
			m[key] = "ok"
		}
	}
	return m
}

func (c *ChainConfig) Init() {
	cs := []*Chain{}
	if c.Configs != "" {
		if err := json.Unmarshal([]byte(c.Configs), &cs); err != nil {
			panic(err)
		}
	}

	cm := make(map[enums.ChainName]*Chain)
	cidm := make(map[uint64]*Chain)
	for _, c := range cs {
		cm[c.Name] = c
		if c.ChainID != 0 {
			cidm[c.ChainID] = c
		}
	}
	c.Chains = cm
	c.ChainIDs = cidm
}

func (c *ChainConfig) GetChain(chainID uint64, chainName enums.ChainName) (*Chain, bool) {
	r, ok := c.ChainIDs[chainID]
	if ok {
		return r, ok
	}
	r, ok = c.Chains[chainName]
	return r, ok
}

func (c *ChainConfig) chainLivenessCheck(chain *Chain) error {
	if chain.IsSolana() {
		cli := client.NewClient(chain.Endpoint)
		_, err := cli.GetLatestBlockhash(context.Background())
		return err
	} else if chain.IsEth() {
		cli, err := ethclient.Dial(chain.Endpoint)
		if err != nil {
			return err
		}
		chainID, err := cli.ChainID(context.Background())
		if err != nil {
			return err
		}
		if chainID.Uint64() != chain.ChainID {
			return errors.Errorf("chainID mismatch, want %d, got %d", chain.ChainID, chainID.Uint64())
		}
	}
	return nil
}

// aliases from base/types
type (
	SFID                     = types.SFID
	SFIDs                    = types.SFIDs
	EthAddress               = types.EthAddress
	Timestamp                = types.Timestamp
	Initializer              = types.Initializer
	ValidatedInitializer     = types.ValidatedInitializer
	InitializerWith          = types.InitializerWith
	ValidatedInitializerWith = types.ValidatedInitializerWith
)

type EthAddressWhiteList []string

func (v *EthAddressWhiteList) Init() {
	lst := EthAddressWhiteList{}
	for _, addr := range *v {
		if err := strfmt.EthAddressValidator.Validate(addr); err == nil {
			lst = append(lst, strings.ToLower(addr))
		}
	}
	*v = lst
}

func (v *EthAddressWhiteList) Validate(address string) bool {
	if v == nil || len(*v) == 0 {
		return true
	}
	for _, addr := range *v {
		if addr == strings.ToLower(address) {
			return true
		}
	}
	return false
}

type StrategyResult struct {
	ProjectName string         `json:"projectName" db:"f_prj_name"`
	AppletID    types.SFID     `json:"appletID"    db:"f_app_id"`
	AppletName  string         `json:"appletName"  db:"f_app_name"`
	InstanceID  types.SFID     `json:"instanceID"  db:"f_ins_id"`
	Handler     string         `json:"handler"     db:"f_hdl"`
	EventType   string         `json:"eventType"   db:"f_evt"`
	AutoCollect datatypes.Bool `json:"autoCollect" db:"f_auto_collect"`
}

type WasmDBConfig struct {
	Endpoint        types.Endpoint
	MaxConnection   int
	PoolSize        int
	ConnMaxLifetime types.Duration
}

func (c *WasmDBConfig) SetDefault() {
	if c.MaxConnection == 0 {
		c.MaxConnection = 2
	}
	if c.PoolSize == 0 {
		c.PoolSize = 2
	}
	if c.ConnMaxLifetime == 0 {
		c.ConnMaxLifetime = *types.AsDuration(time.Second * 20)
	}
}

type MetricsCenterConfig struct {
	Endpoint      string `env:""`
	ClickHouseDSN string `env:""`
}

type RobotNotifierConfig struct {
	Vendor string   `env:""` // Vendor robot vendor eg: `Lark` `Wechat Work` `DingTalk`
	Env    string   `env:""` // Env Service env, eg: dev-staging, prod
	URL    string   `env:""` // URL webhook url
	Secret string   `env:""` // Secret message secret
	PINs   []string `env:""` // PINs pin someone

	SignFn func(int64) (string, error) `env:"-"`
}

func (c *RobotNotifierConfig) IsZero() bool { return c == nil || c.URL == "" }

func (c *RobotNotifierConfig) Init() {
	if c.Secret != "" {
		c.SignFn = func(ts int64) (string, error) {
			payload := fmt.Sprintf("%v", ts) + "\n" + c.Secret

			var data []byte
			h := hmac.New(sha256.New, []byte(payload))
			_, err := h.Write(data)
			if err != nil {
				return "", err
			}

			signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
			return signature, nil
		}
	}
}

type Risc0Config struct {
	Endpoint        string
	CreateProofPath string
}

func (r *Risc0Config) LivenessCheck() map[string]string {
	m := map[string]string{}

	_, err := http.NewRequest("GET", fmt.Sprintf("http://%s", r.Endpoint), nil)
	if err != nil {
		m[r.Endpoint] = err.Error()
	} else {
		m[r.Endpoint] = "ok"
	}

	return m
}
