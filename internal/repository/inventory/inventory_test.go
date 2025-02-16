package inventory_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/suite"

	"AvitoTask/internal/repository/inventory"
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
			switch val := v.(type) {
			case int:
				*d = int64(val)
			case int64:
				*d = val
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
			return fmt.Errorf("unsupported destination type at index %d", i)
		}
	}
	return nil
}

type fakeTx struct {
	queryRowFunc func(ctx context.Context, query string, args ...any) pgx.Row
	execFunc     func(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
	queryFunc    func(ctx context.Context, query string, args ...any) (pgx.Rows, error)
}

func (ft *fakeTx) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	//TODO implement me
	panic("implement me")
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

func (ft *fakeTx) Exec(ctx context.Context, query string, args ...any) (commandTag pgconn.CommandTag, err error) {
	if ft.execFunc != nil {
		return ft.execFunc(ctx, query, args...)
	}
	return pgconn.CommandTag{}, errors.New("Exec not implemented")
}

func (ft *fakeTx) Conn() *pgx.Conn {
	//TODO implement me
	panic("implement me")
}

func (ft *fakeTx) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	if ft.queryRowFunc != nil {
		return ft.queryRowFunc(ctx, query, args...)
	}
	return &fakeRow{err: errors.New("QueryRow not implemented")}
}

func (ft *fakeTx) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	if ft.queryFunc != nil {
		return ft.queryFunc(ctx, query, args...)
	}
	return nil, errors.New("Query not implemented")
}

type fakeRows struct {
	data [][]interface{}
	idx  int
	err  error
}

func (r *fakeRows) CommandTag() pgconn.CommandTag {
	//TODO implement me
	panic("implement me")
}

func (r *fakeRows) FieldDescriptions() []pgconn.FieldDescription {
	//TODO implement me
	panic("implement me")
}

func (r *fakeRows) Values() ([]any, error) {
	//TODO implement me
	panic("implement me")
}

func (r *fakeRows) RawValues() [][]byte {
	//TODO implement me
	panic("implement me")
}

func (r *fakeRows) Conn() *pgx.Conn {
	//TODO implement me
	panic("implement me")
}

func (r *fakeRows) Next() bool {
	return r.idx < len(r.data)
}

func (r *fakeRows) Scan(dest ...interface{}) error {
	// Если количество столбцов не соответствует ожидаемому, вернём ошибку
	if r.idx >= len(r.data) {
		return errors.New("no more rows")
	}
	row := r.data[r.idx]
	if len(dest) != len(row) {
		return fmt.Errorf("expected %d destination(s), got %d", len(row), len(dest))
	}
	for i, v := range row {
		switch d := dest[i].(type) {
		case *string:
			strVal, ok := v.(string)
			if !ok {
				return fmt.Errorf("expected string for column %d", i)
			}
			*d = strVal
		case *int64:
			switch val := v.(type) {
			case int:
				*d = int64(val)
			case int64:
				*d = val
			default:
				return fmt.Errorf("expected int or int64 for column %d", i)
			}
		default:
			return fmt.Errorf("unsupported destination type at index %d", i)
		}
	}
	r.idx++
	return nil
}

func (r *fakeRows) Close() {
	// no-op
}

func (r *fakeRows) Err() error {
	return r.err
}

// ---------------- Тестовый сьют для репозитория inventory ----------------

type InventoryRepoTestSuite struct {
	suite.Suite
	// Для тестов, использующих транзакцию, поле pool не требуется,
	// поэтому можно передать nil при создании репозитория.
	repo *inventory.Repository
}

func (s *InventoryRepoTestSuite) SetupTest() {
	// Если функции не используют r.pool (а работают через переданный tx),
	// можно передать nil.
	s.repo = inventory.NewInsertRepo(nil)
}

// Тест для GetInventoryItem: успешное получение количества предметов.
func (s *InventoryRepoTestSuite) TestGetInventoryItem_Success() {
	ctx := context.Background()
	userID := "user-123"
	itemType := "potion"
	expectedQuantity := int64(5)

	row := &fakeRow{
		values: []interface{}{expectedQuantity},
	}
	tx := &fakeTx{
		queryRowFunc: func(ctx context.Context, query string, args ...any) pgx.Row {
			s.Contains(query, "SELECT quantity")
			s.Equal(userID, args[0])
			s.Equal(itemType, args[1])
			return row
		},
	}

	quantity, err := s.repo.GetInventoryItem(ctx, tx, userID, itemType)
	s.NoError(err)
	s.Equal(expectedQuantity, quantity)
}

// Тест для GetInventoryItem: ошибка при выполнении запроса (например, ошибка сканирования).
func (s *InventoryRepoTestSuite) TestGetInventoryItem_Error() {
	ctx := context.Background()
	userID := "user-123"
	itemType := "potion"
	expectedErr := errors.New("query error")

	row := &fakeRow{
		err: expectedErr,
	}
	tx := &fakeTx{
		queryRowFunc: func(ctx context.Context, query string, args ...any) pgx.Row {
			return row
		},
	}

	quantity, err := s.repo.GetInventoryItem(ctx, tx, userID, itemType)
	s.Error(err)
	s.Equal(int64(0), quantity)
	s.Equal(expectedErr, err)
}

