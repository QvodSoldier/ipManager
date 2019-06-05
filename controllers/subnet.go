package controllers

import (
	"ipManager/models"
	"log"
	"net/http"
)

type SubnetsController struct{}

func (this *SubnetsController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	client := new(models.Client)

	subnets, err := client.HandleSubnets(client.GetOriginalSubnets())
	if err != nil {
		log.Fatal(err)
	}
	w.Write(subnets)
}

func SubnetView(w http.ResponseWriter, r *http.Request) {
	client := new(models.Client)
	tmp := r.URL.Query()

	Sid := tmp.Get("Id")
	Cidr := r.FormValue("cidr")

	if Sid == "" {
		http.Redirect(w, r, "/nsxt/subnets", 301)
		return
	}

	oldnsxidata := client.GetOriginalSubnetView(Sid)
	subnetview, err := client.HandleSubnetView(oldnsxidata, Cidr)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(subnetview)
}

func NsSubnetView(w http.ResponseWriter, r *http.Request) {
	client := new(models.Client)
	ns := r.FormValue("namespace")

	subnets, err := client.NsHandleSubnets(client.GetOriginalSubnets(), ns)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(subnets)
}
