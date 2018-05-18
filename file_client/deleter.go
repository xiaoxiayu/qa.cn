package file_client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func (this *Server) Delete(server_path string) {
	// "!!!DO NOT DELETE NOW!!!"
	return
	data := make(url.Values)
	data["deletefile"] = []string{"/" + server_path}
	target_url := this.Url + "/normal_delete"

	res, err := http.PostForm(target_url, data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer res.Body.Close()
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s\n", result)
}

func (this *Server) DeleteWithOutDB(server_path string) {
	// "!!!DO NOT DELETE NOW!!!"
	return
	data := make(url.Values)
	data["deletefile"] = []string{"/" + server_path}
	data["option"] = []string{"nodb"}
	target_url := this.Url + "/normal_delete"

	res, err := http.PostForm(target_url, data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer res.Body.Close()
	_, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
