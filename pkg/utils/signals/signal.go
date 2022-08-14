/*
 Copyright © 2022 MicroOps-cn.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package signals

import (
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type StopChan struct {
	stopCh chan struct{}
	reqWg  sync.WaitGroup
	wg     sync.WaitGroup
}

var once = sync.Once{}

func (s *StopChan) WaitRequest() {
	s.reqWg.Wait()
}

func (s *StopChan) DoneRequest() {
	s.reqWg.Done()
}

func (s *StopChan) AddRequest(delta int) {
	s.reqWg.Add(delta)
}

func (s *StopChan) Wait() {
	s.wg.Wait()
}

func (s *StopChan) Done() {
	s.wg.Done()
}

func (s *StopChan) Add(delta int) {
	s.wg.Add(delta)
}

func (s *StopChan) Channel() <-chan struct{} {
	return s.stopCh
}

var stopChan *StopChan

func SetupSignalHandler(logger log.Logger) (stopCh *StopChan) {
	once.Do(func() {
		onlyOneSignalHandler := make(chan struct{})
		close(onlyOneSignalHandler) // panics when called twice
		stopChan = &StopChan{
			stopCh: make(chan struct{}),
		}
		c := make(chan os.Signal, 2)
		signal.Notify(c, shutdownSignals...)

		go func() {
			sig := <-c
			level.Info(logger).Log("msg", fmt.Sprintf("收到信号[%s],进程停止\n", sig))
			close(stopChan.stopCh)
			stopChan.WaitRequest()
			stopChan.Wait()
			os.Exit(1) // second signal. Exit directly.
		}()
	})
	return stopChan
}
