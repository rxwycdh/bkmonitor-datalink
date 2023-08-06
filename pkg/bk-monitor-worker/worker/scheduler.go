// Tencent is pleased to support the open source community by making
// 蓝鲸智云 - 监控平台 (BlueKing - Monitor) available.
// Copyright (C) 2022 THL A29 Limited, a Tencent company. All rights reserved.
// Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at http://opensource.org/licenses/MIT
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package worker

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	cron "github.com/robfig/cron/v3"

	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/broker/redis"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/common"
	t "github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/task"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/utils/logger"
	redisUtils "github.com/TencentBlueKing/bkmonitor-datalink/pkg/utils/register/redis"
)

// A Scheduler kicks off tasks at regular intervals based on the user defined schedule.
//
// Schedulers are safe for concurrent use by multiple goroutines.
type Scheduler struct {
	id string

	client          *Client
	rdb             *redis.RDB
	cron            *cron.Cron
	location        *time.Location
	done            chan struct{}
	wg              sync.WaitGroup
	preEnqueueFunc  func(task *t.Task, opts []t.Option)
	postEnqueueFunc func(info *t.TaskInfo, err error)
	errHandler      func(task *t.Task, opts []t.Option, err error)

	// guards idmap
	mu sync.Mutex
	// idmap maps Scheduler's entry ID to cron.EntryID
	// to avoid using cron.EntryID as the public API of
	// the Scheduler.
	idmap map[string]cron.EntryID
}

// SchedulerOpts specifies scheduler options.
type SchedulerOpts struct {
	Location *time.Location

	// enqueue
	PreEnqueueFunc      func(task *t.Task, opts []t.Option)
	PostEnqueueFunc     func(info *t.TaskInfo, err error)
	EnqueueErrorHandler func(task *t.Task, opts []t.Option, err error)
}

// NewScheduler returns a new Scheduler instance given the redis connection option.
// The parameter opts is optional, defaults will be used if opts is set to nil
func NewScheduler(redisOpt redisUtils.Option, opts *SchedulerOpts) (*Scheduler, error) {
	// make a client
	client, err := NewClient()
	if err != nil {
		return nil, err
	}
	// make a rdb
	rdb, err := redis.NewRDB()
	if err != nil {
		return nil, err
	}
	if opts == nil {
		opts = &SchedulerOpts{}
	}

	// 如果不指定，则使用 utc 时间
	loc := opts.Location
	if loc == nil {
		loc = time.UTC
	}

	return &Scheduler{
		id:              generateSchedulerID(),
		client:          client,
		rdb:             rdb,
		cron:            cron.New(cron.WithLocation(loc)),
		location:        loc,
		done:            make(chan struct{}),
		preEnqueueFunc:  opts.PreEnqueueFunc,
		postEnqueueFunc: opts.PostEnqueueFunc,
		errHandler:      opts.EnqueueErrorHandler,
		idmap:           make(map[string]cron.EntryID),
	}, nil
}

func generateSchedulerID() string {
	host, err := os.Hostname()
	if err != nil {
		host = "unknown-host"
	}
	return fmt.Sprintf("%s:%d:%v", host, os.Getpid(), uuid.New())
}

// enqueueJob encapsulates the job of enqueuing a task and recording the event.
type enqueueJob struct {
	id              uuid.UUID
	cronspec        string
	task            *t.Task
	opts            []t.Option
	location        *time.Location
	client          *Client
	rdb             *redis.RDB
	preEnqueueFunc  func(task *t.Task, opts []t.Option)
	postEnqueueFunc func(info *t.TaskInfo, err error)
	errHandler      func(task *t.Task, opts []t.Option, err error)
}

func (j *enqueueJob) Run() {
	if j.preEnqueueFunc != nil {
		j.preEnqueueFunc(j.task, j.opts)
	}
	info, err := j.client.Enqueue(j.task, j.opts...)
	if j.postEnqueueFunc != nil {
		j.postEnqueueFunc(info, err)
	}
	if err != nil {
		if j.errHandler != nil {
			j.errHandler(j.task, j.opts, err)
		}
		return
	}
	logger.Infof("scheduler enqueued a task: %+v", info)
	event := &common.SchedulerEnqueueEvent{
		TaskID:     info.ID,
		EnqueuedAt: time.Now().In(j.location),
	}
	err = j.rdb.RecordSchedulerEnqueueEvent(j.id.String(), event)
	if err != nil {
		logger.Warnf("scheduler could not record enqueue event of enqueued task %s: %v", info.ID, err)
	}
}

