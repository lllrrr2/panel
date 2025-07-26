package ntp

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/beevik/ntp"

	"github.com/tnborg/panel/pkg/shell"
)

var ErrNotReachable = errors.New("failed to reach NTP server")

var ErrNoAvailableServer = errors.New("no available NTP server found")

var defaultAddresses = []string{
	//"ntp.ntsc.ac.cn",      // 中科院国家授时中心的服务器很快，但是多刷几次就会被封
	"ntp.aliyun.com",   // 阿里云
	"ntp1.aliyun.com",  // 阿里云2
	"ntp.tencent.com",  // 腾讯云
	"time.windows.com", // Windows
	"time.apple.com",   // Apple
}

func Now(address ...string) (time.Time, error) {
	if len(address) > 0 {
		if now, err := ntp.Time(address[0]); err != nil {
			return time.Now(), fmt.Errorf("%w: %s", ErrNotReachable, err)
		} else {
			return now, nil
		}
	}

	best, err := bestServer(defaultAddresses...)
	if err != nil {
		return time.Now(), err
	}

	now, err := ntp.Time(best)
	if err != nil {
		return time.Now(), fmt.Errorf("%w: %s", ErrNotReachable, err)
	}

	return now, nil
}

func UpdateSystemTime(t time.Time) error {
	_, err := shell.Execf(`date -s '%s'`, t.Format(time.DateTime))
	return err
}

func UpdateSystemTimeZone(tz string) error {
	_, err := shell.Execf(`timedatectl set-timezone '%s'`, tz)
	return err
}

// pingServer 计算NTP服务器的延迟
func pingServer(addr string) (time.Duration, error) {
	options := ntp.QueryOptions{Timeout: 1 * time.Second}
	response, err := ntp.QueryWithOptions(addr, options)
	if err != nil {
		return 0, err
	}

	return response.RTT, nil
}

// bestServer 返回延迟最低的NTP服务器
func bestServer(addresses ...string) (string, error) {
	if len(addresses) == 0 {
		addresses = defaultAddresses
	}

	type ntpResult struct {
		address string
		delay   time.Duration
		err     error
	}

	results := make(chan ntpResult, len(addresses))
	var wg sync.WaitGroup

	for _, addr := range addresses {
		wg.Add(1)
		go func(addr string) {
			defer wg.Done()

			delay, err := pingServer(addr)
			results <- ntpResult{address: addr, delay: delay, err: err}
		}(addr)
	}

	wg.Wait()
	close(results)

	var bestAddr string
	var bestDelay time.Duration
	found := false

	for result := range results {
		if result.err != nil {
			continue
		}

		if !found || result.delay < bestDelay {
			bestAddr = result.address
			bestDelay = result.delay
			found = true
		}
	}

	if !found {
		return "", ErrNoAvailableServer
	}

	return bestAddr, nil
}
