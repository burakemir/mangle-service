package service

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/mangle/analysis"
	"github.com/google/mangle/ast"
	"github.com/google/mangle/engine"
	"github.com/google/mangle/factstore"
	"github.com/google/mangle/parse"

	pb "github.com/burakemir/mangle-service/proto"
)

var source = flag.String("source", "", "path to source to evaluate")

type MangleService struct {
	store factstore.ReadOnlyFactStore
	pb.UnimplementedMangleServer
}

func New() (*MangleService, error) {
	if *source == "" {
		return nil, errors.New("no --source given")
	}
	log.Printf("read source from %q", *source)
	sourceBytes, err := os.ReadFile(*source)
	if err != nil {
		return nil, err
	}
	u, err := parse.Unit(strings.NewReader(string(sourceBytes)))
	if err != nil {
		return nil, err
	}
	programInfo, err := analysis.Analyze([]parse.SourceUnit{u}, nil)
	if err != nil {
		return nil, err
	}
	strata, predToStratum, err := analysis.Stratify(analysis.Program{
		EdbPredicates: programInfo.EdbPredicates,
		IdbPredicates: programInfo.IdbPredicates,
		Rules:         programInfo.Rules,
	})
	if err != nil {
		return nil, err
	}
	store := factstore.NewSimpleInMemoryStore()
	stats, err := engine.EvalStratifiedProgramWithStats(programInfo, strata, predToStratum, store)
	if err != nil {
		return nil, err
	}
	log.Printf("evaluation finished. stats: %v\n", stats)
	return &MangleService{store: store}, nil
}

func (m *MangleService) Query(req *pb.QueryRequest, stream pb.Mangle_QueryServer) error {
	program := req.GetProgram()
	var p *parse.SourceUnit
	if program != "" {
		u, err := parse.Unit(strings.NewReader(program))
		if err != nil {
			return err
		}
		p = &u
	}
	fmt.Println("program: %v", p)

	query := req.GetQuery()
	u, err := parse.Atom(query)
	if err != nil {
		return err
	}

	fmt.Println("querying store with query %v", u)
	err = m.store.GetFacts(u, func(a ast.Atom) error {
		answer := &pb.QueryAnswer{
			Answer: a.String(),
		}
		if err := stream.Send(answer); err != nil {
			fmt.Println("got send err: %v", err)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
