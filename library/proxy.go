package library

import (
	"regexp"
	"strings"
	"sync"

	"github.com/dong568789/go-forward/library/util"
)

type Addr struct {
	Listen string
	To     string
}

type Proxies struct {
	List []*Addr
}

func NewProxy() *Proxies {
	proxies := &Proxies{}
	return proxies
}

func (p *Proxies) Run() {
	wg := &sync.WaitGroup{}
	for _, item := range p.List {
		wg.Add(1)
		go ListenAndServer(item.Listen, item.To, wg)
	}
	wg.Wait()
}

func (p *Proxies) ParseAddr(s []string) *Proxies {
	for _, v := range s {
		item := strings.Split(v, " ")
		if len(item) != 2 || !p.checkAddr(item[0]) || !p.checkAddr(item[1]) {
			util.Log().Panic("config fail: %v", v)
		}
		p.List = append(p.List, &Addr{Listen: item[0], To: item[1]})
	}
	return p
}

func (p Proxies) checkAddr(addr string) bool {
	regex := regexp.MustCompile(`^([\d]{1,3}\.){3}[\d]{1,3}:[\d]{2,5}$`)
	if regex == nil {
		util.Log().Panic("ip fail: %v", addr)
	}

	if regex.MatchString(addr) {
		return true
	}
	return false
}