// Тест для InsertInventoryItem: успешная вставка нового предмета.
func (s *InventoryRepoTestSuite) TestInsertInventoryItem_Success() {
	ctx := context.Background()
	id := "item-123"
	userID := "user-123"
	itemType := "potion"

	tx := &fakeTx{
		execFunc: func(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
			s.Contains(query, "INSERT INTO inventory")
			s.Equal(id, args[0])
			s.Equal(userID, args[1])
			s.Equal(itemType, args[2])
			// Имитация успешной вставки.
			return pgconn.CommandTag{}, nil
		},
	}

	err := s.repo.InsertInventoryItem(ctx, tx, id, userID, itemType)
	s.NoError(err)
}

// Тест для InsertInventoryItem: ошибка при вставке.
func (s *InventoryRepoTestSuite) TestInsertInventoryItem_Error() {
	ctx := context.Background()
	id := "item-123"
	userID := "user-123"
	itemType := "potion"
	expectedErr := errors.New("exec error")

	tx := &fakeTx{
		execFunc: func(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
			return pgconn.CommandTag{}, expectedErr
		},
	}

	err := s.repo.InsertInventoryItem(ctx, tx, id, userID, itemType)
	s.Error(err)
	s.Contains(err.Error(), fmt.Sprintf("failed to insert new item '%s' for user %s", itemType, userID))
}

// Тест для UpdateInventoryItem: успешное обновление количества предмета.
func (s *InventoryRepoTestSuite) TestUpdateInventoryItem_Success() {
	ctx := context.Background()
	userID := "user-123"
	itemType := "potion"
	newQuantity := int64(10)

	tx := &fakeTx{
		execFunc: func(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
			s.Contains(query, "UPDATE inventory")
			s.Equal(newQuantity, args[0])
			s.Equal(userID, args[1])
			s.Equal(itemType, args[2])
			return pgconn.CommandTag{}, nil
		},
	}

	err := s.repo.UpdateInventoryItem(ctx, tx, userID, itemType, newQuantity)
	s.NoError(err)
}

// Тест для UpdateInventoryItem: ошибка при обновлении.
func (s *InventoryRepoTestSuite) TestUpdateInventoryItem_Error() {
	ctx := context.Background()
	userID := "user-123"
	itemType := "potion"
	newQuantity := int64(10)
	expectedErr := errors.New("update error")

	tx := &fakeTx{
		execFunc: func(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
			return pgconn.CommandTag{}, expectedErr
		},
	}

	err := s.repo.UpdateInventoryItem(ctx, tx, userID, itemType, newQuantity)
	s.Error(err)
	s.Contains(err.Error(), fmt.Sprintf("failed to update item '%s' for user %s", itemType, userID))
}

func (s *InventoryRepoTestSuite) TestGetUserInventory_Success() {
	ctx := context.Background()
	userID := "user-123"

	fRows := &fakeRows{
		data: [][]interface{}{
			{"potion", int64(5)},
			{"elixir", int64(3)},
		},
		idx: 0,
		err: nil,
	}

	tx := &fakeTx{
		queryFunc: func(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
			s.Contains(query, `
        SELECT item_type, quantity
        FROM inventory
        WHERE user_id = $1
    `)
			s.Equal(userID, args[0])
			return fRows, nil
		},
	}

	items, err := s.repo.GetUserInventory(ctx, tx, userID)
	s.NoError(err)
	s.Len(items, 2)
	s.Equal("potion", items[0].ItemType)
	s.Equal(int64(5), items[0].Quantity)
	s.Equal("elixir", items[1].ItemType)
	s.Equal(int64(3), items[1].Quantity)
}

// Тест для GetUserInventory: ошибка запроса (например, ошибка выполнения Query).
func (s *InventoryRepoTestSuite) TestGetUserInventory_QueryError() {
	ctx := context.Background()
	userID := "user-123"
	expectedErr := errors.New("query error")

	tx := &fakeTx{
		queryFunc: func(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
			return nil, expectedErr
		},
	}

	items, err := s.repo.GetUserInventory(ctx, tx, userID)
	s.Error(err)
	s.Nil(items)
	s.Contains(err.Error(), "failed to query inventory")
}

// Тест для GetUserInventory: ошибка сканирования (например, неверное число столбцов).
func (s *InventoryRepoTestSuite) TestGetUserInventory_ScanError() {
	ctx := context.Background()
	userID := "user-123"

	// Передадим строку с недостаточным количеством столбцов.
	fRows := &fakeRows{
		data: [][]interface{}{
			{"potion"}, // отсутствует колонка quantity
		},
		idx: 0,
		err: nil,
	}

	tx := &fakeTx{
		queryFunc: func(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
			return fRows, nil
		},
	}

	items, err := s.repo.GetUserInventory(ctx, tx, userID)
	s.Error(err)
	s.Nil(items)
	s.Contains(err.Error(), "expected")
}

// Тест для GetUserInventory: ошибка во время итерации (rows.Err() возвращает ошибку).
func (s *InventoryRepoTestSuite) TestGetUserInventory_RowsErr() {
	ctx := context.Background()
	userID := "user-123"
	expectedErr := errors.New("rows error")
	fRows := &fakeRows{
		data: [][]interface{}{
			{"potion", int64(5)},
		},
		idx: 0,
		err: expectedErr,
	}

	tx := &fakeTx{
		queryFunc: func(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
			return fRows, nil
		},
	}

	items, err := s.repo.GetUserInventory(ctx, tx, userID)
	s.Error(err)
	s.Nil(items)
	s.Contains(err.Error(), "error during rows iteration")
}

func TestInventoryRepoTestSuite(t *testing.T) {
	suite.Run(t, new(InventoryRepoTestSuite))
}
