package postgres_test

import (
	"Hermes/internal/config"
	"Hermes/internal/errs"
	"Hermes/internal/logger"
	"Hermes/internal/models"
	"Hermes/internal/repository/postgres"
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/lib/pq"
	wbf "github.com/wb-go/wbf/config"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/retry"
)

var testStorage *postgres.Storage

func TestMain(m *testing.M) {

	if err := wbf.New().LoadEnvFiles("../../../.env"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cfg := config.Storage{
		Host:     "postgres-test",
		Port:     "5432",
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   "hermes_test",
		SSLMode:  "disable",
		QueryRetryStrategy: config.RetryStrategy{
			Attempts: 3,
			Delay:    100 * time.Millisecond,
			Backoff:  1.5,
		},
	}

	logger, _ := logger.NewLogger(config.Logger{Debug: true})

	db, err := dbpg.New(fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode), nil, &dbpg.Options{})
	if err != nil {
		logger.LogFatal("postgres_test — failed to connect to test DB", err, "layer", "repository.postgres_test")
	}

	testStorage = postgres.NewStorage(logger, cfg, db)

	exitCode := m.Run()
	testStorage.Close()
	os.Exit(exitCode)

}

func setupTest(t *testing.T) {

	ctx := context.Background()
	_, err := testStorage.DB().ExecWithRetry(ctx, retry.Strategy{Attempts: 3, Delay: 100 * time.Millisecond, Backoff: 1.5}, `
	
	TRUNCATE TABLE comments 
	RESTART IDENTITY`)

	if err != nil {
		t.Fatalf("failed to truncate comments: %v", err)
	}

}

func TestCreateComment_Errors(t *testing.T) {

	setupTest(t)

	comment := models.Comment{ParentID: ptr(int64(999999999)), Content: "Invalid parent", Author: "test"}

	_, err := testStorage.CreateComment(context.Background(), comment)
	if err == nil {
		t.Fatalf("expected error for invalid parent, got nil")
	}

	var pqErr *pq.Error
	if errors.As(err, &pqErr) && string(pqErr.Code) == "23503" {
		t.Logf("expected foreign key error captured: %v", err)
	} else {
		t.Fatalf("unexpected error: %v", err)
	}

}

func TestDeleteComment(t *testing.T) {

	setupTest(t)

	ctx := context.Background()
	comment := models.Comment{Content: "To be deleted", Author: "test"}

	id, err := testStorage.CreateComment(ctx, comment)
	if err != nil {
		t.Fatalf("CreateComment failed: %v", err)
	}

	if err := testStorage.DeleteComment(ctx, id); err != nil {
		t.Fatalf("DeleteComment failed: %v", err)
	}

	var count int
	err = testStorage.DB().Master.QueryRowContext(ctx, `
	
	SELECT COUNT(*) 
	FROM comments 
	WHERE id=$1`, id).Scan(&count)

	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if count != 0 {
		t.Fatalf("expected 0 comments, got %d", count)
	}

	err = testStorage.DeleteComment(ctx, 999999)
	if err != errs.ErrCommentNotFound {
		t.Fatalf("expected ErrCommentNotFound, got %v", err)
	}

}

func TestGetCommentTree(t *testing.T) {

	setupTest(t)

	ctx := context.Background()
	root := models.Comment{Content: "Root", Author: "test"}

	rootID, err := testStorage.CreateComment(ctx, root)
	if err != nil {
		t.Fatalf("CreateComment failed: %v", err)
	}

	time.Sleep(100 * time.Millisecond)
	child := models.Comment{ParentID: &rootID, Content: "Child", Author: "test"}

	childID, err := testStorage.CreateComment(ctx, child)
	if err != nil {
		t.Fatalf("CreateComment failed: %v", err)
	}

	time.Sleep(100 * time.Millisecond)
	grandchild := models.Comment{ParentID: &childID, Content: "Grandchild", Author: "test"}

	grandchildID, err := testStorage.CreateComment(ctx, grandchild)
	if err != nil {
		t.Fatalf("CreateComment failed: %v", err)
	}

	tree, err := testStorage.GetCommentTree(ctx, rootID)
	if err != nil {
		t.Fatalf("GetCommentTree failed: %v", err)
	}

	if len(tree) != 3 {
		t.Fatalf("expected 3 comments, got %d", len(tree))
	}

	if tree[0].ID != rootID || tree[1].ID != childID || tree[2].ID != grandchildID {
		t.Fatalf("unexpected order or IDs: %+v", tree)
	}

	tree, err = testStorage.GetCommentTree(ctx, 999999)
	if err != nil {
		t.Fatalf("GetCommentTree failed: %v", err)
	}

	if len(tree) != 0 {
		t.Fatalf("expected 0 comments for non-existent, got %d", len(tree))
	}

}

