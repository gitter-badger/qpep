package api

import (
	"fmt"
	"log"
	"strings"
	"sync"
)

const (
	TOTAL_CONNECTIONS string = "connections"   // total connections open on the server at this time
	PERF_CONN         string = "perf-conn"     // number of current connections for a particular client
	PERF_UP_SPEED     string = "perf-upspeed"  // current upload speed for a particular client
	PERF_DW_SPEED     string = "perf-dwspeed"  // current download speed for a particular client
	PERF_UP_TOTAL     string = "perf-uptotal"  // total number of bytes uploaded by a particular client
	PERF_DW_TOTAL     string = "perf-dwtotal"  // total number of bytes downloaded by a particular client
	INFO_PLATFORM     string = "info-platform" // platform used by the client, as communicated in api echo
	INFO_UPDATE       string = "info-update"   // last time server received an echo from the client
)

var Statistics = &statistics{}

func init() {
	Statistics.Reset()
}

type statistics struct {
	semCounters *sync.RWMutex
	semState    *sync.RWMutex

	counters map[string]float64
	state    map[string]string
	hosts    []string
}

func (s *statistics) init() {
	if s.semCounters != nil && s.semState != nil {
		return
	}

	s.semCounters = &sync.RWMutex{}
	s.semState = &sync.RWMutex{}
	s.hosts = make([]string, 0, 32)
}

func (s *statistics) Reset() {
	s.semCounters = nil
	s.semState = nil
	s.init()

	log.Println("Server statistics reset.")
	s.counters = make(map[string]float64)
	s.state = make(map[string]string)
}

func (s *statistics) AsKey(prefix string, values ...string) string {
	return strings.ToLower(prefix + "-" + strings.Join(values, "-"))
}

// ---- Counters ---- //
func (s *statistics) GetCounter(prefix string, keyparts ...string) float64 {
	key := s.AsKey(prefix, keyparts...)
	if len(key) == 0 {
		return -1
	}

	s.init()
	s.semCounters.RLock()
	defer s.semCounters.RUnlock()

	if val, ok := s.counters[key]; ok {
		return val
	}
	return -1
}

func (s *statistics) SetCounter(value float64, prefix string, keyparts ...string) float64 {
	key := s.AsKey(prefix, keyparts...)
	if len(key) == 0 {
		return -1
	}
	if value < 0 {
		panic(fmt.Sprintf("Will not track negative values: (%s: %f)", key, value))
	}

	s.init()
	s.semCounters.Lock()
	defer s.semCounters.Unlock()

	s.counters[key] = value
	return value
}

func (s *statistics) IncrementCounter(prefix string, keyparts ...string) float64 {

	key := s.AsKey(prefix, keyparts...)
	if len(key) == 0 {
		log.Printf("counter: %s = -1\n", key)
		return -1
	}

	s.init()
	s.semCounters.Lock()
	defer s.semCounters.Unlock()

	value, ok := s.counters[key]
	if !ok {
		log.Printf("counter: %s = *%d\n", key, 1)
		s.counters[key] = 1
		return 1
	}

	log.Printf("counter: %s = %d\n", key, value+1)
	s.counters[key] = value + 1
	return value + 1.0
}

func (s *statistics) DecrementCounter(prefix string, keyparts ...string) float64 {
	key := s.AsKey(prefix, keyparts...)
	if len(key) == 0 {
		return -1.0
	}

	s.init()
	s.semCounters.Lock()
	defer s.semCounters.Unlock()

	value, ok := s.counters[key]
	if !ok || value-1 <= 0.0 {
		s.counters[key] = 0.0
		return 0.0
	}

	log.Printf("counter: %s = %d\n", key, value-1.0)
	s.counters[key] = value - 1
	return value - 1.0
}

// ---- State ---- //
func (s *statistics) GetState(prefix string, keyparts ...string) string {
	key := s.AsKey(prefix, keyparts...)
	if len(key) == 0 {
		return ""
	}

	s.init()
	s.semState.RLock()
	defer s.semState.RUnlock()

	if val, ok := s.state[key]; ok {
		return val
	}
	return ""
}

func (s *statistics) SetState(value string, prefix string, keyparts ...string) string {
	key := s.AsKey(prefix, keyparts...)
	if len(key) == 0 {
		return ""
	}

	s.init()
	s.semState.Lock()
	defer s.semState.Unlock()

	s.state[key] = value
	return value
}

// ---- address mapping ---- //
func (s *statistics) GetMappedAddress(source string) string {
	s.init()
	s.semState.RLock()
	defer s.semState.RUnlock()

	if val, ok := s.state[source]; ok {
		return val
	}
	return ""
}

func (s *statistics) SetMappedAddress(source string, dest string) {
	s.init()
	s.semState.Lock()
	defer s.semState.Unlock()

	if _, ok := s.state[source]; !ok {
		s.hosts = append(s.hosts, dest)
	}
	s.state[source] = dest
}

func (s *statistics) DeleteMappedAddress(source string) {
	s.init()
	s.semState.Lock()
	defer s.semState.Unlock()

	if _, ok := s.state[source]; ok {
		mapped := s.state[source]
		for i := 0; i < len(s.hosts); i++ {
			if !strings.EqualFold(s.hosts[i], mapped) {
				continue
			}
			s.hosts = append(s.hosts[:i], s.hosts[i+1:]...)
			break
		}
	}
	delete(s.state, source)
}

// ---- hosts ---- //
func (s *statistics) GetHosts() []string {
	s.init()
	s.semState.RLock()
	defer s.semState.RUnlock()

	v := append([]string{}, "127.0.0.1") // for test
	v = append(v, s.hosts...)
	return v
}
