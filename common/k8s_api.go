package common

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Get node ip from metadata -> labels -> kubernetes.io/hostname
func GetNodes(rest_url, label string) ([]string, error) {
	node_ips := []string{}
	tem_lable := strings.Split(label, ":")

	if len(tem_lable) != 2 {
		return nil, fmt.Errorf("nodelabel set error.")
	}
	label_key := tem_lable[0]
	label_val := tem_lable[1]

	res, err := HTTPGet(rest_url)
	if err != nil {
		return nil, err
	}

	var k8data interface{}
	err = json.Unmarshal(res, &k8data)
	if err != nil {
		return nil, err
	}
	m := k8data.(map[string]interface{})

	for _, v := range m {
		switch vv := v.(type) {
		case []interface{}:
			for _, v1 := range vv {
				vvv := v1.(map[string]interface{})
				metadata_if := vvv["metadata"]
				metadata := metadata_if.(map[string]interface{})

				labels_if := metadata["labels"]
				labels := labels_if.(map[string]interface{})
				if labels[label_key] != label_val {
					continue
				}
				node_ips = append(node_ips, labels["kubernetes.io/hostname"].(string))
			}
		}
	}
	if len(node_ips) == 0 {
		return node_ips, fmt.Errorf("Nodes is empty.")
	}
	return node_ips, nil
}

// Get node ip from metadata -> labels -> kubernetes.io/hostname
func GetServicePortFromLabel(rest_url, label string) (float64, error) {
	var service_port float64 = 0
	tem_lable := strings.Split(label, ":")
	if len(tem_lable) != 2 {
		return -1, fmt.Errorf("nodelabel set error.")
	}
	label_key := tem_lable[0]
	label_val := tem_lable[1]

	res, err := HTTPGet(rest_url)
	if err != nil {
		return -1, err
	}

	var k8data interface{}
	err = json.Unmarshal(res, &k8data)
	if err != nil {
		return -1, err
	}
	m := k8data.(map[string]interface{})

	for _, v := range m {
		switch vv := v.(type) {
		case []interface{}:
			for _, v1 := range vv {
				vvv := v1.(map[string]interface{})
				metadata_if := vvv["metadata"]
				metadata := metadata_if.(map[string]interface{})

				labels_if := metadata["labels"]
				labels := labels_if.(map[string]interface{})
				if labels[label_key] != label_val {
					continue
				}

				spec_if := vvv["spec"]
				spce := spec_if.(map[string]interface{})

				ports_if := spce["ports"]
				ports := ports_if.([]interface{})

				nodeport_if := ports[0].(map[string]interface{})
				nodeport := nodeport_if["nodePort"].(float64)
				service_port = nodeport
			}
		}
	}

	return service_port, nil
}

func GetServicePortFromSelector(rest_url, selector string) (float64, error) {
	var service_port float64 = 0
	fmt.Println(selector)
	tem_lable := strings.Split(selector, ":")
	if len(tem_lable) != 2 {
		return -1, fmt.Errorf("Service set error.")
	}
	selector_key := tem_lable[0]
	selector_val := tem_lable[1]

	res, err := HTTPGet(rest_url)
	if err != nil {
		return -1, err
	}

	var k8data interface{}
	err = json.Unmarshal(res, &k8data)
	if err != nil {
		return -1, err
	}
	m := k8data.(map[string]interface{})

	for _, v := range m {
		switch vv := v.(type) {
		case []interface{}:
			for _, v1 := range vv {
				vvv := v1.(map[string]interface{})
				spec_if := vvv["spec"]
				spec := spec_if.(map[string]interface{})

				selector_if := spec["selector"]
				if selector_if == nil {
					continue
				}
				selector := selector_if.(map[string]interface{})
				if selector[selector_key] != selector_val {
					continue
				}

				ports_if := spec["ports"]
				ports := ports_if.([]interface{})

				nodeport_if := ports[0].(map[string]interface{})
				nodeport := nodeport_if["nodePort"].(float64)
				service_port = nodeport
			}
		}
	}
	return service_port, nil
}
