package main

import (
  "context"
  "flag"
  "fmt"
  "log"
  "os"

  "github.com/ethereum/go-ethereum/common"
  "github.com/ethereum/go-ethereum/core/types"
  "github.com/ethereum/go-ethereum/ethclient"
  "github.com/pkg/errors"
)

const t = ""

var (
  clientDial = flag.String(
    "client_dial", t, "could be websocket or IPC",
  )
)

func etherscan(h common.Hash) string {
  return "https://etherscan.io/tx/" + h.Hex()
}

func program() error {
  flag.Parse()

  handle, err := ethclient.Dial(*clientDial)
  if err != nil {
    return err
  }

  ch := make(chan *types.Header)
  sub, err := handle.SubscribeNewHead(context.Background(), ch)
  for {
    select {
    case e := <-sub.Err():
      log.Fatal(e)
    case h := <-ch:
      block, err := handle.BlockByNumber(context.Background(), h.Number)
      if err != nil {
        log.Fatal(errors.Wrapf(err, "block by hash issue"))
      }
      for _, tx := range block.Transactions() {
        // new contract?
        if t := tx.To(); t == nil && len(tx.Data()) >= 4 {
          hash := tx.Hash()
          fmt.Println("New contract at block: ", h.Number, "-", etherscan(hash));
          _, err := handle.TransactionReceipt(context.Background(), hash)
          if err != nil {
            log.Fatal(err)
          }
        }
      }
    }
  }
  return nil
}

func main() {
  if err := program(); err != nil {
    fmt.Printf("FATAL: %+v\n", err)
    os.Exit(1)
  }
}
