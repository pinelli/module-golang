package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	//raftfastlog "github.com/tidwall/raft-fastlog"
)

func setRaftDir() {

}

func (s *Store) Open(localID string) error {
	// Setup Raft configuration.
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(localID)

	s.m = make(map[string]string)
	// Setup Raft communication.
	addr, err := net.ResolveTCPAddr("tcp", s.RaftBind)
	if err != nil {
		return err
	}

	transport, err := raft.NewTCPTransport(s.RaftBind, addr, 3, 10*time.Second, os.Stderr)
	if err != nil {
		return err
	}

	// Create the snapshot store. This allows the Raft to truncate the log.
	snapshots, err := raft.NewFileSnapshotStore(s.RaftDir, 2, os.Stderr)
	if err != nil {
		return fmt.Errorf("file snapshot store: %s", err)
	}

	// Create the log store and stable store.
	var logStore raft.LogStore
	var stableStore raft.StableStore
	boltDB, err := raftboltdb.NewBoltStore(filepath.Join(s.RaftDir, "raft.db"))
	if err != nil {
		return fmt.Errorf("new bolt store: %s", err)
	}
	logStore = boltDB
	stableStore = boltDB
	//logStore, errlog := raftfastlog.NewFastLogStore("kv/raftfastlog2/log92.log", raftfastlog.Low, os.Stdout)//boltDB
	//stableStore, errStable := raftfastlog.NewFastLogStore("kv/raftfastlog2/stable92.log", raftfastlog.Low, os.Stdout)//boltDB

	// if errlog != nil{
	// 	fmt.Println("LOG ERROR --------------------------", errlog)
	// }
	// if errStable != nil{
	// 	fmt.Println(errStable)
	// }
	// Instantiate the Raft systems.
	ra, err := raft.NewRaft(config, (*fsm)(s), logStore, stableStore, snapshots, transport)
	if err != nil {
		return fmt.Errorf("new raft: %s", err)
	}
	s.raft = ra

	// f := s.raft.AddVoter(raft.ServerID("nodeID1"), raft.ServerAddress("localhost:9090"), 0, 0)
	// if f.Error() == nil {
	// 	fmt.Println("AddNode: ")
	// }else{
	// 	fmt.Println("Node added: ", "addr")
	// }

	// f = s.raft.AddVoter(raft.ServerID("nodeID2"), raft.ServerAddress("localhost:9092"), 0, 0)
	// if f.Error() == nil {
	// 	fmt.Println("AddNode: ")
	// }else{
	// 	fmt.Println("Node added: ", "addr")
	// }

	configuration := raft.Configuration{
		Servers: []raft.Server{
			{
				ID:      config.LocalID,
				Address: transport.LocalAddr(),
			},
		},
	}
	ra.BootstrapCluster(configuration)

	return nil
}

func KeyValue() {

}

func main() {

	server := CreateServer("localhost:8080")
	server.Run()

	select {}
	store := Store{}
	store.RaftDir = "kv/raftfastlog/db/1"
	store.RaftBind = "localhost:9091"

	store.Open("0")

	time.Sleep(6 * time.Second)

	//store.Join("1", "localhost:9090")

	err := store.Set("a", "5")
	// fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>", err)

	val, err := store.Get("a")
	//  //err := store.Delete("a")
	//   //fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>", err)
	fmt.Println(">>>>>>>>>>>>>>", err, ":", val)

	select {}
}
