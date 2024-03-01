// Package services contains an implementation of the grpc MangleServer interface.
package service

import (
	"io"
	"log"
	"strings"

	"github.com/google/mangle/analysis"
	"github.com/google/mangle/ast"
	"github.com/google/mangle/engine"
	"github.com/google/mangle/factstore"
	"github.com/google/mangle/parse"

	pb "github.com/burakemir/mangle-service/proto"
)

var programInfo *analysis.ProgramInfo

func copyDecl(decls map[ast.PredicateSym]*ast.Decl) map[ast.PredicateSym]ast.Decl {
	m := make(map[ast.PredicateSym]ast.Decl, len(decls))
	for k, v := range decls {
		m[k] = *v
	}
	return m
}

type MangleService struct {
	store factstore.FactStore
	pb.UnimplementedMangleServer
}

func New() *MangleService {
	return &MangleService{store: factstore.NewSimpleInMemoryStore()}
}

// Parses, analyzes and evaluates source, using current store.
func (m *MangleService) UpdateFromSource(reader io.Reader) error {
	u, err := parse.Unit(reader)
	if err != nil {
		return err
	}
	info, err := analysis.Analyze([]parse.SourceUnit{u}, nil)
	if err != nil {
		return err
	}
	programInfo = info
	strata, predToStratum, err := analysis.Stratify(analysis.Program{
		EdbPredicates: programInfo.EdbPredicates,
		IdbPredicates: programInfo.IdbPredicates,
		Rules:         programInfo.Rules,
	})
	if err != nil {
		return err
	}
	stats, err := engine.EvalStratifiedProgramWithStats(programInfo, strata, predToStratum, m.store)
	if err != nil {
		return err
	}
	log.Printf("service.go:UpdateFromSource: evaluation finished. stats: %v\n", stats)
	return nil
}

func (m *MangleService) Query(req *pb.QueryRequest, stream pb.Mangle_QueryServer) error {
	var store = m.store
	if program := req.GetProgram(); program != "" {
		u, err := parse.Unit(strings.NewReader(program))
		if err != nil {
			return err
		}
		programInfo, err := analysis.Analyze([]parse.SourceUnit{u}, copyDecl(programInfo.Decls))
		if err != nil {
			return err
		}
		store = factstore.NewTeeingStore(store)
		stats, err := engine.EvalProgramWithStats(programInfo, store)
		if err != nil {
			return err
		}
		log.Printf("service.go:Query evaluation of request program finished. stats: %v\n", stats)
		log.Printf("service.go:Query store predicates: %s\n", store.ListPredicates())
	}

	query := req.GetQuery()
	u, err := parse.Atom(query)
	if err != nil {
		log.Printf("service.go:Query parse %q (query) failed: %v\n", query, err)
		return err
	}

	log.Printf("querying store with query %v", u)
	err = store.GetFacts(u, func(a ast.Atom) error {
		answer := &pb.QueryAnswer{
			Answer: a.String(),
		}
		if err := stream.Send(answer); err != nil {
			log.Printf("service.go: got send err: %v", err)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
