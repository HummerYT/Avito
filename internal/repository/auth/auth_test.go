package auth_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/suite"

	"AvitoTask/internal/models"
	"AvitoTask/internal/repository/auth"
	"AvitoTask/internal/repository/auth/mocks"
)

type fakeRow struct {
	values []interface{}
	err    error
}

func (r *fakeRow) Scan(dest ...interface{}) error {
	if r.err != nil {
		return r.err
	}
	if len(dest) != len(r.values) {
		return fmt.Errorf("expected %d destination(s), got %d", len(r.values), len(dest))
	}
	for i, v := range r.values {
		switch d := dest[i].(type) {
		case *int64:
			strVal, ok := v.(int64)
			if !ok {
				return fmt.Errorf("expected string value at index %d", i)
			}
			*d = strVal
		case *int:
			switch val := v.(type) {
			case int:
				*d = val
			case int64:
				*d = int(val)
			default:
				return fmt.Errorf("expected int or int64 value at index %d", i)
			}
		case *string:
			strVal, ok := v.(string)
			if !ok {
				return fmt.Errorf("expected string value at index %d", i)
			}
			*d = strVal
		default:
			return fmt.Errorf("unsupported type for destination %d", i)
		}
	}
	return nil
}

// RepositoryTestSuite объединяет тесты для репозитория.
type RepositoryTestSuite struct {
	suite.Suite
	ctrl     *gomock.Controller
	mockPool *mocks.Mockpool
	repo     *auth.Repository
}

func (s *RepositoryTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockPool = mocks.NewMockpool(s.ctrl)
	s.repo = auth.NewInsertRepo(s.mockPool)
}

func (s *RepositoryTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *RepositoryTestSuite) TestIsUserExists_UserExists() {
	ctx := context.Background()
	user := models.User{Username: "testuser"}

	// Эмулируем, что запрос возвращает count = 1
	row := &fakeRow{
		values: []interface{}{1},
	}
	s.mockPool.
		EXPECT().
		QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE username = $1", user.Username).
		Return(row)

	exists, err := s.repo.IsUserExists(ctx, user)
	s.NoError(err)
	s.True(exists)
}

func (s *RepositoryTestSuite) TestIsUserExists_UserDoesNotExist() {
	ctx := context.Background()
	user := models.User{Username: "testuser"}

	// Эмулируем, что запрос возвращает count = 0
	row := &fakeRow{
		values: []interface{}{0},
	}
	s.mockPool.
		EXPECT().
		QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE username = $1", user.Username).
		Return(row)

	exists, err := s.repo.IsUserExists(ctx, user)
	s.NoError(err)
	s.False(exists)
}

func (s *RepositoryTestSuite) TestIsUserExists_QueryError() {
	ctx := context.Background()
	user := models.User{Username: "testuser"}

	row := &fakeRow{
		err: errors.New("scan error"),
	}
	s.mockPool.
		EXPECT().
		QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE username = $1", user.Username).
		Return(row)

	exists, err := s.repo.IsUserExists(ctx, user)
	s.Error(err)
	s.False(exists)
	s.Contains(err.Error(), "failed to check user existence")
}

func (s *RepositoryTestSuite) TestGetUserByLogin_Success() {
	ctx := context.Background()
	login := "testuser"

	id := uuid.New().String()
	row := &fakeRow{
		values: []interface{}{id, "testuser", "hashedpassword", int64(100)},
	}
	s.mockPool.
		EXPECT().
		QueryRow(ctx, "SELECT id, username, password, coins FROM users WHERE username = $1 LIMIT 1", login).
		Return(row)

	user, err := s.repo.GetUserByLogin(ctx, login)
	s.NoError(err)
	s.Equal(id, user.ID)
	s.Equal("testuser", user.Username)
	s.Equal("hashedpassword", user.Password)
	s.Equal(int64(100), user.Coins)
}

// TestGetUserByLogin: ошибка при сканировании результата
func (s *RepositoryTestSuite) TestGetUserByLogin_ScanError() {
	ctx := context.Background()
	login := "testuser"
	row := &fakeRow{
		err: errors.New("scan error"),
	}
	s.mockPool.
		EXPECT().
		QueryRow(ctx, "SELECT id, username, password, coins FROM users WHERE username = $1 LIMIT 1", login).
		Return(row)

	user, err := s.repo.GetUserByLogin(ctx, login)
	s.Error(err)
	s.Contains(err.Error(), "failed to scan user")
	s.Equal(models.User{}, user)
}

