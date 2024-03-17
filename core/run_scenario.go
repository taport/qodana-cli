/*
 * Copyright 2021-2024 JetBrains s.r.o.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package core

import (
	"github.com/JetBrains/qodana-cli/v2024/platform"
	log "github.com/sirupsen/logrus"
	"strings"
)

const (
	runScenarioDefault      = "default"
	runScenarioFullHistory  = "full-history"
	runScenarioLocalChanges = "local-changes"
	runScenarioScoped       = "scoped"
)

type RunScenario = string

func (o *QodanaOptions) determineRunScenario(hasStartHash bool) RunScenario {
	switch o.ForceIncrementalScript {
	case "":
		isDotNet := Prod.Code == platform.QDNET || strings.Contains(o.Linter, "dotnet")
		switch {
		case o.FullHistory:
			return runScenarioFullHistory
		case hasStartHash && isDotNet:
			return runScenarioScoped
		case hasStartHash && !isDotNet:
			return runScenarioLocalChanges
		default:
			return runScenarioDefault
		}
	case platform.IncrementalScriptLocalChanges:
		if !hasStartHash {
			log.Fatal("Cannot run any incremental script without --diff-start/--commit")
		}
		return runScenarioLocalChanges
	case platform.IncrementalScriptScope:
		if !hasStartHash {
			log.Fatal("Cannot run any incremental script without --diff-start/--commit")
		}
		return runScenarioScoped
	default:
		log.Fatalf("Unknown forced incremental script %s", o.ForceIncrementalScript)
		panic("Unreachable")
	}
}
