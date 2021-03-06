package eos

import (
	"encoding/json"
	"io/ioutil"

	"github.com/eoscanada/eos-go/ecc"
)

// System contract functions

// NewAccount belongs to a `system` or `chain` package.. it wraps
// certain actions and makes it easy to use.. but doesn't belong
// top-level.. since it's not an API call per se.
//
// NewAccount pushes a `newaccount` transaction on the `eosio
func (api *API) NewAccount(creator, newAccount AccountName, publicKey ecc.PublicKey) (out *PushTransactionFullResp, err error) {
	a := &Action{
		Account: AccountName("eosio"),
		Name:    ActionName("newaccount"),
		Authorization: []PermissionLevel{
			{creator, PermissionName("active")},
		},
		Data: NewActionData(NewAccount{
			Creator: creator,
			Name:    newAccount,
			Owner: Authority{
				Threshold: 1,
				Keys: []KeyWeight{
					KeyWeight{
						PublicKey: publicKey,
						Weight:    1,
					},
				},
			},
			Active: Authority{
				Threshold: 1,
				Keys: []KeyWeight{
					KeyWeight{
						PublicKey: publicKey,
						Weight:    1,
					},
				},
			},
			Recovery: Authority{
				Threshold: 1,
				Accounts: []PermissionLevelWeight{
					PermissionLevelWeight{
						Permission: PermissionLevel{creator, PermissionName("active")},
						Weight:     1,
					},
				},
			},
		}),
	}

	return api.SignPushActions(a)
}

// SetCode applies the given `wasm` file to an account.  Once this is done, the account's code cannot be changed.
//
// EOS.IO Software uses an older version of the WAST file forma
// (breaks with the introduction of
// https://github.com/WebAssembly/wabt/commit/500b617b1c8ea88a2cf46f60205071da9c7569bc)
// so trying to convert .wast to .wasm with standard tooling will
// fail.
//
// Over here, we use the `wasm` file directly.. so it is your
// responsibility to provide a compiled file.
func (api *API) SetCode(forAccount AccountName, wasmPath, abiPath string) (out *PushTransactionFullResp, err error) {
	codeContent, err := ioutil.ReadFile(wasmPath)
	if err != nil {
		return nil, err
	}

	abiContent, err := ioutil.ReadFile(abiPath)
	if err != nil {
		return nil, err
	}

	var abiDef ABI
	if err := json.Unmarshal(abiContent, &abiDef); err != nil {
		return nil, err
	}

	actions := []*Action{
		{
			Account: AccountName("eosio"),
			Name:    ActionName("setcode"),
			Authorization: []PermissionLevel{
				{forAccount, PermissionName("active")},
			},
			Data: NewActionData(SetCode{
				Account:   forAccount,
				VMType:    0,
				VMVersion: 0,
				Code:      HexBytes(codeContent),
			}),
		},
		{
			Account: AccountName("eosio"),
			Name:    ActionName("setabi"),
			Authorization: []PermissionLevel{
				{forAccount, PermissionName("active")},
			},
			Data: NewActionData(SetABI{
				Account: forAccount,
				ABI:     abiDef,
			}),
		},
	}

	return api.SignPushActions(actions...)
}
