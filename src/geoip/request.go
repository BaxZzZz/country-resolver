package geoip

import (
	"container/list"
	"errors"
	"time"
)

// Current provider information
type providerInfo struct {
	lastSwitchTime time.Time
	provider       Provider
	requestCount   uint
}

// GeoIP request
type Request struct {
	limitRequests uint
	timeInterval  time.Duration
	providers     *list.List
	currProvInfo  *providerInfo
}

// Get front provider from list
func (req *Request) getFrontProvider() *providerInfo {
	return req.providers.Front().Value.(*providerInfo)
}

// Get next provider from list
func (req *Request) nextProvider() {
	element := req.providers.Front()
	req.currProvInfo = element.Value.(*providerInfo)
	req.currProvInfo.lastSwitchTime = time.Now()
	req.currProvInfo.requestCount = 0
	req.providers.MoveToBack(element)
}

// Reception of the information about the IP address with ability of even distribution
// of requests between providers.
func (req *Request) GetIPInfo(address string) (*IPInfo, error) {
	if req.currProvInfo == nil {
		req.nextProvider()
	} else if time.Since(req.currProvInfo.lastSwitchTime) > req.timeInterval {
		req.currProvInfo.requestCount = 0
		req.currProvInfo.lastSwitchTime = time.Now()
	} else if req.currProvInfo.requestCount >= req.limitRequests {
		providerInfo := req.getFrontProvider()
		if time.Since(providerInfo.lastSwitchTime) < req.timeInterval {
			return nil, errors.New("GeoIP providers unavailable")
		}
		req.nextProvider()
	}

	info, err := req.currProvInfo.provider.GetIPInfo(address)

	if err != nil {
		return nil, err
	}
	req.currProvInfo.requestCount++
	return info, nil
}

// Creates new request instance
func NewRequest(providers []Provider, limitRequests uint, timeInterval time.Duration) (*Request, error) {
	if len(providers) == 0 {
		return nil, errors.New("Empty GeoIP providers")
	}

	req := &Request{
		limitRequests: limitRequests,
		timeInterval:  timeInterval,
		providers:     list.New(),
	}

	for _, provider := range providers {
		provInfo := &providerInfo{
			provider: provider,
		}
		req.providers.PushBack(provInfo)
	}

	return req, nil
}
