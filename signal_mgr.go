package util

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"hotpu.cn/xkefu/common/log"
)

// A SignalHandler is an handler of signal receive
// return {1} to quit after handler execution
type SignalHandler func(s os.Signal) int

func defaultSigHandler(s os.Signal) int {
	log.Info("common.util", "defaultSigHandler", "defaultSigHandler invoke ...")
	return 1
}

const quitMark = 1

// A SignalMgr is an basic wrap for system signal control
// In default app quit when receive signal [syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSTOP]
// It is a singleton implementation for defaulted
type SignalMgr struct {
	funcMap map[os.Signal]SignalHandler
}

// The singleton instance variable
var singletonMgr *SignalMgr

// Constructor
func NewSignalMgr() *SignalMgr {
	if singletonMgr == nil {
		singletonMgr = new(SignalMgr)

		// init
		singletonMgr.funcMap = make(map[os.Signal]SignalHandler)
	}

	return singletonMgr
}

func (s *SignalMgr) Register(sig os.Signal, handler SignalHandler) {
	if _, found := s.funcMap[sig]; !found {
		s.funcMap[sig] = handler
	}
}

func (s *SignalMgr) checkDefault() {
	// get empty func map, auto register default signals
	if len(s.funcMap) == 0 {
		s.Register(syscall.SIGHUP, defaultSigHandler)
		s.Register(syscall.SIGQUIT, defaultSigHandler)
		s.Register(syscall.SIGTERM, defaultSigHandler)
		s.Register(syscall.SIGINT, defaultSigHandler)
		s.Register(syscall.SIGSTOP, defaultSigHandler)
	}
}

func (s *SignalMgr) Run() {
	s.checkDefault()

	c := make(chan os.Signal, 1)
	for sig := range s.funcMap {
		signal.Notify(c, sig)
	}

	for {
		recvSig := <-c
		log.Warn("common.util", "SignalMgr.Run", "Run get signal:%v", recvSig.String())

		mark, err := s.handle(recvSig)
		if mark == quitMark {
			if err != nil {
				log.Error("common.util", "SignalMgr.Run", "Run quit handler get error:%v with signal:%v", err, recvSig)
			}
			return
		}
	}
}

func (s *SignalMgr) handle(sig os.Signal) (int, error) {
	if _, found := s.funcMap[sig]; found {
		ret := s.funcMap[sig](sig)
		return ret, nil
	} else {
		return quitMark, fmt.Errorf("no handler for signal %v", sig)
	}
	panic("cannot be here")
}
