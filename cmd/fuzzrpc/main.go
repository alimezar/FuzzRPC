// cmd/fuzzrpc/main.go
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	execpkg "github.com/alimezar/FuzzRPC/pkg/exec"
	mutatepkg "github.com/alimezar/FuzzRPC/pkg/mutate"
	reflectpkg "github.com/alimezar/FuzzRPC/pkg/reflect"
	reportpkg "github.com/alimezar/FuzzRPC/pkg/report"
	seedpkg "github.com/alimezar/FuzzRPC/pkg/seed"
)

func main() {
	// CLI flags
	target := flag.String("target", "", "host:port of gRPC server (e.g. localhost:50051)")
	timeout := flag.Duration("timeout", 5*time.Second, "dial timeout")
	reportJSON := flag.String("report-json", "", "path to write JSON report")
	reportHTML := flag.String("report-html", "", "path to write HTML report")
	reportTmpl := flag.String("report-template", "templates/report.html", "HTML template path")
	baselinePath := flag.String("baseline", "", "previous out.json to diff against")
	flag.Parse()

	if *target == "" {
		log.Fatal("❌ --target is required")
	}

	var baseline []reportpkg.Finding
	if *baselinePath != "" {
		if data, err := os.ReadFile(*baselinePath); err == nil {
			_ = json.Unmarshal(data, &baseline) // best-effort load
		} else {
			log.Printf("⚠️  could not read baseline: %v", err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	conn, err := grpc.DialContext(ctx, *target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("❌ failed to connect: %v", err)
	}
	defer conn.Close()

	// Phase 1: Reflection
	services, err := reflectpkg.ListServices(ctx, conn)
	if err != nil {
		log.Fatalf("❌ reflection error: %v", err)
	}

	// collect all findings
	var findings []reportpkg.Finding

	for _, sd := range services {
		svcName := sd.GetFullyQualifiedName()
		fmt.Printf("Service: %s\n", svcName)
		for _, mdesc := range sd.GetMethods() {
			// skip streaming
			if mdesc.IsClientStreaming() || mdesc.IsServerStreaming() {
				fmt.Printf("  Skipping streaming: %s\n", mdesc.GetName())
				continue
			}
			fmt.Printf("  Method: %s\n", mdesc.GetName())

			// Phase 2: Seed generation
			seedMsg := seedpkg.BuildSeed(mdesc.GetInputType())

			// Phase 3: Mutation
			muts, err := mutatepkg.MutateSeed(seedMsg)
			if err != nil {
				log.Printf("    mutation error: %v", err)
				continue
			}

			// Phase 3: Execution & collect
			execpkg.ExecuteFuzz(conn, svcName, mdesc, muts, func(f reportpkg.Finding) {
				findings = append(findings, f)
			})
		}
	}

	reportpkg.ApplyBaseline(findings, baseline)

	// Phase 4: Reporting
	if *reportJSON != "" {
		if err := reportpkg.WriteJSON(findings, *reportJSON); err != nil {
			log.Fatalf("❌ could not write JSON report: %v", err)
		}
		fmt.Printf("✅ JSON report saved to %s\n", *reportJSON)
	}
	if *reportHTML != "" {
		if err := reportpkg.WriteHTML(findings, *reportTmpl, *reportHTML); err != nil {
			log.Fatalf("❌ could not write HTML report: %v", err)
		}
		fmt.Printf("✅ HTML report saved to %s\n", *reportHTML)
	}
}
