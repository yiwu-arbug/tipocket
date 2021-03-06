// Copyright 2020 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"flag"

	"github.com/pingcap/tipocket/cmd/util"
	"github.com/pingcap/tipocket/pkg/cluster"
	"github.com/pingcap/tipocket/pkg/control"
	"github.com/pingcap/tipocket/pkg/pocket/config"
	"github.com/pingcap/tipocket/pkg/pocket/creator"
	test_infra "github.com/pingcap/tipocket/pkg/test-infra"
	"github.com/pingcap/tipocket/pkg/test-infra/fixture"
)

var (
	configPath = flag.String("config", "", "config file path")
)

func main() {
	flag.Parse()
	cfg := control.Config{
		Mode:        control.ModeSelfScheduled,
		ClientCount: 1,
		RunTime:     fixture.Context.RunTime,
		RunRound:    1,
	}
	pocketConfig := config.Init()
	pocketConfig.Options.Serialize = false
	pocketConfig.Options.Path = fixture.Context.ABTestConfig.LogPath
	pocketConfig.Options.Concurrency = fixture.Context.ABTestConfig.Concurrency
	pocketConfig.Options.GeneralLog = fixture.Context.ABTestConfig.GeneralLog
	pocketConfig.Options.SyncTimeout.Duration = fixture.Context.BinlogConfig.SyncTimeout
	pocketConfig.Options.EnableHint = fixture.Context.EnableHint
	c := fixture.Context
	suit := util.Suit{
		Config:      &cfg,
		Provisioner: cluster.NewK8sProvisioner(),
		ClientCreator: creator.PocketCreator{
			Config: creator.Config{
				ConfigPath: *configPath,
				Mode:       "binlog",
				Config:     pocketConfig,
			},
		},
		NemesisGens:      util.ParseNemesisGenerators(fixture.Context.Nemesis),
		ClientRequestGen: util.OnClientLoop,
		ClusterDefs:      test_infra.NewBinlogCluster(c.Namespace, c.Namespace, c.TiDBClusterConfig),
	}
	suit.Run(context.Background())
}