func TestGetRootComments(t *testing.T) {

	setupTest(t)

	ctx := context.Background()
	root1 := models.Comment{Content: "Root1", Author: "test"}

	root1ID, err := testStorage.CreateComment(ctx, root1)
	if err != nil {
		t.Fatalf("CreateComment failed: %v", err)
	}

	time.Sleep(100 * time.Millisecond)
	root2 := models.Comment{Content: "Root2", Author: "test"}

	root2ID, err := testStorage.CreateComment(ctx, root2)
	if err != nil {
		t.Fatalf("CreateComment failed: %v", err)
	}

	time.Sleep(100 * time.Millisecond)
	root3 := models.Comment{Content: "Root3", Author: "test"}

	root3ID, err := testStorage.CreateComment(ctx, root3)
	if err != nil {
		t.Fatalf("CreateComment failed: %v", err)
	}

	child := models.Comment{ParentID: &root1ID, Content: "Child", Author: "test"}

	_, err = testStorage.CreateComment(ctx, child)
	if err != nil {
		t.Fatalf("CreateComment failed: %v", err)
	}

	params := models.QueryParams{Limit: 2, Offset: 0, Sort: ""}

	roots, err := testStorage.GetRootComments(ctx, params)
	if err != nil {
		t.Fatalf("GetRootComments failed: %v", err)
	}

	if len(roots) != 2 || roots[0].ID != root3ID || roots[1].ID != root2ID {
		t.Fatalf("unexpected roots (desc): %+v", roots)
	}

	params.Sort = "created_at_asc"

	roots, err = testStorage.GetRootComments(ctx, params)
	if err != nil {
		t.Fatalf("GetRootComments failed: %v", err)
	}

	if len(roots) != 2 || roots[0].ID != root1ID || roots[1].ID != root2ID {
		t.Fatalf("unexpected roots (asc): %+v", roots)
	}

	parentID := root2ID
	params.ParentID = &parentID
	roots, err = testStorage.GetRootComments(ctx, params)
	if err != nil {
		t.Fatalf("GetRootComments failed: %v", err)
	}

	if len(roots) != 1 || roots[0].ID != root2ID {
		t.Fatalf("unexpected specific: %+v", roots)
	}

	parentID = 999999
	params.ParentID = &parentID

	roots, err = testStorage.GetRootComments(ctx, params)
	if err != nil {
		t.Fatalf("GetRootComments failed: %v", err)
	}

	if len(roots) != 0 {
		t.Fatalf("expected 0 for non-existent")
	}

	_, err = testStorage.DB().QueryWithRetry(ctx, retry.Strategy{Attempts: 1, Delay: 1, Backoff: 1}, `
	
	SELECT * 
	FROM NonExistentTable;`)

	if err != nil {
		t.Logf("SQL error captured as expected: %v", err)
	} else {
		t.Fatalf("expected SQL error, got nil")
	}

}

func TestClose(t *testing.T) {
	log, _ := logger.NewLogger(config.Logger{Debug: true})
	db, _ := dbpg.New(fmt.Sprintf("host=postgres-test port=5432 user=%s password=%s dbname=hermes_test sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD")), nil, &dbpg.Options{})
	st := postgres.NewStorage(log, config.Storage{}, db)
	st.Close()
}

func ptr(i int64) *int64 {
	return &i
}
