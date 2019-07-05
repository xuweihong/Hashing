package main

import (
	"encoding/binary"
	"fmt"
	"github.com/emirpasic/gods/maps/treemap"
	"hash/fnv"
	"strconv"
)

// 虚拟节点配置
const virtualNode = 5

var nodeList *treemap.Map

// 初始化数据
func init() {
	// 有序map
	nodeList = treemap.NewWithIntComparator()

	// 获取配置
	ips := GetIps()

	// 虚拟节点
	virtualNodeList := make([]string, 0)
	for _, v := range ips {
		for i := 0; i < virtualNode; i++ {
			virtualNodeList = append(virtualNodeList, v+"&&VN"+strconv.Itoa(i))
		}
	}

	// 填充节点数据
	for k, v := range virtualNodeList {
		object := fnv.New32a()
		_, err := object.Write([]byte(v))
		result := binary.LittleEndian.Uint32(object.Sum(nil))
		if err != nil {
			fmt.Print(err)
		}
		fmt.Printf("ip:%s ", v)
		fmt.Printf("node%d:%d\n", k+1, result)
		nodeList.Put(int(result), v)
	}
}

// 配置ip
func GetIps() []string {
	return []string{"127.0.0.1:3306", "127.0.0.2:3306", "127.0.0.3:3306", "127.0.0.4:3306"}
}

func main() {
	key := "username"
	bit := binary.LittleEndian.Uint32([]byte(key))
	node := getNode(int(bit))
	fmt.Printf("key[%d]放入的节点是[%s]\n", bit, node)
}

// 获取节点
func getNode(key int) (node string) {
	for _, k := range nodeList.Keys() {
		if k.(int) >= key {
			result, _ := nodeList.Get(k)
			node = result.(string)
			break
		}
	}

	// 因为是哈希环，当key获取不到节点的时候取第一个节点
	if node == "" {
		result, _ := nodeList.Get(nodeList.Keys()[0].(int))
		node = result.(string)
	}

	return
}
