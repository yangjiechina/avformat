package libhls

import (
	"fmt"
	"io/ioutil"
	"net/http"
	url2 "net/url"
	"strings"
	"time"
)

type Puller struct {
	url         string
	prefix      string
	client      *http.Client
	lastList    map[string]byte
	segmentPool []string
}

func NewPuller() *Puller {
	return &Puller{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		lastList: make(map[string]byte, 64),
	}
}

func (p *Puller) doGet(url string) ([]byte, error) {
	resp, err := p.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return all, nil
}

func (p *Puller) downloadSegment(entry string) ([]byte, error) {
	return p.doGet(p.prefix + entry)
}

func (p *Puller) getPlayList(url string) ([]string, bool, error) {
	data, err := p.doGet(url)
	if err != nil {
		return nil, false, err
	}

	entries := strings.Split(string(data), "\n#")
	var master bool
	var masterPlayList []string
	var mediaPlayList []string
	for _, entry := range entries {
		split := strings.Split(strings.TrimSuffix(entry, "\n"), "\n")
		if len(split) < 2 {
			continue
		}

		if strings.HasPrefix(split[0], "EXT-X-STREAM-INF") {
			master = true
			masterPlayList = append(masterPlayList, split[1:]...)
		} else if strings.HasPrefix(split[0], "EXTINF") {
			mediaPlayList = append(mediaPlayList, split[1:]...)
		}
	}

	if master {
		return masterPlayList, true, nil
	} else {
		return mediaPlayList, false, nil
	}
}

func (p *Puller) Open(url string) error {
	if _, err := url2.Parse(url); err != nil {
		return err
	}

	p.url = url
	index := strings.LastIndex(url, "/")
	p.prefix = url[:index+1]

	_, err := p.doGet(url)
	return err
}

func (p *Puller) Read() ([]byte, error) {
	if len(p.segmentPool) <= 0 {
		list, master, err := p.getPlayList(p.url)
		if master {
			//for _, s := range list {
			//}
			list, master, err = p.getPlayList(p.prefix + list[0])
		}
		if err != nil {
			return nil, err
		}
		if len(list) == 0 {
			return nil, fmt.Errorf("failed to pull stream")
		}

		for _, s := range list {
			//去重
			if _, ok := p.lastList[s]; !ok {
				p.segmentPool = append(p.segmentPool, s)
			}
		}

		p.lastList = map[string]byte{}
		for _, s := range list {
			p.lastList[s] = 0
		}
		if len(p.segmentPool) == 0 {
			return nil, nil
		}
	}

	segment, err := p.downloadSegment(p.segmentPool[0])
	if err != nil {
		return nil, err
	}

	p.segmentPool = p.segmentPool[1:]
	return segment, err
}

func (p *Puller) Close() {
	p.lastList = map[string]byte{}
	p.segmentPool = p.segmentPool[0:0]
}