func (s *RepositoryTestSuite) TestInsertUser_Success() {
	ctx := context.Background()
	newUser := models.User{ID: "user-id-123", Username: "testuser", Password: "hashedpassword"}

	row := &fakeRow{
		values: []interface{}{newUser.ID},
	}

	s.mockPool.
		EXPECT().
		QueryRow(ctx, gomock.Any(), newUser.ID, newUser.Username, newUser.Password).
		DoAndReturn(func(ctx context.Context, query string, args ...any) pgx.Row {
			s.True(strings.Contains(query, "INSERT INTO users"))
			return row
		})

	id, err := s.repo.InsertUser(ctx, newUser)
	s.NoError(err)
	s.Equal(newUser.ID, id)
}

func (s *RepositoryTestSuite) TestInsertUser_Error() {
	ctx := context.Background()
	newUser := models.User{ID: "user-id-123", Username: "testuser", Password: "hashedpassword"}

	row := &fakeRow{
		err: errors.New("insert error"),
	}
	s.mockPool.
		EXPECT().
		QueryRow(ctx, gomock.Any(), newUser.ID, newUser.Username, newUser.Password).
		DoAndReturn(func(ctx context.Context, query string, args ...any) pgx.Row {
			s.True(strings.Contains(query, "INSERT INTO users"))
			return row
		})

	id, err := s.repo.InsertUser(ctx, newUser)
	s.Error(err)
	s.Equal("", id)
	s.Contains(err.Error(), "failed to insert user")
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

type fakeTx struct {
	queryRowFunc func(ctx context.Context, query string, args ...any) pgx.Row
	execFunc     func(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
}

func (ft *fakeTx) Begin(ctx context.Context) (pgx.Tx, error) {
	//TODO implement me
	panic("implement me")
}

func (ft *fakeTx) Commit(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (ft *fakeTx) Rollback(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (ft *fakeTx) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (ft *fakeTx) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	//TODO implement me
	panic("implement me")
}

func (ft *fakeTx) LargeObjects() pgx.LargeObjects {
	//TODO implement me
	panic("implement me")
}

func (ft *fakeTx) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	//TODO implement me
	panic("implement me")
}

func (ft *fakeTx) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	//TODO implement me
	panic("implement me")
}

func (ft *fakeTx) Conn() *pgx.Conn {
	//TODO implement me
	panic("implement me")
}

func (ft *fakeTx) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	if ft.queryRowFunc != nil {
		return ft.queryRowFunc(ctx, query, args...)
	}
	return nil
}

func (ft *fakeTx) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	if ft.execFunc != nil {
		return ft.execFunc(ctx, query, args...)
	}
	return pgconn.CommandTag{}, errors.New("exec not implemented")
}

type fakePool struct {
	beginFunc    func(ctx context.Context) (pgx.Tx, error)
	queryRowFunc func(ctx context.Context, sql string, args ...any) pgx.Row
}

func (fp *fakePool) Begin(ctx context.Context) (pgx.Tx, error) {
	if fp.beginFunc != nil {
		return fp.beginFunc(ctx)
	}
	return nil, errors.New("begin not implemented")
}

func (fp *fakePool) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	if fp.queryRowFunc != nil {
		return fp.queryRowFunc(ctx, sql, args...)
	}
	return nil
}

type TxTestSuite struct {
	suite.Suite
	repo     *auth.Repository
	fakePool *fakePool
}

func (s *TxTestSuite) SetupTest() {
	s.fakePool = &fakePool{}
	s.repo = auth.NewInsertRepo(s.fakePool)
}

func (s *TxTestSuite) TestBeginTx_Success() {
	ctx := context.Background()
	ftx := &fakeTx{}
	s.fakePool.beginFunc = func(ctx context.Context) (pgx.Tx, error) {
		return ftx, nil
	}

	tx, err := s.repo.BeginTx(ctx)
	s.NoError(err)
	s.Equal(ftx, tx)
}

