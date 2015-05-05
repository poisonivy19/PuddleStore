package puddlestore

import (
	"fmt"
	"net/rpc"
)

var connMap = make(map[string]*rpc.Client)

type ConnectRequest struct {
	FromNode PuddleAddr
}

type ConnectReply struct {
	Ok bool
}

func ConnectRPC(remotenode *PuddleAddr, request ConnectRequest) (*ConnectReply, error) {
	fmt.Println("(Puddlestore) RPC Connect")
	var reply ConnectReply

	err := makeRemoteCall(remotenode, "ConnectImpl", request, &reply)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &reply, nil
}

type lsRequest struct {
	curdir   string
	FromNode PuddleAddr
}

type lsReply struct {
	elements []string
	Ok       bool
}

func lsRPC(remotenode *PuddleAddr, request lsRequest) (*lsReply, error) {
	var reply lsReply

	err := makeRemoteCall(remotenode, "lsImpl", request, &reply)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &reply, nil
}

type cdRequest struct {
	curdir   string
	FromNode PuddleAddr
}

type cdReply struct {
	Ok bool
}

func cdRPC(remotenode *PuddleAddr, request lsRequest) (*cdReply, error) {
	var reply cdReply

	err := makeRemoteCall(remotenode, "cdImpl", request, &reply)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &reply, nil
}

/* Helper function to make a call to a remote node */
func makeRemoteCall(remoteNode *PuddleAddr, method string, req interface{}, rsp interface{}) error {
	// Dial the server if we don't already have a connection to it
	remoteNodeAddrStr := remoteNode.Addr
	var err error
	client, ok := connMap[remoteNodeAddrStr]
	if !ok {
		client, err = rpc.Dial("tcp", remoteNode.Addr)
		if err != nil {
			return err
		}
		connMap[remoteNodeAddrStr] = client
	}

	// Make the request
	uniqueMethodName := fmt.Sprintf("%v.%v", remoteNodeAddrStr, method)
	err = client.Call(uniqueMethodName, req, rsp)
	if err != nil {
		client.Close()
		delete(connMap, remoteNodeAddrStr)
		return err
	}

	return nil
}
