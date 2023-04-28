/*
 Copyright © 2023 MicroOps-cn.

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

package service

import (
	"context"
	"fmt"
	"time"

	"github.com/MicroOps-cn/idas/config"
	"github.com/MicroOps-cn/idas/pkg/service/gormservice"
	"github.com/MicroOps-cn/idas/pkg/service/models"
)

type LoggingService interface {
	migrator
	PostEventLog(ctx context.Context, eventId, userId, username, clientIP, loc, action, message string, status bool, took time.Duration, log ...interface{}) error
	GetEvents(ctx context.Context, filters map[string]string, keywords string, startTime time.Time, endTime time.Time, current int64, size int64) (count int64, event []*models.Event, err error)
	GetEventLogs(ctx context.Context, filters map[string]string, keywords string, current int64, size int64) (count int64, event []*models.EventLog, err error)
}

func NewLoggingService(ctx context.Context) LoggingService {
	var loggingService LoggingService
	loggingStorage := config.Get().GetStorage().GetLogging()
	switch loggingSource := loggingStorage.GetStorageSource().(type) {
	case *config.Storage_Mysql:
		loggingService = gormservice.NewLoggingService(ctx, loggingStorage.Name, loggingSource.Mysql.Client)
	case *config.Storage_Sqlite:
		loggingService = gormservice.NewLoggingService(ctx, loggingStorage.Name, loggingSource.Sqlite.Client)
	default:
		panic(fmt.Sprintf("failed to initialize LoggingService: unknown data source: %T", loggingSource))
	}
	return loggingService
}
