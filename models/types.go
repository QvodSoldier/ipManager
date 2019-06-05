package models

//######################################################################
// 接收NSX-T返回Subnet的数据
// GET /api/v1/pools/ip-pools
type OriginalData struct {
	Results []Result
}

type Result struct {
	Pool_usage   Pool_usage
	Subnets      []OldSubnets
	Tags         []map[string]string
	Display_name string
	Id           string
}

type Pool_usage struct {
	Total_ids     int
	Allocated_ids int
	Free_ids      int
}

type OldSubnets struct {
	Cidr string
	// Allocation_ranges []map[string]string
}

// 处理后返回给接口的数据
type NewData struct {
	Subnets []NewSubnet
}

type NewSubnet struct {
	SubnetName  string
	NameSpace   string
	Cidr        string
	TotalIp     int
	FreeIp      int
	AllocatedIp int
	Id          string
	// IpStart     string
	// IpEnd       string
}

//######################################################################
type NewData2 struct {
	Subnets []NewSubnet2
}

type NewSubnet2 struct {
	SubnetName string
	Cidr       string
}

//######################################################################
// OSV Old Subnet view
type OSV struct {
	Results []LogicalSwitchPort
}

type LogicalSwitchPort struct {
	Address_bindings []Podip
	Display_name     string
}

type Podip struct {
	Ip_address string
}

// NSV new Subnet view
type NSV struct {
	SubnetViews []SubnetView
}

type SubnetView struct {
	Ip      string
	Status  bool
	PodName string
}

//######################################################################
//LSS means logical switches
type LSS struct {
	Results []LS
}

// LS means a logical switche
type LS struct {
	Ip_pool_id string
	Id         string
}

// ####################################################################
type Nss struct {
	NsSubnets []NsSubnet
}

type NsSubnet struct {
	SubnetName string
	Cidr       string
}

// 判断IP是否已被分配
func Judge(ip string, s []LogicalSwitchPort) (string, bool) {
	for _, v := range s {
		if v.Address_bindings[0].Ip_address == ip {
			return v.Display_name, true
			break
		}
	}
	return "", false
}
