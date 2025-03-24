package memberlist

import "encoding/json"

type Meta struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	BindAddr       string `json:"bind_addr"`
	MemberlistPort int    `json:"memberlist_port"`
	GRPCPort       int    `json:"grpc_port"`
	IsCommander    bool   `json:"is_commander"`
}

func MetaFromJSON(data []byte) *Meta {
	meta := &Meta{}
	_ = json.Unmarshal(data, meta)

	return meta
}

func MetaToJSON(meta *Meta) []byte {
	data, _ := json.Marshal(meta)

	return data
}

type MemberListDelegate struct {
	Meta *Meta
}

func (d *MemberListDelegate) NodeMeta(limit int) []byte {
	if d.Meta == nil {
		return nil
	}

	return MetaToJSON(d.Meta)
}

func (d *MemberListDelegate) NotifyMsg(msg []byte) {
}

func (d *MemberListDelegate) GetBroadcasts(overhead, limit int) [][]byte {
	return nil
}

func (d *MemberListDelegate) LocalState(join bool) []byte {
	return nil
}

func (d *MemberListDelegate) MergeRemoteState(buf []byte, join bool) {
}
