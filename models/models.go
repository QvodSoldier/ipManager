package models

import (
	"encoding/json"
	"fmt"
	"ipManager/config"
	"ipManager/tools"
	"log"
	"strconv"
	"strings"
)

type Models interface {
	// 从NSX-T获取subnet的原始JSON
	GetOriginalSubnets() []byte

	// 处理从NSX-T返回的原始JSON，并返回一个新subnet JSON
	HandleSubnets(data []byte) ([]byte, error)

	// 从NSX-T查询一个Subnet里IP的使用情况
	GetOriginalSubnetView(sid string) []byte

	// 处理从NSX-T返回的原始JSON，并返回一个新json
	HandleSubnetView(data []byte, Cidr string) ([]byte, error)

	//处理从NSX-T返回的原始JSON，并返回一个通过命名空间过滤后的json
	NsHandleSubnets(data []byte, ns string) ([]byte, error)
}

var cfg = config.Cfg

type Client struct{}

func (this *Client) GetOriginalSubnets() []byte {
	// 从NSX-T获取所有的IP池信息
	req := tools.NewRequest()
	req.Url = "https://" + cfg.Nsxmanager[0]["host1"] + "/api/v1/pools/ip-pools"
	data := req.ClientGet()
	return data
}

func (this *Client) GetOriginalSubnetView(sid string) []byte {
	// 从NSX-T获取该Subnet对应的逻辑交换机ID
	var LsId string
	LSlist, _ := HandleLS(GetOriginalLS())
	for _, v := range LSlist {
		if v.Ip_pool_id == sid {
			LsId = v.Id
			break
		}
	}
	// 从NSX-T获取该交换机下所有的容器端口
	req := tools.NewRequest()
	req.Url = "https://" + cfg.Nsxmanager[0]["host1"] + "/api/v1/logical-ports?container_ports_only=true&logical_switch_id=" + LsId
	data := req.ClientGet()
	return data
}

func (this *Client) HandleSubnets(data []byte) ([]byte, error) {
	// 从所有的IP池中，筛选出容器网络的Subnet，并组织数据返回
	var originaldata OriginalData
	var newdata NewData

	err := json.Unmarshal(data, &originaldata)
	if err != nil {
		return nil, err
	}

	for _, v := range originaldata.Results {
		if v.Tags != nil && v.Tags[0]["scope"] == "ncp/subnet" {
			// 获取命名空间名字
			tmpV := strings.FieldsFunc(v.Display_name, func(c rune) bool {
				return c == '-'
			})
			tmpV = tmpV[:len(tmpV)-1]
			nsn := strings.Replace(strings.Join(tmpV, " "), " ", "-", -1)
			// 传入一个新Subnet
			var Subnet = NewSubnet{
				SubnetName:  v.Display_name,
				NameSpace:   nsn,
				Cidr:        v.Subnets[0].Cidr,
				TotalIp:     v.Pool_usage.Total_ids,
				FreeIp:      v.Pool_usage.Free_ids,
				AllocatedIp: v.Pool_usage.Allocated_ids,
				Id:          v.Id,
				// IpStart:     v.Subnets[0].Allocation_ranges[0]["start"],
				// IpEnd:       v.Subnets[0].Allocation_ranges[0]["end"],
			}
			newdata.Subnets = append(newdata.Subnets, Subnet)
		}
	}
	// json序列化并返回
	data2, err := json.Marshal(newdata)
	if err != nil {
		return nil, err
	}
	return data2, nil
}

func (this *Client) HandleSubnetView(data []byte, Cidr string) ([]byte, error) {
	// 接受并筛选该子网的逻辑交换机下端口数据
	var osv OSV
	err := json.Unmarshal(data, &osv)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// 获取Cidr完整的Ip列表-iplist
	cidr := tools.NewCidr(Cidr)
	MaxIp := cidr.Ip2Long(cidr.GetCidrIpRange().Max)
	MinIp := cidr.Ip2Long(cidr.GetCidrIpRange().Min)
	Count, _ := strconv.Atoi(cidr.GetCidrHostNum().Count)
	iplist := make([]string, 0, Count)
	for i := MinIp; i < MaxIp; i++ {
		i := int64(i)
		iplist = append(iplist, cidr.BacktoIP4(i))
	}

	var nsv NSV
	var sview SubnetView
	for _, v := range iplist {
		// 判断Ip是否已被Pod分配,并返回Pod名字
		Pn, j := Judge(v, osv.Results)
		if j == true {
			sview = SubnetView{
				Ip:      v,
				Status:  true,
				PodName: Pn,
			}
		} else {
			sview = SubnetView{
				Ip:      v,
				Status:  false,
				PodName: "",
			}
		}
		// 生成完整的IP列表及使用状态
		nsv.SubnetViews = append(nsv.SubnetViews, sview)
	}
	// json序列化并返回
	Svs, err := json.Marshal(nsv)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return Svs, nil
}

func (this *Client) NsHandleSubnets(data []byte, ns string) ([]byte, error) {
	// 从所有的IP池中，筛选出容器网络的Subnet，并组织数据返回
	var originaldata OriginalData
	var newdata NewData2

	err := json.Unmarshal(data, &originaldata)
	if err != nil {
		return nil, err
	}

	for _, v := range originaldata.Results {
		if v.Tags != nil && v.Tags[0]["scope"] == "ncp/subnet" {
			// 获取命名空间名字
			tmpV := strings.FieldsFunc(v.Display_name, func(c rune) bool {
				return c == '-'
			})
			tmpV = tmpV[:len(tmpV)-1]
			nsn := strings.Replace(strings.Join(tmpV, " "), " ", "-", -1)
			if nsn == ns {
				var Subnet = NewSubnet2{
					SubnetName: v.Display_name,
					Cidr:       v.Subnets[0].Cidr,
				}
				newdata.Subnets = append(newdata.Subnets, Subnet)
			}
		}
	}
	data2, err := json.Marshal(newdata)
	if err != nil {
		return nil, err
	}
	return data2, nil
}

// 从NSX-T获得一个逻辑交换机列表
// GetOriginalLS() []byte
func GetOriginalLS() []byte {
	// Get subnets from nsx-t
	req := tools.NewRequest()
	req.Url = "https://" + cfg.Nsxmanager[0]["host1"] + "/api/v1/logical-switches"

	data := req.ClientGet()
	fmt.Println(data)
	return data
}

// 处理从NSX-T返回的逻辑交换机列表
// HandleLS() []byte
func HandleLS(data []byte) ([]LS, error) {
	var tmp LSS
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return tmp.Results, nil
}