// Register registers a task to be enqueued on the given schedule specified by the cronspec.
// It returns an ID of the newly registered entry.
func (s *Scheduler) Register(cronspec string, task *t.Task, opts ...t.Option) (entryID string, err error) {
	job := &enqueueJob{
		id:              uuid.New(),
		cronspec:        cronspec,
		task:            task,
		opts:            opts,
		location:        s.location,
		client:          s.client,
		rdb:             s.rdb,
		preEnqueueFunc:  s.preEnqueueFunc,
		postEnqueueFunc: s.postEnqueueFunc,
		errHandler:      s.errHandler,
	}
	cronID, err := s.cron.AddJob(cronspec, job)
	if err != nil {
		return "", err
	}
	s.mu.Lock()
	s.idmap[job.id.String()] = cronID
	s.mu.Unlock()
	return job.id.String(), nil
}

// Unregister removes a registered entry by entry ID.
// Unregister returns a non-nil error if no entries were found for the given entryID.
func (s *Scheduler) Unregister(entryID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	cronID, ok := s.idmap[entryID]
	if !ok {
		return fmt.Errorf("no scheduler entry found")
	}
	delete(s.idmap, entryID)
	s.cron.Remove(cronID)
	return nil
}

// Run starts the scheduler until an os signal to exit the program is received.
// It returns an error if scheduler is already running or has been shutdown.
func (s *Scheduler) Run() error {
	if err := s.Start(); err != nil {
		return err
	}
	s.Shutdown()
	return nil
}

// Start starts the scheduler.
// It returns an error if the scheduler is already running or has been shutdown.
func (s *Scheduler) Start() error {
	logger.Info("Scheduler starting")
	logger.Infof("Scheduler timezone is set to %v", s.location)
	s.cron.Start()
	s.wg.Add(1)
	go s.runHeartbeater()
	return nil
}

// Shutdown stops and shuts down the scheduler.
func (s *Scheduler) Shutdown() {
	logger.Info("Scheduler shutting down")
	close(s.done) // signal heartbeater to stop
	ctx := s.cron.Stop()
	<-ctx.Done()
	s.wg.Wait()

	s.clearHistory()
	s.client.Close()
	s.rdb.Close()
	logger.Info("Scheduler stopped")
}

func (s *Scheduler) runHeartbeater() {
	defer s.wg.Done()
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-s.done:
			logger.Debugf("Scheduler heatbeater shutting down")
			s.rdb.ClearSchedulerEntries(s.id)
			ticker.Stop()
			return
		case <-ticker.C:
			s.beat()
		}
	}
}

// beat writes a snapshot of entries to redis.
func (s *Scheduler) beat() {
	var entries []*common.SchedulerEntry
	for _, entry := range s.cron.Entries() {
		job := entry.Job.(*enqueueJob)
		e := &common.SchedulerEntry{
			ID:      job.id.String(),
			Spec:    job.cronspec,
			Kind:    job.task.Kind,
			Payload: job.task.Payload,
			Opts:    stringifyOptions(job.opts),
			Next:    entry.Next,
			Prev:    entry.Prev,
		}
		entries = append(entries, e)
	}
	logger.Debugf("Writing entries %v", entries)
	if err := s.rdb.WriteSchedulerEntries(s.id, entries, 5*time.Second); err != nil {
		logger.Warnf("Scheduler could not write heartbeat data: %v", err)
	}
}

func stringifyOptions(opts []t.Option) []string {
	var res []string
	for _, opt := range opts {
		res = append(res, opt.String())
	}
	return res
}

func (s *Scheduler) clearHistory() {
	for _, entry := range s.cron.Entries() {
		job := entry.Job.(*enqueueJob)
		if err := s.rdb.ClearSchedulerHistory(job.id.String()); err != nil {
			logger.Warnf("Could not clear scheduler history for entry %q: %v", job.id.String(), err)
		}
	}
}
