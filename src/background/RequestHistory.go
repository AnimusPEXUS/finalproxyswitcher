package main

import (
	"encoding/json"
	"log"
	"net"
	"net/url"
	"strings"
	"syscall/js"
	"time"

	"github.com/AnimusPEXUS/utils/set"
)

type RequestHistoryItem struct {
	Date        time.Time
	TabId       int
	FrameId     int
	RequestId   string
	DocumentURL *string
	Host        string
	URL         string
}

type RequestHistory struct {
	items []*RequestHistoryItem

	// TODO: add mutex
}

func NewRequestHistory() *RequestHistory {
	self := &RequestHistory{}
	return self
}

func (self *RequestHistory) ClearTabHistory(tabId int) {
	for i := len(self.items) - 1; i != -1; i += -1 {
		if self.items[i].TabId == tabId {
			self.items = append(self.items[:i], self.items[i+1:]...)
		}
	}
}

func (self *RequestHistory) AddFromMozillaObject(
	obj js.Value,
	clear_tab_history_on_documentUrl_undefined bool,
) error {

	date := time.Now().UTC()

	req_url := obj.Get("url").String()

	u, err := url.Parse(req_url)
	if err != nil {
		return err
	}

	var doc_url *string

	tabId := obj.Get("tabId").Int()

	if documentUrl := obj.Get("documentUrl"); documentUrl != js.Undefined() {
		t := documentUrl.String()
		doc_url = &t
	} else {
		self.ClearTabHistory(tabId)
	}

	new_item := &RequestHistoryItem{
		Date:        date,
		TabId:       tabId,
		FrameId:     obj.Get("frameId").Int(),
		RequestId:   obj.Get("requestId").String(),
		Host:        strings.ToLower(u.Host),
		URL:         req_url,
		DocumentURL: doc_url,
	}

	b, err := json.MarshalIndent(new_item, "  ", "  ")
	if err != nil {
		return err
	}

	log.Println("new_item", string(b))

	self.items = append(self.items, new_item)

	return nil
}

func (self *RequestHistory) ComputeTabHosts(tabId int) []string {

	ret := set.NewSetString()

	for _, i := range self.items {

		if i.TabId != tabId {
			continue
		}

		h := i.Host

		if parsed_ip := net.ParseIP(h); parsed_ip != nil {
			ret.AddStrings(h)
		} else {

			splitted_host := strings.Split(h, ".")
			splitted_host_len := len(splitted_host)

			for j := 1; j != splitted_host_len+1; j++ {
				host := strings.Join(splitted_host[splitted_host_len-j:], ".")
				ret.AddStrings(host)
			}
		}

	}

	return ret.ListStringsSorted()
}