func (s *TxTestSuite) TestBeginTx_Error() {
	ctx := context.Background()
	expectedErr := errors.New("begin error")
	s.fakePool.beginFunc = func(ctx context.Context) (pgx.Tx, error) {
		return nil, expectedErr
	}

	tx, err := s.repo.BeginTx(ctx)
	s.Error(err)
	s.Nil(tx)
	s.Equal(expectedErr, err)
}

func (s *TxTestSuite) TestGetUserById_Success() {
	ctx := context.Background()
	userID := "user-id-123"
	expectedUser := models.User{
		ID:       userID,
		Username: "testuser",
		Password: "hashedpassword",
		Coins:    100,
	}

	row := &fakeRow{
		values: []interface{}{expectedUser.ID, expectedUser.Username, expectedUser.Password, int64(expectedUser.Coins)},
	}
	tx := &fakeTx{
		queryRowFunc: func(ctx context.Context, query string, args ...any) pgx.Row {
			return row
		},
	}

	user, err := s.repo.GetUserById(ctx, tx, userID)
	s.NoError(err)
	s.Equal(expectedUser, user)
}

func (s *TxTestSuite) TestGetUserById_NoRows() {
	ctx := context.Background()
	userID := "user-id-123"

	row := &fakeRow{
		err: pgx.ErrNoRows,
	}
	tx := &fakeTx{
		queryRowFunc: func(ctx context.Context, query string, args ...any) pgx.Row {
			return row
		},
	}

	user, err := s.repo.GetUserById(ctx, tx, userID)
	s.Error(err)
	s.Equal(auth.ErrNoUserExist, err)
	s.Equal(models.User{}, user)
}

func (s *TxTestSuite) TestGetUserById_OtherError() {
	ctx := context.Background()
	userID := "user-id-123"

	row := &fakeRow{
		err: errors.New("scan error"),
	}
	tx := &fakeTx{
		queryRowFunc: func(ctx context.Context, query string, args ...any) pgx.Row {
			return row
		},
	}

	user, err := s.repo.GetUserById(ctx, tx, userID)
	s.Error(err)
	s.Contains(err.Error(), fmt.Sprintf("cannot find user '%s'", userID))
	s.Equal(models.User{}, user)
}

func (s *TxTestSuite) TestUpdateUserCoins_Success() {
	ctx := context.Background()
	userID := "user-id-123"
	newCoins := int64(200)

	tx := &fakeTx{
		execFunc: func(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
			return pgconn.CommandTag{}, nil
		},
	}

	err := s.repo.UpdateUserCoins(ctx, tx, userID, newCoins)
	s.NoError(err)
}

func (s *TxTestSuite) TestUpdateUserCoins_Error() {
	ctx := context.Background()
	userID := "user-id-123"
	newCoins := int64(200)
	expectedErr := errors.New("update error")

	tx := &fakeTx{
		execFunc: func(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
			return pgconn.CommandTag{}, expectedErr
		},
	}

	err := s.repo.UpdateUserCoins(ctx, tx, userID, newCoins)
	s.Error(err)
	s.Contains(err.Error(), fmt.Sprintf("failed to update coins for userID=%s", userID))
}

func (s *TxTestSuite) TestGetUserCoins_Success() {
	ctx := context.Background()
	userID := "user-id-123"
	expectedCoins := int64(300)

	row := &fakeRow{
		values: []interface{}{expectedCoins},
	}
	tx := &fakeTx{
		queryRowFunc: func(ctx context.Context, query string, args ...any) pgx.Row {
			return row
		},
	}

	coins, err := s.repo.GetUserCoins(ctx, tx, userID)
	s.NoError(err)
	s.Equal(expectedCoins, coins)
}

func (s *TxTestSuite) TestGetUserCoins_Error() {
	ctx := context.Background()
	userID := "user-id-123"
	expectedErr := errors.New("scan error")

	row := &fakeRow{
		err: expectedErr,
	}
	tx := &fakeTx{
		queryRowFunc: func(ctx context.Context, query string, args ...any) pgx.Row {
			return row
		},
	}

	coins, err := s.repo.GetUserCoins(ctx, tx, userID)
	s.Error(err)
	s.Equal(int64(0), coins)
	s.Contains(err.Error(), fmt.Sprintf("failed to get user coins (userID=%s)", userID))
}

func TestTxTestSuite(t *testing.T) {
	suite.Run(t, new(TxTestSuite))
}
