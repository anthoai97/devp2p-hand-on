package main

import (
	"encoding/binary"
	"fmt"

	"github.com/anthoai97/devp2p-hand-on/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

// At the start block 0, we deploy contract -> block 2
// the contract address is 0x568B6B552311415cb5a4139324F8d2Ef43F06086
// we call set_new_value -> block 3.
// This means block2's test_int will be 123 and block3's test_int will be 456.
// Our objective is to get storage slot 0 of account 0x568B... at a specific block number from levelDB.

func OpenDatabase(name string, cache, handles int, namespace string, readonly bool) (ethdb.Database, error) {
	var db ethdb.Database
	var err error

	db, err = rawdb.Open(rawdb.OpenOptions{
		Type:      "pebble",
		Directory: "./chaindata",
		Namespace: namespace,
		Cache:     cache,
		Handles:   handles,
		ReadOnly:  readonly,
	})

	return db, err
}

var (
	BLOCK_NUMBER     = uint64(3)
	CONTRACT_ADDRESS = "0xB740431df1aBBe6197dCD9aFDF702C88b918b094"
)

// func (bc *BlockChain) GetBlockByNumber(number uint64) *types.Block {
// 	hash := rawdb.ReadCanonicalHash(bc.db, number)
// 	if hash == (common.Hash{}) {
// 		return nil
// 	}
// 	return bc.GetBlock(hash, number)
// }

// ReadCanonicalHash retrieves the hash assigned to a canonical block number.
// func ReadCanonicalHash(db ethdb.Reader, number uint64) common.Hash {
// 	var data []byte
// 	db.ReadAncients(func(reader ethdb.AncientReaderOp) error {
// 		data, _ = reader.Ancient("hashes", number)
// 		if len(data) == 0 {
// 			// Get it by hash from leveldb
// 			data, _ = db.Get(headerHashKey(number))
// 		}
// 		return nil
// 	})
// 	return common.BytesToHash(data)
// }

// encodeBlockNumber encodes a block number as big endian uint64
func encodeBlockNumber(number uint64) []byte {
	enc := make([]byte, 8)
	binary.BigEndian.PutUint64(enc, number)
	return enc
}

// Get header hashkey for get block hash
func headerHashKey(number uint64) []byte {
	return append(append([]byte("h"), encodeBlockNumber(number)...), []byte("n")...)
}

// headerKey = headerPrefix + num (uint64 big endian) + hash
func headerKey(number uint64, hash common.Hash) []byte {
	return append(append([]byte("h"), encodeBlockNumber(number)...), hash.Bytes()...)
}

// ReadHeaderRLP retrieves a block header in its raw RLP database encoding.
// func ReadHeaderRLP(db ethdb.Reader, hash common.Hash, number uint64) rlp.RawValue {
// 	var data []byte
// 	db.ReadAncients(func(reader ethdb.AncientReaderOp) error {
// 		// First try to look up the data in ancient database. Extra hash
// 		// comparison is necessary since ancient database only maintains
// 		// the canonical data.
// 		data, _ = reader.Ancient("headers", number)
// 		if len(data) > 0 && crypto.Keccak256Hash(data) == hash {
// 			return nil
// 		}
// 		// If not, try reading from leveldb
// 		data, _ = db.Get(headerKey(number, hash))
// 		return nil
// 	})
// 	return data
// }

func main() {
	db, err := OpenDatabase("hello", 0, 0, "", true)
	if err != nil {
		utils.Log.Crit("OpenDatabase fail", "err", err)
	}

	// get block hash from block Number
	// data => Block hash
	var data []byte
	db.ReadAncients(func(aro ethdb.AncientReaderOp) error {
		data, _ = db.Ancient("hashes", uint64(BLOCK_NUMBER))
		if len(data) == 0 {
			// Get it by hash from leveldb
			utils.Log.Info("len(data) == 0")
			data, _ = db.Get(headerHashKey(uint64(BLOCK_NUMBER)))
		}
		utils.Log.Info("ReadAncients data", "data", common.BytesToHash(data))
		return nil
	})

	fmt.Printf("Block hash after ReadAncients data %s\n", common.BytesToHash(data))

	// get block header key from block hash
	blockHeaderKey := headerKey(uint64(BLOCK_NUMBER), common.BytesToHash(data))

	var headerDataRawRPL []byte
	db.ReadAncients(func(aro ethdb.AncientReaderOp) error {
		// First try to look up the data in ancient database. Extra hash
		// comparison is necessary since ancient database only maintains
		// the canonical data.
		headerDataRawRPL, _ = db.Ancient("headers", BLOCK_NUMBER)
		if len(headerDataRawRPL) > 0 && crypto.Keccak256Hash(data) == common.BytesToHash(blockHeaderKey) {
			return nil
		}
		// If not, try reading from leveldb
		headerDataRawRPL, _ = db.Get(blockHeaderKey)
		return nil
	})

	fmt.Printf("Block header Data in RawRPL after ReadAncients from %s\n", common.BytesToHash(headerDataRawRPL))
	fmt.Printf("Length Block header Data in RawRPL after ReadAncients from %d\n", len(headerDataRawRPL))

	// convert to header
	header := new(types.Header)
	if err := rlp.DecodeBytes(headerDataRawRPL, header); err != nil {
		log.Error("Invalid block header RLP", "hash", headerDataRawRPL, "err", err)
		return
	}

	fmt.Printf("Block header stateRootKey %s\n", header.Root)

	// GET trie and traverse to get blockchain data from RootStateKey hash
	stateDB, err := state.New(header.Root, state.NewDatabase(db), nil)
	if err != nil {
		utils.Log.Crit("stateDB error", "err", err.Error())
	}

	storageRoot := stateDB.GetStorageRoot(common.HexToAddress(CONTRACT_ADDRESS))
	fmt.Printf("Contract Storage Root hash %s\n", storageRoot)
	// Get storage rootHash

	obj := stateDB.GetOrNewStateObject(common.HexToAddress(CONTRACT_ADDRESS))
	fmt.Printf("GetOrNewStateObject Root hash %s\n", obj.Root())

	data, err = db.Get(storageRoot.Bytes())
	if err != nil {
		utils.Log.Crit("Get(storageRoot.Bytes()) error", "err", err.Error())
	}

	fmt.Print("data from contract storage Root hash ", data)

}
