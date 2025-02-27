/*
 *     Copyright 2020 The Dragonfly Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package peer

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/go-http-utils/headers"
	"github.com/golang/mock/gomock"
	testifyassert "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"
	"golang.org/x/time/rate"

	commonv1 "d7y.io/api/v2/pkg/apis/common/v1"
	schedulerv1 "d7y.io/api/v2/pkg/apis/scheduler/v1"

	"d7y.io/dragonfly/v2/client/config"
	"d7y.io/dragonfly/v2/client/daemon/storage"
	"d7y.io/dragonfly/v2/client/daemon/test"
	clientutil "d7y.io/dragonfly/v2/client/util"
	logger "d7y.io/dragonfly/v2/internal/dflog"
	"d7y.io/dragonfly/v2/internal/util"
	"d7y.io/dragonfly/v2/pkg/digest"
	nethttp "d7y.io/dragonfly/v2/pkg/net/http"
	_ "d7y.io/dragonfly/v2/pkg/rpc/dfdaemon/server"
	"d7y.io/dragonfly/v2/pkg/source"
	"d7y.io/dragonfly/v2/pkg/source/clients/httpprotocol"
)

func TestPieceManager_DownloadSource(t *testing.T) {
	assert := testifyassert.New(t)
	ctrl := gomock.NewController(t)
	source.UnRegister("http")
	require.Nil(t, source.Register("http", httpprotocol.NewHTTPSourceClient(), httpprotocol.Adapter))
	defer source.UnRegister("http")
	testBytes, err := os.ReadFile(test.File)
	assert.Nil(err, "load test file")

	var (
		peerID = "peer0"
		taskID = "task0"
		output = "../test/testdata/test.output"
	)

	pieceDownloadTimeout := 30 * time.Second
	storageManager, _ := storage.NewStorageManager(
		config.SimpleLocalTaskStoreStrategy,
		&config.StorageOption{
			DataPath: t.TempDir(),
			TaskExpireTime: clientutil.Duration{
				Duration: -1 * time.Second,
			},
		}, func(request storage.CommonTaskRequest) {}, os.FileMode(0755))

	hash := md5.New()
	hash.Write(testBytes)
	digest := digest.New(digest.AlgorithmMD5, hex.EncodeToString(hash.Sum(nil)[:16]))

	testCases := []struct {
		name               string
		pieceSize          uint32
		withContentLength  bool
		checkDigest        bool
		recordDownloadTime bool
		bandwidth          clientutil.Size
		concurrentOption   *config.ConcurrentOption
	}{
		{
			name:              "multiple pieces with content length, check digest",
			pieceSize:         1024,
			checkDigest:       true,
			withContentLength: true,
		},
		{
			name:               "multiple pieces with content length",
			pieceSize:          1024,
			checkDigest:        false,
			withContentLength:  true,
			recordDownloadTime: true,
		},
		{
			name:               "multiple pieces with content length, concurrent download with 2 goroutines",
			pieceSize:          2048,
			checkDigest:        false,
			withContentLength:  true,
			recordDownloadTime: true,
			concurrentOption: &config.ConcurrentOption{
				GoroutineCount: 2,
				ThresholdSize: clientutil.Size{
					Limit: 1024,
				},
			},
		},
		{
			name:               "multiple pieces with content length, concurrent download with 4 goroutines",
			pieceSize:          2048,
			checkDigest:        false,
			withContentLength:  true,
			recordDownloadTime: true,
			concurrentOption: &config.ConcurrentOption{
				GoroutineCount: 4,
				ThresholdSize: clientutil.Size{
					Limit: 1024,
				},
			},
		},
		{
			name:               "multiple pieces with content length, concurrent download with 8 goroutines",
			pieceSize:          2048,
			checkDigest:        false,
			withContentLength:  true,
			recordDownloadTime: true,
			concurrentOption: &config.ConcurrentOption{
				GoroutineCount: 8,
				ThresholdSize: clientutil.Size{
					Limit: 1024,
				},
			},
		},
		{
			name:               "multiple pieces with content length, concurrent download with 16 goroutines",
			pieceSize:          2048,
			checkDigest:        false,
			withContentLength:  true,
			recordDownloadTime: true,
			concurrentOption: &config.ConcurrentOption{
				GoroutineCount: 16,
				ThresholdSize: clientutil.Size{
					Limit: 1024,
				},
			},
		},
		{
			name:               "multiple pieces with content length, single-thread download with download bandwidth 2KB/s",
			pieceSize:          2048,
			checkDigest:        false,
			withContentLength:  true,
			recordDownloadTime: true,
			bandwidth:          clientutil.Size{Limit: 2048},
		},
		{
			name:               "multiple pieces with content length, concurrent download with 2 goroutines, download bandwidth 2KB/s",
			pieceSize:          2048,
			checkDigest:        false,
			withContentLength:  true,
			recordDownloadTime: true,
			bandwidth:          clientutil.Size{Limit: 2048},
			concurrentOption: &config.ConcurrentOption{
				GoroutineCount: 2,
				ThresholdSize: clientutil.Size{
					Limit: 1024 * 1024,
				},
				ThresholdSpeed: 8192,
			},
		},
		{
			name:               "multiple pieces with content length, concurrent download with 4 goroutines, download bandwidth 2KB/s",
			pieceSize:          2048,
			checkDigest:        false,
			withContentLength:  true,
			recordDownloadTime: true,
			bandwidth:          clientutil.Size{Limit: 2048},
			concurrentOption: &config.ConcurrentOption{
				GoroutineCount: 4,
				ThresholdSize: clientutil.Size{
					Limit: 1024 * 1024,
				},
				ThresholdSpeed: 8192,
			},
		},
		{
			name:               "multiple pieces with content length, concurrent download with 8 goroutines, download bandwidth 2KB/s",
			pieceSize:          2048,
			checkDigest:        false,
			withContentLength:  true,
			recordDownloadTime: true,
			bandwidth:          clientutil.Size{Limit: 2048},
			concurrentOption: &config.ConcurrentOption{
				GoroutineCount: 8,
				ThresholdSize: clientutil.Size{
					Limit: 1024 * 1024,
				},
				ThresholdSpeed: 8192,
			},
		},
		{
			name:               "multiple pieces with content length, concurrent download with 16 goroutines, download bandwidth 2KB/s",
			pieceSize:          2048,
			checkDigest:        false,
			withContentLength:  true,
			recordDownloadTime: true,
			bandwidth:          clientutil.Size{Limit: 2048},
			concurrentOption: &config.ConcurrentOption{
				GoroutineCount: 16,
				ThresholdSize: clientutil.Size{
					Limit: 1024 * 1024,
				},
				ThresholdSpeed: 8192,
			},
		},
		{
			name:               "multiple pieces with content length, single-thread download with download bandwidth 4KB/s",
			pieceSize:          2048,
			checkDigest:        false,
			withContentLength:  true,
			recordDownloadTime: true,
			bandwidth:          clientutil.Size{Limit: 4096},
		},
		{
			name:               "multiple pieces with content length, concurrent download with download bandwidth 4KB/s",
			pieceSize:          2048,
			checkDigest:        false,
			withContentLength:  true,
			recordDownloadTime: true,
			bandwidth:          clientutil.Size{Limit: 4096},
			concurrentOption: &config.ConcurrentOption{
				GoroutineCount: 4,
				ThresholdSize: clientutil.Size{
					Limit: 1024 * 1024,
				},
				ThresholdSpeed: 8192,
			},
		},
		{
			name:               "multiple pieces with content length, single-thread download with download bandwidth 8KB/s",
			pieceSize:          2048,
			checkDigest:        false,
			withContentLength:  true,
			recordDownloadTime: true,
			bandwidth:          clientutil.Size{Limit: 8192},
		},
		{
			name:               "multiple pieces with content length, concurrent download with download bandwidth 8KB/s",
			pieceSize:          2048,
			checkDigest:        false,
			withContentLength:  true,
			recordDownloadTime: true,
			bandwidth:          clientutil.Size{Limit: 8192},
			concurrentOption: &config.ConcurrentOption{
				GoroutineCount: 4,
				ThresholdSize: clientutil.Size{
					Limit: 1024 * 1024,
				},
				ThresholdSpeed: 8192,
			},
		},
		{
			name:               "multiple pieces with content length, single-thread download with download bandwidth 16KB/s",
			pieceSize:          2048,
			checkDigest:        false,
			withContentLength:  true,
			recordDownloadTime: true,
			bandwidth:          clientutil.Size{Limit: 16384},
		},
		{
			name:               "multiple pieces with content length, concurrent download with download bandwidth 16KB/s",
			pieceSize:          2048,
			checkDigest:        false,
			withContentLength:  true,
			recordDownloadTime: true,
			bandwidth:          clientutil.Size{Limit: 16384},
			concurrentOption: &config.ConcurrentOption{
				GoroutineCount: 4,
				ThresholdSize: clientutil.Size{
					Limit: 1024 * 1024,
				},
				ThresholdSpeed: 8192,
			},
		},
		{
			name:              "multiple pieces without content length, check digest",
			pieceSize:         1024,
			checkDigest:       true,
			withContentLength: false,
		},
		{
			name:              "multiple pieces without content length",
			pieceSize:         1024,
			checkDigest:       false,
			withContentLength: false,
		},
		{
			name:              "one pieces with content length case 1",
			pieceSize:         uint32(len(testBytes)),
			withContentLength: true,
		},
		{
			name:              "one pieces without content length case 1",
			pieceSize:         uint32(len(testBytes)),
			withContentLength: false,
		},
		{
			name:              "one pieces with content length case 2",
			pieceSize:         uint32(len(testBytes)) + 1,
			withContentLength: true,
		},
		{
			name:              "one pieces without content length case 2",
			pieceSize:         uint32(len(testBytes)) + 1,
			withContentLength: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			/********** prepare test start **********/
			mockPeerTask := NewMockTask(ctrl)
			var (
				totalPieces = &atomic.Int32{}
				taskStorage storage.TaskStorageDriver
			)
			mockPeerTask.EXPECT().SetContentLength(gomock.Any()).AnyTimes().DoAndReturn(
				func(arg0 int64) error {
					return nil
				})
			mockPeerTask.EXPECT().SetTotalPieces(gomock.Any()).AnyTimes().DoAndReturn(
				func(arg0 int32) {
					totalPieces.Store(arg0)
				})
			mockPeerTask.EXPECT().GetTotalPieces().AnyTimes().DoAndReturn(
				func() int32 {
					return totalPieces.Load()
				})
			mockPeerTask.EXPECT().GetPeerID().AnyTimes().DoAndReturn(
				func() string {
					return peerID
				})
			mockPeerTask.EXPECT().GetTaskID().AnyTimes().DoAndReturn(
				func() string {
					return taskID
				})
			mockPeerTask.EXPECT().GetStorage().AnyTimes().DoAndReturn(
				func() storage.TaskStorageDriver {
					return taskStorage
				})
			mockPeerTask.EXPECT().AddTraffic(gomock.Any()).AnyTimes().DoAndReturn(func(int642 uint64) {})
			mockPeerTask.EXPECT().ReportPieceResult(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().DoAndReturn(
				func(*DownloadPieceRequest, *DownloadPieceResult, error) {

				})
			mockPeerTask.EXPECT().PublishPieceInfo(gomock.Any(), gomock.Any()).AnyTimes().DoAndReturn(
				func(pieceNum int32, size uint32) {

				})
			mockPeerTask.EXPECT().Context().AnyTimes().DoAndReturn(func() context.Context {
				return context.Background()
			})
			mockPeerTask.EXPECT().Log().AnyTimes().DoAndReturn(func() *logger.SugaredLoggerOnWith {
				return logger.With("test case", tc.name)
			})
			taskStorage, err = storageManager.RegisterTask(context.Background(),
				&storage.RegisterTaskRequest{
					PeerTaskMetadata: storage.PeerTaskMetadata{
						PeerID: mockPeerTask.GetPeerID(),
						TaskID: mockPeerTask.GetTaskID(),
					},
					DesiredLocation: output,
					ContentLength:   int64(len(testBytes)),
				})
			assert.Nil(err)
			defer storageManager.CleanUp()
			defer os.Remove(output)
			/********** prepare test end **********/
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tc.withContentLength {
					w.Header().Set(headers.ContentLength,
						fmt.Sprintf("%d", len(testBytes)))
				}
				var (
					buf *bytes.Buffer
					l   int
				)
				rg := r.Header.Get(headers.Range)
				if rg == "" {
					buf = bytes.NewBuffer(testBytes)
					l = len(testBytes)
				} else {
					parsedRange, err := nethttp.ParseRange(rg, int64(len(testBytes)))
					assert.Nil(err)
					w.Header().Set(headers.ContentRange,
						fmt.Sprintf("bytes %d-%d/%d",
							parsedRange[0].Start,
							parsedRange[0].Start+parsedRange[0].Length-1,
							len(testBytes)))
					w.WriteHeader(http.StatusPartialContent)

					buf = bytes.NewBuffer(testBytes[parsedRange[0].Start : parsedRange[0].Start+parsedRange[0].Length])
					l = int(parsedRange[0].Length)
				}
				if tc.bandwidth.Limit > 0 {
					ctx := context.Background()
					limiter := rate.NewLimiter(tc.bandwidth.Limit, int(tc.bandwidth.Limit))
					limiter.AllowN(time.Now(), int(tc.bandwidth.Limit))
					maxPieceNum := util.ComputePieceCount(int64(l), uint32(tc.bandwidth.Limit))
					for pieceNum := int32(0); pieceNum < maxPieceNum; pieceNum++ {
						size := int(tc.bandwidth.Limit)
						offset := uint64(pieceNum) * uint64(size)
						// calculate piece size for last piece
						if int64(offset)+int64(size) > int64(l) {
							size = int(int64(l) - int64(offset))
						}
						err := limiter.WaitN(ctx, size)
						assert.Nil(err)
						n, err := io.CopyN(w, buf, int64(size))
						if err != nil {
							// broken pipe due to auto switch to concurrent download, which will close conn
							return
						}
						assert.Equal(int64(size), n)
					}
				} else {
					n, err := io.Copy(w, buf)
					assert.Nil(err)
					assert.Equal(int64(l), n)
				}
			}))
			defer ts.Close()

			pm, err := NewPieceManager(pieceDownloadTimeout, WithConcurrentOption(tc.concurrentOption))
			assert.Nil(err)
			pm.(*pieceManager).computePieceSize = func(length int64) uint32 {
				return tc.pieceSize
			}

			request := &schedulerv1.PeerTaskRequest{
				Url: ts.URL,
				UrlMeta: &commonv1.UrlMeta{
					Digest: "",
					Range:  "",
					Header: nil,
				},
			}
			if tc.checkDigest {
				request.UrlMeta.Digest = digest.String()
			}
			var start time.Time
			if tc.recordDownloadTime {
				start = time.Now()
			}
			err = pm.DownloadSource(context.Background(), mockPeerTask, request, nil)
			assert.Nil(err)
			if tc.recordDownloadTime {
				elapsed := time.Since(start)
				log := mockPeerTask.Log()
				log.Infof("download took %s", elapsed)
			}

			err = storageManager.Store(context.Background(),
				&storage.StoreRequest{
					CommonTaskRequest: storage.CommonTaskRequest{
						PeerID:      peerID,
						TaskID:      taskID,
						Destination: output,
					},
				})
			assert.Nil(err)

			outputBytes, err := os.ReadFile(output)
			assert.Nil(err, "load output file")
			assert.Equal(testBytes, outputBytes, "output and desired output must match")
		})
	}
}

func TestDetectBackSourceError(t *testing.T) {
	assert := testifyassert.New(t)
	testCases := []struct {
		name              string
		genError          func() error
		detectError       bool
		isBackSourceError bool
	}{
		{
			name: "is back source error - connect error",
			genError: func() error {
				_, err := http.Get("http://127.0.0.1:12345")
				return err
			},
			detectError:       true,
			isBackSourceError: true,
		},
		{
			name: "is back source error - timeout",
			genError: func() error {
				client := http.Client{
					Timeout: 10 * time.Millisecond,
				}
				request, _ := http.NewRequest(http.MethodGet, "http://127.0.0.2:12345", nil)
				_, err := client.Do(request)
				return err
			},
			detectError:       true,
			isBackSourceError: true,
		},
		{
			name: "not back source error",
			genError: func() error {
				return fmt.Errorf("test")
			},
			detectError:       true,
			isBackSourceError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.genError()

			err = detectBackSourceError(err)
			if tc.detectError {
				assert.NotNil(err)
			} else {
				assert.Nil(err)
			}

			assert.Equal(tc.isBackSourceError, isBackSourceError(err))
		})
	}
}
